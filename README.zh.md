## 关于 Satellity (Pre Alpha 版)

Satellity 是一个开源论坛。

## 相关技术

1. 前后端分离，因为 Golang 模板相当不好用的原因。
2. 后端主要是 Golang, 没有用框架，主要是引用 package 来解决。
3. 前端 React, 反正也没有太多的选择，就相对随意一点，另外一点就是有可能会写点 React Navive 啥的。
4. 数据存储 Postgres, 关系型数据库，而且开源，足够强大了。

## 目录结构

1. `./web` 下是所有的前端代码，请参照 `package.json` 下的 scripts。
2. `./internal` 是 Golang 相关的代码，运行参照 `Makefile`。
3. 其它部署示例，配置相关。

## 本地运行

1. 本地运行，需要准备数据库 `./internal/models` 下。
2. './web' 下 `.env.example` 需要准备一个测试用的 `Github Client Id`，目前只支持 github 登录。

## 生产环境部署

现在功能并不全，暂时空着吧。。。
