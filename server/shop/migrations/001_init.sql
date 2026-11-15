-- ============================================================
-- 题目01 Shop 商城小型交易系统 · 初始化建表
-- 目标库：shop（MySQL 8.0）
-- 可重复执行：全部使用 IF NOT EXISTS
-- ============================================================

-- 环境设置：SET NAMES utf8mb4;
SET NAMES utf8mb4;

-- ---------- 账号 ----------
CREATE TABLE IF NOT EXISTS users (
-- SQL：id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：username      VARCHAR(64)  NOT NULL,
  username      VARCHAR(64)  NOT NULL,
-- SQL：password_hash VARCHAR(160) NOT NULL DEFAULT '',
  password_hash VARCHAR(160) NOT NULL DEFAULT '',
-- SQL：role          ENUM('buyer','ops','warehouse','cs') NOT NULL,
  role          ENUM('buyer','ops','warehouse','cs') NOT NULL,
-- SQL：display_name  VARCHAR(64)  NOT NULL DEFAULT '',
  display_name  VARCHAR(64)  NOT NULL DEFAULT '',
-- SQL：created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：UNIQUE KEY uk_users_username (username)
  UNIQUE KEY uk_users_username (username)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='买家与员工账号';

-- 建表/建对象：CREATE TABLE IF NOT EXISTS sessions (
CREATE TABLE IF NOT EXISTS sessions (
-- SQL：token      CHAR(64)        NOT NULL,
  token      CHAR(64)        NOT NULL,
-- SQL：user_id    BIGINT UNSIGNED NOT NULL,
  user_id    BIGINT UNSIGNED NOT NULL,
-- SQL：role       ENUM('buyer','ops','warehouse','cs') NOT NULL,
  role       ENUM('buyer','ops','warehouse','cs') NOT NULL,
-- SQL：expires_at DATETIME(3)     NOT NULL,
  expires_at DATETIME(3)     NOT NULL,
-- SQL：created_at DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (token),
  PRIMARY KEY (token),
-- SQL：KEY idx_sessions_user (user_id),
  KEY idx_sessions_user (user_id),
-- SQL：KEY idx_sessions_expire (expires_at),
  KEY idx_sessions_expire (expires_at),
-- SQL：CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id)
  CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录会话';

-- ---------- 商品与库存 ----------
CREATE TABLE IF NOT EXISTS products (
-- SQL：id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：sku           VARCHAR(32)  NOT NULL,
  sku           VARCHAR(32)  NOT NULL,
-- SQL：name          VARCHAR(128) NOT NULL,
  name          VARCHAR(128) NOT NULL,
-- SQL：price_cents   BIGINT       NOT NULL COMMENT '单价，单位：分',
  price_cents   BIGINT       NOT NULL COMMENT '单价，单位：分',
-- SQL：available     INT          NOT NULL COMMENT '当前可售数量',
  available     INT          NOT NULL COMMENT '当前可售数量',
-- SQL：initial_stock INT          NOT NULL COMMENT '初始库存，用于守恒校验',
  initial_stock INT          NOT NULL COMMENT '初始库存，用于守恒校验',
-- SQL：on_sale       TINYINT(1)   NOT NULL DEFAULT 1,
  on_sale       TINYINT(1)   NOT NULL DEFAULT 1,
-- SQL：created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：updated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  updated_at    DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：UNIQUE KEY uk_products_sku (sku),
  UNIQUE KEY uk_products_sku (sku),
-- SQL：CONSTRAINT ck_products_available CHECK (available >= 0),
  CONSTRAINT ck_products_available CHECK (available >= 0),
-- SQL：CONSTRAINT ck_products_price     CHECK (price_cents > 0),
  CONSTRAINT ck_products_price     CHECK (price_cents > 0),
-- SQL：CONSTRAINT ck_products_initial   CHECK (initial_stock >= 0)
  CONSTRAINT ck_products_initial   CHECK (initial_stock >= 0)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品与可售数量';

-- ---------- 订单 ----------
CREATE TABLE IF NOT EXISTS orders (
-- SQL：id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：order_no       CHAR(20)     NOT NULL,
  order_no       CHAR(20)     NOT NULL,
-- SQL：buyer_id       BIGINT UNSIGNED NOT NULL,
  buyer_id       BIGINT UNSIGNED NOT NULL,
-- SQL：product_id     BIGINT UNSIGNED NOT NULL,
  product_id     BIGINT UNSIGNED NOT NULL,
-- SQL：product_name   VARCHAR(128) NOT NULL COMMENT '下单时商品名快照',
  product_name   VARCHAR(128) NOT NULL COMMENT '下单时商品名快照',
-- SQL：price_cents    BIGINT       NOT NULL COMMENT '下单时单价快照',
  price_cents    BIGINT       NOT NULL COMMENT '下单时单价快照',
-- SQL：qty            INT          NOT NULL,
  qty            INT          NOT NULL,
-- SQL：amount_cents   BIGINT       NOT NULL COMMENT '金额快照 = price*qty',
  amount_cents   BIGINT       NOT NULL COMMENT '金额快照 = price*qty',
-- SQL：status         ENUM('pending_payment','cancelled','closed','paid','shipped',
  status         ENUM('pending_payment','cancelled','closed','paid','shipped',
-- SQL：'refund_requested','refund_failed','refunded') NOT NULL,
                      'refund_requested','refund_failed','refunded') NOT NULL,
-- SQL：expire_at      DATETIME(3)  NOT NULL COMMENT '服务端计算的付款截止时间',
  expire_at      DATETIME(3)  NOT NULL COMMENT '服务端计算的付款截止时间',
-- SQL：paid_at        DATETIME(3)  NULL,
  paid_at        DATETIME(3)  NULL,
-- SQL：shipped_at     DATETIME(3)  NULL,
  shipped_at     DATETIME(3)  NULL,
-- SQL：closed_at      DATETIME(3)  NULL,
  closed_at      DATETIME(3)  NULL,
-- SQL：created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：updated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  updated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：UNIQUE KEY uk_orders_no (order_no),
  UNIQUE KEY uk_orders_no (order_no),
-- SQL：KEY idx_orders_buyer (buyer_id, created_at),
  KEY idx_orders_buyer (buyer_id, created_at),
-- SQL：KEY idx_orders_status_expire (status, expire_at),
  KEY idx_orders_status_expire (status, expire_at),
-- SQL：CONSTRAINT fk_orders_buyer   FOREIGN KEY (buyer_id)   REFERENCES users(id),
  CONSTRAINT fk_orders_buyer   FOREIGN KEY (buyer_id)   REFERENCES users(id),
-- SQL：CONSTRAINT fk_orders_product FOREIGN KEY (product_id) REFERENCES products(id),
  CONSTRAINT fk_orders_product FOREIGN KEY (product_id) REFERENCES products(id),
-- SQL：CONSTRAINT ck_orders_qty    CHECK (qty BETWEEN 1 AND 2),
  CONSTRAINT ck_orders_qty    CHECK (qty BETWEEN 1 AND 2),
-- SQL：CONSTRAINT ck_orders_amount CHECK (amount_cents > 0)
  CONSTRAINT ck_orders_amount CHECK (amount_cents > 0)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单（含商品与价格快照）';

-- ---------- 支付与通知 ----------
CREATE TABLE IF NOT EXISTS payments (
-- SQL：id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：order_no       CHAR(20)     NOT NULL,
  order_no       CHAR(20)     NOT NULL,
-- SQL：channel_txn_id VARCHAR(64)  NOT NULL COMMENT '渠道交易号',
  channel_txn_id VARCHAR(64)  NOT NULL COMMENT '渠道交易号',
-- SQL：amount_cents   BIGINT       NOT NULL,
  amount_cents   BIGINT       NOT NULL,
-- SQL：paid_at        DATETIME(3)  NOT NULL,
  paid_at        DATETIME(3)  NOT NULL,
-- SQL：created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：UNIQUE KEY uk_payments_txn   (channel_txn_id) COMMENT '同一渠道交易号不得归属两个订单',
  UNIQUE KEY uk_payments_txn   (channel_txn_id) COMMENT '同一渠道交易号不得归属两个订单',
-- SQL：UNIQUE KEY uk_payments_order (order_no)       COMMENT '一单只记款一次',
  UNIQUE KEY uk_payments_order (order_no)       COMMENT '一单只记款一次',
-- SQL：CONSTRAINT fk_payments_order FOREIGN KEY (order_no) REFERENCES orders(order_no)
  CONSTRAINT fk_payments_order FOREIGN KEY (order_no) REFERENCES orders(order_no)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='成功记款（幂等基石）';

-- 建表/建对象：CREATE TABLE IF NOT EXISTS payment_notices (
CREATE TABLE IF NOT EXISTS payment_notices (
-- SQL：id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：order_no       VARCHAR(32)  NOT NULL DEFAULT '',
  order_no       VARCHAR(32)  NOT NULL DEFAULT '',
-- SQL：channel_txn_id VARCHAR(64)  NOT NULL DEFAULT '',
  channel_txn_id VARCHAR(64)  NOT NULL DEFAULT '',
-- SQL：amount_cents   BIGINT       NOT NULL DEFAULT 0,
  amount_cents   BIGINT       NOT NULL DEFAULT 0,
-- SQL：sign_ok        TINYINT(1)   NOT NULL DEFAULT 0,
  sign_ok        TINYINT(1)   NOT NULL DEFAULT 0,
-- SQL：result         ENUM('accepted','rejected','duplicated','abnormal') NOT NULL,
  result         ENUM('accepted','rejected','duplicated','abnormal') NOT NULL,
-- SQL：reason         VARCHAR(255) NOT NULL DEFAULT '',
  reason         VARCHAR(255) NOT NULL DEFAULT '',
-- SQL：raw_body       TEXT         NULL,
  raw_body       TEXT         NULL,
-- SQL：created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：KEY idx_notices_order  (order_no),
  KEY idx_notices_order  (order_no),
-- SQL：KEY idx_notices_result (result, created_at)
  KEY idx_notices_result (result, created_at)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付通知原始留痕（含被拒）';

-- ---------- 退款 ----------
CREATE TABLE IF NOT EXISTS refunds (
-- SQL：id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：refund_no      CHAR(20)     NOT NULL,
  refund_no      CHAR(20)     NOT NULL,
-- SQL：order_no       CHAR(20)     NOT NULL,
  order_no       CHAR(20)     NOT NULL,
-- SQL：source         ENUM('presale','abnormal_late') NOT NULL COMMENT '售前退款/迟到支付退款',
  source         ENUM('presale','abnormal_late') NOT NULL COMMENT '售前退款/迟到支付退款',
-- SQL：amount_cents   BIGINT       NOT NULL,
  amount_cents   BIGINT       NOT NULL,
-- SQL：audit_status   ENUM('pending','approved','rejected') NOT NULL DEFAULT 'pending',
  audit_status   ENUM('pending','approved','rejected') NOT NULL DEFAULT 'pending',
-- SQL：audit_reason   VARCHAR(255) NOT NULL DEFAULT '',
  audit_reason   VARCHAR(255) NOT NULL DEFAULT '',
-- SQL：audit_by       BIGINT UNSIGNED NULL,
  audit_by       BIGINT UNSIGNED NULL,
-- SQL：audit_at       DATETIME(3)  NULL,
  audit_at       DATETIME(3)  NULL,
-- SQL：channel_status ENUM('pending','success','failed') NOT NULL DEFAULT 'pending',
  channel_status ENUM('pending','success','failed') NOT NULL DEFAULT 'pending',
-- SQL：retry_count    INT          NOT NULL DEFAULT 0,
  retry_count    INT          NOT NULL DEFAULT 0,
-- SQL：restored       TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否已恢复可售数量',
  restored       TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否已恢复可售数量',
-- SQL：created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：updated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  updated_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：UNIQUE KEY uk_refunds_no    (refund_no),
  UNIQUE KEY uk_refunds_no    (refund_no),
-- SQL：UNIQUE KEY uk_refunds_order (order_no) COMMENT '整单退款，一单一张',
  UNIQUE KEY uk_refunds_order (order_no) COMMENT '整单退款，一单一张',
-- SQL：KEY idx_refunds_audit (audit_status, channel_status),
  KEY idx_refunds_audit (audit_status, channel_status),
-- SQL：CONSTRAINT fk_refunds_order FOREIGN KEY (order_no) REFERENCES orders(order_no)
  CONSTRAINT fk_refunds_order FOREIGN KEY (order_no) REFERENCES orders(order_no)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款单';

-- ---------- 发货 ----------
CREATE TABLE IF NOT EXISTS shipments (
-- SQL：id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：order_no    CHAR(20)     NOT NULL,
  order_no    CHAR(20)     NOT NULL,
-- SQL：carrier     VARCHAR(32)  NOT NULL,
  carrier     VARCHAR(32)  NOT NULL,
-- SQL：tracking_no VARCHAR(64)  NOT NULL,
  tracking_no VARCHAR(64)  NOT NULL,
-- SQL：created_by  BIGINT UNSIGNED NULL,
  created_by  BIGINT UNSIGNED NULL,
-- SQL：created_at  DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at  DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：UNIQUE KEY uk_shipments_order (order_no) COMMENT '合法订单只能登记一次发货',
  UNIQUE KEY uk_shipments_order (order_no) COMMENT '合法订单只能登记一次发货',
-- SQL：CONSTRAINT fk_shipments_order FOREIGN KEY (order_no) REFERENCES orders(order_no)
  CONSTRAINT fk_shipments_order FOREIGN KEY (order_no) REFERENCES orders(order_no)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货登记';

-- ---------- 异常支付 ----------
CREATE TABLE IF NOT EXISTS abnormal_payments (
-- SQL：id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：order_no       VARCHAR(32)  NOT NULL,
  order_no       VARCHAR(32)  NOT NULL,
-- SQL：channel_txn_id VARCHAR(64)  NOT NULL,
  channel_txn_id VARCHAR(64)  NOT NULL,
-- SQL：amount_cents   BIGINT       NOT NULL,
  amount_cents   BIGINT       NOT NULL,
-- SQL：reason         VARCHAR(255) NOT NULL DEFAULT '',
  reason         VARCHAR(255) NOT NULL DEFAULT '',
-- SQL：status         ENUM('pending','handled') NOT NULL DEFAULT 'pending',
  status         ENUM('pending','handled') NOT NULL DEFAULT 'pending',
-- SQL：handled_by     BIGINT UNSIGNED NULL,
  handled_by     BIGINT UNSIGNED NULL,
-- SQL：handled_at     DATETIME(3)  NULL,
  handled_at     DATETIME(3)  NULL,
-- SQL：refund_id      BIGINT UNSIGNED NULL,
  refund_id      BIGINT UNSIGNED NULL,
-- SQL：created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at     DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：UNIQUE KEY uk_abnormal_txn (channel_txn_id),
  UNIQUE KEY uk_abnormal_txn (channel_txn_id),
-- SQL：KEY idx_abnormal_status (status, created_at)
  KEY idx_abnormal_status (status, created_at)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='迟到支付等异常收款待处理';

-- ---------- 库存流水与操作留痕 ----------
CREATE TABLE IF NOT EXISTS inventory_ledger (
-- SQL：id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：product_id      BIGINT UNSIGNED NOT NULL,
  product_id      BIGINT UNSIGNED NOT NULL,
-- SQL：order_no        VARCHAR(32)  NOT NULL DEFAULT '',
  order_no        VARCHAR(32)  NOT NULL DEFAULT '',
-- SQL：change_type     ENUM('occupy','release','restore','ship_out') NOT NULL,
  change_type     ENUM('occupy','release','restore','ship_out') NOT NULL,
-- SQL：delta           INT          NOT NULL COMMENT '对可售数量的增减（负数为扣减）',
  delta           INT          NOT NULL COMMENT '对可售数量的增减（负数为扣减）',
-- SQL：available_after INT          NOT NULL,
  available_after INT          NOT NULL,
-- SQL：created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at      DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：KEY idx_ledger_product (product_id, id),
  KEY idx_ledger_product (product_id, id),
-- SQL：KEY idx_ledger_order   (order_no)
  KEY idx_ledger_order   (order_no)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存变动流水（守恒校验依据）';

-- 建表/建对象：CREATE TABLE IF NOT EXISTS operation_logs (
CREATE TABLE IF NOT EXISTS operation_logs (
-- SQL：id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
-- SQL：order_no   VARCHAR(32)  NULL,
  order_no   VARCHAR(32)  NULL,
-- SQL：actor_id   BIGINT UNSIGNED NULL,
  actor_id   BIGINT UNSIGNED NULL,
-- SQL：actor_role VARCHAR(16)  NOT NULL DEFAULT '',
  actor_role VARCHAR(16)  NOT NULL DEFAULT '',
-- SQL：action     VARCHAR(48)  NOT NULL,
  action     VARCHAR(48)  NOT NULL,
-- SQL：detail     VARCHAR(500) NOT NULL DEFAULT '',
  detail     VARCHAR(500) NOT NULL DEFAULT '',
-- SQL：created_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  created_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
-- SQL：PRIMARY KEY (id),
  PRIMARY KEY (id),
-- SQL：KEY idx_oplogs_order (order_no, id),
  KEY idx_oplogs_order (order_no, id),
-- SQL：KEY idx_oplogs_time  (created_at)
  KEY idx_oplogs_time  (created_at)
-- 语句结束
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作留痕（追加写，不提供修改与删除接口）';
