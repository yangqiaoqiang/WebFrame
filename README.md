# WebFrame Go

## 1.http server 基础框架

net/http库
http.Handler接口、重写ServeHttp方法

## 2.context

用context封装HandlerFunc
将router提出

## 3.前缀树动态路由

将url分割为前缀树形式
动态匹配":"  "*"

## 4.分组控制

增加routerGroup组便于拓展中间件

原engine添加路由的实现方法、实例改为routerGroup

## 5.middleware

在RouterGroup实例后增加可拓展的middleware

内容保存在context中

## 6.panic recover

以中间件形式添加recovery中间件完成对panic、recover的捕获
