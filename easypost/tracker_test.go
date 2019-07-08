// Copyright 2019 RetailNext, Inc.
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
	"reflect"
	"strings"
	"testing"
)

var (
	testServer *httptest.Server
	testClient = NewClient("")
)

func setup() {
	m := http.NewServeMux()
	m.HandleFunc("/trackers", getTestTrackers)
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
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error: ErrorMessage{
					Message: "not found",
					FieldErrors: []FieldError{
						{
							Field:   "tracking_code",
							Message: "not found",
						},
					},
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

func TestCreateTracker(t *testing.T) {
	setup()

	trackingCode := "EZ3000000003"
	b, err := readTestTrackerFile(trackingCode)
	if err != nil {
		t.Fatalf("error reading file: %s", err)
	}

	expectedTracker := Tracker{}
	if err := json.Unmarshal(b, &expectedTracker); err != nil {
		t.Fatalf("expected tracker build error: %s", err)
	}

	gotTracker, err := testClient.GetTracker(trackingCode, "")
	if err != nil {
		t.Fatalf("not success response: %s", err)
	}
	if !reflect.DeepEqual(gotTracker, &expectedTracker) {
		t.Fatalf("trackers: \nexpected %+v\n     got %+v", &expectedTracker, gotTracker)
	}

	expectedError := ProcessingError{
		msg:     "not found",
		details: map[string]string{"tracking_code": "not found"},
	}
	_, err = testClient.GetTracker("EZ3000000002", "")
	if !reflect.DeepEqual(expectedError, err) {
		t.Fatalf("error:\nexpected: %v \ngot: %v", expectedError, err)
	}

	_, err = testClient.GetTracker(paymentError.Error(), "")
	if _, ok := err.(PaymentRequiredError); !ok {
		t.Fatalf("payment error expected: %T (%s)", err, err)
	}

	_, err = testClient.GetTracker(unauthorizedError.Error(), "")
	if _, ok := err.(UnauthorizedError); !ok {
		t.Fatalf("unauthorized error expected: %T (%s)", err, err)
	}
}
