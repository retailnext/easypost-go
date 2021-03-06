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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var apiURL = "https://api.easypost.com/v2"

const (
	trackerURL = "trackers"
	addressURL = "addresses"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Client struct {
	c           http.Client
	apiKey      string
	errorLogger Logger
}

func (c *Client) SetErrorLog(l Logger) {
	c.errorLogger = l
}

func (c Client) errorf(f string, attr ...interface{}) {
	if c.errorLogger != nil {
		c.errorLogger.Printf(f, attr)
	}
}

func NewClient(apiKey string) *Client {
	return &Client{
		c:      http.Client{},
		apiKey: apiKey,
	}
}

func (c *Client) post(objectURL string, parameters url.Values) ([]byte, error) {
	if c == nil {
		panic("client is not initialized")
	}
	rawURL, err := url.ParseRequestURI(fmt.Sprintf("%s/%s?%s", apiURL, objectURL, parameters.Encode()))
	if err != nil {
		panic(err)
	}

	r, err := http.NewRequest("POST", rawURL.String(), nil)
	if err != nil {
		panic(err)
	}
	r.SetBasicAuth(c.apiKey, "")

	response, err := c.c.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 {
		return nil, fmt.Errorf("request can't be processed by server")
	}

	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			err := fmt.Errorf("error read response: %s", err)
			c.errorf("%s\n", err)
			return nil, err
		}
		return body, nil
	}

	err = c.processErrorResponse(response)
	if err != nil {
		c.errorf("%s\n", err)
	}
	return nil, err
}

func (c Client) processErrorResponse(response *http.Response) error {
	switch response.StatusCode {
	case http.StatusUnauthorized:
		return unauthorizedError
	case http.StatusPaymentRequired:
		return paymentError
	case http.StatusNotFound:
		return errors.New("resource is not reachable")
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading error reponse: %s", err)
	}
	c.errorf("error response body: %s\n", b)

	errorResponse := ErrorResponse{}
	if err := json.Unmarshal(b, &errorResponse); err != nil {
		return fmt.Errorf("error parse not success response: %s", err)
	}
	errorMessage := errorResponse.Error
	return ProcessingError{
		msg:     errorMessage.Message,
		code:    errorMessage.Code,
		details: errorMessage.FieldErrors,
	}
}
