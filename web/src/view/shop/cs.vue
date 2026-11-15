<!-- <template> -->
<template>
<!-- <div class="shop-page"> -->
  <div class="shop-page">
<!-- <el-card shadow="never" class="mb8"> -->
    <el-card shadow="never" class="mb8">
<!-- <template #header> -->
      <template #header>
<!-- <div class="flex-between"> -->
        <div class="flex-between">
<!-- <span>退款审核（同意后调渠道；失败保留记录可重试；申请期间发货暂停）</span> -->
          <span>退款审核（同意后调渠道；失败保留记录可重试；申请期间发货暂停）</span>
<!-- <div> -->
          <div>
<!-- <el-button size="small" type="danger" @click="setMode('fail')">首次失败（fail）</el-button> -->
            <el-button size="small" type="danger" @click="setMode('fail')">首次失败（fail）</el-button>
<!-- <el-button size="small" type="success" @click="setMode('success')">退款成功（success）</el-button> -->
            <el-button size="small" type="success" @click="setMode('success')">退款成功（success）</el-button>
<!-- <span class="muted">当前：{{ channel.refundMode }} · 累计 {{ channel.refundCalls || 0 }} 次</span> -->
            <span class="muted">当前：{{ channel.refundMode }} · 累计 {{ channel.refundCalls || 0 }} 次</span>
<!-- </div> -->
          </div>
<!-- </div> -->
        </div>
<!-- </template> -->
      </template>
<!-- <el-table :data="refunds" size="small"> -->
      <el-table :data="refunds" size="small">
<!-- <el-table-column prop="refundNo" label="退款单" width="190" /> -->
        <el-table-column prop="refundNo" label="退款单" width="190" />
<!-- <el-table-column prop="orderNo" label="订单" width="190" /> -->
        <el-table-column prop="orderNo" label="订单" width="190" />
<!-- <el-table-column prop="source" label="来源" width="90" /> -->
        <el-table-column prop="source" label="来源" width="90" />
<!-- <el-table-column label="金额" width="90"><template #default="{ row }">{{ fen(row.amountCents) }}</template></el-... -->
        <el-table-column label="金额" width="90"><template #default="{ row }">{{ fen(row.amountCents) }}</template></el-table-column>
<!-- <el-table-column label="审核" width="90"> -->
        <el-table-column label="审核" width="90">
<!-- <template #default="{ row }"><el-tag size="small" :type="tagType(row.auditStatus)">{{ row.auditStatus }}</el-t... -->
          <template #default="{ row }"><el-tag size="small" :type="tagType(row.auditStatus)">{{ row.auditStatus }}</el-tag></template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- <el-table-column label="渠道" width="90"> -->
        <el-table-column label="渠道" width="90">
<!-- <template #default="{ row }"><el-tag size="small" :type="tagType(row.channelStatus)">{{ row.channelStatus }}</... -->
          <template #default="{ row }"><el-tag size="small" :type="tagType(row.channelStatus)">{{ row.channelStatus }}</el-tag></template>
<!-- </el-table-column> -->
        </el-table-column>
<!-- <el-table-column prop="retryCount" label="重试" width="60" /> -->
        <el-table-column prop="retryCount" label="重试" width="60" />
<!-- <el-table-column prop="auditReason" label="原因" min-width="120" /> -->
        <el-table-column prop="auditReason" label="原因" min-width="120" />
<!-- <el-table-column label="操作" width="170"> -->
        <el-table-column label="操作" width="170">
<!-- <template #default="{ row }"> -->
          <template #default="{ row }">
<!-- <el-button v-if="row.auditStatus === 'pending'" size="small" type="success" @click="audit(row, true)">同意</el-b... -->
            <el-button v-if="row.auditStatus === 'pending'" size="small" type="success" @click="audit(row, true)">同意</el-button>
<!-- <el-button v-if="row.auditStatus === 'pending'" size="small" type="danger" @click="audit(row, false)">拒绝</el-b... -->
            <el-button v-if="row.auditStatus === 'pending'" size="small" type="danger" @click="audit(row, false)">拒绝</el-button>
<!-- <el-button v-if="row.auditStatus === 'approved' && row.channelStatus === 'failed'" size="small" type="warning"... -->
            <el-button v-if="row.auditStatus === 'approved' && row.channelStatus === 'failed'" size="small" type="warning" @click="retry(row)">重试退款</el-button>
