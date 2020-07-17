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
	"net/url"
)

type Address struct {
	ID              string         `json:"id"`
	Object          RecordType     `json:"object"`
	Mode            *string        `json:"mode"`
	Street1         string         `json:"street1"`
	Street2         string         `json:"street2"`
	City            string         `json:"city"`
	State           string         `json:"state"`
	Zip             string         `json:"zip"`
	Country         string         `json:"country"`
	Residential     bool           `json:"residential"`
	CarrierFacility *string        `json:"carrier_facility"`
	Name            *string        `json:"name"`
	Company         *string        `json:"company"`
	Phone           *string        `json:"phone"`
	Email           *string        `json:"email"`
	FederalTaxID    *string        `json:"federal_tax_id"`
	StateTaxID      *string        `json:"state_tax_id"`
	Verifications   *Verifications `json:"verifications"`
}

func (a Address) Geocode() (float32, float32, bool) {
	if a.Verifications == nil {
		return 0, 0, false
	}
	if a.Verifications.Zip4 == nil && a.Verifications.Delivery == nil {
		return 0, 0, false
	}
	if a.Verifications.Delivery != nil {
		if a.Verifications.Delivery.Details == nil {
			return 0, 0, false
		}
		return a.Verifications.Delivery.Details.Latitude, a.Verifications.Delivery.Details.Longitude, true
	}
	if a.Verifications.Zip4.Details == nil {
		return 0, 0, false
	}
	return a.Verifications.Zip4.Details.Latitude, a.Verifications.Zip4.Details.Longitude, true
}

type AddressVerificationError struct {
	Message    string  `json:"message"`
	Suggestion *string `json:"suggestion"`
	Field      string  `json:"field"`
	Code       string  `json:"code"`
}

type VerificationType string

const (
	DeliveryVerification VerificationType = "delivery"
	Zip4Verification     VerificationType = "zip4"
)

type Verifications struct {
	Zip4     *Verification `json:"zip4"`
	Delivery *Verification `json:"delivery"`
}

type Verification struct {
	Success bool                 `json:"success"`
	Errors  []FieldError         `json:"errors"`
	Details *VerificationDetails `json:"details"`
}

type VerificationDetails struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	TimeZone  string  `json:"time_zone"`
}

func (c *Client) VerifyAndCreateAddress(address Address, verificationType VerificationType) (*Address, error) {
	parameters := url.Values{}
	parameters.Set("verify_strict[]", string(verificationType))
	parameters.Set("address[country]", address.Country)
	parameters.Set("address[city]", address.City)
	parameters.Set("address[street1]", address.Street1)
	if address.State != "" {
		parameters.Set("address[state]", address.State)
	}
	if address.Zip != "" {
		parameters.Set("address[zip]", address.State)
	}
	if address.Street2 != "" {
		parameters.Set("address[street2]", address.Street2)
	}
	if address.Company != nil {
		parameters.Set("address[company]", *address.Company)
	}
	if address.Name != nil {
		parameters.Set("address[name]", *address.Name)
	}
	if address.Phone != nil {
		parameters.Set("address[phone]", *address.Name)
	}
	if address.Email != nil {
		parameters.Set("address[email]", *address.Name)
	}

	responseBody, err := c.post(addressURL, parameters)
	if err != nil {
		return nil, err
	}
	var verifiedAddress Address
	if err := json.Unmarshal(responseBody, &verifiedAddress); err != nil {
		return nil, fmt.Errorf("error decode response: %s", err)
	}
	return &verifiedAddress, nil
}
