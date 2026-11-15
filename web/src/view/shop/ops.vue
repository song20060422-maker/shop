<!-- <template> -->
<template>
<!-- <div class="shop-page"> -->
  <div class="shop-page">
<!-- <el-card shadow="never" class="mb8"> -->
    <el-card shadow="never" class="mb8">
<!-- <template #header><span>交易汇总</span></template> -->
      <template #header><span>交易汇总</span></template>
<!-- <el-row :gutter="16"> -->
      <el-row :gutter="16">
<!-- <el-col :span="6"><el-statistic title="订单总数" :value="summary.orderCount || 0" /></el-col> -->
        <el-col :span="6"><el-statistic title="订单总数" :value="summary.orderCount || 0" /></el-col>
<!-- <el-col :span="6"><el-statistic title="已收款（元）" :value="fen(summary.totalPaidCents)" :precision="2" /></el-col> -->
        <el-col :span="6"><el-statistic title="已收款（元）" :value="fen(summary.totalPaidCents)" :precision="2" /></el-col>
<!-- <el-col :span="6"><el-statistic title="已退款（元）" :value="fen(summary.totalRefundedCents)" :precision="2" /></el-... -->
        <el-col :span="6"><el-statistic title="已退款（元）" :value="fen(summary.totalRefundedCents)" :precision="2" /></el-col>
<!-- <el-col :span="6"><el-statistic title="净收款（元）" :value="fen((summary.totalPaidCents || 0) - (summary.totalRefun... -->
        <el-col :span="6"><el-statistic title="净收款（元）" :value="fen((summary.totalPaidCents || 0) - (summary.totalRefundedCents || 0))" :precision="2" /></el-col>
<!-- </el-row> -->
      </el-row>
<!-- </el-card> -->
    </el-card>

<!-- <el-card shadow="never" class="mb8"> -->
    <el-card shadow="never" class="mb8">
<!-- <template #header><span>商品与库存（可售 = 初始库存 - 未释放占用；改价不影响已有订单快照）</span></template> -->
      <template #header><span>商品与库存（可售 = 初始库存 - 未释放占用；改价不影响已有订单快照）</span></template>
<!-- <el-table :data="products" size="small"> -->
      <el-table :data="products" size="small">
<!-- <el-table-column prop="sku" label="SKU" width="90" /> -->
        <el-table-column prop="sku" label="SKU" width="90" />
<!-- <el-table-column prop="name" label="商品" min-width="120" /> -->
        <el-table-column prop="name" label="商品" min-width="120" />
<!-- <el-table-column label="单价(分)" width="120"> -->
        <el-table-column label="单价(分)" width="120">
<!-- <template #default="{ row }"> -->
          <template #default="{ row }">
<!-- <el-input v-model="row._price" size="small" style="width: 90px" /> -->
            <el-input v-model="row._price" size="small" style="width: 90px" />
<!-- </template> -->
          </template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- <el-table-column label="可售" width="110"> -->
        <el-table-column label="可售" width="110">
<!-- <template #default="{ row }"> -->
          <template #default="{ row }">
<!-- <el-input v-model="row._available" size="small" style="width: 80px" /> -->
            <el-input v-model="row._available" size="small" style="width: 80px" />
<!-- </template> -->
          </template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- <el-table-column prop="initialStock" label="初始" width="70" /> -->
        <el-table-column prop="initialStock" label="初始" width="70" />
<!-- <el-table-column label="在售" width="80"> -->
        <el-table-column label="在售" width="80">
<!-- <template #default="{ row }"> -->
          <template #default="{ row }">
<!-- <el-tag :type="row.onSale ? 'success' : 'info'" size="small">{{ row.onSale ? '在售' : '下架' }}</el-tag> -->
            <el-tag :type="row.onSale ? 'success' : 'info'" size="small">{{ row.onSale ? '在售' : '下架' }}</el-tag>
<!-- </template> -->
          </template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- <el-table-column label="操作" width="220"> -->
        <el-table-column label="操作" width="220">
<!-- <template #default="{ row }"> -->
          <template #default="{ row }">
<!-- <el-button size="small" @click="saveProduct(row)">保存调价</el-button> -->
            <el-button size="small" @click="saveProduct(row)">保存调价</el-button>
<!-- <el-button size="small" :type="row.onSale ? 'danger' : 'success'" @click="toggleSale(row)"> -->
            <el-button size="small" :type="row.onSale ? 'danger' : 'success'" @click="toggleSale(row)">
