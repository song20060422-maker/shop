# 阶段5：S4 迟到支付闭环 + S5 操作竞争 + B5 越权 + R8 运营核查与对账
. D:\Users\Lenovo\Desktop\shop\shop\acceptance-r2\_common.ps1
Write-Output "########## PHASE5: S4 / S5 / B5 / R8 ##########"

Reset
$b1 = Login 'buyer1'; $b2 = Login 'buyer2'; $b3 = Login 'buyer3'
$ops = Login 'ops1'; $wh = Login 'wh1'; $cs = Login 'cs1'
SetRefundMode 'success' $cs | Out-Null

# ---- S4 迟到支付：订单关闭 → 成功通知 → 客服退款 → 对账 ----
$o1 = Json (NewOrder $b1 1 1)
$o1no = $o1.data.order.orderNo
ExpireAll $b1 | Out-Null
Start-Sleep -Seconds 1
$availBefore = Avail
$r = Notify $o1no 'TXN-S4' 10000 (PaySign $o1no 'TXN-S4' 10000)
$j = Json $r
Check 'S4-late-notify-abnormal' ($j.data.result -eq 'abnormal') ('result=' + $j.data.result)
$r = Req Post '/api/cs/abnormals/handle' ('{"channelTxnId":"TXN-S4","reason":"late-pay-refund"}') $cs
$j = Json $r
Write-Output ('S4-HANDLE-DEBUG: http=' + $r.code + ' body=' + $r.body)
Check 'S4-cs-refund-done' (($r.code -eq 200 -or $r.code -eq 409) -and ($j.data.channelStatus -eq 'success')) ('http=' + $r.code + ' channelStatus=' + $j.data.channelStatus)
# 迟到支付退款不恢复数量（已释放过）
Check 'S4-no-double-restore' ((Avail) -eq $availBefore) ('available=' + (Avail) + ' expect=' + $availBefore)
# 运营确认无重复退款：该订单退款单仅 1 张且成功
$rfs = @((Json (Req Get '/api/cs/refunds' $null $cs)).data)
$rf4 = @($rfs | Where-Object { $_.orderNo -eq $o1no })
$st4 = $rf4[0].channelStatus; if (-not $st4) { $st4 = $rf4[0].'channel_status' }
Write-Output ('S4-REFUND-DEBUG: total refunds=' + $rfs.Count + ' for-order=' + $rf4.Count + ' status=' + $st4)
Check 'S4-no-double-refund' ($rf4.Count -eq 1 -and $st4 -eq 'success') ('refund records=' + $rf4.Count + ' status=' + $st4)

# ---- S5-1 退款申请与发货并发：只允许一种合法结果 ----
$o2 = Json (NewOrder $b2 1 1)
$o2no = $o2.data.order.orderNo
Notify $o2no 'TXN-S5' 10000 (PaySign $o2no 'TXN-S5' 10000) | Out-Null
$jA = Start-Job -ScriptBlock {
  param($o)
  $l = Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/login' -Method Post -ContentType 'application/json' -Body '{"username":"buyer2","password":"shop123"}'
  try {
    $r = Invoke-RestMethod -Uri ('http://127.0.0.1:8080/api/orders/' + $o + '/refund') -Method Post -ContentType 'application/json' -Headers @{Authorization = 'Bearer ' + $l.data.token} -Body '{"reason":"concurrent-refund"}'
    "REFUND-OK"
  } catch { "REFUND-ERR " + [int]$_.Exception.Response.StatusCode }
} -ArgumentList $o2no
$jB = Start-Job -ScriptBlock {
  param($o)
  $l = Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/login' -Method Post -ContentType 'application/json' -Body '{"username":"wh1","password":"shop123"}'
  try {
    $r = Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/warehouse/ship' -Method Post -ContentType 'application/json' -Headers @{Authorization = 'Bearer ' + $l.data.token} -Body ('{"orderNo":"' + $o + '","carrier":"SF","trackingNo":"SF-CONC"}')
    "SHIP-OK"
  } catch { "SHIP-ERR " + [int]$_.Exception.Response.StatusCode }
} -ArgumentList $o2no
Wait-Job $jA, $jB | Out-Null
$rA = Receive-Job $jA; $rB = Receive-Job $jB
Remove-Job $jA, $jB -Force
Write-Output ("CONCURRENT RESULT: " + $rA + " / " + $rB)
$lst = (Json (Req Get '/api/orders' $null $b2)).data
$ord2 = $lst | Where-Object { $_.orderNo -eq $o2no } | Select-Object -First 1
# 合法结果 = 要么发货成功且退款被拒，要么退款挂起且发货被拒；不允许“已发货且已退款”
$legal = (($rB -eq 'SHIP-OK') -and ($ord2.status -eq 'shipped')) -or (($rA -eq 'REFUND-OK') -and ($ord2.status -ne 'shipped'))
Check 'S5-1-concurrent-one-legal-outcome' ($legal) ('refund=' + $rA + ' ship=' + $rB + ' finalStatus=' + $ord2.status)
Check 'S5-1-not-both-succeed' (-not (($rA -eq 'REFUND-OK') -and ($rB -eq 'SHIP-OK'))) ('refund=' + $rA + ' ship=' + $rB)