<!-- </template> -->
          </template>
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
<!-- <template #header><span>异常收款（迟到支付，核实后发起退款）</span></template> -->
          <template #header><span>异常收款（迟到支付，核实后发起退款）</span></template>
<!-- <el-table :data="abnormals" size="small"> -->
          <el-table :data="abnormals" size="small">
<!-- <el-table-column prop="channelTxnId" label="交易号" width="170" /> -->
            <el-table-column prop="channelTxnId" label="交易号" width="170" />
<!-- <el-table-column prop="orderNo" label="订单" width="180" /> -->
            <el-table-column prop="orderNo" label="订单" width="180" />
<!-- <el-table-column label="金额" width="80"><template #default="{ row }">{{ fen(row.amountCents) }}</template></el-... -->
            <el-table-column label="金额" width="80"><template #default="{ row }">{{ fen(row.amountCents) }}</template></el-table-column>
<!-- <el-table-column prop="status" label="状态" width="90" /> -->
            <el-table-column prop="status" label="状态" width="90" />
<!-- <el-table-column prop="reason" label="原因" min-width="110" /> -->
            <el-table-column prop="reason" label="原因" min-width="110" />
<!-- <el-table-column label="操作" width="110"> -->
            <el-table-column label="操作" width="110">
<!-- <template #default="{ row }"> -->
              <template #default="{ row }">
<!-- <el-button v-if="row.status === 'pending'" size="small" type="warning" @click="handleAbn(row)">发起退款</el-button... -->
                <el-button v-if="row.status === 'pending'" size="small" type="warning" @click="handleAbn(row)">发起退款</el-button>
<!-- </template> -->
              </template>
<!-- </el-table-column> -->
            </el-table-column>
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
<!-- <template #header><span>支付通知处理记录（只增不改）</span></template> -->
          <template #header><span>支付通知处理记录（只增不改）</span></template>
<!-- <el-table :data="notices" size="small" max-height="360"> -->
          <el-table :data="notices" size="small" max-height="360">
<!-- <el-table-column prop="orderNo" label="订单" width="180" /> -->
            <el-table-column prop="orderNo" label="订单" width="180" />
<!-- <el-table-column prop="channelTxnId" label="交易号" width="160" /> -->
            <el-table-column prop="channelTxnId" label="交易号" width="160" />
<!-- <el-table-column label="验签" width="60"><template #default="{ row }">{{ row.signOk ? '通过' : '失败' }}</template><... -->
            <el-table-column label="验签" width="60"><template #default="{ row }">{{ row.signOk ? '通过' : '失败' }}</template></el-table-column>
<!-- <el-table-column prop="result" label="结果" width="100" /> -->
            <el-table-column prop="result" label="结果" width="100" />
<!-- <el-table-column prop="reason" label="原因" min-width="110" /> -->
            <el-table-column prop="reason" label="原因" min-width="110" />
<!-- </el-table> -->
          </el-table>
<!-- </el-card> -->
        </el-card>
<!-- </el-col> -->
      </el-col>
<!-- </el-row> -->
    </el-row>
<!-- </div> -->
  </div>
<!-- </template> -->
</template>

<script setup name="ShopCs">
// 执行：import { ref, onMounted } from 'vue'
import { ref, onMounted } from 'vue'
// 执行：import { ElMessage, ElMessageBox } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
// 执行：import {
import {
// 执行：csRefunds, csRefundAudit, csRefundRetry, csAbnormals, csHandleAbnormal, csNotices, csChannelMode
  csRefunds, csRefundAudit, csRefundRetry, csAbnormals, csHandleAbnormal, csNotices, csChannelMode
// 结束当前代码块
} from '@/api/shop'

// 声明：const refunds = ref([])
const refunds = ref([])
// 声明：const abnormals = ref([])
const abnormals = ref([])
// 声明：const notices = ref([])
const notices = ref([])
// 声明：const channel = ref({})
const channel = ref({})

// 声明：const fen = (v) => ((v || 0) / 100).toFixed(2)
const fen = (v) => ((v || 0) / 100).toFixed(2)
// 声明：const tagType = (s) => (['success', 'approved', 'refunded', 'handled', 'accepted'].includes(s) ? 'success'
const tagType = (s) => (['success', 'approved', 'refunded', 'handled', 'accepted'].includes(s) ? 'success'
// 执行：: (['failed', 'rejected', 'abnormal'].includes(s) ? 'danger' : 'warning'))
  : (['failed', 'rejected', 'abnormal'].includes(s) ? 'danger' : 'warning'))

