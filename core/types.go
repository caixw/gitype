// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

// ErrorResult 用于描述简单的错误返回信息。
type ErrorResult struct {
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail,omitempty"`
}

func (errs *ErrorResult) HasErrors() bool {
	return len(errs.Detail) > 0
}

func (errs *ErrorResult) Add(key, message string) {
	errs.Detail[key] = message
}
