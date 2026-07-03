# 怪兽充电宝共享租赁系统 - 答辩文档（30分钟）

---

## 一、项目概述（3分钟）

### 1.1 项目背景
共享充电宝已成为出行刚需，用户随时随地需要为手机补充电量。本项目旨在构建一套完整的共享充电宝租赁平台，涵盖用户端借还流程、支付结算、站点管理等功能。

`### 1.2 技术栈
- **后端**：Go语言 + Kratos微服务框架
- **前端**：Vue3 + Vant移动端UI库
- **数据库**：MySQL（GORM ORM）+ Redis缓存
- **依赖注入**：Google Wire
- **通信协议**：HTTP + gRPC
- **支付渠道**：支付宝、微信支付`

### 1.3 模块划分
| 模块 | 功能 |
|------|------|
| 用户模块 | 微信登录、手机验证码登录、个人信息管理 |
| 订单模块 | 借还充电宝、费用计算、续租、赔付 |
| 支付模块 | 支付宝/微信支付集成、异步回调处理 |
| 钱包模块 | 余额管理、押金冻结解冻、资金流水 |
| 站点模块 | 附近站点搜索、机柜扫码、仓位管理 |

---

## 二、系统架构（5分钟）

### 2.1 分层架构（DDD思想）
采用领域驱动设计的分层架构，每层职责清晰、可独立测试：

```
HTTP/gRPC 接口层（server）
    ↓
服务编排层（service）—— 参数校验、DTO转换
    ↓
业务逻辑层（biz）—— 核心领域模型与用例
    ↓
数据访问层（data）—— MySQL + Redis 持久化
```

**各层核心职责：**
- **server层**：注册HTTP/gRPC路由，配置中间件
- **service层**：处理请求参数，调用业务用例，返回响应
- **biz层**：核心领域逻辑，定义接口契约，不含任何框架依赖
- **data层**：实现数据持久化，ORM模型与业务模型转换

### 2.2 依赖注入（Wire）
通过Google Wire在编译期完成依赖注入，无需手动管理对象创建顺序：
- server注入service -> service注入biz -> biz注入data
- 核心链路：`main -> wireApp -> NewHTTPServer -> NewOrderService -> NewOrderUsecase -> NewOrderRepo -> NewData(DB+Redis)`

### 2.3 双协议支持
同时提供HTTP JSON和gRPC两种接口协议，protobuf统一定义：
```protobuf
service Order {
  rpc CreateOrder (CreateOrderRequest) returns (CreateOrderResponse);
  rpc ReturnPowerBank (ReturnPowerBankRequest) returns (ReturnPowerBankResponse);
  // 其余接口...
}
```

### 2.4 数据模型关系
```
用户（users）1:N 订单（orders）
  |                      |
  v                      v
钱包（wallets）       充电宝（power_banks）
  |                      |
  v                      v
资金流水（wallet_transactions）  仓位（slots）
                                    |
                                    v
                                机柜（cabinets）
                                    |
                                    v
                                站点（stations）
```
           
### 2.5 核心接口清单
```
用户: Login、LoginPhone、SendSmsCode、BindPhone、GetProfile
钱包: GetWallet、PayDeposit、RefundDeposit、ListTransactions
订单: PreCheckBorrow、CreateOrder、ReturnPowerBank、GetRealTimeFee
      ExtendRental、GetCurrentOrder、ListOrders、CreateCompensationPayment
站点: ListNearbyStations、ScanCabinet、GetCabinetList、ListReturnCabinets
支付: CreatePayment、PaymentCallback
```
      
---

## 三、核心业务流程（10分钟）

### 3.1 用户认证流程（1分钟）
```
登录流程：
用户输入手机号 -> 发送验证码（Redis存储） ->
输入验证码 -> 校验通过 -> 查询/创建用户 ->
生成JWT Token -> 返回用户信息

JWT结构：user_id + open_id + phone + 过期时间
中间件拦截器：白名单放行登录接口，其余接口需携带Token
```

### 3.2 租借充电宝完整流程（4分钟）
```
Step 1 - 预检：
用户扫码 -> 查当前是否有未完成订单 ->
查仓位状态 -> 查充电宝状态 -> 查余额是否足够押金

Step 2 - 加分布式锁：
Redis SETNX 实现 -> 防止同一用户并发下单 ->
锁定10秒 -> 超时自动释放

Step 3 - 正式下单（数据库事务）：
创建订单记录 -> 冻结押金（余额->冻结） ->
释放仓位（空闲->已借出）-> 充电宝状态更新 ->
清除充电宝位置信息

Step 4 - 归还充电宝（数据库事务）：
核验订单归属 -> 自动分配空仓位 ->
计算租借时长和费用 -> 应用优惠券 ->
从余额扣费 -> 不足扣押金 -> 退还剩余押金 ->
更新订单为已完成 -> 更新仓位和充电宝状态

Step 5 - 实时费用查询：
根据已用时长在线计算费用 ->
起步价（2元/60分钟）-> 超时按小时计费（1元/小时）->
每日封顶20元 -> 展示扣费明细
```

### 3.3 费用计算算法（1.5分钟）
```
function calculateFee(duration, startFee, startMins, hourlyFee, dailyCap):
    totalMinutes = duration.Minutes()
    if totalMinutes <= startMins:       // 起步时长内
        return startFee
    
    extraMinutes = totalMinutes - startMins
    extraHours = ceil(extraMinutes / 60) // 超出部分按小时向上取整
    total = startFee + extraHours * hourlyFee
    
    if dailyCap > 0 && total > dailyCap: // 每日封顶
        total = dailyCap
    return total
```

