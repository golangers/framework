/*

Copyright 2013 Golanger.com. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


Golanger Web Framework
=======================================================================

Golanger is a lightweight framework for writing web applications in Go. 


## Features
 * Configuration (Support runtime modify with immediate effect)
 * Routing
 * Controllers
 * Templates
 * Session
 * Plugins
 * UrlManager
 * DataValidation
 * Log
 * Debug
 * RunTime Environment (Develop、Pre-release、Produce)

## Wishlist
 * Further optimize performance
 * Further enhance the stability
 * Architectural refactoring achieve interface package

## Quick Start
 * 首先要确保你已经安装了Go，如果没有请参考<a href="http://golang.org/doc/install" target="_blank">Go Installation</a>
 * 设置GOPATH...
 * 安装Golanger

```
   go get -u golanger.com/framework
```

 * 如果想了解Golanger可以先从<a href="https://github.com/golangers/samples" target="_blank">Samples</a>开始

## Golanger 框架约定：
 * 框架本身不实现任何模式
 * 框架可组合实现多种模式，如mvc,restful等

## Golanger 运行时环境说明

运行程序时通过 -env=\<environment\> 命令行参数指定

 * 开发环境 - develop

    ```
!一般用于本地开发时设定。
内部会有另一个进程在监听 某对象，比如你改了啥html文件，你改了啥go文件，都可以实现实时(其实是伪实时，因为重新编译了程序，然后做逐步替换处理)生效
    ```

 * 预发布环境 - prerelease

    ```
!一般用于线上进入beta状态时使用
只监听模板内文件，不监听go文件，就是说，不会重新编译，但是会在请求时重新加载模板文件
    ```

 * 生产环境 - produce

    ```
!一般用于几乎不变动的稳定生产环境使用
不监听任何，包括模板和go文件。
    ```

## Golanger 应用目录结构说明（配置为MVC模式）

```
构建新应用时可以将samples/applicationTemplate/web/website.zip解压，作为程序应用模板来扩展功能
```

Golanger 框架支持配置为MVC模式的实现(三层架构模式)(Model-View-Controller), 它是软件工程中的一种软件架构模式，把软件系统分为三个基本部分：模型(Model)、视图(View)和控制器(Controller)

Golanger 应用目录的命名规则（可根据自己需求自行变化）：

 * 控制器(Controller): 存放在controllers目录中, 负责转发请求，对请求进行处理.
 * 模型(Model): 存放在models目录中, 程序员编写程序应有的功能(实现算法等等)、数据管理和数据库设计(可以实现具体的功能).
 * 视图(View): 存放在views目录中, 界面设计人员进行图形界面设计. 
 * 资源文件放在assets目录中
 * 配置文件放在config目录中
 * 自定义助手类放在helper目录中
 * 自定义模板函数放在templateFunc目录中
 * add-on存放独立此应用的第三方库，默认死把GOPATH设置到此目录

本文以 samples/applicationTemplate/web/website.zip 项目为例，说说Golanger的MVC模式的目录结构: 

```
 ~/website <master>$ tree
.
└── src                                           // 项目源码目录
    ├── main.go                                   // 主程序文件，配置监听端口，配置文件，可以通过此实现MVC、RestFul等模式
    ├── config                                    // 项目配置信息目录，支持程序运行时修改并实时生效，可以通过配置实现MVC、RestFul等模式，采用增强型json的格式(支持注释)
    │   ├── site                                  // 应用的相关站点配置
    │   ├── assets                                // 资源型文件相关配置
    │   ├── custom                                // 自定义变量等相关设置
    │   ├── database                              // 数据库等相关配置
    │   ├── environment                           // 应用程序内部环境变量
    │   ├── html                                  // 静态html缓存等相关配置
    │   ├── i18n                                  // i18n多国语言的相关配置
    │   ├── log                                   // 日志相关的配置，支持log分级的组合形式，log的输出形式(控制台，文件)
    │   ├── session                               // Session相关的配置，如MemorySession、FileSession、CookieSession等
    │   ├── template                              // 模板相关的配置
    │   ├── urlmanage                             // Url的管理，支持url的重写，兼容apache的重写标记
    │   └── locale                                // i18n的语言文件目录
    │       └── zh-cn                             // 中文语种文件
    ├── controllers                               // 控制器(Controller)模块， 负责所有业务的请求，转发等业务逻辑
    │   ├── 404.go                                // 404逻辑处理页面
    │   ├── app.go                                // 控制器(Controller)模块的初始化， 每一个逻辑处理页面都要注册app.go
    │   └── index.go                              // 
    │                                       // 1. 习惯性的将文件名与要处理的url路径名相同(通过App.RegisterController注册路径的Handle)，来清晰的划分功能模块， 
    │                                       // 注册响应要处理的url目录路径，然后处理子页面的逻辑。 
    │                                       // 2. 文件名和App.RegisterController注册的路径的Handle可以不同，主要依据App.RegisterController注册的路径的Handle。 
    │                                       // 3. 每个模块用一个文件名，用来功能性分离， 让结构更清晰而已。
    ├── data                                      // 存放App自有数据的目录
    ├── doc.go                                    // Go的文档文件
    ├── helper                                    // 助手模块, add-on里面的第三方库的包含, helper可以对它做二次封装。比如: website-admin 就对mgo进行了二次封装，方便以后的扩展。
    │                                       //可以关闭 App.HandleFavicon() 和 App.HandleStatic()。支持无需停止程序，动态修改生效(除特殊说明的配置项)
    ├── models                                    // 模型(Model)，编写程序应有的功能(实现算法等等)、数据管理和数据库设计(可以实现具体的功能) 
    │   └── table_name.go                         // 示例文件，基于mongodb的模型的。可自行根据自己使用的数据库系统来构建
    ├── readme                                    // 项目说明文件
    ├── templateFunc                              // 模板函数扩展目录
    │   └── operator.go                           // 一个简单的模板运算的库
    ├── tmp                                       // 临时文件目录
    ├── views                                     // 视图(View)， 支持动态修改模板内容
    │   └── theme                                 // 主题目录
    │       └── default                           // 网站的默认主题，可以在config/site中进行配置
    │           ├── _global                       // golanger会自动判断此目录是否存在， 如果存在会自动加载相关的模板文件
    │           │   ├── footer.html               // footer 模板
    │           │   └── header.html               // header 模板
    │           ├── index                         // App.RegisterController注册路径的Handle， golanger会自动判断此目录是否存在
    │           │   └── index.html                // 如果请求相关的注册路径，会自动加载相关的模板信息
    │           └── _notfound
    │               └── 404.html
    ├── assets                                    // 资源文件目录
    │   ├── static                                // 静态文件目录， 如果需要nginx之类的程序处理静态文件，可以关闭config中asset最静态目录的支持
    │   │   ├── add-on                            // 静态文件的第三方库
    │   │   ├── theme                             // 主题目录
    │   │   │   └── default                       // 网站的默认主题，可以在config/site中进行配置
    │   │   │       ├── css                       // 样式文件目录
    │   │   │       │   ├── global                // golanger会自动判断此目录是否存在
    │   │   │       │   │   └── global.css        // 如果目录存在，golanger会自动判断global.css是否存在，如果存在，会自动相应的模板变量，并在全局包含
    │   │   │       │   └── index                 // App.RegisterController注册路径的Handle， golanger会自动判断此目录是否存在
    │   │   │       │       ├── global.css        // 如果注册路径存在，golanger会自动判断global.css是否存在，如果存在，会自动相应的模板变量，并在全模块下包含
    │   │   │       │       └── index.css         // 如果注册路径存在，golanger会自动判断“注册路径”.css(index.css)是否存在，如果存在，会自动相应的模板变量，并在当前同名请求的页面包含
    │   │   │       ├── img
    │   │   │       │   ├── global                // golanger会自动判断此目录是否存在
    │   │   │       │   └── index                 // App.RegisterController注册路径的Handle， golanger会自动判断此目录是否存在
    │   │   │       └── js
    │   │   │           ├── global                // golanger会自动判断此目录是否存在
    │   │   │           │   └── global.js         // 如果目录存在，golanger会自动判断global.js是否存在，如果存在，会自动相应的模板变量，并在全局包含
    │   │   │           └── index                 // App.RegisterController注册路径的Handle， golanger会自动判断此目录是否存在
    │   │   │               ├── global.js         // 如果注册路径存在，golanger会自动判断global.js是否存在，如果存在，会自动相应的模板变量，并在全模块下包含
    │   │   │               └── index.js          // 如果注册路径存在，golanger会自动判断“注册路径”.css(index.css)是否存在，如果存在，会自动相应的模板变量，并在当前同名请求的页面包含
    │   │   └── upload                            // 上传文件路径
    │   └── html                                  // 生成静态html文件时存放的目录
    ├── add-on                                    // 第三方扩展目录， 默认GOPATH目录
    ├── build.bat                                 // windows 平台编译脚本
    ├── build.sh                                  // Linux/Mac 平台编译脚本
    └── website                                   // 生成的执行文件
```

## Samples Online
 * <a href="https://github.com/golangers/samples/tree/master/website/api/src" target="_blank">WebServer服务</a>
 * <a href="https://github.com/golangers/samples/tree/master/website/chatroom/src" target="_blank">聊天室(chatroom)</a> - <a href="http://chatroom.golanger.com" target="_blank">在线demo</a>
 * <a href="https://github.com/golangers/samples/tree/master/website/guestbook/src" target="_blank">记事本(guestbook)</a> - <a href="http://guestbook.golanger.com" target="_blank">在线demo</a>
 * <a href="https://github.com/golangers/samples/tree/master/website/helloworld/src" target="_blank">Helloworld</a> - <a href="http://helloworld.golanger.com" target="_blank">在线demo</a>
 * <a href="https://github.com/golangers/samples/tree/master/website/pinterest/src" target="_blank">图片分享(pinterest)</a> - <a href="http://pinterest.golanger.com" target="_blank">在线demo</a>
 * <a href="https://github.com/golangers/samples/tree/master/website/play/src" target="_blank">Golang Play</a> - <a href="http://play.golanger.com" target="_blank">在线demo</a>
 * <a href="https://github.com/golangers/samples/tree/master/website/todo/src" target="_blank">Todo List</a> - <a href="http://todo.golanger.com" target="_blank">在线demo</a>
 * <a href="https://github.com/golangers/samples/tree/master/website/website-admin/src" target="_blank">权限管理(website-admin)</a> - <a href="http://website-admin.golanger.com" target="_blank">在线demo</a>
   * User: testgolanger
   * Password: testgolanger 

## 性能测试
 * [Golanger Using Apache Bench for Simple Load Testing](https://github.com/golangers/framework/wiki/Golanger-Using-Apache-Bench-for-Simple-Load-Testing)

## 开发团队成员清单
 * Li Wei <lee@leetaifook.com> 网名：海意/LeeTaiFook/李大福， QQ：20660991，新浪微博：weibo.com/leetaifook
 * Jiang Bian <borderj@gmail.com>

## 联系方式
### WebSite
 * <a href="http://golanger.com/framework" target="_blank">Golanger Web Framework</a>
 * <a href="http://golanger.com" target="_blank">Golanger</a>
 * <a href="http://weibo.com/golanger" target="_blank">新浪微博</a>
 * QQ群: 29994666

### 邮件列表
 * <a href="https://groups.google.com/group/golanger" target="_blank">Golanger邮件列表</a>
 * <golanger@googlegroups.com>

*/

package documentation
