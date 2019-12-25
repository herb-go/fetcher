# Fetcher HTTP请求库
一个易于配置和扩充的HTTP请求库。

主要目的为
* 创建易于标准化配置和加载的客户端以及访问点
* 提供一个构建请求的标准化框架

## Fetcher 请求配置

每一个Fetcher对象包括了创建一个HTTP请求必须的参数，包括：
* URL 请求的地址
* Method 请求的方式
* Header 请求头
* Body 请求正文
* Builders 其他在http.Request对象创建后进一步设置的构建器

原则上在使用本库时不应该手工创建Fetch数据

通过Fetcher.Raw方法可以将Fetcher转化为http.Request对象和Doer请求器。

## Command 请求命令

请求命令的作用是修改Fethcer的配置参数，最后确定创建出的http请求。
大部分的扩展功能应该以Command的形式来实现

### 常用命令

本库定义了一些常用命令如下

* URL 请求地址命令
* Method 请求方式命令
* Replace 替换地址中路径部分命令
* PathPrefix 将路径加入指定前缀命令
* Body 指定请求正文命令
* JSONBody 将对象以JSON格式序列化为正文命令
* Header 添加请求头命令
* SetDoer 设置请求器命令
* SetQuery 设置查询字符串命令
* BasicAuth 设置Basic auth命令
* RequestBuilder 设置请求构建器命令
* HeaderBuilder 设置请求头构建器命令
* MethodBuilder 设置请求方式建器命令

## Preset 预设

Preset是一系列有书虚的命令的集合。
Preset提供了一系列快速操作的方法以便维护。
Preset应该是使用本库的主要方式。

## 配置

本库预先提供了两种常用的易与反序列化的配置结构

* ServerInfo 通过URL，Method,Header来定义需要创建的请求。
* Server 通过ServerInfo和Client来定义需要创建的请求。

## Response响应

Response结构是对 http.Response的简单封装。

提供了BodyContent方法来供反复读取数据。

提供了直接作为error对象的能力

能够通过传入一个code参数，直接生成带code的api错误。

## Parser 解析器

解析器是能够对响应结果进行进一步处理的接口。

能够在正常的请求形式中快速的解析需要的数据。

当解析过程中出错时，则认为整个请求都出错。

### 预定义的解析器

* Should200 判断请求状态码是否为200.不是的的话将请求当错误抛出。是的话继续执行传入的解析器
* ShouldSuccess 判断请求状态码是否为成功(<300).不是的的话将请求当错误抛出。是的话继续执行传入的解析器
* ShouldNoError 判断请求状态码是否不是服务器错误(<500).不是的的话将请求当错误抛出。是的话继续执行传入的解析器
* AsBytes 将响应内容当成字节切片读出
* AsString 将响应内容当成字符串读出
* AsJSON 将响应内容按JSON格式反序列化

## Doer 请求器

用于发起请求的接口，为空时使用http.DefaultClient发起请求

### Client

一个易于反序列化的请求起配置结构

可配置属性如下:

* TimeoutInSecond 按秒计算的超时属性，默认120
* MaxIdleConns 最大空闲链接，默认20
* IdleConnTimeoutInSecond 以秒计算的空闲超时，默认120
* TLSHandshakeTimeoutInSecond int64 TLS握手超时
* Proxy  URL形式的代理地址
