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
	"net/http"
	"time"
)

const (
	EventStatusCompleted EventStatus = "completed"
	EventStatusFailed    EventStatus = "failed"
	EventStatusInQueue   EventStatus = "in_queue"
	EventStatusRetrying  EventStatus = "retrying"
)

type EventStatus string

type Event struct {
	Object             RecordType      `json:"object"`
	ID                 string          `json:"id"`
	Mode               string          `json:"mode"`
	Description        string          `json:"description"`
	PreviousAttributes json.RawMessage `json:"previous_attributes,omitempty"`
	Result             json.RawMessage `json:"result,omitempty"`
	Status             EventStatus     `json:"status"`
	PendingURLs        []string        `json:"pending_urls"`
	CompletedURLs      []string        `json:"completed_urls"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

func (e Event) GetResult() (interface{}, error) {
	if len(e.Result) == 0 {
		return nil, nil
	}

	record := Record{}
	if err := json.Unmarshal(e.Result, &record); err != nil {
		return nil, fmt.Errorf("error getting result type: %s", err)
	}

	var result interface{}
	switch record.Object {
	case RecordTypeTracker:
		result = &Tracker{}
	default:
		return nil, NotSupportedRecordError{record.Object}
	}

	if err := json.Unmarshal(e.Result, result); err != nil {
		return nil, fmt.Errorf("error getting %s as result: %s", record.Object, err)
	}
	return result, nil
}

type WebHookHandler func(r *http.Request) (*Event, error)

func NewWebHookHandler(apiKey, keySecret string) WebHookHandler {
	return func(r *http.Request) (*Event, error) {
		username, password, ok := r.BasicAuth()
		if ok {
			if username != apiKey || password != keySecret {
				return nil, unauthorizedError
			}
		} else if apiKey != "" {
			return nil, unauthorizedError
		}

		event := Event{}
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			return nil, fmt.Errorf("error reading easypost response: %s", err)
		}
		return &event, nil
	}
}
