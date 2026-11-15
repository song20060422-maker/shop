# 阶段4：S2 缺货与弃单（20 并发争抢 5 件）+ 释放后可复购
. D:\Users\Lenovo\Desktop\shop\shop\acceptance-r2\_common.ps1
Write-Output "########## PHASE4: S2 concurrency 20 ##########"

Reset
# 20 个并发下单请求（轮换 3 个买家，各购 SKU-A 1 件，争抢 5 件库存）
$jobs = @()
1..20 | ForEach-Object {
  $i = $_
  $buyer = 'buyer' + (($i % 3) + 1)
  $jobs += Start-Job -ScriptBlock {
    param($buyer, $idx)
    $login = Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/login' -Method Post -ContentType 'application/json' -Body ('{"username":"' + $buyer + '","password":"shop123"}')
    $tok = $login.data.token
    try {
      $r = Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/orders' -Method Post -ContentType 'application/json' -Headers @{Authorization = 'Bearer ' + $tok} -Body '{"productId":1,"qty":1}'
      "OK " + $r.data.order.orderNo
    } catch {
      $sc = [int]$_.Exception.Response.StatusCode
      "ERR " + $sc
    }
  } -ArgumentList $buyer, $i
}
$jobs | Wait-Job | Out-Null
$results = $jobs | Receive-Job
$jobs | Remove-Job -Force
$oks = ($results | Where-Object { $_ -like 'OK*' })
$errs = ($results | Where-Object { $_ -like 'ERR*' })
Write-Output ("RAW RESULTS: success=" + $oks.Count + " failed=" + $errs.Count)
Write-Output ($results -join ' | ')
# 无超卖：成功单数 = 5（每单 1 件），失败 15
Check 'S2-20-concurrent-exactly-5-win' ($oks.Count -eq 5 -and $errs.Count -eq 15) ('success=' + $oks.Count + ' failed=' + $errs.Count)
Check 'S2-no-oversell-stock-zero' ((Avail) -eq 0) ('available=' + (Avail))
# 部分订单未付款关闭 → 释放（循环触发直至 5 件全部释放，模拟陆续关闭）
$b1t = Login 'buyer1'
for ($k = 0; $k -lt 8 -and (Avail) -lt 5; $k++) { ExpireAll $b1t | Out-Null; Start-Sleep -Milliseconds 600 }
Check 'S2-close-release-5' ((Avail) -eq 5) ('available after close=' + (Avail))
# 其他买家购买释放的商品
$b2 = Login 'buyer2'
$r = NewOrder $b2 1 2
$j = Json $r
Check 'S2-repurchase-released-goods' ($r.code -eq 200 -and (Avail) -eq 3) ('http=' + $r.code + ' orderNo=' + $j.data.order.orderNo + ' available=' + (Avail))
Req Post ('/api/orders/' + $j.data.order.orderNo + '/cancel') '{}' $b2 | Out-Null
Check 'S2-cancel-back-to-5' ((Avail) -eq 5) ('available=' + (Avail))
Summary 'PHASE4'
