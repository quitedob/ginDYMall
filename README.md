```markdown
# 抖音商城个人小项目后端服务 🛍️

本项目是一个基于 Go 语言 Gin 框架开发的抖音商城后端服务，提供了一系列电商核心功能。

## ✨ 项目特点

* **模块化设计** 🧱：清晰的项目结构，分为 API 接口层、服务逻辑层、数据访问层。
* **RESTful API** 🌐：提供标准的 RESTful API 接口。
* **用户管理** 👤：
    * 用户注册与登录 [cite: 36, 37]
    * 用户注销 [cite: 38]
    * 密码修改 [cite: 39]
    * 昵称修改 [cite: 40]
    * 用户信息更新与展示 [cite: 41, 42]
* **商品管理** 🛍️：
    * 商品创建、查询、更新、删除 [cite: 28, 29, 30, 31, 32, 33, 34]
    * 商品列表查询（支持分页） [cite: 35]
* **购物车管理** 🛒：
    * 创建购物车、获取购物车信息 [cite: 16, 18]
    * 清空购物车、添加/更新购物车商品 [cite: 19, 20, 21, 22]
* **订单管理** 🧾：
    * 创建订单、更新订单信息 [cite: 24, 25, 26, 27]
* **结算流程** 💳：
    * 订单结算功能 [cite: 23]
* **认证与授权** 🔑：
    * 基于 JWT (JSON Web Tokens) 的用户认证机制 [cite: 81, 83]
    * Auth 中间件进行接口访问权限控制 [cite: 65]
* **配置管理** ⚙️：
    * 使用 Viper 进行灵活的配置管理 (`config/config.yaml`) [cite: 54, 55]
* **数据库操作** 🗄️：
    * 使用 GORM 作为 ORM 框架操作 MySQL 数据库 [cite: 1]
* **缓存机制** ⚡：
    * 集成 Redis 用于缓存及 JWT 令牌存储 [cite: 45, 82]
* **日志系统** 📝：
    * 使用 Logrus 进行结构化日志记录 [cite: 90]
* **分布式追踪** 📡：
    * 集成 Jaeger 和 OpenTracing 实现分布式链路追踪 [cite: 91]
* **中间件支持** 🔗：
    * CORS 跨域处理 [cite: 64]
    * 统一的响应格式化与错误处理 [cite: 74, 75]

## 🛠️ 技术栈与主要依赖

* **Go**: 1.18+ [cite: 1]
* **Gin**: Web 框架 [cite: 1]
* **GORM**: ORM 框架 [cite: 1]
* **MySQL**: 关系型数据库 [cite: 1]
* **Redis**: 缓存数据库 (go-redis/v9) [cite: 1]
* **JWT-Go**: JSON Web Token 实现 [cite: 1]
* **Viper**: 配置管理 [cite: 1]
* **Logrus**: 日志记录 [cite: 1]
* **Jaeger Client**: 分布式追踪 [cite: 1]
* **OpenTracing**: 分布式追踪标准接口 [cite: 1]

## 🚀 快速开始

### 先决条件

确保您的开发环境中已安装以下软件：

* Go (版本 1.18 或更高版本) [cite: 1]
* MySQL 数据库
* Redis 服务

### 配置文件

项目启动前，需要正确配置 `config/config.yaml` 文件。请根据您的实际环境修改以下关键配置项：

* **MySQL 数据库连接信息** [cite: 57]：
    * `mysql.default.dbHost`
    * `mysql.default.dbPort`
    * `mysql.default.dbName` (默认为 `douyin_mall`)
    * `mysql.default.userName`
    * `mysql.default.password`
* **Redis 连接信息** [cite: 58]：
    * `redis.redisHost`
    * `redis.redisPort`
    * `redis.redisPwd` (如果需要密码)
    * `redis.redisDbName`
* **JWT 密钥**：
    * `encryptSecret.jwtSecret` (默认为 `DouyinSecret`) [cite: 62, 63]
* **服务端口**：
    * `system.HttpPort` (默认为 `:5001`) [cite: 56]

**示例 `config/config.yaml` 结构** [cite: 56, 57, 58, 59]：
```yaml
system:
  HttpPort: ":5001"
  # ... 其他系统配置

mysql:
  default:
    dbHost: "127.0.0.1"
    dbPort: "3306"
    dbName: "douyin_mall"
    userName: "root"
    password: "your_mysql_password"
    charset: "utf8mb4"

redis:
  redisHost: "127.0.0.1"
  redisPort: "6379"
  redisPwd: ""
  redisDbName: 0

encryptSecret:
  jwtSecret: "DouyinSecret" # 强烈建议修改为更安全的密钥

