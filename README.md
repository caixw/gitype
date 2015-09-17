# typing

简单的博客系统，具有以下特性：
1. 单用户；



### 支持的数据库
1. mysql
1. sqlite3
1. postgresql



### 安装

1. 下载代码:`go get github.com/caixw/typing`；
1. 将typing复制到指定目录；
1. 执行`typing -init=config`输出基本的配置内容；
1. 修改`config`下的相关配置；
1. 执行`typing -init=db`初始化数据库配置；
1. 将源文件目录下的static复制到config/app.json指定的目录中；
1. 登录后台，作一些自定义设置，默认登录密码为`123`；


### 配置文件
提供了三种类型的配置：
- 源码级别：在main.go的顶部，可以修改配置文件名称等内容，修改这些内容需要重新编译程序；
- 程序级别：在config/app.json中，修改这些内容需要重新启动程序；
- 内置：通过后能可更改的配置项，更改之后即时生效；

### 更改配置
- 修改AdminAPIPrefix之后，记得同时修改static/admin/app.js中的apiPrefix变量
