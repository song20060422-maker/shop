# 阶段2：R3 支付通知有效性 + R4 迟到支付
. D:\Users\Lenovo\Desktop\shop\shop\acceptance-r2\_common.ps1
Write-Output "########## PHASE2: R3 / R4 ##########"

Reset
$b1 = Login 'buyer1'; $ops = Login 'ops1'; $cs = Login 'cs1'
$o1 = Json (NewOrder $b1 1 1)
$o1no = $o1.data.order.orderNo
# R3-1：伪造签名被拒
$r = Notify $o1no 'TXN-BAD' 10000 'deadbeef'
$j = Json $r
$rejected = ($r.code -ge 400) -or ($j.data.result -eq 'rejected')
Check 'R3-1-bad-sign-rejected' ($rejected) ('http=' + $r.code + ' result=' + $j.data.result + ' reason=' + $j.data.reason)
# 处理结果可查询（运营通知列表）
$notices = (Json (Req Get '/api/ops/notices' $null $ops)).data
$bad = $notices | Where-Object { $_.channelTxnId -eq 'TXN-BAD' } | Select-Object -First 1
Check 'R3-1-bad-sign-queryable' ($null -ne $bad) ('notice result=' + $bad.result + ' reason=' + $bad.reason)
# R3-2：金额不符（签名对金额 9999 有效，但与订单 10000 不匹配）
$r = Notify $o1no 'TXN-AMT' 9999 (PaySign $o1no 'TXN-AMT' 9999)
$j = Json $r
$rejected = ($r.code -ge 400) -or ($j.data.result -eq 'rejected')
Check 'R3-2-amount-mismatch-rejected' ($rejected) ('http=' + $r.code + ' result=' + $j.data.result)
# R3-2：未知订单号
$r = Notify 'NO-SUCH-ORDER' 'TXN-UNK' 10000 (PaySign 'NO-SUCH-ORDER' 'TXN-UNK' 10000)
$j = Json $r
$rejected = ($r.code -ge 400) -or ($j.data.result -eq 'rejected')
Check 'R3-2-unknown-order-rejected' ($rejected) ('http=' + $r.code + ' result=' + $j.data.result)
$notices = (Json (Req Get '/api/ops/notices' $null $ops)).data
$unk = $notices | Where-Object { $_.channelTxnId -eq 'TXN-UNK' } | Select-Object -First 1
$amt = $notices | Where-Object { $_.channelTxnId -eq 'TXN-AMT' } | Select-Object -First 1
Check 'R3-2-illegal-queryable' ($null -ne $unk -and $null -ne $amt) ('unk=' + $unk.result + ' amt=' + $amt.result)
# R3-3：有效通知受理 → 订单待发货；重复通知幂等
$r = Notify $o1no 'TXN-OK1' 10000 (PaySign $o1no 'TXN-OK1' 10000)
$j = Json $r
Check 'R3-3-valid-notify-accepted' (($r.code -eq 200) -and ($j.data.result -eq 'accepted')) ('http=' + $r.code + ' result=' + $j.data.result + ' orderStatus=' + $j.data.orderStatus)
$lst = (Json (Req Get '/api/orders' $null $b1)).data
$ord1 = $lst | Where-Object { $_.orderNo -eq $o1no } | Select-Object -First 1
Check 'R4-1-paid-status' ($ord1.status -eq 'paid') ('status=' + $ord1.status)
$r = Notify $o1no 'TXN-OK1' 10000 (PaySign $o1no 'TXN-OK1' 10000)
$j = Json $r
$dup = ($j.data.result -eq 'duplicated') -or ($r.code -eq 409)
Check 'R3-3-replay-duplicated' ($dup) ('http=' + $r.code + ' result=' + $j.data.result)
$payRes = Req Get ('/api/ops/payments?order_no=' + $o1no) $null $ops
Write-Output ('PAYMENTS-DEBUG: http=' + $payRes.code + ' len=' + $payRes.body.Length + ' head=' + $payRes.body.Substring(0, [Math]::Min(120, $payRes.body.Length)))
$paysRaw = (Json $payRes).data
$pays = @($paysRaw)
Write-Output ('PAYS-DBG: rawType=' + $paysRaw.GetType().Name + ' paysCount=' + @($pays).Count + ' item0src=' + $pays[0].source + ' item0res=' + $pays[0].result)
$payRows = @($pays | Where-Object { $_.source -eq 'payment' }).Count
$dupCnt = @($pays | Where-Object { $_.result -eq 'duplicated' }).Count
Check 'R3-3-one-payment-record' ($payRows -eq 1 -and $dupCnt -ge 1) ('payment rows=' + $payRows + ' duplicated notices=' + $dupCnt)
# R3-4：同一渠道交易号归属另一订单
$o2 = Json (NewOrder $b1 1 1)
$o2no = $o2.data.order.orderNo
$r = Notify $o2no 'TXN-OK1' 10000 (PaySign $o2no 'TXN-OK1' 10000)
$j = Json $r
$r34 = ($r.code -ge 400) -or ($j.data.result -ne 'accepted')
Check 'R3-4-same-txn-other-order-rejected' ($r34) ('http=' + $r.code + ' result=' + $j.data.result)
# R4-2：订单超时关闭后收到成功通知 → 不恢复订单/不重占库存/进异常列表
ExpireAll $b1 | Out-Null
Start-Sleep -Seconds 1
$availBefore = Avail
$r = Notify $o2no 'TXN-LATE' 10000 (PaySign $o2no 'TXN-LATE' 10000)
$j = Json $r
$lst = (Json (Req Get '/api/orders' $null $b1)).data
$ord2 = $lst | Where-Object { $_.orderNo -eq $o2no } | Select-Object -First 1
$abns = (Json (Req Get '/api/cs/abnormals' $null $cs)).data
$abn = $abns | Where-Object { $_.channelTxnId -eq 'TXN-LATE' } | Select-Object -First 1
Check 'R4-2-order-stays-closed' ($ord2.status -eq 'closed') ('status=' + $ord2.status)
Check 'R4-2-no-reoccupy' ((Avail) -eq $availBefore) ('available before=' + $availBefore + ' after=' + (Avail))
Check 'R4-2-in-abnormal-list' ($null -ne $abn) ('abnormal found, status=' + $abn.status)
# R4-3：客服核实后发起退款（先切模拟渠道退款为 success）
$sr = SetRefundMode 'success' $cs
Write-Output ('SET-REFUND-MODE: http=' + $sr.code + ' body=' + $sr.body)
$r = Req Post '/api/cs/abnormals/handle' ('{"channelTxnId":"TXN-LATE","reason":"late-pay-refund"}') $cs
$j = Json $r
Check 'R4-3-cs-handle-refund' ($r.code -eq 200 -and $j.data.channelStatus -eq 'success') ('http=' + $r.code + ' body=' + $r.body + ' channelStatus=' + $j.data.channelStatus)
# R6-4：迟到支付退款不得恢复库存
Check 'R6-4-late-refund-no-restore' ((Avail) -eq $availBefore) ('available=' + (Avail) + ' (expect ' + $availBefore + ')')
# 重复处理同一异常 → 拒绝（防重复退款）
$r = Req Post '/api/cs/abnormals/handle' ('{"channelTxnId":"TXN-LATE","reason":"dup-handle"}') $cs
Check 'R6-4-double-handle-rejected' ($r.code -ge 400) ('http=' + $r.code + ' body=' + $r.body)
SetRefundMode 'external' $cs | Out-Null
Summary 'PHASE2'
