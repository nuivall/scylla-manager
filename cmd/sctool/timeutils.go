// Copyright (C) 2017 ScyllaDB

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"
	"github.com/scylladb/mermaid/timeutc"
)

const nowSafety = 30 * time.Second

func parseStartDate(value string) (strfmt.DateTime, error) {
	now := timeutc.Now()

	if !strings.HasPrefix(value, "now") {
		t, err := timeutc.Parse(time.RFC3339, value)
		if err != nil {
			return strfmt.DateTime(t), err
		}
		if t.Before(now) {
			return strfmt.DateTime(time.Time{}), errors.New("start date cannot be in the past")
		}
		if t.Before(now.Add(nowSafety)) {
			return strfmt.DateTime(time.Time{}), errors.Errorf("start date must be at least in %s", nowSafety)
		}
	}

	if value != "now" {
		d, err := time.ParseDuration(value[3:])
		if err != nil {
			return strfmt.DateTime{}, err
		}
		if d < 0 {
			return strfmt.DateTime(time.Time{}), errors.New("start date cannot be in the past")
		}
		if d < nowSafety {
			return strfmt.DateTime(time.Time{}), errors.Errorf("start date must be at least in %s", nowSafety)
		}
		now = now.Add(d)
	} else {
		now = now.Add(nowSafety)
	}

	return strfmt.DateTime(now), nil
}

func formatTime(t strfmt.DateTime) string {
	if isZero(t) {
		return ""
	}
	return time.Time(t).Format(time.RFC822)
}

func isZero(t strfmt.DateTime) bool {
	return time.Time(t).IsZero()
}

func duration(t0 strfmt.DateTime, t1 strfmt.DateTime) string {
	var d time.Duration
	if isZero(t1) {
		d = timeutc.Now().Sub(time.Time(t0))
	} else {
		d = time.Time(t1).Sub(time.Time(t0))
	}

	return fmt.Sprint(d.Truncate(time.Second))
}
