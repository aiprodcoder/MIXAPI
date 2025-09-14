<p align="right">
   <strong>中文</strong> | <a href="./README.en.md">English</a>
</p>
<div align="center">



# MIXAPI

🍥新一代大模型网关,聚合大模型API调用，通过标准的 OpenAI API 格式访问所有的大模型，开箱即用

 
<p align="center">
  <a href="https://raw.githubusercontent.com/Calcium-Ion/new-api/main/LICENSE">
    <img src="https://img.shields.io/github/license/Calcium-Ion/new-api?color=brightgreen" alt="license">
  </a>
 
</p>
</div>

## 📝 项目说明

**全新AI大模型接口管理与API聚合分发系统**，支持将多种大模型转换成统一的OpenAI兼容接口,Claude接口,Gemini接口，可供个人或者企业内部大模型API
统一管理和渠道分发使用(key管理与二次分发)，支持国际国内所有主流大模型，gemini,claude,qwen3,kimi-k2,豆包等，提供单可执行文件，
docker镜像，一键部署，开箱即用，完全开源，自主可控！<br>
 * MixAPI基于New-API和One-API，整合了NewAPI,OneAPI所有重要功能及问题改进优化，内置众多第三方插件为一身，成为名副其实的全能六边形战士！  
 * 超高性能优化，用魔法打败魔法，对主要转发通路的代码用AI大模型进行多轮特别优化，写出了人类想象不出的高效代码，提高大流量高并发场景性能50%以上！  
 * **MixAPI-PRO版本**，专门针对公司企业客户，去掉收费，充值，用户注册等繁琐环节,聚焦公司内部大模型API统一管理，安全审计，隐私防泄漏，大模型调用次数限制，频率限制等企业内部管理核心问题,详情点这里 MixAPI-PRO(https://github.com/aiprodcoder/MixAPI-PRO) 


 <div align="center"> <img src="/img/mixapi-info.jpg" width = "960" height = "520" alt="mixapi" /> </div>

## ✨ 主要特性

MIXAPI提供了丰富的功能：

* 🎨 全新的UI界面
* 🌍 多语言支持
* 💰 支持在线充值功能（易支付）
* 🔍 支持用key查询使用额度(已经内置)
* 🔄 兼容原版One API的数据库
* 💵 支持模型按次数收费
* ⚖️ 支持渠道加权随机
* 📈 数据看板（控制台）
* 🔒 令牌分组、模型限制
* 🤖 支持更多授权登陆方式（LinuxDO,Telegram、OIDC）
* 🔄 支持Rerank模型（Cohere和Jina）
* ⚡ 支持OpenAI Realtime API（包括Azure渠道）
* ⚡ 支持Claude Messages 格式
* 💵 支持使用路由/chat2link进入聊天界面
* 🔄 针对用户的模型限流功能
* 💰 缓存计费支持，开启后可以在缓存命中时按照设定的比例计费：
* 🔄 新增对token令牌的控制，可控制分钟请求次数限制和日请求次数限制
  <div align="center"> <img src="/img/git1.png" width = "960" height = "460" alt="mixapi" /> </div>

* 📊 新增用量日统计
  <div align="center"> <img src="/img/git2.png" width = "960" height = "460" alt="mixapi" /> </div>


* 📊 新增用量月统计
 <div align="center"> <img src="/img/git3.png" width = "960" height = "460" alt="mixapi" /> </div>



* 📋 新增令牌管理显示该令牌的今日次数和总次数
 <div align="center"> <img src="/img/git4.png" width = "960" height = "460" alt="mixapi" /> </div>



* 📝 新增通过令牌请求的内容记录显示
 <div align="center"> <img src="/img/git5.png" width = "960" height = "460" alt="mixapi" /> </div>



* 📝 支持通过令牌直接查询余额，无需登录 
 <div align="center"> <img src="/img/git6.png" width = "960" height = "460" alt="mixapi" /> </div>

* 📝 支持令牌中选择渠道组,设定该令牌只有渠道组下面的渠道可用,并结合选择渠道组,选择模型达到限制需求
* 📝 API令牌中增加控制使用总次数限制,当达到总次数限制时返回额度已用完
 <div align="center"> <img src="/img/git7.png" width = "960" height = "460" alt="mixapi" /> </div>

## 部署

详细部署指南请参考[安装指南-部署方式](https://docs.newapi.pro/installation)：

### 部署要求
- 本地数据库（默认）：SQLite（Docker部署必须挂载`/data`目录）
- 远程数据库：MySQL版本 >= 5.7.8，PgSQL版本 >= 9.6 (非必须)

### 部署方式
#### 下载二进制程序双击运行 (小白推荐)
 windows对应下载release里面的.exe文件双击运行,下载好.exe程序, 双击运行,运行起来后通过浏览器访问
```shell
http://localhost:3000
```
#### 本地运行方式
下载本项目源码  安装好go环境, 然后在根目录运行命令 , 可用于本地开发测试
```shell
git clone https://github.com/aiprodcoder/MIXAPI
cd MIXAPI
go run main.go
```

#### 自行构建docker镜像，容器运行
下载本项目Dockerfile文件，自行构建docker镜像,容器运行，可用于测试和正式运行
```shell
wget https://github.com/aiprodcoder/MIXAPI/blob/main/Dockerfile
docker build -t mixapi .    ( ARM平台 docker build --build-arg TARGETPLATFORM=linux/arm64 -t mixapi . ）

//测试运行命令
docker run -it --rm  -p 3000:3000  -v $PWD/logs:/app/logs mixapi:latest    ($PWD为当前目录)

```

#### 通过docker 容器运行(docker推荐)
直接下载docker镜像运行,把镜像通过docker load加载上,文件在release里面
```shell
# 使用SQLite 无需配置数据库连接直接运行如下命令
docker run --name mixapicon \
--restart always \
-e TZ="Asia/Shanghai" \
-v /home/ubuntu/data/mixapi:/data  \
-d  -p 3000:3000  mixapi:v1

# 使用MySQL  需要对应修改mysql的账号密码
docker run --name mixapi -d --restart always -p 3000:3000 -e SQL_DSN="root:123456@tcp(localhost:3306)/oneapi" -e TZ=Asia/Shanghai -v /home/ubuntu/data/mixapi:/data mixapi:v1
```

## 渠道重试与缓存
渠道重试功能已经实现，可以在`设置->运营设置->通用设置`设置重试次数，**建议开启缓存**功能。

### 缓存设置方法
1. `REDIS_CONN_STRING`：设置Redis作为缓存
2. `MEMORY_CACHE_ENABLED`：启用内存缓存（设置了Redis则无需手动设置）

## 接口文档

详细接口文档请参考[接口文档](https://docs.newapi.pro/api)：

- [聊天接口（Chat）](https://docs.newapi.pro/api/openai-chat)
- [图像接口（Image）](https://docs.newapi.pro/api/openai-image)
- [重排序接口（Rerank）](https://docs.newapi.pro/api/jinaai-rerank)
- [实时对话接口（Realtime）](https://docs.newapi.pro/api/openai-realtime)
- [Claude聊天接口（messages）](https://docs.newapi.pro/api/anthropic-chat)

## 相关项目
- [One API](https://github.com/songquanpeng/one-api)：原版项目
- [Midjourney-Proxy](https://github.com/novicezk/midjourney-proxy)：Midjourney接口支持
- [chatnio](https://github.com/Deeptrain-Community/chatnio)：下一代AI一站式B/C端解决方案
- [neko-api-key-tool](https://github.com/Calcium-Ion/neko-api-key-tool)：用key查询使用额度

其他基于New API的项目：
- [new-api-horizon](https://github.com/Calcium-Ion/new-api-horizon)：New API高性能优化版

## 帮助支持

 ### 社区交流 ，问题反馈 ，软件定制，软件修改，商务咨询  请扫描加入微信群，我们一起把MixAPI做到最好
 <div align="center"> <img src="/img/wx.jpg" width = "300" height = "344" alt="mixapi" /> </div>  

