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
	"reflect"
	"testing"
	"time"
)

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

	b, err = json.Marshal([]FieldError{
		{
			Field:   "tracking_code",
			Message: "not found",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedError := ProcessingError{
		msg:     "not found",
		details: b,
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

func TestJSONCarrierDetails(t *testing.T) {
	raw := []byte(`{
	  "est_delivery_date_local": "2022-12-08",
	  "est_delivery_time_local": "20:11:53"
	}`)
	var c CarrierDetails
	if err := json.Unmarshal(raw, &c); err != nil {
		t.Fatalf("error unmarshalling details: %s", err)
	}
	if c.EstimatedDeliveryTime() == nil || !time.Date(2022, 12, 8, 20, 11, 53, 0, time.UTC).Equal(*c.EstimatedDeliveryTime()) {
		t.Errorf("unexpected time, expected: 2022-12-08 20:11:53, got: %s", c.EstimatedDeliveryTime())
	}
}
