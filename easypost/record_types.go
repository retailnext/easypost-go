// Copyright 2019 RetailNext, Inc.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package easypost

const (
	RecordTypeCarrierDetail    RecordType = "CarrierDetail"
	RecordTypeEvent            RecordType = "Event"
	RecordTypeFee              RecordType = "Fee"
	RecordTypeTracker          RecordType = "Tracker"
	RecordTypeTrackingDetail   RecordType = "TrackingDetail"
	RecordTypeTrackingLocation RecordType = "TrackingLocation"
)

type RecordType string

func (t RecordType) String() string { return string(t) }

type Record struct {
	ID     string     `json:"id"`
	Object RecordType `json:"object"`
}