# ---- S5-2 重复支付通知（已有 PHASE2 R3-3）+ 越权 ----
# B5：买家2 访问/操作 买家1 的订单
$o3 = Json (NewOrder $b1 1 1)
$o3no = $o3.data.order.orderNo
$r = Req Get ('/api/orders/' + $o3no) $null $b2
Check 'B5-cross-view-rejected' ($r.code -eq 403 -or $r.code -eq 404) ('http=' + $r.code)
$r = Req Post ('/api/orders/' + $o3no + '/cancel') '{}' $b2
Check 'B5-cross-cancel-rejected' ($r.code -ge 400) ('http=' + $r.code)
$r = Req Post ('/api/orders/' + $o3no + '/refund') '{"reason":"cross-access"}' $b2
Check 'B5-cross-refund-rejected' ($r.code -ge 400) ('http=' + $r.code)

# ---- R8 运营核查：逐单四类记录 + 汇总对账 ----
# 先给 o3 付款并发货，保证收款/发货记录同时在册
Notify $o3no 'TXN-R8' 10000 (PaySign $o3no 'TXN-R8' 10000) | Out-Null
$r = Req Post '/api/warehouse/ship' ('{"orderNo":"' + $o3no + '","carrier":"YT","trackingNo":"YT-888"}') $wh
Write-Output ('R8-SHIP: http=' + $r.code)
$sum = (Json (Req Get '/api/ops/summary' $null $ops)).data
Write-Output ("OPS SUMMARY: " + ($sum | ConvertTo-Json -Depth 4 -Compress))
$allOrders = @((Json (Req Get '/api/ops/orders' $null $ops)).data)
$pays = @((Json (Req Get ('/api/ops/payments?order_no=' + $o3no) $null $ops)).data)
$pays2 = @((Json (Req Get ('/api/ops/payments?order_no=' + $o1no) $null $ops)).data)
$rfs = @((Json (Req Get '/api/cs/refunds' $null $cs)).data)
$ships = @((Json (Req Get '/api/warehouse/shipments' $null $wh)).data)
$abns = @((Json (Req Get '/api/ops/abnormals' $null $ops)).data)
$allPays = @($pays + $pays2)
$payRowsAll = @($allPays | Where-Object { $_.source -eq 'payment' })
$okRefunds = @($rfs | Where-Object { ($_.channelStatus -eq 'success') -or ($_.'channel_status' -eq 'success') })
Check 'R8-1-four-record-types-visible' (($payRowsAll.Count -ge 1) -and ($okRefunds.Count -ge 1) -and $abns.Count -ge 1 -and $ships.Count -ge 1) ('payments=' + $allPays.Count + ' payRows=' + $payRowsAll.Count + ' refunds=' + $rfs.Count + ' shipments=' + $ships.Count + ' abnormals=' + $abns.Count + ' orders=' + $allOrders.Count)
# R8-2 对账：已收款 = 记款行 + 迟到支付异常款；已退款 = 成功退款金额和；均与汇总一致
$paidSum = ($payRowsAll | Measure-Object -Property amountCents -Sum).Sum
$latePaid = (@($abns) | Where-Object { $_.status -ne 'pending' } | Measure-Object -Property amountCents -Sum).Sum
$refSum = ($okRefunds | Measure-Object -Property amountCents -Sum).Sum
Write-Output ("RECON: paidRows=" + $paidSum + " latePaid=" + $latePaid + " refundSum=" + $refSum + " summaryPaid=" + $sum.totalPaidCents + " summaryRefund=" + $sum.totalRefundedCents)
$ledger = (Json (Req Get '/api/ops/ledger?productId=1' $null $ops)).data
Write-Output ("LEDGER ENTRIES: " + @($ledger).Count)
Check 'R8-2-summary-struct-ok' ((($paidSum + $latePaid) -eq $sum.totalPaidCents) -and ($refSum -eq $sum.totalRefundedCents)) ('paidSum+late=' + ($paidSum + $latePaid) + ' refundSum=' + $refSum)
# R8-3：无直接修改订单状态的旁路接口（404 即不存在）
$r = Req Post ('/api/ops/orders/' + $o3no + '/status') '{"status":"refunded"}' $ops
Check 'R8-3-no-status-bypass-endpoint' ($r.code -eq 404 -or $r.code -eq 405 -or $r.code -ge 400) ('http=' + $r.code)
SetRefundMode 'external' $cs | Out-Null
Summary 'PHASE5'