**计费规则示例：**
| 时长 | 费用 | 说明 |
|------|------|------|
| 30分钟 | 2元 | 起步价 |
| 90分钟 | 3元 | 起步+1小时 |
| 5小时 | 6元 | 起步+4小时 |
| 10小时 | 20元 | 每日封顶 |

### 3.4 支付系统设计（2分钟）
```
支付统一接口：
CreatePayment(bizType, bizNo, channel, amount)
  -> 创建支付单（状态pending）
  -> 调用支付渠道（支付宝/微信）
  -> 返回支付参数给前端

异步回调处理（幂等设计）：
支付渠道回调 -> 验签（支付宝SDK内置） ->
查询支付单 -> 校验状态（已处理则跳过） ->
状态流转 success/failed -> 触发业务逻辑 ->
发送支付成功通知

业务类型处理：
deposit   -> 钱包余额增加
rental    -> 订单标记完成
compensation -> 订单标记已赔付
```

### 3.5 钱包与押金体系（1.5分钟）
```
余额 = 可用余额 + 冻结余额（押金）

状态流转：
租借时：余额->押金（冻结）
归还时：扣除费用 -> 解冻剩余押金
退还时：冻结->余额（解冻）

安全设计：
- FreezeBalance 用 CAS 条件更新（WHERE balance >= amount）
- 所有资金操作在数据库事务中执行
- 资金流水记录每一笔变动（记账式设计）
```

---

## 四、关键技术亮点（5分钟）

### 4.1 分布式锁
```go
// 基于 Redis SETNX + Lua 脚本实现
Lock(key, ttl)     -> SETNX key token EX ttl
Unlock(key, token) -> EVAL "if GET key==token then DEL key end"
```
**应用场景**：防止用户并发下单，确保同一时刻只有一个订单创建请求被处理。锁超时10秒避免死锁。

### 4.2 JWT认证中间件
```go
Auth中间件 -> 解析Bearer Token -> 验证签名 -> 提取user_id -> 注入Context
白名单机制 -> 登录、短信、回调接口跳过认证
Context传递 -> UserIDFromContext(ctx) 随处获取当前用户
```

### 4.3 支付渠道策略模式
```go
// 统一的 Channel 接口
type Channel interface {
    Name() string
    Pay(ctx, req) (*PaymentResponse, error)
    Refund(ctx, req) (*RefundResponse, error)
    VerifyNotify(ctx, params) (*PaymentResponse, error)
}
// 支付宝 / 微信分别实现
Manager -> 按名称注册 -> 运行时按渠道路由
```
**优势**：新增支付渠道只需实现Channel接口，Manager注册即可，开闭原则。

### 4.4 附近站点空间搜索
```go
// 经纬度范围查询 + Haversine距离计算
latDelta = radius / 111320.0
lngDelta = radius / (111320.0 * cos(lat))
WHERE latitude BETWEEN minLat AND maxLat
  AND longitude BETWEEN minLng AND maxLng
// Haversine公式精确计算距离 -> 按距离排序
```

### 4.5 数据库事务保证一致性
```go
// 租借操作在单事务中完成：
// 创建订单 + 冻结押金 + 释放仓位 + 更新充电宝
Transaction(func(ctx) {
    repo.Create(order)
    wallet.FreezeBalance()
    repo.ReleaseSlot()
    repo.UpdatePowerBankStatus()
})
// 任一失败 -> 全部回滚
```

### 4.6 Redis缓存应用总结
| 用途 | TTL | Key模式 |
|------|-----|---------|
| 短信验证码 | 5分钟 | sms_code:{phone} |
| 分布式锁 | 10秒 | lock:order:create:{userId} |
| 热点数据缓存 | 可配置 | cache:{key} |

---

## 五、测试与部署（5分钟）

### 5.1 测试策略
- **单元测试**：覆盖biz层核心业务逻辑
- **Mock测试**：使用gomock模拟repo层接口
- **集成测试**：真实数据库验证事务一致性
- **压测场景**：并发扣库存、高并发下单

### 5.2 性能优化措施
- 数据库连接池：最大100连接
- Redis缓存热点数据，减少DB查询
- 慢查询加索引（order_no、user_id、status）
- 地理位置查询范围过滤 + 内存排序

### 5.3 部署架构
```
Kratos应用 -> HTTP端口 + gRPC端口
MySQL持久化 + Redis缓存
配置文件 -> YAML管理多环境配置（开发/测试/生产）
CI/CD流水线 -> 自动构建、测试、部署
```

### 5.4 待完善方向
- 微信支付V3完整对接
- 熔断限流保护
- 操作日志审计
- 消息队列异步解耦
- 数据分库分表（扩展性）

---

## 六、总结（2分钟）

### 6.1 项目成果
- 完成共享充电宝租赁核心完整闭环
- 三层架构清晰，业务与基础设施解耦
- 支付宝/微信支付双渠道支持
- 分布式锁保障高并发数据一致性
- 实时费用计算引擎
- 地理位置搜索就近推荐站点

### 6.2 个人收获
- 深入理解DDD分层架构在微服务中的实践
- 掌握Kratos框架 + Wire依赖注入的最佳实践
- 支付系统幂等设计与异步回调处理经验
- Redis分布式锁、JWT认证等中间件研发
- 复杂订单状态机的事务一致性保障

---

## 答辩时间分配

| 环节 | 时间 |
|------|------|
| 项目概述 | 3分钟 |
| 系统架构 | 5分钟 |
| 核心业务流程 | 10分钟 |
| 技术亮点 | 5分钟 |
| 测试与部署 | 5分钟 |
| 总结 | 2分钟 |
| 合计 | 30分钟 |
