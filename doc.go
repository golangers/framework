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


### 约定的命名规则：

    Model: 存放在models目录中
    Controller: 存放在controllers目录中
    View: 存放在views目录中
    静态文件放在static目录中

## 运行方法


### 安装：

1.下载go安装包，部署go的编译环境

2.安装相应的扩展包<非必须>

3.执行初步命令

``` bash
cd [path]/example/helloworld/src
chmod +x ./build.sh
```

4.编译并执行程序

``` bash
./build.sh
```

5.打开浏览器访问

    根据输出的端口 [port]
    http://localhost:[port]

## Wishlist

 * Validation -- 帮助验证管理
 * Hot Compile -- 代码或者模板修改后自动编译项目，你只需要刷新的你浏览器
 * Debug -- 快速定位问题


## 主要开发者
```
Li Wei <lee#leetaifook.com>
Jiang Bian <borderj#gmail.com>
```


## 联系方式


### WebSite

```
http://wWw.GoLangEr.Com
```

### 微博

```
新浪：http://weibo.com/golanger
```

### IM

```
QQ : 20660991
QQ群 : 29994666
Gtalk : lee#leetaifook.com 
```

### Email

```
lee#leetaifook.com
borderj#gmail.com
```

### 邮件列表

```
https://groups.google.com/group/golang-china/
golang-china@googlegroups.com
```

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
