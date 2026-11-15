<!-- <template> -->
<template>
<!-- <div class="shop-page"> -->
  <div class="shop-page">
<!-- <el-card shadow="never" class="mb8"> -->
    <el-card shadow="never" class="mb8">
<!-- <template #header><span>可发货订单（仅已付款且无进行中退款；一单只能登记一次）</span></template> -->
      <template #header><span>可发货订单（仅已付款且无进行中退款；一单只能登记一次）</span></template>
<!-- <el-table :data="shipable" size="small"> -->
      <el-table :data="shipable" size="small">
<!-- <el-table-column prop="orderNo" label="订单号" width="190" /> -->
        <el-table-column prop="orderNo" label="订单号" width="190" />
<!-- <el-table-column prop="productName" label="商品" min-width="120" /> -->
        <el-table-column prop="productName" label="商品" min-width="120" />
<!-- <el-table-column prop="qty" label="数量" width="60" /> -->
        <el-table-column prop="qty" label="数量" width="60" />
<!-- <el-table-column label="金额" width="90"><template #default="{ row }">{{ fen(row.amountCents) }}</template></el-... -->
        <el-table-column label="金额" width="90"><template #default="{ row }">{{ fen(row.amountCents) }}</template></el-table-column>
<!-- <el-table-column label="付款时间" width="160"><template #default="{ row }">{{ fmt(row.paidAt) }}</template></el-ta... -->
        <el-table-column label="付款时间" width="160"><template #default="{ row }">{{ fmt(row.paidAt) }}</template></el-table-column>
<!-- <el-table-column label="承运方 / 运单号" min-width="220"> -->
        <el-table-column label="承运方 / 运单号" min-width="220">
<!-- <template #default="{ row }"> -->
          <template #default="{ row }">
<!-- <el-input v-model="row._carrier" size="small" placeholder="承运方" style="width: 100px" /> -->
            <el-input v-model="row._carrier" size="small" placeholder="承运方" style="width: 100px" />
<!-- <el-input v-model="row._trackingNo" size="small" placeholder="运单号" style="width: 120px" /> -->
            <el-input v-model="row._trackingNo" size="small" placeholder="运单号" style="width: 120px" />
<!-- </template> -->
          </template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- <el-table-column label="操作" width="150"> -->
        <el-table-column label="操作" width="150">
<!-- <template #default="{ row }"> -->
          <template #default="{ row }">
<!-- <el-button size="small" type="primary" @click="ship(row)">登记发货</el-button> -->
            <el-button size="small" type="primary" @click="ship(row)">登记发货</el-button>
<!-- <el-button size="small" @click="openTrace(row.orderNo)">详情</el-button> -->
            <el-button size="small" @click="openTrace(row.orderNo)">详情</el-button>
<!-- </template> -->
          </template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- </el-table> -->
      </el-table>
<!-- </el-card> -->
    </el-card>

<!-- <el-card shadow="never"> -->
    <el-card shadow="never">
<!-- <template #header><span>发货记录</span></template> -->
      <template #header><span>发货记录</span></template>
<!-- <el-table :data="shipments" size="small"> -->
      <el-table :data="shipments" size="small">
<!-- <el-table-column prop="orderNo" label="订单号" width="190" /> -->
        <el-table-column prop="orderNo" label="订单号" width="190" />
<!-- <el-table-column prop="carrier" label="承运方" width="140" /> -->
        <el-table-column prop="carrier" label="承运方" width="140" />
<!-- <el-table-column prop="trackingNo" label="运单号" min-width="140" /> -->
        <el-table-column prop="trackingNo" label="运单号" min-width="140" />
<!-- <el-table-column label="登记时间" width="160"><template #default="{ row }">{{ fmt(row.createdAt) }}</template></el... -->
        <el-table-column label="登记时间" width="160"><template #default="{ row }">{{ fmt(row.createdAt) }}</template></el-table-column>
<!-- </el-table> -->
      </el-table>
<!-- </el-card> -->
    </el-card>

<!-- <el-dialog v-model="traceVisible" title="订单详情" width="680px"> -->
    <el-dialog v-model="traceVisible" title="订单详情" width="680px">
<!-- <el-descriptions v-if="traceData" :column="2" size="small" border> -->
      <el-descriptions v-if="traceData" :column="2" size="small" border>
