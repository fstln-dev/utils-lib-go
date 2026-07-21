# utils-lib-go

Go 微服务通用工具库，基于 [Kratos v2](https://github.com/go-kratos/kratos) 生态。

## 安装

```shell
go get github.com/fstln-dev/utils-lib-go
```

## 包概览

### 基础设施

| 包       | 说明                              |
|----------|-----------------------------------|
| `cache`  | BigCache 内存缓存，支持过期时间和分片配置 |
| `config` | etcd 配置源，支持 Watch 热更新     |
| `data`   | MySQL (GORM) / Redis 客户端工厂   |
| `etcd`   | etcd 客户端连接封装               |
| `zap`    | Zap 结构化日志，支持按天滚动       |

### 中间件

| 包        | 说明                              |
|-----------|-----------------------------------|
| `alert`   | 飞书卡片告警 + 企微机器人告警      |
| `grpc`    | gRPC 服务发现与元数据传递中间件    |
| `http`    | HTTP 客户端封装、CORS、Kratos 请求解码 |
| `logging` | 服务端/客户端日志中间件            |
| `plugin`  | JWT 鉴权、参数校验、重定向中间件    |

### 安全

| 包       | 说明                              |
|----------|-----------------------------------|
| `encrypt`| AES-256-CBC / BCrypt / HS256      |
| `jwt`    | JWT 签发、验签、解码               |

### 消息队列

| 包         | 说明                              |
|------------|-----------------------------------|
| `mq/rabbit`| RabbitMQ 生产者/消费者，支持自动重连 |

### 通用工具

| 包          | 说明                              |
|-------------|-----------------------------------|
| `graphql`   | 通用 GraphQL 客户端与查询构造器    |
| `goshopify` | Shopify GraphQL API 完整封装      |
| `pool`      | 协程池与协程限制器                 |
| `rand`      | 加密安全的随机字符串生成           |
| `trie`      | 路径前缀树路由                     |

## 许可

MIT