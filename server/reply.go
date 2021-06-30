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
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func replyOK(w http.ResponseWriter, body []byte, mimeType string) {
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		log.Printf("i/o error: %v", err) // Already sent StatusOK, so just log.
	}
}

func replyServerError(w http.ResponseWriter, err error) {
	log.Print(err)
	replyError(w, http.StatusInternalServerError)
}

func replyClientError(w http.ResponseWriter, err error) {
	log.Print(err)
	replyError(w, http.StatusBadRequest)
}

func replyClientErrorSilent(w http.ResponseWriter) {
	replyError(w, http.StatusBadRequest)
}

func replyError(w http.ResponseWriter, code int) {
	http.Error(w, fmt.Sprintf("%d %s", code, http.StatusText(code)), code)
}
