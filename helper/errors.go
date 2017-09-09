// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package helper

import "fmt"

// FieldError 表示加载文件出错时的具体的错误信息
type FieldError struct {
	File    string // 所在文件
	Message string // 错误信息
	Field   string // 所在的字段
}

func (err *FieldError) Error() string {
	return fmt.Sprintf("在文件 %s 中的 %s 字段发生错误：%s", err.File, err.Field, err.Message)
}
