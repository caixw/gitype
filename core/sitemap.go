// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"time"
)

const (
	Always  = "always"
	Hourly  = "hourly"
	Daily   = "daily"
	Weekly  = "weekly"
	Monthly = "monthly"
	Yearly  = "yearly"
	Never   = "never"
)

type URL struct {
	Loc        string
	Lastmod    time.Time
	Changefreq string
	Priority   float32
}
