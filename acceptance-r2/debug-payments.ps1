. D:\Users\Lenovo\Desktop\shop\shop\acceptance-r2\_common.ps1
$ops = Login 'ops1'
$orders = (Json (Req Get '/api/ops/orders?status=paid' $null $ops)).data
$o1 = $orders | Select-Object -First 1
$res = Req Get ('/api/ops/payments?order_no=' + $o1.orderNo) $null $ops
$paysRaw = (Json $res).data
'rawType=' + $paysRaw.GetType().Name + ' rawCount=' + @($paysRaw).Count
if ($paysRaw.list) { $pays = $paysRaw.list; 'USED .list' } else { $pays = $paysRaw; 'USED raw' }
'paysType=' + $pays.GetType().Name + ' paysCount=' + @($pays).Count
$payRows = @($pays | Where-Object { $_.source -eq 'payment' }).Count
$dupCnt = @($pays | Where-Object { $_.result -eq 'duplicated' }).Count
'payRows=' + $payRows + ' dupCnt=' + $dupCnt
$pays | ForEach-Object { 'ITEM source=' + $_.source + ' result=' + $_.result + ' status=' + $_.status }
