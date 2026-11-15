# 阶段3：R5 退款与发货互斥 + R6 退款结果与数量恢复 + R7 发货约束（覆盖场景 S3）
. D:\Users\Lenovo\Desktop\shop\shop\acceptance-r2\_common.ps1
Write-Output "########## PHASE3: R5 / R6 / R7 (S3) ##########"

Reset
$b1 = Login 'buyer1'; $wh = Login 'wh1'; $cs = Login 'cs1'
# 两笔已付款订单：o1 走“批准退款（首次失败→重试成功）”，o2 走“拒绝退款后发货”
$o1 = Json (NewOrder $b1 1 1)
$o1no = $o1.data.order.orderNo
Notify $o1no 'TXN-P1' 10000 (PaySign $o1no 'TXN-P1' 10000) | Out-Null
$o2 = Json (NewOrder $b1 1 1)
$o2no = $o2.data.order.orderNo
Notify $o2no 'TXN-P2' 10000 (PaySign $o2no 'TXN-P2' 10000) | Out-Null
Check 'S3-setup-two-paid' ((Avail) -eq 3) ('available after 2 paid orders=' + (Avail))

# ---- o1：售前退款（批准 + 首次渠道失败 + 重试成功）----
# 买家申请整单退款
$r = Req Post ('/api/orders/' + $o1no + '/refund') '{"reason":"no-longer-want"}' $b1
$j = Json $r
$rn = $j.data.refund.refundNo
Check 'R5-1-refund-applied' (($r.code -eq 200) -and ($j.data.refundPending -eq $true)) ('http=' + $r.code + ' refundNo=' + $rn + ' pending=' + $j.data.refundPending)
# R5-1：申请期间发货被拒
$r = Req Post '/api/warehouse/ship' ('{"orderNo":"' + $o1no + '","carrier":"SF","trackingNo":"SF-BLOCK"}') $wh
Check 'R5-1-ship-blocked-while-refunding' ($r.code -ge 400) ('http=' + $r.code + ' body=' + $r.body)
# 先把模拟渠道退款结果切成 fail，审核时会立即发起第一次渠道退款
SetRefundMode 'fail' $cs | Out-Null
# R5-2：客服批准（写明原因）；渠道失败 → 审核生效且退款单保留 FAILED
$r = Req Post '/api/cs/refunds/audit' ('{"orderNo":"' + $o1no + '","approve":true,"reason":"agree-pre-ship-refund"}') $cs
$j = Json $r
$st = $j.data.auditStatus
$cs1 = $j.data.channelStatus
if (-not $st) { $st = $j.data.refund.auditStatus; $cs1 = $j.data.refund.channelStatus }
Check 'R5-2-audit-approved' ($st -eq 'approved') ('http=' + $r.code + ' auditStatus=' + $st)
Check 'R6-3-first-refund-fail-kept' ($cs1 -eq 'failed') ('channelStatus=' + $cs1)
$availMid = Avail
# 渠道仍为 fail：重试再次失败，记录保留且不恢复库存
$r = Req Post '/api/cs/refunds/retry' ('{"refundNo":"' + $rn + '"}') $cs
$j = Json $r
$cs1b = $j.data.channelStatus
if (-not $cs1b) { $cs1b = $j.data.refund.channelStatus }
Check 'R6-3-retry-under-fail-still-failed' ($cs1b -eq 'failed') ('http=' + $r.code + ' channelStatus=' + $cs1b)
Check 'R6-3-fail-no-restore' ((Avail) -eq $availMid) ('available=' + (Avail) + ' expect=' + $availMid)
# 渠道切 success：重试成功，恢复 1 件（仅一次）
SetRefundMode 'success' $cs | Out-Null
$r = Req Post '/api/cs/refunds/retry' ('{"refundNo":"' + $rn + '"}') $cs
$j = Json $r
$cs1c = $j.data.channelStatus
if (-not $cs1c) { $cs1c = $j.data.refund.channelStatus }
Check 'R6-3-retry-success' ($cs1c -eq 'success') ('http=' + $r.code + ' channelStatus=' + $cs1c)
# R6-2：售前退款恢复 1 件（只恢复一次）
Check 'R6-2-stock-restored-once' ((Avail) -eq ($availMid + 1)) ('before=' + $availMid + ' after=' + (Avail))
# R6-1：买家可见退款结果
$rf = (Json (Req Get ('/api/orders/' + $o1no + '/refund') $null $b1)).data
Check 'R6-1-buyer-sees-refund-success' ($rf.channelStatus -eq 'success') ('channelStatus=' + $rf.channelStatus)
$lst = (Json (Req Get '/api/orders' $null $b1)).data
$ord1 = $lst | Where-Object { $_.orderNo -eq $o1no } | Select-Object -First 1
Check 'S3-order-refunded' ($ord1.status -eq 'refunded') ('status=' + $ord1.status)
# 重复重试退款 → 拒绝（防重复退款，库存不再变动）
$r = Req Post '/api/cs/refunds/retry' ('{"refundNo":"' + $rn + '"}') $cs
Check 'R6-3-triple-retry-rejected' ($r.code -ge 400) ('http=' + $r.code + ' body=' + $r.body)
Check 'R6-3-no-extra-restore' ((Avail) -eq ($availMid + 1)) ('available=' + (Avail))

# ---- o2：拒绝退款 → 恢复发货 → 已发货不可退款 ----
$r = Req Post ('/api/orders/' + $o2no + '/refund') '{"reason":"mis-order"}' $b1
Check 'R5-o2-applied' ($r.code -eq 200) ('http=' + $r.code)
# R5-2：拒绝并写明原因
$r = Req Post '/api/cs/refunds/audit' ('{"orderNo":"' + $o2no + '","approve":false,"reason":"reject-ship-soon"}') $cs
$j = Json $r
$st2 = $j.data.auditStatus
if (-not $st2) { $st2 = $j.data.refund.auditStatus }
Check 'R5-2-audit-rejected-with-reason' ($st2 -eq 'rejected') ('auditStatus=' + $st2)
# R5-3：拒绝后允许继续发货
$r = Req Post '/api/warehouse/ship' ('{"orderNo":"' + $o2no + '","carrier":"SF","trackingNo":"SF-002"}') $wh
$j = Json $r
Check 'R5-3-ship-after-reject-ok' (($r.code -eq 200) -and ($j.data.status -eq 'shipped')) ('http=' + $r.code + ' status=' + $j.data.status)
# R7-1：重复发货被拒
$r = Req Post '/api/warehouse/ship' ('{"orderNo":"' + $o2no + '","carrier":"SF","trackingNo":"SF-003"}') $wh
Check 'R7-1-double-ship-rejected' ($r.code -ge 400) ('http=' + $r.code + ' body=' + $r.body)
# R7-2：缺承运方/运单号被拒
$o3 = Json (NewOrder $b1 1 1)
$o3no = $o3.data.order.orderNo
Notify $o3no 'TXN-P3' 10000 (PaySign $o3no 'TXN-P3' 10000) | Out-Null
$r = Req Post '/api/warehouse/ship' ('{"orderNo":"' + $o3no + '","carrier":"","trackingNo":""}') $wh
Check 'R7-2-missing-fields-rejected' ($r.code -ge 400) ('http=' + $r.code)
# R7-3：已发货订单不受理售前退款
$r = Req Post ('/api/orders/' + $o2no + '/refund') '{"reason":"already-shipped"}' $b1
Check 'R7-3-shipped-refund-rejected' ($r.code -ge 400) ('http=' + $r.code + ' body=' + $r.body)
SetRefundMode 'external' $cs | Out-Null
Summary 'PHASE3'
