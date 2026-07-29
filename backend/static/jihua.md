

# Juchuan（菊传）项目当前状态交接

仓库：

```
https://github.com/boss44944/juchuan
```

分支：

```
main
```

目标：

局域网 AirDrop 类工具：

* Go 后端
* Vue3 + Element Plus 前端
* 浏览器客户端（手机无需 App）
* WebSocket 实时通信
* 文件/文字传输
* 多设备管理
* 多语言

---

# 一、已经完成并提交 GitHub 的部分

## 1. Vue3 前端工程 ✅

已经存在：

```
frontend/
```

技术栈：

```
Vue3
Vite
TypeScript
Element Plus
Pinia
Vue Router
vue-i18n
Axios
```

package.json 已经有：

* vue
* vue-router
* pinia
* element-plus
* vue-i18n
* vite

---

## 2. 前端目录结构

目前：

```
frontend/src

├── api
│   ├── index.ts
│   ├── auth.ts
│   └── config.ts
│
├── router
│   └── index.ts
│
├── stores
│   ├── device.ts
│   ├── message.ts
│   └── auth.ts
│
├── websocket
│   └── client.ts
│
├── i18n
│   └── index.ts
│
├── views
│   ├── Login.vue
│   ├── Devices.vue
│   ├── Messages.vue
│   ├── Send.vue
│   └── Config.vue
│
├── App.vue
└── main.ts
```

---

# 二、前端已经完成

## App 主布局 ✅

已经从：

```
Vue3迁移中
```

改成：

Element Plus Layout。

---

## Router ✅

已有：

```
/login

/devices

/messages

/send

/config
```

---

## i18n ✅

已经支持：

```
zh-CN

en-US

ja-JP
```

目前语言包：

菜单：

* 首页
* 历史
* 设备
* 配置
* 退出

错误：

* AUTH_PASSWORD_INVALID
* FILE_TOO_LARGE

---

## Devices.vue ✅

已经完成基础页面：

显示：

* 设备名称
* 状态
* 平台

基于：

Element Plus Table

---

## Messages.vue ✅

已经完成基础页面：

支持：

TEXT

FILE

显示：

文字消息

文件消息

---

## Send.vue ✅

已经完成基础框架：

支持：

* 选择设备
* 发送文字
* 发送文件

---

## Config.vue ⚠️

只有页面框架。

API 已创建：

但是页面还没有完整绑定。

---

## Login.vue ⚠️

页面已经创建。

支持：

* 设备 ID
* 密码

但是登录流程还没有完全联调。

---

# 三、后端已经完成

## 设备体系 ✅

完成：

### Device

支持：

```
id

display_name

platform

browser

device_secret

last_seen

status
```

---

## 设备注册 ✅

流程：

```
浏览器打开

↓

注册设备

↓

SQLite保存

↓

EventBus

↓

WebSocket通知
```

---

## Heartbeat ✅

支持：

```
LastSeen更新
```

---

## 在线/离线状态 ✅

支持：

```
online

offline
```

离线检测：

120 秒。

---

## Device Event ✅

事件：

```
DEVICE_ONLINE

DEVICE_OFFLINE
```

---

## EventBus ✅

已经存在。

---

# 四、消息系统完成情况

## Message Model ✅

已有：

```
messages

message_targets
```

支持：

```
TEXT

FILE
```

---

## 多设备发送设计 ✅

支持：

例如：

```
Mac

 |
 + iphone

 |
 + 小米手机

 |
 + Windows
```

---

## 文字消息 API ⚠️

已经有：

```
POST /api/message/text
```

但是：

需要实际联调。

---

## 文件消息 ⚠️

已经有：

```
POST /api/message/file
```

但是：

完整文件链路还没验证。

---

# 五、文件系统状态

目前：

## 已完成

模型：

```
FileRecord
```

files 表设计。

---

## 未完成

缺：

### 文件上传完整实现

需要：

```
POST /api/file/upload
```

完成：

multipart

保存磁盘

写 files 表

---

### 文件下载完整实现

需要：

```
GET /api/file/download/{id}
```

---

# 六、最大的问题（下一阶段重点）

现在不要继续加功能。

应该进入：

# 阶段1：跑通闭环

目标：

真正测试：

## 场景：

电脑：

```
MacBook
```

手机：

```
iPhone
```

---

流程：

### 1

两个设备打开网页。

### 2

设备列表：

看到：

```
MacBook 在线

iPhone 在线
```

### 3

Mac发送文字：

```
hello
```

手机收到：

```
hello

复制
```

### 4

Mac发送文件：

```
test.zip
```

手机：

收到：

```
test.zip

下载
```

### 5

状态：

```
CREATED

DELIVERED

READ
```

---

# 七、目前还没有完成的功能列表

按优先级：

---

# P0 必须完成

## 1. Vue 登录联调

缺：

```
Login.vue

↓

auth API

↓

session

↓

进入主页
```

---

## 2. 文件上传

缺：

```
multipart upload
```

---

## 3. 文件下载

缺：

```
file_id

↓

真实文件
```

---

## 4. WebSocket 全链路验证

需要确认：

```
Go

↓

EventBus

↓

WebSocket

↓

Vue Store

↓

页面更新
```

---

## 5. Go embed Vue

缺：

```
npm build

↓

frontend/dist

↓

Go embed

↓

单文件运行
```

---

# P1 后续完成

## 设备管理完善

缺：

* 重命名页面按钮
* 删除按钮
* 在线排序

---

## 消息完善

缺：

* 历史消息加载
* 已读状态
* 投递状态显示

---

## 配置页面

缺：

真正配置：

* 密码
* 端口
* 语言
* 文件目录

---

# P2 后续功能

暂时不要做：

* tus断点续传
* 大文件优化
* 图片预览
* 视频
* 托盘国际化
* 自动更新
* mDNS发现

---

# 当前真实完成度

比较准确：

```
后端基础架构      70%

设备体系          85%

消息系统          50%

文件系统          30%

Vue前端           60%

整体可运行闭环    40%左右
```

---

# 新对话继续时第一句话建议：

复制：

> 继续开发 juchuan 项目。之前已经完成 Vue3 前端迁移和部分 Go 后端改造。请先以 GitHub 当前代码为准检查，不要假设已经完成。下一步目标是先跑通闭环：Vue 登录 -> 设备在线 -> WebSocket -> 文字发送 -> 文件上传下载 -> 消息状态。

---

后续开发顺序建议不要再扩功能：

1. 修编译问题
2. 跑通登录
3. 跑通设备在线
4. 跑通文字消息
5. 跑通文件传输
6. Vue build + Go embed

完成这 6 步，项目才算从“开发状态”进入“可用状态”。
