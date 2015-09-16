// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	cfg, err := LoadConfig("./testdata/app.json")
	a.NotError(err).NotNil(cfg)

	a.Equal(cfg.FrontAPIPrefix, "/api").
		Equal(cfg.DBDriver, "sqlite3")
}
