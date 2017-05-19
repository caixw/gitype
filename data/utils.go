// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

func (link *Link) check() *FieldError {
	if len(link.Text) == 0 {
		return &FieldError{Field: "Text", Message: "不能为空"}
	}

	if len(link.URL) == 0 {
		return &FieldError{Field: "URL", Message: "不能为空"}
	}

	return nil
}

func (author *Author) check() *FieldError {
	if len(author.Name) == 0 {
		return &FieldError{Field: "Name", Message: "不能为空"}
	}

	return nil
}
