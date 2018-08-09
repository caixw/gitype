CONTRIBUTING
===


### 代码规范

请使用一个支持 editorconfig 插件的编辑器，在 `.editorconfig` 中定义了除 Go
以外的其它文本文件的编码要求。

Go 代码提交前，请使用 `go fmt` 对代码进行格式化。



### 测试数据

源码目录下的 testdata 为一份完整的测试数据，如果你需要用到测试数据，
可以直接从 tetsdata 中引用，如非必要，不建议重新做部分测试数据。



### 开发环境

OS：windows、linux 和 macOS

Go：1.10 及以上

本地测试可以使用 ./testdata 作为测试数据，配合 [gobuild](https://github.com/caixw/gobuild) 使用很方便：
`gobuild -ext="go,html,css,yaml" -x="-appdir=./testdata"`。
