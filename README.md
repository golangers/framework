/*

Copyright (c) 2012 The Golanger Authors. All rights reserved.


Golanger Web Framework
=======================================================================

Golanger is a lightweight framework for writing web applications in Go. 


## 框架简单实现了MVC的设计方式

### Features
 * Routing
 * Controllers
 * Templates
 * Session
 * Plugins

## Quick Start
 * 首先要确保你已经安装了Go，如果没有请参考<a href="http://golang.org/doc/install" target="_blank">Go Installation</a>
 * 设置GOPATH...
 * 安装Golanger

 ```
   go get -u golanger.com/framework
```

 * 如果想了解Golanger可以先从<a href="https://github.com/golangers/samples" target="_blank">Samples</a>开始

### 约定的命名规则：
 * Model: 存放在models目录中
 * Controller: 存放在controllers目录中
 * View: 存放在views目录中
 * 静态文件放在static目录中

## Golanger目录结构说明

Golanger框架主要实现了MVC模式(三层架构模式)(Model-View-Controller), 它是软件工程中的一种软件架构模式，把软件系统分为三个基本部分：模型(Model)、视图(View)和控制器(Controller)

Golanger约定的命名规则：
 * 控制器(Controller): 存放在controllers目录中, 负责转发请求，对请求进行处理.
 * 模型(Model): 存放在models目录中, 程序员编写程序应有的功能(实现算法等等)、数据管理和数据库设计(可以实现具体的功能).
 * 视图(View): 存放在views目录中, 界面设计人员进行图形界面设计. 
 * 静态文件放在static目录中.
 * add-on存放第三方库文件，默认是把GOPATH设置为这个目录.
 
本文以 samples/applicationTemplate/web/website.zip 项目为例，说说Golanger的目录结构: 
 ```
 ~/website <master> tree
.
└── src                             // 项目源码目录
    ├── add-on                      // 第三方扩展目录， 默认GOPATH目录
    ├── build.bat                   // windows 平台编译脚本
    ├── build.sh                    // Linux/Mac 平台编译脚本
    ├── config                      // 项目配置信息目录
    │   └── site                    // 项目配置文件，采用json的格式，配置了网站的一些基本信息
    ├── controllers                 // 控制器(Controller)模块， 负责所有业务的请求，转发等业务逻辑
    │   ├── 404.go                  // 404逻辑处理页面
    │   ├── app.go                  // 控制器(Controller)模块的初始化， 每一个逻辑处理页面都要注册app.go
    │   └── index.go                // 
                     // 1. 习惯性的将文件名与要处理的url路径名相同(通过App.RegisterController注册路径的Handle)，来清晰的划分功能模块， 
                     // 注册响应要处理的url目录路径，然后处理子页面的逻辑。 
                     // 2. 文件名和App.RegisterController注册的路径的Handle可以不同，主要依据App.RegisterController注册的路径的Handle。 
                     // 3. 每个模块用一个文件名，用来功能性分离， 让结构更清晰而已。
    ├── data                        // 存放App自有数据的目录
    ├── doc.go                      // Go的文档文件
    ├── helper                      // 助手模块, add-on里面的第三方库的包含, helper可以对它做二次封装。比如: website-admin 就对mgo进行了二次封装，方便以后的扩展。
    ├── main.go                     // 项目启动文件，配置监听端口，配置文件，如果需要nginx之类的程序处理静态文件，
                    //可以关闭 App.HandleFavicon() 和 App.HandleStatic()。支持无需停止程序，动态修改生效(除特殊说明的配置项)
    ├── models                      // 模型(Model)，编写程序应有的功能(实现算法等等)、数据管理和数据库设计(可以实现具体的功能) 
    │   └── table_name.go
    ├── readme                      // 项目说明文件
    ├── static                      // 静态文件目录， 如果需要nginx之类的程序处理静态文件，可以关闭main.go里面中的 App.HandleFavicon() 和 App.HandleStatic()。
    │   ├── add-on                  // 静态文件的第三方库
    │   ├── theme                   // 主题目录
    │   │   └── default             // 网站的默认主题，可以在config/site中进行配置
    │   │       ├── css             // 样式文件夹
    │   │       │   ├── global                  // golanger会自动判断此目录是否存在
    │   │       │   │   └── global.css          // 如果目录存在，golanger会自动判断global.css是否存在，如果存在，会自动相应的模板变量，并在全局包含
    │   │       │   └── index                   // App.RegisterController注册路径的Handle， golanger会自动判断此目录是否存在
    │   │       │       ├── global.css          // 如果注册路径存在，golanger会自动判断global.css是否存在，如果存在，会自动相应的模板变量，并在全模块下包含
    │   │       │       └── index.css           // 如果注册路径存在，golanger会自动判断“注册路径”.css(index.css)是否存在，如果存在，会自动相应的模板变量，并在当前同名请求的页面包含
    │   │       ├── img
    │   │       │   ├── global                  // golanger会自动判断此目录是否存在
    │   │       │   └── index                   // App.RegisterController注册路径的Handle， golanger会自动判断此目录是否存在
    │   │       └── js
    │   │           ├── global                  // golanger会自动判断此目录是否存在
    │   │           │   └── global.js           // 如果目录存在，golanger会自动判断global.js是否存在，如果存在，会自动相应的模板变量，并在全局包含
    │   │           └── index                   // App.RegisterController注册路径的Handle， golanger会自动判断此目录是否存在
    │   │               ├── global.js           // 如果注册路径存在，golanger会自动判断global.js是否存在，如果存在，会自动相应的模板变量，并在全模块下包含
    │   │               └── index.js            // 如果注册路径存在，golanger会自动判断“注册路径”.css(index.css)是否存在，如果存在，会自动相应的模板变量，并在当前同名请求的页面包含
    │   └── upload                              // 上传文件路径
    ├── templateFunc                            // 模板函数扩展目录
    │   └── operator.go                         // 一个简单的模板运算的库
    ├── tmp                                     // 临时文件目录
    ├── views                                   // 视图(View)， 支持动态修改模板内容
    │   └── theme                               // 主题目录
    │       └── default                         // 网站的默认主题，可以在config/site中进行配置
    │           ├── _global                     // golanger会自动判断此目录是否存在， 如果存在会自动加载相关的模板文件
    │           │   ├── footer.html             // footer 模板
    │           │   └── header.html             // header 模板
    │           ├── index                       // App.RegisterController注册路径的Handle， golanger会自动判断此目录是否存在
    │           │   └── index.html              // 如果请求相关的注册路径，会自动加载相关的模板信息
    │           └── _notfound
    │               └── 404.html
    └── website                                 // 生成的执行文件
 ```

## Wishlist
 * Validation -- 帮助验证管理
 * Hot Compile -- 代码修改后自动编译项目，你只需要刷新的你浏览器
 * Debug -- 快速定位问题

## Samples Online
 * <a href="http://chatroom.golanger.com" target="_blank">聊天室(chatroom)</a>
 * <a href="http://guestbook.golanger.com" target="_blank">记事本(guestbook)</a>
 * <a href="http://helloworld.golanger.com" target="_blank">Helloworld</a>
 * <a href="http://pinterest.golanger.com" target="_blank">图片分享(pinterest)</a>
 * <a href="http://play.golanger.com" target="_blank">Golang Play</a>
 * <a href="http://todo.golanger.com" target="_blank">Todo List</a>
 * <a href="http://website-admin.golanger.com" target="_blank">权限管理(website-admin)</a>. 
   * User: testgolanger
   * Password: testgolanger 

## 性能测试
 * [Golanger Using Apache Bench for Simple Load Testing](https://github.com/golangers/framework/wiki/Golanger-Using-Apache-Bench-for-Simple-Load-Testing)

## 开发者
 * Li Wei <lee@leetaifook.com>
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

Open Source License
------------------------------------------------------------------------------------------
This version of golanger framework is licensed under the terms of the Open Source GPL 3.0 license. 

http://www.gnu.org/licenses/gpl.html

Alternate Licensing
------------------------------------------------------------------------------------------
Commercial and OEM Licenses are available for an alternate download of golanger framework.
This is the appropriate option if you are creating proprietary applications and you are 
not prepared to distribute and share the source code of your application under the 
GPL v3 license. 

--

This library is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NON-INFRINGEMENT OF THIRD-PARTY INTELLECTUAL PROPERTY RIGHTS.  See the GNU General Public License for more details.


*/

package golanger