<!-- {{ row.onSale ? '下架' : '上架' }} -->
              {{ row.onSale ? '下架' : '上架' }}
<!-- </el-button> -->
            </el-button>
<!-- </template> -->
          </template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- </el-table> -->
      </el-table>
<!-- </el-card> -->
    </el-card>

<!-- <el-card shadow="never" class="mb8"> -->
    <el-card shadow="never" class="mb8">
<!-- <template #header> -->
      <template #header>
<!-- <div class="flex-between"> -->
        <div class="flex-between">
<!-- <span>订单核查</span> -->
          <span>订单核查</span>
<!-- <el-select v-model="statusFilter" size="small" clearable placeholder="按状态筛选" style="width: 160px" @change="loa... -->
          <el-select v-model="statusFilter" size="small" clearable placeholder="按状态筛选" style="width: 160px" @change="loadOrders">
<!-- <el-option v-for="(v, k) in statusMap" :key="k" :label="v" :value="k" /> -->
            <el-option v-for="(v, k) in statusMap" :key="k" :label="v" :value="k" />
<!-- </el-select> -->
          </el-select>
<!-- </div> -->
        </div>
<!-- </template> -->
      </template>
<!-- <el-table :data="orders" size="small"> -->
      <el-table :data="orders" size="small">
<!-- <el-table-column prop="orderNo" label="订单号" width="190" /> -->
        <el-table-column prop="orderNo" label="订单号" width="190" />
<!-- <el-table-column prop="buyerId" label="买家ID" width="80" /> -->
        <el-table-column prop="buyerId" label="买家ID" width="80" />
<!-- <el-table-column prop="productName" label="商品" min-width="120" /> -->
        <el-table-column prop="productName" label="商品" min-width="120" />
<!-- <el-table-column prop="qty" label="数量" width="60" /> -->
        <el-table-column prop="qty" label="数量" width="60" />
<!-- <el-table-column label="金额" width="90"><template #default="{ row }">{{ fen(row.amountCents) }}</template></el-... -->
        <el-table-column label="金额" width="90"><template #default="{ row }">{{ fen(row.amountCents) }}</template></el-table-column>
<!-- <el-table-column label="状态" width="110"> -->
        <el-table-column label="状态" width="110">
<!-- <template #default="{ row }"><el-tag size="small" :type="statusTag(row.status)">{{ statusMap[row.status] || ro... -->
          <template #default="{ row }"><el-tag size="small" :type="statusTag(row.status)">{{ statusMap[row.status] || row.status }}</el-tag></template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- <el-table-column label="创建时间" width="160"><template #default="{ row }">{{ fmt(row.createdAt) }}</template></el... -->
        <el-table-column label="创建时间" width="160"><template #default="{ row }">{{ fmt(row.createdAt) }}</template></el-table-column>
<!-- <el-table-column label="操作" width="90"> -->
        <el-table-column label="操作" width="90">
<!-- <template #default="{ row }"><el-button size="small" @click="openTrace(row.orderNo)">核查</el-button></template> -->
          <template #default="{ row }"><el-button size="small" @click="openTrace(row.orderNo)">核查</el-button></template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- </el-table> -->
      </el-table>
<!-- </el-card> -->
    </el-card>

<!-- <el-row :gutter="12"> -->
    <el-row :gutter="12">
<!-- <el-col :span="12"> -->
      <el-col :span="12">
<!-- <el-card shadow="never" class="mb8"> -->
        <el-card shadow="never" class="mb8">
<!-- <template #header><span>库存流水（占用/释放/恢复守恒）</span></template> -->
          <template #header><span>库存流水（占用/释放/恢复守恒）</span></template>
<!-- <el-table :data="ledger" size="small" max-height="320"> -->
          <el-table :data="ledger" size="small" max-height="320">
<!-- <el-table-column prop="productId" label="商品" width="70" /> -->
            <el-table-column prop="productId" label="商品" width="70" />
<!-- <el-table-column prop="orderNo" label="订单" width="180" /> -->
            <el-table-column prop="orderNo" label="订单" width="180" />
<!-- <el-table-column prop="changeType" label="类型" width="90" /> -->
            <el-table-column prop="changeType" label="类型" width="90" />
<!-- <el-table-column label="变动" width="70"> -->
            <el-table-column label="变动" width="70">