<!-- <el-descriptions-item label="订单号">{{ traceData.order.orderNo }}</el-descriptions-item> -->
        <el-descriptions-item label="订单号">{{ traceData.order.orderNo }}</el-descriptions-item>
<!-- <el-descriptions-item label="状态">{{ traceData.order.status }}</el-descriptions-item> -->
        <el-descriptions-item label="状态">{{ traceData.order.status }}</el-descriptions-item>
<!-- <el-descriptions-item label="商品">{{ traceData.order.productName }} × {{ traceData.order.qty }}</el-description... -->
        <el-descriptions-item label="商品">{{ traceData.order.productName }} × {{ traceData.order.qty }}</el-descriptions-item>
<!-- <el-descriptions-item label="金额">{{ fen(traceData.order.amountCents) }}</el-descriptions-item> -->
        <el-descriptions-item label="金额">{{ fen(traceData.order.amountCents) }}</el-descriptions-item>
<!-- <el-descriptions-item label="可发货">{{ traceData.shipableNow ? '是' : '否' }}</el-descriptions-item> -->
        <el-descriptions-item label="可发货">{{ traceData.shipableNow ? '是' : '否' }}</el-descriptions-item>
<!-- <el-descriptions-item v-if="traceData.shipment" label="已发货">{{ traceData.shipment.carrier }} / {{ traceData.sh... -->
        <el-descriptions-item v-if="traceData.shipment" label="已发货">{{ traceData.shipment.carrier }} / {{ traceData.shipment.trackingNo }}</el-descriptions-item>
<!-- </el-descriptions> -->
      </el-descriptions>
<!-- </el-dialog> -->
    </el-dialog>
<!-- </div> -->
  </div>
<!-- </template> -->
</template>

<script setup name="ShopWarehouse">
// 执行：import { ref, onMounted } from 'vue'
import { ref, onMounted } from 'vue'
// 执行：import { ElMessage } from 'element-plus'
import { ElMessage } from 'element-plus'
// 执行：import { whShipable, whShip, whShipments, whOrderTrace } from '@/api/shop'
import { whShipable, whShip, whShipments, whOrderTrace } from '@/api/shop'

// 声明：const shipable = ref([])
const shipable = ref([])
// 声明：const shipments = ref([])
const shipments = ref([])
// 声明：const traceData = ref(null)
const traceData = ref(null)
// 声明：const traceVisible = ref(false)
const traceVisible = ref(false)

// 声明：const fen = (v) => ((v || 0) / 100).toFixed(2)
const fen = (v) => ((v || 0) / 100).toFixed(2)
// 声明：const fmt = (t) => (t ? new Date(t).toLocaleString() : '-')
const fmt = (t) => (t ? new Date(t).toLocaleString() : '-')

// 声明：const load = async () => {
const load = async () => {
// 声明：const [a, s] = await Promise.all([whShipable(), whShipments()])
  const [a, s] = await Promise.all([whShipable(), whShipments()])
// 执行：shipable.value = (a || []).map(x => ({ ...x, _carrier: '', _trackingNo: '' }))
  shipable.value = (a || []).map(x => ({ ...x, _carrier: '', _trackingNo: '' }))
// 执行：shipments.value = s || []
  shipments.value = s || []
// 结束当前代码块
}
// 声明：const ship = async (row) => {
const ship = async (row) => {
// 条件判断：if (!row._carrier || !row._trackingNo) {
  if (!row._carrier || !row._trackingNo) {
// 执行：ElMessage.warning('请填写承运方与运单号')
    ElMessage.warning('请填写承运方与运单号')
// 返回：return
    return
// 结束当前代码块
  }
// 执行：await whShip({ orderNo: row.orderNo, carrier: row._carrier, trackingNo: row._trackingNo })
  await whShip({ orderNo: row.orderNo, carrier: row._carrier, trackingNo: row._trackingNo })
// 执行：ElMessage.success('发货登记成功')
  ElMessage.success('发货登记成功')
// 执行：await load()
  await load()
// 结束当前代码块
}
// 声明：const openTrace = async (orderNo) => {
const openTrace = async (orderNo) => {
// 执行：traceData.value = await whOrderTrace(orderNo)
  traceData.value = await whOrderTrace(orderNo)
// 执行：traceVisible.value = true
  traceVisible.value = true
// 结束当前代码块
}

// 调用：onMounted(load)
onMounted(load)
</script>

<style scoped>
.mb8 { margin-bottom: 8px; }
</style>