// 声明：const load = async () => {
const load = async () => {
// 声明：const [r, a, n] = await Promise.all([csRefunds({}), csAbnormals({}), csNotices({ limit: 50 })])
  const [r, a, n] = await Promise.all([csRefunds({}), csAbnormals({}), csNotices({ limit: 50 })])
// 执行：refunds.value = r || []
  refunds.value = r || []
// 执行：abnormals.value = a || []
  abnormals.value = a || []
// 执行：notices.value = (n || []).slice(0, 50)
  notices.value = (n || []).slice(0, 50)
// 结束当前代码块
}
// 客服同意/拒绝、异常处理都必须写明原因：取消返回空串（不提交），空白原因就地拦下
const askReason = async (title, tip, def) => {
// 执行：try {
  try {
// 声明：const { value } = await ElMessageBox.prompt(title, tip, {
    const { value } = await ElMessageBox.prompt(title, tip, {
// 执行：inputValue: def,
      inputValue: def,
// 执行：inputValidator: (v) => (v && String(v).trim() ? true : '必须填写原因')
      inputValidator: (v) => (v && String(v).trim() ? true : '必须填写原因')
// 结束当前代码块
    })
// 返回：return String(value || '').trim()
    return String(value || '').trim()
// 结束当前代码块
  } catch {
// 返回：return ''
    return ''
// 结束当前代码块
  }
// 结束当前代码块
}
// 声明：const audit = async (row, approve) => {
const audit = async (row, approve) => {
// 声明：const reason = await askReason(approve ? '审核意见（必填）：' : '拒绝原因（必填）：', '退款审核', approve ? '同意整单退款' : '发货在即，不同意退款')
  const reason = await askReason(approve ? '审核意见（必填）：' : '拒绝原因（必填）：', '退款审核', approve ? '同意整单退款' : '发货在即，不同意退款')
// 条件判断：if (!reason) return
  if (!reason) return
// 执行：await csRefundAudit({ orderNo: row.orderNo, approve, reason })
  await csRefundAudit({ orderNo: row.orderNo, approve, reason })
// 执行：ElMessage.success('审核完成')
  ElMessage.success('审核完成')
// 执行：await load()
  await load()
// 结束当前代码块
}
// 声明：const retry = async (row) => {
const retry = async (row) => {
// 执行：await csRefundRetry({ refundNo: row.refundNo })
  await csRefundRetry({ refundNo: row.refundNo })
// 执行：ElMessage.success('重试完成')
  ElMessage.success('重试完成')
// 执行：await load()
  await load()
// 结束当前代码块
}
// 声明：const handleAbn = async (row) => {
const handleAbn = async (row) => {
// 声明：const reason = await askReason('处理备注（必填）：', '异常收款处理', '核实为订单关闭后到达的支付，发起退款')
  const reason = await askReason('处理备注（必填）：', '异常收款处理', '核实为订单关闭后到达的支付，发起退款')
// 条件判断：if (!reason) return
  if (!reason) return
// 执行：await csHandleAbnormal({ channelTxnId: row.channelTxnId, reason })
  await csHandleAbnormal({ channelTxnId: row.channelTxnId, reason })
// 执行：ElMessage.success('处理完成')
  ElMessage.success('处理完成')
// 执行：await load()
  await load()
// 结束当前代码块
}
// 声明：const setMode = async (refundMode) => {
const setMode = async (refundMode) => {
// 声明：const d = await csChannelMode({ refundMode })
  const d = await csChannelMode({ refundMode })
// 执行：channel.value = d || {}
  channel.value = d || {}
// 执行：ElMessage.success('渠道退款模式已切换为 ' + (d?.refundMode || refundMode))
  ElMessage.success('渠道退款模式已切换为 ' + (d?.refundMode || refundMode))
// 结束当前代码块
}

// 调用：onMounted(async () => { await load() })
onMounted(async () => { await load() })
</script>

<style scoped>
.mb8 { margin-bottom: 8px; }
.flex-between { display: flex; justify-content: space-between; align-items: center; }
.muted { color: #6b7280; font-size: 12px; }
</style>
