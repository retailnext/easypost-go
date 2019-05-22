// Copyright 2019 RetailNext, Inc.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package easypost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewWebHookHandler(t *testing.T) {
	setup()
	username := "user"
	secret := "secret"
	webHookHandler := NewWebHookHandler(username, secret)
	m := http.NewServeMux()
	var (
		testEvent *Event
		testError error
	)
	m.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		testEvent, testError = webHookHandler(r)
	})
	testServer = httptest.NewServer(m)
	c := &http.Client{}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/webhook", testServer.URL), nil)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	request.SetBasicAuth("user", "secret1")
	_, err = c.Do(request)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if testError != unauthorizedError {
		t.Fatalf("unexpected error: %T (%s)", testError, testError)
	}

	request, err = http.NewRequest("POST", fmt.Sprintf("%s/webhook", testServer.URL), nil)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	_, err = c.Do(request)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if testError != unauthorizedError {
		t.Fatalf("unexpected error: %T (%s)", testError, testError)
	}

	trackerBody, err := readTestTrackerFile(TestTrackerCodes[2])
	if err != nil {
		t.Fatalf("error reading tracking: %s", err)
	}

	event := Event{
		Object: RecordTypeEvent,
		Result: trackerBody,
	}

	b, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	request, err = http.NewRequest("POST", fmt.Sprintf("%s/webhook", testServer.URL), bytes.NewReader(b))
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	request.SetBasicAuth(username, secret)
	_, err = c.Do(request)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if testError != nil {
		t.Fatalf("error: %s", err)
	}

	if testEvent == nil {
		t.Fatal("missing event")
	}

	easyPostRecord, err := testEvent.GetResult()
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	tracker, ok := easyPostRecord.(*Tracker)
	if !ok {
		t.Fatalf("unexpected record: %T", easyPostRecord)
	}

	expectedTracker := Tracker{}
	if err := json.Unmarshal(trackerBody, &expectedTracker); err != nil {
		t.Fatalf("error: %s", err)
	}

	if !reflect.DeepEqual(*tracker, expectedTracker) {
		t.Fatalf("trackers: \nexpected %+v\n     got %+v", expectedTracker, *tracker)
	}
}
