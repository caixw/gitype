// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/is"
)

func TestTemplateExtension(t *testing.T) {
	a := assert.New(t)

	a.Equal(TemplateExtension[0], '.')
	a.True(len(TemplateExtension) > 2)
}

func TestURLSuffix(t *testing.T) {
	a := assert.New(t)

	a.Equal(urlSuffix[0], '.')
	a.True(len(urlSuffix) > 2)
}

func TestURL(t *testing.T) {
	a := assert.New(t)

	a.True(is.URL(URL))
}