# ... 其他配置如 email, oss 等
```

### 数据库初始化

1.  确保您的 MySQL 服务已启动并且可以访问。
2.  手动创建名为 `douyin_mall` (或其他在 `config.yaml` 中配置的数据库名) 的数据库。
3.  项目启动时，GORM 会根据定义的模型 (在 `repository/db/model/` 目录下) 自动迁移（创建或更新）数据表。
    * 用户表 (`users`) [cite: 131]
    * 商品表 (`products`) [cite: 130]
    * 订单表 (`orders`) [cite: 124]
    * 订单项表 (`order_items`) [cite: 128]
    * 购物车项表 (`cart_items`) [cite: 120]
    * 支付表 (`payments`) [cite: 122]
    * 分类表 (`categories`) [cite: 121]
    * 商品分类关联表 (`product_categories`)

### 安装依赖

克隆项目到本地后，在项目根目录下打开终端，执行以下命令以下载和安装项目所需的依赖：

```bash
go mod tidy
# 或者
go mod download
```

### 本地开发运行

完成配置和依赖安装后，可以通过以下命令启动后端服务：

```bash
go run cmd/main.go
```
服务将默认在 `config.yaml` 中 `system.HttpPort` 指定的端口（如 `:5001`）启动 [cite: 43] 。

## 📦 项目构建 (生产环境)

当您准备好将项目部署到生产环境时，可以执行以下命令来构建可执行文件：

```bash
go build -o douyin_server cmd/main.go
```
该命令会在项目根目录下生成一个名为 `douyin_server` (或您指定的名称) 的可执行文件。

## 部署

1.  **构建可执行文件**：如上一步所示。
2.  **准备配置文件**：确保生产环境的 `config/config.yaml` 文件已根据实际情况配置完毕。
3.  **部署文件**：将构建好的可执行文件 (`douyin_server`) 和配置文件 (`config/config.yaml` 及其所在目录 `config`) 上传到您的服务器。
4.  **运行服务**：
    * 直接运行：`./douyin_server`
    * 后台运行 (例如使用 `nohup`)：`nohup ./douyin_server &`
    * 使用进程管理工具（如 `systemd`, `Supervisor`）来管理服务，以确保服务的稳定运行和自动重启。

**注意**：

* 确保服务器防火墙已开放服务所需的端口 (例如 `5001`)。
* 确保 MySQL 和 Redis 服务在生产环境中正常运行并可被后端服务访问。
* 为了安全，生产环境中的 `jwtSecret` 等敏感信息应妥善管理，避免硬编码或明文存储在易暴露的地方。

## 📋 API 接口

项目 API 路由定义在 `routes/routes.go` 文件中 [cite: 136, 137] 。主要接口分组如下：

* `/api/v1/user/`：用户相关接口 (注册、登录、信息修改等)
* `/api/v1/product/`：商品相关接口
* `/api/v1/cart/`：购物车相关接口 (需认证)
* `/api/v1/order/`：订单相关接口 (需认证)
* `/api/v1/checkout/`：结算相关接口 (需认证)

所有需要认证的接口，请求时需要在 HTTP Header 中加入 `Authorization: Bearer <your_jwt_token>`。

## 📝 主要目录结构

```
douyin/
├── api/v1/             # API 接口处理层 (Gin Handlers)
│   ├── cart.go
│   ├── checkout.go
│   ├── order.go
│   ├── product.go
│   └── user.go
├── cmd/                  # 程序入口
│   └── main.go
├── config/               # 配置文件及加载逻辑
│   ├── config.go
│   └── config.yaml
├── consts/               # 项目常量定义
├── logs/                 # 日志文件存储目录 (自动创建) [cite: 86, 87, 88]
├── middleware/           # Gin 中间件 (JWT认证, CORS, Jaeger追踪)
├── pkg/                  # 项目工具包
│   ├── error_code/       # 错误码定义
│   ├── utils/            # 通用工具 (ctl, jwt, log, track等)
├── repository/           # 数据存储层
│   ├── cache/            # Redis 缓存操作
│   └── db/               # 数据库操作
│       ├── dao/          # 数据访问对象 (DAO)
│       └── model/        # GORM 数据模型
├── routes/               # 路由定义
│   └── routes.go
├── service/              # 业务逻辑层
│   ├── cart.go
│   ├── checkout.go
│   ├── order.go
│   ├── payment.go
│   ├── product.go
│   └── user.go
├── types/                # API 请求和响应的数据结构体定义
├── go.mod                # Go 模块依赖文件
└── go.sum                # Go 模块依赖校验文件
```