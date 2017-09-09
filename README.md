typing [![Build Status](https://travis-ci.org/caixw/typing.svg?branch=nosql)](https://travis-ci.org/caixw/typing)
======

基于 Git 的博客系统，具有以下特性：

1. 无数据库，通过 Git 管理发布的内容；
1. 无分类，通过标签来归类；
1. 不区分页面和普通文章；
1. 可以实时搜索内容；
1. 自动生成 RSS、Atom、Sitemap 和 Opensearch 等内容；
1. 自定义主题。



演示地址： https://caixw.io



### 使用

1. 下载代码：`go get github.com/caixw/typing`；
1. 运行 `scripts/build.sh` 编译代码（也可以直接执行 `go build` 编译，除了版本号，并无其它差别。）；
1. 执行 `typing -init=to_path` 输出初始的数据内容；
1. 运行 `typing -appdir=to_path`。

*./script 目录下包含了部分平台下的转换成守护进程的脚本*



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
|     |--- app.yaml 程序的配置文件
|
|--- data 程序的数据目录
      |
      |--- meta 博客的一些设置项
      |     |
      |     |--- config.yaml 基本设置项，比如网站名称等
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



#### conf 目录下内容

conf 目录下的为程序级别的配置文件，需要重启才能使更改生效。其中：
- app.yaml 定义了诸如端口，证书等基本数据；
- logs.xml 定义了日志的输出形式和保存路径，具体配置可参考 [logs](https://github.com/issue9/logs) 的相关文档。


##### app.yaml

名称              | 类型        | 描述
:-----------------|:------------|:------
https             | bool        | 是否启用 https
httpState         | string      | 当 https 为 true 时，对 80 端口的处理方式，可以为 disable, redirect, default
certFile          | string      | 当 https 为 true 时，此值为必填
keyFile           | string      | 当 https 为 true 时，此值为必填
port              | string      | 端口，不指定，默认为 80 或是 443
pprof             | bool        | 是否需要在 /debug/pprof 输出调试信息
headers           | map         | 附加的头信息，头信息可能在其它地方被修改
webhook           | webhook     | 与 webhook 相关的设置



###### webhook

名称              | 类型          | 描述
:-----------------|:--------------|:------
url               | string        | webhooks 的接收地址
frequency         | time.Duration | webhooks 的最小更新频率
method            | string        | webhooks 接收地址的接收方法，不指定，则默认为 POST
repoURL           | string        | 远程仓库的地址


#### data 目录下内容


涉及的时间均为 RFC3339 格式：2006-01-02T15:04:05Z07:00。


##### meta/config.yaml

config.yaml 指定了网站的一些基本配置情况：

名称            | 类型        | 描述
:---------------|:------------|:------
title           | string      | 网站标题
subtitle        | string      | 网站副标题
url             | string      | 网站的地址
keywords        | string      | 默认情况下的 keyword 内容
description     | string      | 默认情况下的 descrription 内容
beian           | string      | 备案号
uptime          | string      | 上线时间，字符串表示
pageSize        | int         | 每页显示的数量
longDateFormat  | string      | 长时间的显示格式，Go 的时间格式化方式
shortDateFormat | string      | 短时间的显示格式，Go 的时间格式化方式
theme           | string      | 默认主题
type            | string      | 所有 HTML 页面的 mimetype，默认使用 vars.ContentTypeHTML
icon            | Icon        | 网站的图标
menus           | []Link      | 菜单内容，格式与 links.yaml 的相同
author          | Author      | 文章的默认作者信息
license         | Link        | 文章的默认版权信息
archive         | Archive     | 存档页的相关配置
outdated        | Outdated    | 文章过时提示信息设置
rss             | RSS         | rss 配置，若不需要，则不指定该值即可
atom            | RSS         | atom 配置，若不需要，则不指定该值即可
sitemap         | Sitemap     | sitemap 相关配置，若不需要，则不指定该值即可
opensearch      | Opensearch  | opensearch 相关配置，若不需要，则不指定该值即可

###### Author

名称      | 类型        | 描述
:---------|:------------|:----------
name      | string      | 名称
url       | string      | 网站地址
email     | string      | 邮箱
avatar    | string      | 头像


###### Outdated

名称      | 类型        | 描述
:---------|:------------|:----------
type      | string      | 比较方式，可以 created 或是 modified
duration  | string      | 超过此时间值，显示提示信息，为一个可以被 time.ParseDuration 解析的字符串
content   | string      | 提示内容，可以带上一个 %d 用于表示有多少天未被改过


###### Archive

名称      | 类型        | 描述
:---------|:------------|:----------
order     | string      | 存档的排序方式，可以是：desc(默认) 和 month
type      | string      | 存档的分类方式，可以是按年：year(默认) 或是按月：month
format    | string      | 标题的格式

###### RSS

名称      | 类型        | 描述
:---------|:------------|:----------
title     | string      | 标题
size      | int         | 显示数量
url       | string      | 地址
type      | string      | 当前文件的 mimetype 若不指定，会使用 vars 包中的默认值


###### Sitemap

名称           | 类型        | 描述
:--------------|:------------|:----------
url            | string      | Sitemap 的地址
xslURL         | string      | 为 sitemap.xml 配置的 xsl，可以为空
enableTag      | bool        | 是否把标签放到 Sitemap 中
priority       | float       | 标签页的权重
changefreq     | string      | 标签页的修改频率
postPriority   | float       | 文章页的权重
postChangefreq | string      | 文章页的修改频率
type           | string      | 当前文件的 mimetype 若不指定，会使用 vars 包中的默认值


###### Opensearch

名称      | 类型        | 描述
:-----------|:------------|:----------
url         | string      | opensearch 的地址
title       | string      | 出现于 html>head>link.title 属性中
shortName   | string      | shortName 值
description | string      | description 值
longName    | string      | longName 值
image       | Icon        | image 值
type        | string      | 当前文件的 mimetype 若不指定，会使用 vars 包中的默认值


Icon

名称      | 类型        | 描述
:---------|:------------|:----------
type      | string      | 图标的 mimetype
sizes     | string      | 图标的大小
url       | string      | 图标地址



##### meta/links.yaml

links.yaml 用于指定友情链接，为一个数组。每个元素包含以下字段：

名称      | 类型        | 描述
:---------|:------------|:----------
text      | string      | 字面文字，可以不唯一
url       | string      | 对应的链接地址
title     | string      | a 标签的 title 属性。可以为空
icon      | string      | 一个 URL
rel       | string      | 与该网站的关系，可用 [XFN](https://gmpg.org/xfn/) 的相关定义


##### meta/tags.yaml

tags.yaml 用于指定所有的标签内容。为一个数组，每个元素包含以下字段：

名称      | 类型        | 描述
:---------|:------------|:----------
slug      | string      | 唯一名称，文章引用此值，地址中也使用此值
title     | string      | 字面文字，可以不唯一
color     | string      | 颜色值，在展示所有标签的页面，会以此颜色显示
content   | string      | 用于描述该标签的详细内容，可以是 **HTML**



##### themes

data/themes 下为主题文件，可定义多个主题，通过 config 中的 theme 指定当前使用的主题。
主题模板语法为 [html/template](https://golang.org/pkg/html/template/)。


单一主题下，可以为文章详细页定义多个模板，通过每篇文章的 meta.yaml 可以自定义当前文章使用的模板，
默认情况下，使用 post 模板。


###### 错误模板

400 及以上的错误信息，均可以自定义，方式为在当前主题目录下，新建一个与错误代码相对应的 HTML 文件，
比如 400 错误，会读取 400.html 文件，以此类推。但是只能是纯 HTML 文本，不能包含模板代码。


##### raws

当访问的页面不存在时，会尝试从 raws 下访问相关内容。比如 `/abc.html`，会尝试在查找 `raws/abc.html`
文件是否存在；甚至当 `/post/2016/about.htm` 这样标准的文章路由，如果文章不存在，会也访问 `raws`
目录，查看其下是否在正好相同的文件。 


### 开发

typing 以自用为主，原则上*不支持新功能的 PR*。
BUG 可在[此处](https://github.com/caixw/typing/issues)提交或是直接 PR。



### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