<!-- <template #default="{ row }">{{ row.delta > 0 ? '+' + row.delta : row.delta }}</template> -->
              <template #default="{ row }">{{ row.delta > 0 ? '+' + row.delta : row.delta }}</template>
<!-- </el-table-column> -->
            </el-table-column>
<!-- <el-table-column prop="availableAfter" label="结余" width="70" /> -->
            <el-table-column prop="availableAfter" label="结余" width="70" />
<!-- <el-table-column label="时间"><template #default="{ row }">{{ fmt(row.createdAt) }}</template></el-table-column> -->
            <el-table-column label="时间"><template #default="{ row }">{{ fmt(row.createdAt) }}</template></el-table-column>
<!-- </el-table> -->
          </el-table>
<!-- </el-card> -->
        </el-card>
<!-- </el-col> -->
      </el-col>
<!-- <el-col :span="12"> -->
      <el-col :span="12">
<!-- <el-card shadow="never" class="mb8"> -->
        <el-card shadow="never" class="mb8">
<!-- <template #header><span>支付通知记录（含失败通知，只增不改）</span></template> -->
          <template #header><span>支付通知记录（含失败通知，只增不改）</span></template>
<!-- <el-table :data="notices" size="small" max-height="320"> -->
          <el-table :data="notices" size="small" max-height="320">
<!-- <el-table-column prop="orderNo" label="订单" width="180" /> -->
            <el-table-column prop="orderNo" label="订单" width="180" />
<!-- <el-table-column prop="channelTxnId" label="交易号" width="160" /> -->
            <el-table-column prop="channelTxnId" label="交易号" width="160" />
<!-- <el-table-column label="验签" width="60"> -->
            <el-table-column label="验签" width="60">
<!-- <template #default="{ row }">{{ row.signOk ? '通过' : '失败' }}</template> -->
              <template #default="{ row }">{{ row.signOk ? '通过' : '失败' }}</template>
<!-- </el-table-column> -->
            </el-table-column>
<!-- <el-table-column prop="result" label="结果" width="100" /> -->
            <el-table-column prop="result" label="结果" width="100" />
<!-- <el-table-column prop="reason" label="原因" min-width="120" /> -->
            <el-table-column prop="reason" label="原因" min-width="120" />
<!-- </el-table> -->
          </el-table>
<!-- </el-card> -->
        </el-card>
<!-- </el-col> -->
      </el-col>
<!-- </el-row> -->
    </el-row>

<!-- <el-dialog v-model="traceVisible" title="订单全链路核查" width="760px"> -->
    <el-dialog v-model="traceVisible" title="订单全链路核查" width="760px">
<!-- <template v-if="traceData"> -->
      <template v-if="traceData">
<!-- <el-descriptions :column="2" size="small" border> -->
        <el-descriptions :column="2" size="small" border>
<!-- <el-descriptions-item label="订单号">{{ traceData.order.orderNo }}</el-descriptions-item> -->
          <el-descriptions-item label="订单号">{{ traceData.order.orderNo }}</el-descriptions-item>
<!-- <el-descriptions-item label="状态">{{ statusMap[traceData.order.status] || traceData.order.status }}</el-descrip... -->
          <el-descriptions-item label="状态">{{ statusMap[traceData.order.status] || traceData.order.status }}</el-descriptions-item>
<!-- <el-descriptions-item label="商品">{{ traceData.order.productName }} × {{ traceData.order.qty }}</el-description... -->
          <el-descriptions-item label="商品">{{ traceData.order.productName }} × {{ traceData.order.qty }}</el-descriptions-item>
<!-- <el-descriptions-item label="金额">{{ fen(traceData.order.amountCents) }}</el-descriptions-item> -->
          <el-descriptions-item label="金额">{{ fen(traceData.order.amountCents) }}</el-descriptions-item>
<!-- <el-descriptions-item label="付款时间">{{ traceData.order.paidAt ? fmt(traceData.order.paidAt) : '未付款' }}</el-desc... -->
          <el-descriptions-item label="付款时间">{{ traceData.order.paidAt ? fmt(traceData.order.paidAt) : '未付款' }}</el-descriptions-item>
