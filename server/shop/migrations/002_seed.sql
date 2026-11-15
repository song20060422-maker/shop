-- ============================================================
-- 题目01 Shop 商城小型交易系统 · 预置数据
-- 可重复执行：INSERT ... ON DUPLICATE KEY UPDATE（不改动运行期数值）
-- 说明：密码哈希留空，服务启动时由 ensureSeed 写入 bcrypt('shop123')
-- ============================================================

-- 环境设置：SET NAMES utf8mb4;
SET NAMES utf8mb4;

-- ---------- 员工账号 ----------
INSERT INTO users (username, password_hash, role, display_name) VALUES
-- SQL：('ops1',      '', 'ops',       '运营·张敏'),
  ('ops1',      '', 'ops',       '运营·张敏'),
-- SQL：('wh1',       '', 'warehouse', '仓库·李强'),
  ('wh1',       '', 'warehouse', '仓库·李强'),
-- SQL：('cs1',       '', 'cs',        '客服·王芳')
  ('cs1',       '', 'cs',        '客服·王芳')
-- SQL：ON DUPLICATE KEY UPDATE display_name = VALUES(display_name);
ON DUPLICATE KEY UPDATE display_name = VALUES(display_name);

-- ---------- 买家账号 ----------
INSERT INTO users (username, password_hash, role, display_name) VALUES
-- SQL：('buyer1', '', 'buyer', '买家·陈一'),
  ('buyer1', '', 'buyer', '买家·陈一'),
-- SQL：('buyer2', '', 'buyer', '买家·赵二'),
  ('buyer2', '', 'buyer', '买家·赵二'),
-- SQL：('buyer3', '', 'buyer', '买家·孙三')
  ('buyer3', '', 'buyer', '买家·孙三')
-- SQL：ON DUPLICATE KEY UPDATE display_name = VALUES(display_name);
ON DUPLICATE KEY UPDATE display_name = VALUES(display_name);

-- ---------- 契约账号别名（api-contract 约定账号名为 ops/warehouse/cs，供验收脚本对接） ----------
-- 与 ops1/wh1/cs1 同角色、同密码，仅用户名不同；不改动业务权限模型。
INSERT INTO users (username, password_hash, role, display_name) VALUES
-- SQL：('ops',       '', 'ops',       '运营·张敏（契约别名）'),
  ('ops',       '', 'ops',       '运营·张敏（契约别名）'),
-- SQL：('warehouse', '', 'warehouse', '仓库·李强（契约别名）'),
  ('warehouse', '', 'warehouse', '仓库·李强（契约别名）'),
-- SQL：('cs',        '', 'cs',        '客服·王芳（契约别名）')
  ('cs',        '', 'cs',        '客服·王芳（契约别名）')
-- SQL：ON DUPLICATE KEY UPDATE display_name = VALUES(display_name);
ON DUPLICATE KEY UPDATE display_name = VALUES(display_name);

-- ---------- 商品：A 单价 10000 分、可售 5 件；B 单价 5000 分 ----------
INSERT INTO products (sku, name, price_cents, available, initial_stock, on_sale) VALUES
-- SQL：('SKU-A', '商品 A（标准款）', 10000, 5, 5, 1),
  ('SKU-A', '商品 A（标准款）', 10000, 5, 5, 1),
-- SQL：('SKU-B', '商品 B（配件款）',  5000, 3, 3, 1)
  ('SKU-B', '商品 B（配件款）',  5000, 3, 3, 1)
-- SQL：ON DUPLICATE KEY UPDATE name = VALUES(name);
ON DUPLICATE KEY UPDATE name = VALUES(name);
