# 阶段1：部署核对(D1-D3) + R1 商品与下单约束 + R2 保留与释放
. D:\Users\Lenovo\Desktop\shop\shop\acceptance-r2\_common.ps1
Write-Output "########## PHASE1: D1-D3 / R1 / R2 ##########"

# 重置演示数据，保证从 5 件满库存开始
Reset
$prod = (Json (Req Get '/api/products' $null $null)).data
$skuA = $prod | Where-Object { $_.sku -eq 'SKU-A' } | Select-Object -First 1
$skuB = $prod | Where-Object { $_.sku -eq 'SKU-B' } | Select-Object -First 1
# D2：SKU-A 10000 分 / 5 件
Check 'D2-SKUA-10000fen-5pcs' ($skuA.priceCents -eq 10000 -and $skuA.available -eq 5 -and $skuA.initialStock -eq 5) ("priceCents=" + $skuA.priceCents + " available=" + $skuA.available + " initialStock=" + $skuA.initialStock)
# D2：SKU-B 5000 分 / 3 件
Check 'D2-SKUB-5000fen-3pcs' ($skuB.priceCents -eq 5000 -and $skuB.available -eq 3) ("priceCents=" + $skuB.priceCents + " available=" + $skuB.available)
# D2：3 买家 + 3 员工账号
$users = (Json (Req Get '/api/users' $null $null)).data
$names = ($users | ForEach-Object { $_.username }) -join ','
$missing = @('ops1','wh1','cs1','buyer1','buyer2','buyer3') | Where-Object { $names -notmatch $_ }
Check 'D2-accounts-3buyers-3staff' ($missing.Count -eq 0) ('users=' + $names)
# D3：买家页可访问
$html = Req Get '/shop' $null $null
Check 'D3-buyer-page-200' ($html.code -eq 200) ('http=' + $html.code)
# D1：五角色均可登录（8080 买家/业务服务 + GVA 8085/8099 已另行人工验证）
$b1 = Login 'buyer1'; $ops = Login 'ops1'; $wh = Login 'wh1'; $cs = Login 'cs1'; $b2 = Login 'buyer2'
Check 'D1-login-all-roles' ($b1 -and $ops -and $wh -and $cs -and $b2) 'buyer1/buyer2/ops1/wh1/cs1 all got JWT'

# R1-1：数量超过可售（3 > 5? 不，此处 3<=5 可售，先用满库存测 6 件）
$r = NewOrder $b1 1 6
Check 'R1-1-qty-exceed-rejected' ($r.code -ge 400 -and $r.code -lt 500) ('http=' + $r.code + ' body=' + $r.body)
Check 'R1-1-no-stock-occupied' ((Avail) -eq 5) ('available=' + (Avail))
# R1-1：下架商品拒购
Req Patch '/api/ops/products' '{"productId":2,"onSale":false}' $ops | Out-Null
$r = NewOrder $b1 2 1
Check 'R1-1-offsale-rejected' ($r.code -ge 400) ('http=' + $r.code + ' body=' + $r.body)
Req Patch '/api/ops/products' '{"productId":2,"onSale":true}' $ops | Out-Null
Check 'R1-1-offsale-no-leak' ((Avail 'SKU-B') -eq 3) ('SKU-B available=' + (Avail 'SKU-B'))

# R1-2：订单保存下单时刻价格，改价不改已有订单
$o1 = Json (NewOrder $b1 1 1)
$o1no = $o1.data.order.orderNo
Check 'R1-2-order-created' ($o1.code -eq 0 -and $o1no) ('orderNo=' + $o1no + ' amount=' + $o1.data.order.amountCents)
Req Patch '/api/ops/products' '{"productId":1,"priceCents":12000}' $ops | Out-Null
$lst = (Json (Req Get '/api/orders' $null $b1)).data
$ord1 = $lst | Where-Object { $_.orderNo -eq $o1no } | Select-Object -First 1
Check 'R1-2-price-snapshot-10000' ($ord1.amountCents -eq 10000) ('after price->12000, order amountCents=' + $ord1.amountCents)
Req Patch '/api/ops/products' '{"productId":1,"priceCents":10000}' $ops | Out-Null
# R1-3：下单 1 件后可售立即减 1
Check 'R1-3-occupy-reflected' ((Avail) -eq 4) ('after qty1 order, available=' + (Avail))
# R1-4：数量越界（0 / -1 / 3 件在剩 4 件时合法，用上限 2 约束测 3 件）
$r0 = NewOrder $b1 1 0
$rm = NewOrder $b1 1 -1
$r3 = NewOrder $b1 1 3
Check 'R1-4-qty0-rejected' ($r0.code -ge 400) ('http=' + $r0.code)
Check 'R1-4-qty-negative-rejected' ($rm.code -ge 400) ('http=' + $rm.code)
Check 'R1-4-qty3-over-limit-rejected' ($r3.code -ge 400) ('http=' + $r3.code + ' body=' + $r3.body)

# R2-1：待付款保留 10 分钟（服务端时间；配置走查 order-timeout-minutes: 10）
$o2 = Json (NewOrder $b1 1 2)
$o2no = $o2.data.order.orderNo
$ea = [datetime]::Parse($o2.data.order.expireAt)
$diffMin = ($ea - (Get-Date)).TotalMinutes
Check 'R2-1-expire-in-10min' ($diffMin -ge 8.5 -and $diffMin -le 11.5) ('orderNo=' + $o2no + ' expireIn=' + [math]::Round($diffMin,1) + 'min')
# R1-3 定量：下单 2 件后 5-1-2=2
Check 'R1-3-occupy-5-1-2-eq-2' ((Avail) -eq 2) ('available=' + (Avail))
# R2-2：买家主动取消释放占用
$r = Req Post ('/api/orders/' + $o2no + '/cancel') '{}' $b1
Check 'R2-2-cancel-ok' ($r.code -eq 200) ('http=' + $r.code)
Check 'R2-2-stock-restored-4' ((Avail) -eq 4) ('available=' + (Avail))
# R2-3 前置：重复取消不得重复释放
$r2 = Req Post ('/api/orders/' + $o2no + '/cancel') '{}' $b1
Check 'R2-3-double-cancel-rejected' ($r2.code -ge 400) ('http=' + $r2.code + ' body=' + $r2.body)
Check 'R2-3-no-double-release' ((Avail) -eq 4) ('available=' + (Avail))
# 先取消 o1，让 R2-3 的库存核算完全确定（无其他未决订单）
Req Post ('/api/orders/' + $o1no + '/cancel') '{}' $b1 | Out-Null
# R2-3：超时关闭释放且重复关单幂等（此时仅 o3 一单占用 2 件，基线 3）
$o3 = Json (NewOrder $b1 1 2)
$o3no = $o3.data.order.orderNo
Check 'R2-3-baseline-3' ((Avail) -eq 3) ('available before expire=' + (Avail))
ExpireAll $b1 | Out-Null
Start-Sleep -Seconds 1
$lst = (Json (Req Get '/api/orders' $null $b1)).data
$ord3 = $lst | Where-Object { $_.orderNo -eq $o3no } | Select-Object -First 1
Check 'R2-3-closed-by-expire' ($ord3.status -eq 'closed') ('status=' + $ord3.status)
Check 'R2-3-expire-restored-5' ((Avail) -eq 5) ('available=' + (Avail))
ExpireAll $b1 | Out-Null
Start-Sleep -Seconds 1
Check 'R2-3-second-close-idempotent' ((Avail) -eq 5) ('available after 2nd close=' + (Avail))
Summary 'PHASE1'
