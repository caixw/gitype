typing [![Build Status](https://travis-ci.org/caixw/typing.svg?branch=nosql)](https://travis-ci.org/caixw/typing)
======

简单的半静态博客系统，具有以下特性：

1. 无数据库，通过 Git 管理发布的内容；
1. 无分类，通过标签来归类；
1. 不区分页面和普通文章；
1. 可以实时搜索内容。



### 安装

1. 下载代码:`go get github.com/caixw/typing`；
1. 运行程序，使用 appdir 参数指定程序的工作目录；

*源码目录下的 testdata 为一个完整的工作目录内容，
用户可根据自己的需求决定是否直接使用此目录，或是在其基础上作修改。*



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

conf 目录下的为程序级别的配置文件，在程序加载之后，无法再次更改。其中：
- app.json 定义了诸如端口，证书等基本数据；
- logs.xml 定义了日志的输出形式和保存路径，具体配置可参考 [logs](https://github.com/issue9/logs) 的相关文档。


##### app.json

名称          | 类型        | 描述
:-------------|:------------|:------
adminURL      | string      | 后台管理的地址
adminPassword | string      | 后台管理密码
https         | bool        | 是否启用 https
httpState     | string      | 当 https 为 true 时，对 80 端口的处理方式，可以为 disable, redirect, default
certFile      | string      | 当 https 为 true 时，此值为必填
keyFile       | string      | 当 https 为 true 时，此值为必填
port          | string      | 端口，不指定，默认为 80 或是 443
headers       | map         | 附加的头信息，头信息可能在其它地方被修改
pprof         | bool        | 是否需要 /debug/pprof




#### data 目录下内容


涉及的时间均为 RFC3339 格式：2006-01-02T15:04:05Z07:00。


##### config.yaml

config.yaml 指定了网站的一些基本配置情况：

名称            | 类型        | 描述
:---------------|:------------|:------
title           | string      | 网站标题
subtitle        | string      | 网站副标题
url             | string      | 网站的地址
icon            | Icon        | 网站的图标
keywords        | string      | 默认情况下的 keyword 内容
description     | string      | 默认情况下的 descrription 内容
beian           | string      | 备案号
uptimeFormat    | string      | 上线时间，字符串表示
pageSize        | int         | 每页显示的数量
longDateFormat  | string      | 长时间的显示格式，Go 的时间格式化方式
shortDateFormat | string      | 短时间的显示格式，Go 的时间格式化方式
theme           | string      | 默认主题
type            | string      | 所有 html 页面的 mime type，默认使用 vars.ContentTypeHTML
menus           | []Link      | 菜单内容，格式与 links.yaml 的相同
author          | Author      | 默认的作者信息
rss             | RSS         | rss 配置，若不需要，则不指定该值即可
atom            | RSS         | atom 配置，若不需要，则不指定该值即可
sitemap         | Sitemap     | sitemap 相关配置，若不需要，则不指定该值即可
opensearch      | Opensearch  | opensearch 相关配置，若不需要，则不指定该值即可

Author
名称      | 类型        | 描述
:---------|:------------|:----------
name      | string      | 名称
url       | string      | 网站地址
email     | string      | 邮箱
avatar    | string      | 头像


RSS

名称      | 类型        | 描述
:---------|:------------|:----------
title     | string      | 标题
size      | int         | 显示数量
url       | string      | 地址
type      | string      | 当前文件的 mimetype 若不指定，会使用 vars 包中的默认值


Sitemap

名称           | 类型        | 描述
:--------------|:------------|:----------
url            | string      | Sitemap 的地址
xslURL         | string      | 为 sitemap.xml 配置的 xsl，可以为空
enableTag      | bool        | 是否把标签放到 Sitemap 中
tagPriority    | float64     | 标签页的权重
postPriority   | float64     | 文章页的权重
tagChangefreq  | string      | 标签页的修改频率
postChangefreq | string      | 文章页的修改频率
type           | string      | 当前文件的 mimetype 若不指定，会使用 vars 包中的默认值


Opensearch

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
type      | string      | 图标的 mime-type
sizes     | string      | 图标的大小
url       | string      | 图标地址



##### links.yaml

links.yaml 用于指定友情链接，为一个数组。每个元素包含以下字段：

名称      | 类型        | 描述
:---------|:------------|:----------
text      | string      | 字面文字，可以不唯一
url       | string      | 对应的链接地址
title     | string      | a 标签的 title 属性。可以为空
icon      | string      | 一个 URL
rel       | string      | 与该网站的关系，可用 [XFN](https://gmpg.org/xfn/) 的相关定义


##### tags.yaml

tags.yaml 用于指定所有的标签内容。为一个数组，每个元素包含以下字段：

名称      | 类型        | 描述
:---------|:------------|:----------
slug      | string      | 唯一名称，文章引用此值，地址中也使用此值
title     | string      | 字面文字，可以不唯一
color     | string      | 颜色值，在展示所有标签的页面，会以此颜色显示
content   | string      | 用于描述该标签的详细内容，可以是**HTML**



##### 主题

data/themes 下为主题文件，可定义多个主题，通过 config 中的 theme 指定当前使用的主题。
主题模板语法为 [html/template](https://golang.org/pkg/html/template/)。


单一主题下，可以为文章详细页定义多个模板，通过每篇文章的 meta.yaml 可以自定义当前文章使用的模板，
默认情况下，使用 post 模板。


400 及以上的错误信息，均可以自定义，方式为在当前主题目录下，新建一个与错误代码相对应的 html 文件，
比如 400 错误，会读取 400.html 文件，以此类推。



### 开发

typing 以自用为主，暂时*不支持新功能的 PR*。
BUG 可在[此处](https://github.com/caixw/typing/issues)提交或是直接 PR。



### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
