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

   go get -u golanger.com/framework

 * 如果想了解Golanger可以先从<a href="https://github.com/golangers/samples" target="_blank">Samples</a>开始

### 约定的命名规则：
 * Model: 存放在models目录中
 * Controller: 存放在controllers目录中
 * View: 存放在views目录中
 * 静态文件放在static目录中

## Wishlist
 * Validation -- 帮助验证管理
 * Hot Compile -- 代码修改后自动编译项目，你只需要刷新的你浏览器
 * Debug -- 快速定位问题

## Samples Online
 * <a href="http://chatroom.golanger.com" target="_blank">聊天室(chatroom)</a>
 * <a href="http://guestbook.golanger.com" target="_blank">记事本(guestbook)</a>
 * <a href="http://helloworld.golanger.com" target="_blank">(helloworld)</a>
 * <a href="http://pinterest.golanger.com" target="_blank">图片分享(pinterest)</a>
 * <a href="http://play.golanger.com" target="_blank">Golang Play</a>
 * <a href="http://todo.golanger.com" target="_blank">Todo List</a>
 * <a href="http://website-admin.golanger.com" target="_blank">权限管理(website-admin)</a>

   User: testgolanger
   Password: testgolanger 


## 开发者
 * Li Wei <lee#leetaifook.com>
 * Jiang Bian <borderj#gmail.com>

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
