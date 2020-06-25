// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/google/webpackager"
	"github.com/google/webpackager/certchain/certmanager"
	"github.com/google/webpackager/fetch"
	"github.com/google/webpackager/internal/timeutil"
	"github.com/google/webpackager/internal/urlutil"
	"github.com/google/webpackager/processor/preverify"
	"github.com/google/webpackager/server/tomlconfig"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/xerrors"
)

const (
	mimeTypeCertChain = "application/cert-chain+cbor"
	mimeTypeExchange  = "application/signed-exchange"
	mimeTypeValidity  = "application/cbor"
)

var (
	emptyMapCBOR = []byte{0xa0}
)

// Handler handles HTTP requests. See the package GoDoc for details.
type Handler struct {
	mux *http.ServeMux
	Config
}

var _ http.Handler = (*Handler)(nil)

// Config holds the parameters to NewHandler.
type Config struct {
	// Packager is used to produce signed exchanges. ExchangeFactory should
	// be an ExchangeMetaFactory set with CertManager (the following field)
	// to keep the signing certificate and the cert-url parameter consistent
	// with this Handler.
	Packager *webpackager.Packager

	// CertManager provides the AugmentedChain to serve from this Handler.
	CertManager *certmanager.Manager

	// ServerConfig specifies the endpoints. All fields must contain a valid
	// value as described in cmd/webpkgserver/webpkgserver.example.toml.
	tomlconfig.ServerConfig
}

// NewHandler creates and initializes a new Handler.
func NewHandler(c Config) *Handler {
	// Remove the trailing slash.
	c.DocPath = path.Clean(c.DocPath)
	c.CertPath = path.Clean(c.CertPath)
	c.ValidityPath = path.Clean(c.ValidityPath)

	h := &Handler{new(http.ServeMux), c}

	h.mux.HandleFunc(c.CertPath+"/", h.handleCert)
	h.mux.HandleFunc(c.DocPath, h.handleDoc)
	h.mux.HandleFunc(c.ValidityPath, h.handleValidity)

	return h
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// All handlers assume GET requests.
	if req.Method != http.MethodGet {
		replyError(w, http.StatusMethodNotAllowed)
		return
	}

	path := req.URL.EscapedPath()

	// http.ServeMux normalizes the URL path and causes multiple issues:
	//   - "https://..." is reduced to "https:/...".
	//   - ".." can be used to replace the authority
	//     (e.g. "/priv/doc/https://www.example.com/../bad.example.com/").
	// To work around it, we handle this case specially.
	if url := strings.TrimPrefix(path, h.DocPath+"/"); len(url) < len(path) {
		if req.URL.RawQuery != "" {
			url += "?" + req.URL.RawQuery
		}
		h.handleDocImpl(w, req, url)
	} else {
		h.mux.ServeHTTP(w, req)
	}
}

func (h *Handler) handleCert(w http.ResponseWriter, req *http.Request) {
	ac := h.CertManager.GetAugmentedChain()

	// TODO(yuizumi): Keep the previous AugmentedChains for a while; serve
	// them when requested (through past signed exchanges).
	if req.URL.Path != path.Join(h.CertPath, ac.Digest) {
		replyError(w, http.StatusNotFound)
		return
	}
	var body bytes.Buffer
	if err := ac.WriteCBOR(&body); err != nil {
		replyServerError(w, xerrors.Errorf("serializing cert-chain: %w", err))
		return
	}
	replyOK(w, body.Bytes(), mimeTypeCertChain)
}

func (h *Handler) handleDoc(w http.ResponseWriter, req *http.Request) {
	h.handleDocImpl(w, req, req.URL.Query().Get(h.SignParam))
}

func (h *Handler) handleDocImpl(w http.ResponseWriter, req *http.Request, signURL string) {
	if err := verifyAcceptHeader(req); err != nil {
		replyClientError(w, err)
		return
	}
	u, err := parseSignURL(signURL)
	if err != nil {
		replyClientError(w, xerrors.Errorf("invalid sign url: %w", err))
		return
	}
	// TODO(yuizumi): Copy some request headers from req.
	newReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		replyServerError(w, err)
		return
	}
	if err := h.Packager.RunForRequest(newReq, timeutil.Now()); err != nil {
		err = filterError(err, u.String())
		// TODO(banaag): ideally, we should pass through that error response
		// from the upstream.
		if httpErr, ok := err.(*preverify.HTTPStatusError); ok {
			replyError(w, httpErr.StatusCode)
			return
		}
		if xerrors.Is(err, fetch.ErrURLMismatch) {
			replyClientError(w, err)
			return
		}
		if err != nil {
			replyServerError(w, xerrors.Errorf("Packager.RunForRequest: %w", err))
			return
		}
	}
	r, err := h.Packager.ResourceCache.Lookup(newReq)
	if err != nil {
		replyServerError(w, xerrors.Errorf("ResourceCache.Lookup: %w", err))
		return
	}
	var body bytes.Buffer
	if err := r.Exchange.Write(&body); err != nil {
		replyServerError(w, xerrors.Errorf("serializing exchange: %w", err))
		return
	}
	replyOK(w, body.Bytes(), r.Exchange.Version.MimeType())
}

func (h *Handler) handleValidity(w http.ResponseWriter, req *http.Request) {
	replyOK(w, emptyMapCBOR, mimeTypeValidity)
}

func filterError(err error, url string) error {
	switch err := err.(type) {
	case *webpackager.Error:
		if err.URL.String() != url {
			return nil
		}
		return err

	case *multierror.Error:
		var errs *multierror.Error
		for _, e := range err.Errors {
			errs = multierror.Append(errs, filterError(e, url))
		}
		if len(errs.Errors) == 1 {
			return errs.Errors[0]
		}
		return errs.ErrorOrNil()

	default:
		return err // TODO(yuizumi): Should this be nil?
	}
}

func parseSignURL(rawurl string) (*url.URL, error) {
	if rawurl == "" {
		return nil, errors.New("must be non-empty")
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, errors.New("must start with https://")
	}
	if u.User != nil {
		return nil, errors.New("must not have user:pass@")
	}
	if u.Fragment != "" {
		return nil, errors.New("must not have #fragment")
	}

	// Prevent malformed URLs from eluding the PathRE protections.
	u.Path = urlutil.GetCleanPath(u)
	// Escape special characters in the query component such as "<" or "|"
	// (but not "&" or "=").
	u.RawQuery = url.PathEscape(u.RawQuery)

	return u, nil
}

func verifyAcceptHeader(req *http.Request) error {
	// TODO(yuizumi): Parse the Accept header properly. If SXG has a lower
	// q value (say, than "text/html"), behave like a reverse proxy.
	// For now, we just verify it contains application/signed-exchange and
	// and reject the request otherwise for minimal safety.
	for _, v := range req.Header["Accept"] {
		if strings.Contains(v, mimeTypeExchange) {
			return nil
		}
	}
	return xerrors.Errorf("Accept header missing %q", mimeTypeExchange)
}
