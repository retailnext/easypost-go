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
)

var (
	paymentError      PaymentRequiredError
	unauthorizedError UnauthorizedError
)

type ErrorCode string

type UnauthorizedError struct{}

func (e UnauthorizedError) Error() string {
	return "unauthorized error"
}

type PaymentRequiredError struct{}

func (e PaymentRequiredError) Error() string {
	return "payment required"
}

type ProcessingError struct {
	code    string
	msg     string
	details json.RawMessage
}

func (e ProcessingError) Error() string {
	return e.msg
}

func (e ProcessingError) Details(target interface{}) error {
	return json.Unmarshal(e.details, target)
}

type errorMessage struct {
	Code        string          `json:"code"`
	Message     string          `json:"message"`
	FieldErrors json.RawMessage `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error errorMessage `json:"error"`
}

type NotSupportedRecordError struct {
	recordType RecordType
}

func (e NotSupportedRecordError) Error() string {
	return fmt.Sprintf("not supported record: %s", e.recordType)
}
