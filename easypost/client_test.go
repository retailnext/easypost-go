// Copyright 2020 RetailNext, Inc.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package easypost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
)

var (
	testServer *httptest.Server
	testClient = NewClient("")
)

func setup() {
	m := http.NewServeMux()
	m.HandleFunc("/trackers", getTestTrackers)
	m.HandleFunc("/addresses", validateTestAddress)
	testServer = httptest.NewServer(m)
	apiURL = testServer.URL
}

func readTestTrackerFile(trackingCode string) ([]byte, error) {
	f, err := os.Open(path.Join("./test/trackers", fmt.Sprintf("%s.json", strings.ToUpper(trackingCode))))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func getTestTrackers(w http.ResponseWriter, r *http.Request) {
	trackingCode := r.FormValue("tracker[tracking_code]")
	switch trackingCode {
	case paymentError.Error():
		w.WriteHeader(http.StatusPaymentRequired)
		return
	case unauthorizedError.Error():
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	b, err := readTestTrackerFile(trackingCode)
	if err != nil {
		if os.IsNotExist(err) {
			b, err := json.Marshal([]FieldError{
				{
					Field:   "tracking_code",
					Message: "not found",
				},
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: errorMessage{
					Message:     "not found",
					FieldErrors: b,
				},
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func validateTestAddress(w http.ResponseWriter, r *http.Request) {
	var (
		addressFileName string
		responseCode    int
	)
	streetOne := r.FormValue("address[street1]")
	switch streetOne {
	case "Valid Street Name":
		addressFileName = "valid_address.json"
		responseCode = http.StatusOK
	default:
		addressFileName = "invalid_address.json"
		responseCode = http.StatusUnprocessableEntity
	}

	f, err := os.Open(path.Join("./test/addresses", addressFileName))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(responseCode)
	w.Write(b)
}
