// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	conf, err := loadConfig("./testdata/app.yaml")
	a.NotError(err).NotNil(conf)
	a.Equal(conf.Port, ":8080")
	a.Equal(conf.Webhook.Frequency, time.Minute)
}
