// Copyright 2022 RetailNext, Inc.
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
	"time"
)

type TrackerStatus string

const (
	TrackerStatusUnknown            TrackerStatus = "unknown"
	TrackerStatusPreTransit         TrackerStatus = "pre_transit"
	TrackerStatusInTransit          TrackerStatus = "in_transit"
	TrackerStatusOutForDelivery     TrackerStatus = "out_for_delivery"
	TrackerStatusDelivered          TrackerStatus = "delivered"
	TrackerStatusAvailableForPickup TrackerStatus = "available_for_pickup"
	TrackerStatusReturnToSender     TrackerStatus = "return_to_sender"
	TrackerStatusFailure            TrackerStatus = "failure"
	TrackerStatusCancelled          TrackerStatus = "cancelled"
	TrackerStatusError              TrackerStatus = "error"
)

var (
	TestTrackerCodes = []string{"EZ1000000001", "EZ2000000002", "EZ3000000003", "EZ4000000004", "EZ5000000005", "EZ6000000006", "EZ7000000007"}
)

type Carrier string

func (c Carrier) String() string { return string(c) }

type Tracker struct {
	ID              string            `json:"id"`
	Object          RecordType        `json:"object"`
	Mode            string            `json:"mode"`
	TrackingCode    string            `json:"tracking_code"`
	Status          TrackerStatus     `json:"status"`
	SignedBy        string            `json:"signed_by"`
	Weight          float64           `json:"weight"`
	EstDeliveryDate *DateTime         `json:"est_delivery_date"`
	ShipmentID      string            `json:"shipment_id"`
	Carrier         string            `json:"carrier"`
	TrackingDetails []TrackingDetails `json:"tracking_details"`
	CarrierDetail   CarrierDetails    `json:"carrier_detail"`
	PublicURL       string            `json:"public_url"`
	Fees            []Fee             `json:"fees"`
	CreatedAt       DateTime          `json:"created_at"`
	UpdatedAt       DateTime          `json:"updated_at"`
}

type Fee struct {
	Object   RecordType `json:"object"`
	Type     string     `json:"type"`
	Amount   string     `json:"amount"`
	Charged  bool       `json:"charged"`
	Refunded bool       `json:"refunded"`
}

type TrackingDetails struct {
	Object           RecordType       `json:"object"`
	Message          string           `json:"message"`
	Status           TrackerStatus    `json:"status"`
	Datetime         DateTime         `json:"datetime"`
	Source           string           `json:"source"`
	TrackingLocation TrackingLocation `json:"tracking_location"`
}

type TrackingLocation struct {
	Object  RecordType `json:"object"`
	City    string     `json:"city"`
	State   string     `json:"state"`
	Country string     `json:"country"`
	Zip     string     `json:"zip"`
}

type CarrierDetails struct {
	Object                      RecordType        `json:"object"`
	Service                     string            `json:"service"`
	ContainerType               string            `json:"container_type"`
	EstDeliveryDateLocal        *DateTime         `json:"est_delivery_date_local,omitempty"`
	EstDeliveryTimeLocal        *DateTime         `json:"est_delivery_time_local,omitempty"`
	OriginLocation              string            `json:"origin_location"`
	OriginTrackingLocation      *TrackingLocation `json:"origin_tracking_location,omitempty"`
	DestinationLocation         string            `json:"destination_location"`
	DestinationTrackingLocation *TrackingLocation `json:"destination_tracking_location,omitempty"`
	GuaranteedDeliveryDate      *DateTime         `json:"guaranteed_delivery_date,omitempty"`
	AlternateIdentifier         string            `json:"alternate_identifier"`
	InitialDeliveryAttempt      DateTime          `json:"initial_delivery_attempt"`
}

func (c CarrierDetails) EstimatedDeliveryTime() *time.Time {
	if c.EstDeliveryDateLocal == nil {
		return nil
	}
	t := c.EstDeliveryDateLocal.Time
	if c.EstDeliveryTimeLocal != nil {
		t = time.Date(t.Year(), t.Month(), t.Day(), c.EstDeliveryTimeLocal.Hour(), c.EstDeliveryTimeLocal.Minute(), 0, 0, time.UTC)
	}
	return &t
}

func (c *Client) GetTracker(trackingCode string, carrier Carrier) (*Tracker, error) {
	parameters := url.Values{}
	parameters.Set("tracker[tracking_code]", trackingCode)
	if carrier != "" {
		parameters.Set("tracker[carrier]", carrier.String())
	}
	parameters.Encode()
	responseBody, err := c.post(trackerURL, parameters)
	if err != nil {
		return nil, err
	}

	tracker := &Tracker{}
	if err := json.Unmarshal(responseBody, tracker); err != nil {
		return nil, fmt.Errorf("error decode response: %s", err)
	}
	return tracker, nil
}
