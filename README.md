typing [![Build Status](https://travis-ci.org/caixw/typing.svg?branch=nosql)](https://travis-ci.org/caixw/typing)
======
 
简单的半静态博客系统，具有以下特性：
 
1. 无数据库，通过 git 管理发布的内容；
1. 无分类，通过标签来归类；
1. 不区分页面和普通文章；
 
 
 
### 安装
 
1. 下载代码:`go get github.com/caixw/typing`；
1. 运行程序，使用 appdir 参数指定程序的工作目录；

*源码目录下的 testdata 为一个完整的工作目录内容，用户可根
据自己的需求决定是否直接使用此目录，或是在其基础上作修改*



### 目录结构

appdir 的目录结构是固定的。
其中 conf 为程序的配置相关内容，包含了后台更新界面的密码，不能对外公开；
data 为博客的实际内容，包含了文章，标签，友链以及网站名称等相关的配置，
所有针对博客内容的相关设置和内容发布，都直接体现在此目录下。

```
|--- conf 程序的配置文件
|     |
|     |--- logs.xml 日志的配置文件
|     |
|     |--- app.json 程序的配置文件
|
|--- data 程序的数据目录
      |
      |--- meta 博客的一些设置项
      |     |
      |     |--- config.yaml 基本设置项，比如网站名称等
      |     |
      |     |--- urls.yaml 自定义路由项的一些设置
      |     |
      |     |--- tags.yaml 标签的定义
      |     |
      |     |--- links.yaml 友情链接
      |
      |--- posts 文章所在的目录
      |
      |--- raws 其它可直接通过地址访问的内容可直接放在此处
      |
      |--- themes 自定义的主题目录
            |
            |--- default 默认的主题
```
 
 
 
#### conf
 
conf 目录下的为程序级别的配置文件，在程序加载之后，无法再次更改。其中
app.json 定义了诸如端口，证书等基本数据；
logs.xml 定义了日志的输出形式和保存路径，具体配置可参考[logs](https://github.com/issue9/logs)的相关文档。


##### app.json

名称          | 类型        | 描述
:-------------|:------------|:------
adminURL      | string      | 后台管理的地址
adminPassword | string      | 后台管理密码
core          | web.Core
core.https    | bool        | 是否启用https
core.certFile | string      | 当https为true时，此值为必填
core.keyFile  | string      | 当https为true时，此值为必填
core.port     | string      | 端口，不指定，默认为80或是443
core.headers  | map         | 附加的头信息，头信息可能在其它地方被修改
core.pprof    | string      | 指定pprof地址，输出net/pprof中指定的一些信息




#### 网站内容


##### config.yaml

config.yaml 指定了网站的一些基本配置情况：

名称            | 类型        | 描述
:-------------  |:------------|:------
title           | string      | 网站标题
subtitle        | string      | 网站副标题
url             | string      | 网站的地址
keywords        | string      | 默认情况下的 keyword 内容
description     | string      | 默认情况下的 descrription 内容
beian           | string      | 备案号
uptimeFormat    | string      | 上线时间，字符串表示
pageSize        | int         | 每页显示的数量
longDateFormat  | string      | 长时间的显示格式，Go 的时间格式化方式
shortDateFormat | string      | 短时间的显示格式，Go 的时间格式化方式
theme           | string      | 默认主题
menus           | []Link      | 菜单内容，格式与 links.yaml 的相同
author          | Author      | 默认的作者信息
rss             | RSS         | rss 配置，若不需要，则不指定该值即可
atom            | RSS         | atom 配置，若不需要，则不指定该值即可
sitemap         | Sitemap     | sitemap 相关配置，若不需要，则不指定该值即可


##### links.yaml

links.yaml 用于指定友情链接。目前包含以下字段：

名称      | 类型        | 描述
:---------|:------------|:----------
text      | string      | 字面文字，可以不唯一
url       | string      | 对应的链接地址
title     | string      | a 标签的 title 属性。可以为空
icon      | string      | 图标名称，图标名称为 [fontawesome](http://fontawesome.io) 下的图标


##### tags.yaml

tags.yaml 用于指定所有的标签内容。为一个数组，每个元素包含以下字段：

名称      | 类型        | 描述
:---------|:------------|:----------
slug      | string      | 唯一名称，文章引用此值，地址中也使用此值
title     | string      | 字面文字，可以不唯一
color     | string      | 颜色值，在展示所有标签的页面，会以此颜色显示
content   | string      | 用于描述该标签的详细内容，可以是**HTML**


##### urls.yaml

用于自定义URL

名称      | 类型        | 描述
:---------|:------------|:----------
root      | string      | 根地址，可为空，表示使用根网址
suffix    | string      | 地址后缀
posts     | string      | 列表页地址
post      | string      | 文章详细页地址
tags      | string      | 标签列表页地址
tag       | string      | 标签详细页地址
themes    | string      | 主题地址


##### 主题

data/themes 下为主题文件，可定义多个主题，通过 config 中的 theme 指定当前使用的主题。
主题模板为 [Go 官方模板](https://golang.org/pkg/html/template/)。

单一主题下，可以为文章详细页定义多个模板，通过 meta.yaml 可以自定义当前文章使用的模板，
默认情况下，使用 post 模板。


400 及以上的错误信息，均可以自定义，方式为在当前主题目录下，新一个与错误代码相对应的 html 文件，
比如 400 错误，会读取 400.html 文件，以此类推。

 
 
 
 
### 开发
 
typing 以自用为主，暂时*不支持新功能的 PR*。
BUG 可在[此处](https://github.com/caixw/typing/issues)提交或是直接 PR。
 
详细的开发文档可在[DEV](DEV.md)中找到。
 
 
 
### 版权
 
本项目采用[MIT](https://opensource.org/licenses/MIT)开源授权许可证，完整的授权说明可在[LICENSE](LICENSE)文件中找到。
