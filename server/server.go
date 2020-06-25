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
	"net"
	"net/http"
)

// Server encapsulates http.Server and Config so it can start and stop
// CertManager automatically in Serve.
type Server struct {
	*http.Server
	Config
}

// NewServer creates a new Server. s.Handler is replaced with NewHandler(c).
func NewServer(s *http.Server, c Config) *Server {
	s.Handler = NewHandler(c)
	return &Server{s, c}
}

// ListenAndServe wraps s.Server.ListenAndServe to start/stop s.CertManager
// automatically.
func (s *Server) ListenAndServe() error {
	if err := s.CertManager.Start(); err != nil {
		return err
	}
	defer s.CertManager.Stop()
	return s.Server.ListenAndServe()
}

// ListenAndServeTLS wraps s.Server.ListenAndServeTLS to start/stop
// s.CertManager automatically.
func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if err := s.CertManager.Start(); err != nil {
		return err
	}
	defer s.CertManager.Stop()
	return s.Server.ListenAndServeTLS(certFile, keyFile)
}

// Serve wraps s.Server.Serve to start/stop s.CertManager automatically.
func (s *Server) Serve(l net.Listener) error {
	if err := s.CertManager.Start(); err != nil {
		return err
	}
	defer s.CertManager.Stop()
	return s.Server.Serve(l)
}

// ServeTLS wraps s.Server.ServeTLS to start/stop s.CertManager automatically.
func (s *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {
	if err := s.CertManager.Start(); err != nil {
		return err
	}
	defer s.CertManager.Stop()
	return s.Server.ServeTLS(l, certFile, keyFile)
}
