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
	"os"
	"reflect"
	"testing"
)

func TestValidateAddress(t *testing.T) {
	setup()

	address, err := testClient.VerifyAndCreateAddress(Address{
		Street1: "Valid Street Name",
	}, DeliveryVerification)
	if err != nil {
		t.Fatalf("unxepected error: %s", err)
	}

	f, err := os.Open("./test/addresses/valid_address.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var expectedAddress Address
	if err := json.NewDecoder(f).Decode(&expectedAddress); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(address, &expectedAddress) {
		t.Fatalf("address: \nexpected %+v\n     got %+v", &expectedAddress, address)
	}
}

func TestInvalidAddress(t *testing.T) {
	setup()

	_, err := testClient.VerifyAndCreateAddress(Address{
		Street1: "address",
	}, DeliveryVerification)
	if err == nil {
		t.Fatal("error expected")
	}

	processingError, ok := err.(ProcessingError)
	if !ok {
		t.Fatalf("expected ProcessingError, got: %T(%s)", err, err)
	}

	var details []AddressVerificationError
	if err := processingError.Details(&details); err != nil {
		t.Fatal(err)
	}

	expectedDetails := []AddressVerificationError{
		{
			Code:    "E.ADDRESS.NOT_FOUND",
			Field:   "address",
			Message: "Address not found",
		}, {
			Code:    "E.HOUSE_NUMBER.MISSING",
			Field:   "street1",
			Message: "House number is missing",
		},
	}
	if !reflect.DeepEqual(details, expectedDetails) {
		t.Fatalf("address error: \nexpected %+v\n     got %+v", expectedDetails, details)
	}
}
