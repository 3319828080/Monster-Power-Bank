-- 1. users - 用户表
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    open_id VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信open_id',
    union_id VARCHAR(64) NOT NULL DEFAULT '' COMMENT '微信union_id',
    nickname VARCHAR(64) NOT NULL DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像',
    phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号',
    gender TINYINT NOT NULL DEFAULT 0 COMMENT '性别 0-未知 1-男 2-女',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1-正常 0-禁用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    UNIQUE INDEX idx_open_id (open_id),
    INDEX idx_phone (phone),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 2. wallets - 钱包表
CREATE TABLE wallets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    balance BIGINT NOT NULL DEFAULT 0 COMMENT '余额(分)',
    frozen BIGINT NOT NULL DEFAULT 0 COMMENT '冻结金额(分)',
    total_recharge BIGINT NOT NULL DEFAULT 0 COMMENT '累计充值(分)',
    total_consume BIGINT NOT NULL DEFAULT 0 COMMENT '累计消费(分)',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1-正常 0-冻结',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='钱包表';

-- 3. wallet_transactions - 钱包交易流水表
CREATE TABLE wallet_transactions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    wallet_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '钱包ID',
    order_id VARCHAR(64) NOT NULL DEFAULT '' COMMENT '关联订单号',
    trade_type VARCHAR(32) NOT NULL DEFAULT '' COMMENT '交易类型: recharge-充值 consume-消费 refund-退款 frozen-冻结 unfrozen-解冻',
    amount BIGINT NOT NULL DEFAULT 0 COMMENT '交易金额(分) 正数收入 负数支出',
    balance_before BIGINT NOT NULL DEFAULT 0 COMMENT '交易前余额',
    balance_after BIGINT NOT NULL DEFAULT 0 COMMENT '交易后余额',
    remark VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_wallet_id (wallet_id),
    INDEX idx_order_id (order_id),
    INDEX idx_trade_type (trade_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='钱包交易流水表';

-- 4. stations - 充电站表
CREATE TABLE stations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '站点名称',
    address VARCHAR(255) NOT NULL DEFAULT '' COMMENT '详细地址',
    latitude DECIMAL(10,7) NOT NULL DEFAULT 0 COMMENT '纬度',
    longitude DECIMAL(10,7) NOT NULL DEFAULT 0 COMMENT '经度',
    province VARCHAR(32) NOT NULL DEFAULT '' COMMENT '省',
    city VARCHAR(32) NOT NULL DEFAULT '' COMMENT '市',
    district VARCHAR(32) NOT NULL DEFAULT '' COMMENT '区',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1-营业 0-关闭',
    open_time VARCHAR(32) NOT NULL DEFAULT '00:00-24:00' COMMENT '营业时间',
    total_cabinets INT NOT NULL DEFAULT 0 COMMENT '充电柜数量',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_city (city),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    SPATIAL INDEX idx_location (latitude, longitude)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='充电站表';

-- 5. cabinets - 充电柜表
CREATE TABLE cabinets (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    station_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '站点ID',
    cabinet_no VARCHAR(32) NOT NULL DEFAULT '' COMMENT '柜编号',
    total_slots INT NOT NULL DEFAULT 0 COMMENT '总槽位数',
    available_slots INT NOT NULL DEFAULT 0 COMMENT '可用槽位数',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1-在线 0-离线 2-故障',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_station_id (station_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='充电柜表';

-- 6. cabinet_slots - 充电柜仓位表
CREATE TABLE cabinet_slots (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    cabinet_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '充电柜ID',
    slot_no VARCHAR(16) NOT NULL DEFAULT '' COMMENT '仓位编号',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态 0-空 1-占用 2-故障 3-预约',
    power_bank_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '当前充电宝ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_cabinet_id (cabinet_id),
    INDEX idx_status (status),
    INDEX idx_power_bank_id (power_bank_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='充电柜仓位表';

-- 7. power_banks - 充电宝表
CREATE TABLE power_banks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    cabinet_slot_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '仓位ID',
    cabinet_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '充电柜ID',
    station_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '站点ID',
    power_bank_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '充电宝编号',
    battery_level INT NOT NULL DEFAULT 0 COMMENT '电量百分比',
    power VARCHAR(32) NOT NULL DEFAULT '' COMMENT '功率规格',
    total_charge_count INT NOT NULL DEFAULT 0 COMMENT '总充电次数',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1-可用 0-租出 2-充电 3-故障 4-报废',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    UNIQUE INDEX idx_power_bank_no (power_bank_no),
    INDEX idx_cabinet_slot_id (cabinet_slot_id),
    INDEX idx_cabinet_id (cabinet_id),
    INDEX idx_station_id (station_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='充电宝表';

-- 8. orders - 租借订单表
CREATE TABLE orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '订单号',
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    power_bank_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '充电宝ID',
    power_bank_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '充电宝编号',
    borrow_station_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '借出站点ID',
    borrow_cabinet_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '借出柜ID',
    borrow_slot_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '借出仓位ID',
    return_station_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '归还站点ID',
    return_cabinet_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '归还柜ID',
    return_slot_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '归还仓位ID',
    borrow_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '借出时间',
    return_time DATETIME DEFAULT NULL COMMENT '归还时间',
    duration_minutes INT NOT NULL DEFAULT 0 COMMENT '租借时长(分钟)',
    start_fee BIGINT NOT NULL DEFAULT 0 COMMENT '起步价(分)',
    hourly_fee BIGINT NOT NULL DEFAULT 0 COMMENT '每小时费用(分)',
    total_amount BIGINT NOT NULL DEFAULT 0 COMMENT '总金额(分)',
    paid_amount BIGINT NOT NULL DEFAULT 0 COMMENT '已支付金额(分)',
    status VARCHAR(32) NOT NULL DEFAULT 'borrowed' COMMENT '订单状态: borrowed-租借中 returned-已归还 paying-待支付 paid-已支付 refunding-退款中 refunded-已退款 closed-已关闭',
    remark VARCHAR(255) NOT NULL DEFAULT '' COMMENT '备注',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    UNIQUE INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_power_bank_id (power_bank_id),
    INDEX idx_status (status),
    INDEX idx_borrow_time (borrow_time),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租借订单表';

-- 9. order_items - 订单费用明细表
CREATE TABLE order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单ID',
    order_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '订单号',
    fee_type VARCHAR(32) NOT NULL DEFAULT '' COMMENT '费用类型: rent-租金 overtime-超时费 damage-损坏费 discount-优惠',
    amount BIGINT NOT NULL DEFAULT 0 COMMENT '金额(分)',
    description VARCHAR(255) NOT NULL DEFAULT '' COMMENT '描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order_id (order_id),
    INDEX idx_order_no (order_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单费用明细表';

-- 10. payments - 支付记录表
CREATE TABLE payments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    payment_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '支付单号',
    order_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '订单号',
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    channel VARCHAR(32) NOT NULL DEFAULT '' COMMENT '支付渠道: wechat-微信 alipay-支付宝 balance-余额',
    amount BIGINT NOT NULL DEFAULT 0 COMMENT '支付金额(分)',
    currency VARCHAR(8) NOT NULL DEFAULT 'CNY' COMMENT '币种',
    status VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '支付状态: pending-待支付 processing-支付中 success-成功 failed-失败',
    channel_transaction_no VARCHAR(128) NOT NULL DEFAULT '' COMMENT '渠道交易号',
    paid_at DATETIME DEFAULT NULL COMMENT '支付时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_payment_no (payment_no),
    INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_channel (channel),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';

-- 11. refunds - 退款记录表
CREATE TABLE refunds (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    refund_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '退款单号',
    payment_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '支付单号',
    order_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '订单号',
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    amount BIGINT NOT NULL DEFAULT 0 COMMENT '退款金额(分)',
    reason VARCHAR(255) NOT NULL DEFAULT '' COMMENT '退款原因',
    status VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '退款状态: pending-待处理 processing-处理中 success-成功 failed-失败',
    channel_transaction_no VARCHAR(128) NOT NULL DEFAULT '' COMMENT '渠道退款交易号',
    refunded_at DATETIME DEFAULT NULL COMMENT '退款完成时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_refund_no (refund_no),
    INDEX idx_payment_no (payment_no),
    INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款记录表';

-- 12. coupons - 优惠券定义表
CREATE TABLE coupons (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '优惠券名称',
    coupon_type VARCHAR(32) NOT NULL DEFAULT '' COMMENT '类型: full_reduction-满减 discount-折扣 free_rent-免租金',
    condition_amount BIGINT NOT NULL DEFAULT 0 COMMENT '满减条件金额(分) 0-无限制',
    discount_amount BIGINT NOT NULL DEFAULT 0 COMMENT '减免金额(分)',
    discount_rate INT NOT NULL DEFAULT 100 COMMENT '折扣率(百分比) 如80表示8折',
    total_count INT NOT NULL DEFAULT 0 COMMENT '发行总量 0-不限',
    used_count INT NOT NULL DEFAULT 0 COMMENT '已使用量',
    max_per_user INT NOT NULL DEFAULT 1 COMMENT '每人限领',
    start_time DATETIME DEFAULT NULL COMMENT '有效开始时间',
    end_time DATETIME DEFAULT NULL COMMENT '有效结束时间',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1-上架 0-下架',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_coupon_type (coupon_type),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券定义表';

-- 13. user_coupons - 用户优惠券表
CREATE TABLE user_coupons (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    coupon_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '优惠券ID',
    order_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '使用订单号',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态 0-未使用 1-已使用 2-已过期',
    used_at DATETIME DEFAULT NULL COMMENT '使用时间',
    expired_at DATETIME DEFAULT NULL COMMENT '过期时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_coupon_id (coupon_id),
    INDEX idx_order_no (order_no),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户优惠券表';

-- 14. price_rules - 计费规则表
CREATE TABLE price_rules (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '规则名称',
    rule_type VARCHAR(32) NOT NULL DEFAULT '' COMMENT '规则类型: default-默认规则 station-站点规则 vip-VIP规则',
    station_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联站点ID 0-全局',
    start_fee BIGINT NOT NULL DEFAULT 0 COMMENT '起步价(分)',
    start_duration INT NOT NULL DEFAULT 60 COMMENT '起步时长(分钟)',
    hourly_fee BIGINT NOT NULL DEFAULT 0 COMMENT '每小时费用(分)',
    daily_cap BIGINT NOT NULL DEFAULT 0 COMMENT '每日封顶(分) 0-不封顶',
    free_duration INT NOT NULL DEFAULT 0 COMMENT '免费时长(分钟)',
    overtime_fee BIGINT NOT NULL DEFAULT 0 COMMENT '超时费用(分/小时)',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1-启用 0-禁用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_rule_type (rule_type),
    INDEX idx_station_id (station_id),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='计费规则表';

-- 15. user_vip - 用户会员表
CREATE TABLE user_vip (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    level INT NOT NULL DEFAULT 0 COMMENT '会员等级',
    start_time DATETIME DEFAULT NULL COMMENT '生效时间',
    end_time DATETIME DEFAULT NULL COMMENT '到期时间',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态 0-过期 1-有效',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户会员表';

-- 16. feedback - 用户反馈表
CREATE TABLE feedback (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    order_no VARCHAR(64) NOT NULL DEFAULT '' COMMENT '关联订单号',
    feedback_type VARCHAR(32) NOT NULL DEFAULT '' COMMENT '反馈类型: complaint-投诉 suggestion-建议 repair-报修 other-其他',
    content TEXT COMMENT '反馈内容',
    images VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '图片JSON数组',
    contact_phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '联系电话',
    status TINYINT NOT NULL DEFAULT 0 COMMENT '状态 0-待处理 1-已处理 2-已驳回',
    reply TEXT COMMENT '回复内容',
    handled_at DATETIME DEFAULT NULL COMMENT '处理时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_order_no (order_no),
    INDEX idx_feedback_type (feedback_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户反馈表';

-- 17. login_logs - 登录日志表
CREATE TABLE login_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    login_type VARCHAR(32) NOT NULL DEFAULT '' COMMENT '登录类型: wechat-微信登录 phone-手机验证码',
    ip VARCHAR(64) NOT NULL DEFAULT '' COMMENT '登录IP',
    device_info VARCHAR(255) NOT NULL DEFAULT '' COMMENT '设备信息',
    status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1-成功 0-失败',
    fail_reason VARCHAR(255) NOT NULL DEFAULT '' COMMENT '失败原因',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_login_type (login_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';
