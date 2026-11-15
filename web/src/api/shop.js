// 执行：import service from '@/utils/request'
import service from '@/utils/request'
// Shop 商城后台管理（运营/仓库/客服）接口封装
// 后端：server/shop/adminapi（挂在 GVA PrivateGroup，走 gva 登录 + Casbin）

// ---------- 运营 ----------
export const opsSummary = () => service({ url: '/shop/ops/summary', method: 'get' })
// 执行：export const opsProducts = () => service({ url: '/shop/ops/products', method: 'get' })
export const opsProducts = () => service({ url: '/shop/ops/products', method: 'get' })
// 执行：export const opsUpdateProduct = (id, data) => service({ url: `/shop/ops/products/${id}`, method: 'patch', data...
export const opsUpdateProduct = (id, data) => service({ url: `/shop/ops/products/${id}`, method: 'patch', data })
// 执行：export const opsOrders = (params) => service({ url: '/shop/ops/orders', method: 'get', params })
export const opsOrders = (params) => service({ url: '/shop/ops/orders', method: 'get', params })
// 执行：export const opsOrderTrace = (orderNo) => service({ url: `/shop/ops/orders/${orderNo}`, method: 'get' })
export const opsOrderTrace = (orderNo) => service({ url: `/shop/ops/orders/${orderNo}`, method: 'get' })
// 执行：export const opsNotices = (params) => service({ url: '/shop/ops/notices', method: 'get', params })
export const opsNotices = (params) => service({ url: '/shop/ops/notices', method: 'get', params })
// 执行：export const opsAbnormals = (params) => service({ url: '/shop/ops/abnormals', method: 'get', params })
export const opsAbnormals = (params) => service({ url: '/shop/ops/abnormals', method: 'get', params })
// 执行：export const opsLedger = (params) => service({ url: '/shop/ops/ledger', method: 'get', params })
export const opsLedger = (params) => service({ url: '/shop/ops/ledger', method: 'get', params })
// 执行：export const opsTrail = (params) => service({ url: '/shop/ops/trail', method: 'get', params })
export const opsTrail = (params) => service({ url: '/shop/ops/trail', method: 'get', params })

// ---------- 仓库 ----------
export const whShipable = () => service({ url: '/shop/warehouse/shipable', method: 'get' })
// 执行：export const whShip = (data) => service({ url: '/shop/warehouse/ship', method: 'post', data })
export const whShip = (data) => service({ url: '/shop/warehouse/ship', method: 'post', data })
// 执行：export const whShipments = () => service({ url: '/shop/warehouse/shipments', method: 'get' })
export const whShipments = () => service({ url: '/shop/warehouse/shipments', method: 'get' })
// 执行：export const whOrderTrace = (orderNo) => service({ url: `/shop/warehouse/orders/${orderNo}`, method: 'get' })
export const whOrderTrace = (orderNo) => service({ url: `/shop/warehouse/orders/${orderNo}`, method: 'get' })

// ---------- 客服 ----------
export const csRefunds = (params) => service({ url: '/shop/cs/refunds', method: 'get', params })
// 执行：export const csRefundAudit = (data) => service({ url: '/shop/cs/refunds/audit', method: 'post', data })
export const csRefundAudit = (data) => service({ url: '/shop/cs/refunds/audit', method: 'post', data })
// 执行：export const csRefundRetry = (data) => service({ url: '/shop/cs/refunds/retry', method: 'post', data })
export const csRefundRetry = (data) => service({ url: '/shop/cs/refunds/retry', method: 'post', data })
// 执行：export const csAbnormals = (params) => service({ url: '/shop/cs/abnormals', method: 'get', params })
export const csAbnormals = (params) => service({ url: '/shop/cs/abnormals', method: 'get', params })
// 执行：export const csHandleAbnormal = (data) => service({ url: '/shop/cs/abnormals/handle', method: 'post', data })
export const csHandleAbnormal = (data) => service({ url: '/shop/cs/abnormals/handle', method: 'post', data })
// 执行：export const csNotices = (params) => service({ url: '/shop/cs/notices', method: 'get', params })
export const csNotices = (params) => service({ url: '/shop/cs/notices', method: 'get', params })
// 执行：export const csOrderTrace = (orderNo) => service({ url: `/shop/cs/orders/${orderNo}`, method: 'get' })
export const csOrderTrace = (orderNo) => service({ url: `/shop/cs/orders/${orderNo}`, method: 'get' })
// 执行：export const csChannelMode = (data) => service({ url: '/shop/cs/channel/mode', method: 'post', data })
export const csChannelMode = (data) => service({ url: '/shop/cs/channel/mode', method: 'post', data })
