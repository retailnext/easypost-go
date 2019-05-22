// Copyright 2019 RetailNext, Inc.
//
// Licensed under the BSD 3-Clause License (the "License");
// you may not use this file except in compliance with the License.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package easypost

import "fmt"

var (
	paymentError      PaymentRequiredError
	unauthorizedError UnauthorizedError
)

type UnauthorizedError struct{}

func (e UnauthorizedError) Error() string {
	return "unauthorized error"
}

type PaymentRequiredError struct{}

func (e PaymentRequiredError) Error() string {
	return "payment required"
}

type ProcessingError struct {
	msg     string
	details map[string]string
}

func (e ProcessingError) Error() string {
	return e.msg
}

func (e ProcessingError) Details() map[string]string {
	return e.details
}

type ErrorResponse struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type NotSupportedRecordError struct {
	recordType RecordType
}

func (e NotSupportedRecordError) Error() string {
	return fmt.Sprintf("not supported record: %s", e.recordType)
}
