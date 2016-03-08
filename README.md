typing [![Build Status](https://travis-ci.org/caixw/typing.svg?branch=master)](https://travis-ci.org/caixw/typing)
======
 
简单的半静态博客系统，具有以下特性：
 
1. 无数据库，通过git管理发布的内容；
1. 无分类，通过标签来归类；
1. 不区分页面和普通文章；
 
 
 
### 安装
 
1. 下载代码:`go get github.com/caixw/typing`；
1. 编译并运行代码，使用appdir指定数据地址；
** 源码目录下的testdata为测试用数据，同时也是个完整的数据内容，可以根据些目录下的内容作为数据的初始内容 **



### 目录结构

```
|--- conf 程序的配置文件
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
      |
      |--- posts 文章所在的目录
      |
      |--- raws 其它可直接通过地址访问的内容可直接放在此处
      |
      |--- themes 自定义的主题目录
```
 
 
 
### 设置
 
conf目录下的为程序级别的配置文件，在程序加载之后，无法再次更改。其中
app.json定义了诸如端口，证书等基本数据；
logs.xml定义了日志的输出形式和保存路径，具体配置可参考[logs](https://github.com/issue9/logs)的相关文档。


##### app.json

值        | 类型        | 描述
:---------|:------------|:------
https     | bool        | 是否启用https
certFile  | string      | 当https为true时，此值为必填
keyFile   | string      | 当https为true时，此值为必填
port      | string      | 端口，不指定，默认为80或是443
headers   | map         | 附加的头信息，头信息可能在其它地方被修改
pprof     | string      | 指定pprof地址，输出net/pprof中指定的一些信息


##### 自定义路由


 
 
 
 
### 开发
 
typing以自用为主，暂时*不支持新功能的PR*。
bug可在[此处](https://github.com/caixw/typing/issues)提交或是直接PR。
 
详细的开发文档可在[DEV](DEV.md)中找到。
 
 
 
### 版权
 
本项目采用[MIT](http://opensource.org/licenses/MIT)开源授权许可证，完整的授权说明可在[LICENSE](LICENSE)文件中找到。
