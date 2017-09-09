// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"testing"

	"github.com/caixw/typing/vars"
	"github.com/issue9/assert"
)

func TestLoadLinks(t *testing.T) {
	a := assert.New(t)
	p := vars.NewPath("./testdata")

	links, err := loadLinks(p)
	a.NotError(err).NotNil(links)

	a.True(len(links) > 0)
	a.Equal(links[0].Text, "text0")
	a.Equal(links[0].URL, "url0")
	a.Equal(links[1].Text, "text1")
	a.Equal(links[1].URL, "url1")
}

func TestOutdatedConfig_sanitize(t *testing.T) {
	a := assert.New(t)
	o := &outdatedConfig{}

	a.Error(o.sanitize())

	o.Type = "not exits"
	a.Error(o.sanitize())
	o.Type = outdatedTypeCreated

	a.Error(o.sanitize())
	o.Content = "test"

	a.NotError(o.sanitize())
}

func TestAuthor_sanitize(t *testing.T) {
	a := assert.New(t)

	author := &Author{}
	a.Error(author.sanitize())

	author.Name = ""
	a.Error(author.sanitize())

	author.Name = "caixw"
	a.NotError(author.sanitize())
}
