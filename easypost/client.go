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
)

type Client struct {
	c      http.Client
	apiKey string
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
			return nil, fmt.Errorf("error read response: %s", err)
		}
		return body, nil
	}

	return nil, c.processErrorResponse(response)
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
		return fmt.Errorf("error reading not success reponse: %s", err)
	}
	errorResponse := ErrorResponse{}
	if err := json.Unmarshal(b, &errorResponse); err != nil {
		return fmt.Errorf("error parse not success response: %s", err)
	}
	details := make(map[string]string, len(errorResponse.Errors))
	for _, e := range errorResponse.Errors {
		details[e.Field] = e.Message
	}
	return ProcessingError{
		msg:     errorResponse.Message,
		details: details,
	}
}