<!-- <el-descriptions-item label="可发货">{{ traceData.shipableNow ? '是' : '否' }}{{ traceData.refundPending ? '（退款申请中，... -->
          <el-descriptions-item label="可发货">{{ traceData.shipableNow ? '是' : '否' }}{{ traceData.refundPending ? '（退款申请中，发货暂停）' : '' }}</el-descriptions-item>
<!-- <el-descriptions-item v-if="traceData.payment" label="收款交易号">{{ traceData.payment.channelTxnId }}</el-descript... -->
          <el-descriptions-item v-if="traceData.payment" label="收款交易号">{{ traceData.payment.channelTxnId }}</el-descriptions-item>
<!-- <el-descriptions-item v-if="traceData.shipment" label="发货">{{ traceData.shipment.carrier }} / {{ traceData.shi... -->
          <el-descriptions-item v-if="traceData.shipment" label="发货">{{ traceData.shipment.carrier }} / {{ traceData.shipment.trackingNo }}</el-descriptions-item>
<!-- </el-descriptions> -->
        </el-descriptions>
<!-- <template v-if="traceData.refund"> -->
        <template v-if="traceData.refund">
<!-- <h4>退款：{{ traceData.refund.refundNo }}（审核 {{ traceData.refund.auditStatus }} / 渠道 {{ traceData.refund.channelS... -->
          <h4>退款：{{ traceData.refund.refundNo }}（审核 {{ traceData.refund.auditStatus }} / 渠道 {{ traceData.refund.channelStatus }}）</h4>
<!-- </template> -->
        </template>
<!-- <template v-if="traceData.abnormal && traceData.abnormal.length"> -->
        <template v-if="traceData.abnormal && traceData.abnormal.length">
<!-- <h4>异常收款</h4> -->
          <h4>异常收款</h4>
<!-- <div v-for="a in traceData.abnormal" :key="a.id" class="muted">{{ a.channelTxnId }} · {{ fen(a.amountCents) }}... -->
          <div v-for="a in traceData.abnormal" :key="a.id" class="muted">{{ a.channelTxnId }} · {{ fen(a.amountCents) }} · {{ a.status }} · {{ a.reason }}</div>
<!-- </template> -->
        </template>
<!-- <h4>操作轨迹</h4> -->
        <h4>操作轨迹</h4>
<!-- <el-input type="textarea" :rows="6" readonly :model-value="timelineText" /> -->
        <el-input type="textarea" :rows="6" readonly :model-value="timelineText" />
<!-- <h4>支付通知</h4> -->
        <h4>支付通知</h4>
<!-- <el-input type="textarea" :rows="4" readonly :model-value="noticesText" /> -->
        <el-input type="textarea" :rows="4" readonly :model-value="noticesText" />
<!-- </template> -->
      </template>
<!-- </el-dialog> -->
    </el-dialog>
<!-- </div> -->
  </div>
<!-- </template> -->
</template>

<script setup name="ShopOps">
// 执行：import { ref, computed, onMounted } from 'vue'
import { ref, computed, onMounted } from 'vue'
// 执行：import { ElMessage } from 'element-plus'
import { ElMessage } from 'element-plus'
// 执行：import {
import {
// 执行：opsSummary, opsProducts, opsUpdateProduct, opsOrders, opsOrderTrace,
  opsSummary, opsProducts, opsUpdateProduct, opsOrders, opsOrderTrace,
// 执行：opsNotices, opsAbnormals, opsLedger
  opsNotices, opsAbnormals, opsLedger
// 结束当前代码块
} from '@/api/shop'

// 声明：const summary = ref({})
const summary = ref({})
// 声明：const products = ref([])
const products = ref([])
// 声明：const orders = ref([])
const orders = ref([])
// 声明：const ledger = ref([])
const ledger = ref([])
// 声明：const notices = ref([])
const notices = ref([])
// 声明：const statusFilter = ref('')
const statusFilter = ref('')
// 声明：const traceData = ref(null)
const traceData = ref(null)
// 声明：const traceVisible = ref(false)
const traceVisible = ref(false)

// 声明：const statusMap = {
const statusMap = {
// 执行：pending_payment: '待付款', paid: '待发货', shipped: '已发货', cancelled: '已取消',
  pending_payment: '待付款', paid: '待发货', shipped: '已发货', cancelled: '已取消',
// 执行：closed: '已关闭', refund_requested: '退款申请中', refund_failed: '退款失败待重试', refunded: '已退款'
  closed: '已关闭', refund_requested: '退款申请中', refund_failed: '退款失败待重试', refunded: '已退款'
// 结束当前代码块
}
// 声明：const statusTag = (s) => (['cancelled', 'closed'].includes(s) ? 'info' : ['refunded', 'shipped'].includes(s) ?...
const statusTag = (s) => (['cancelled', 'closed'].includes(s) ? 'info' : ['refunded', 'shipped'].includes(s) ? 'success' : ['refund_failed'].includes(s) ? 'danger' : 'warning')

// 声明：const fen = (v) => ((v || 0) / 100).toFixed(2)
const fen = (v) => ((v || 0) / 100).toFixed(2)
// 声明：const fmt = (t) => (t ? new Date(t).toLocaleString() : '-')
const fmt = (t) => (t ? new Date(t).toLocaleString() : '-')

// 声明：const timelineText = computed(() =>
const timelineText = computed(() =>
// 执行：(traceData.value?.timeline || []).map(t => `${fmt(t.createdAt)} [${t.actorRole}] ${t.action} ${t.detail || ''}...
  (traceData.value?.timeline || []).map(t => `${fmt(t.createdAt)} [${t.actorRole}] ${t.action} ${t.detail || ''}`).join('\n'))
// 声明：const noticesText = computed(() =>
const noticesText = computed(() =>
// 执行：(traceData.value?.notices || []).map(n => `${fmt(n.createdAt)} ${n.channelTxnId} ${n.amountCents}分 验签${n.signO...
  (traceData.value?.notices || []).map(n => `${fmt(n.createdAt)} ${n.channelTxnId} ${n.amountCents}分 验签${n.signOk ? '通过' : '失败'} → ${n.result} ${n.reason || ''}`).join('\n'))

// 声明：const loadData = async () => {
const loadData = async () => {
// 声明：const [s, p, l, n] = await Promise.all([opsSummary(), opsProducts(), opsLedger({}), opsNotices({ limit: 50 })]...
  const [s, p, l, n] = await Promise.all([opsSummary(), opsProducts(), opsLedger({}), opsNotices({ limit: 50 })])
// 执行：summary.value = s
  summary.value = s
// 执行：products.value = (p || []).map(x => ({ ...x, _price: x.priceCents, _available: x.available }))
  products.value = (p || []).map(x => ({ ...x, _price: x.priceCents, _available: x.available }))
// 执行：ledger.value = (l || []).slice(0, 50)
  ledger.value = (l || []).slice(0, 50)
// 执行：notices.value = (n || []).slice(0, 50)
  notices.value = (n || []).slice(0, 50)
// 执行：await loadOrders()
  await loadOrders()
// 结束当前代码块
}
// 声明：const loadOrders = async () => {
const loadOrders = async () => {
// 声明：const list = await opsOrders(statusFilter.value ? { status: statusFilter.value } : {})
  const list = await opsOrders(statusFilter.value ? { status: statusFilter.value } : {})
// 执行：orders.value = list || []
  orders.value = list || []
// 结束当前代码块
}
// 声明：const saveProduct = async (row) => {
const saveProduct = async (row) => {
// 执行：await opsUpdateProduct(row.id, { priceCents: Number(row._price), available: Number(row._available) })
  await opsUpdateProduct(row.id, { priceCents: Number(row._price), available: Number(row._available) })
// 执行：ElMessage.success('已保存（已有订单金额快照不变）')
  ElMessage.success('已保存（已有订单金额快照不变）')
// 执行：await loadData()
  await loadData()
// 结束当前代码块
}
// 声明：const toggleSale = async (row) => {
const toggleSale = async (row) => {
// 执行：await opsUpdateProduct(row.id, { onSale: !row.onSale })
  await opsUpdateProduct(row.id, { onSale: !row.onSale })
// 执行：ElMessage.success(row.onSale ? '已下架' : '已上架')
  ElMessage.success(row.onSale ? '已下架' : '已上架')
// 执行：await loadData()
  await loadData()
// 结束当前代码块
}
// 声明：const openTrace = async (orderNo) => {
const openTrace = async (orderNo) => {
// 执行：traceData.value = await opsOrderTrace(orderNo)
  traceData.value = await opsOrderTrace(orderNo)
// 执行：traceVisible.value = true
  traceVisible.value = true
// 结束当前代码块
}

// 调用：onMounted(loadData)
onMounted(loadData)
</script>

<style scoped>
.mb8 { margin-bottom: 8px; }
.flex-between { display: flex; justify-content: space-between; align-items: center; }
h4 { margin: 12px 0 6px; }
.muted { color: #6b7280; font-size: 12px; line-height: 1.8; }
</style>
