[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$BASE = 'http://127.0.0.1:8080'
$SECRET = 'channel-shared-secret'
$script:passN = 0
$script:failN = 0
$script:failList = @()

function Req {
  param($Method, $Path, $Body, $Token)
  $headers = @{}
  if ($Token) { $headers['Authorization'] = 'Bearer ' + $Token }
  try {
    if ($null -ne $Body) {
      $r = Invoke-WebRequest -UseBasicParsing -Method $Method -Uri ($BASE + $Path) -Body $Body -ContentType 'application/json' -Headers $headers
    } else {
      $r = Invoke-WebRequest -UseBasicParsing -Method $Method -Uri ($BASE + $Path) -Headers $headers
    }
    return @{ code = [int]$r.StatusCode; body = $r.Content }
  } catch {
    $resp = $_.Exception.Response
    if ($null -eq $resp) { return @{ code = -1; body = ($_.Exception.Message) } }
    $sc = [int]$resp.StatusCode
    $txt = $_.ErrorDetails.Message
    if (-not $txt) {
      $stream = $resp.GetResponseStream()
      $reader = New-Object System.IO.StreamReader($stream)
      $txt = $reader.ReadToEnd()
    }
    return @{ code = $sc; body = $txt }
  }
}

function Json($res) { return ($res.body | ConvertFrom-Json) }

function Login { param($u)
  $r = Req Post '/api/login' ('{"username":"' + $u + '","password":"shop123"}') $null
  return ((Json $r).data.token)
}

function PaySign { param($orderNo, $txn, $amount)
  $hmac = New-Object System.Security.Cryptography.HMACSHA256
  $hmac.Key = [System.Text.Encoding]::UTF8.GetBytes($SECRET)
  $bytes = $hmac.ComputeHash([System.Text.Encoding]::UTF8.GetBytes("$orderNo|$txn|$amount"))
  return (-join ($bytes | ForEach-Object { $_.ToString('x2') }))
}

function Notify { param($orderNo, $txn, $amount, $sign)
  return (Req Post '/api/pay/notify' ('{"orderNo":"' + $orderNo + '","channelTxnId":"' + $txn + '","amountCents":' + $amount + ',"sign":"' + $sign + '"}') $null)
}

function NewOrder { param($tok, $prodId, $qty)
  return (Req Post '/api/orders' ('{"productId":' + $prodId + ',"qty":' + $qty + '}') $tok)
}

function Avail { param($sku)
  if (-not $sku) { $sku = 'SKU-A' }
  $p = (Json (Req Get '/api/products' $null $null)).data
  $m = $p | Where-Object { $_.sku -eq $sku } | Select-Object -First 1
  return $m.available
}

function Check {
  param($name, $cond, $detail)
  if ($cond) { $script:passN++; Write-Output ("[PASS] " + $name + " | " + $detail) }
  else { $script:failN++; $script:failList += $name; Write-Output ("[FAIL] " + $name + " | " + $detail) }
}

function Reset { Req Post '/api/demo/reset' '{}' $null | Out-Null }

function SetRefundMode { param($mode, $tok)
  return (Req Post '/api/pay/channel-mode' ('{"mode":"' + $mode + '"}') $tok)
}

function ExpireAll { param($tok) return (Req Post '/api/test/expire-orders' '{}' $tok) }

function Summary { param($title)
  Write-Output ("=== " + $title + " :: PASS=" + $script:passN + " FAIL=" + $script:failN)
  if ($script:failList.Count -gt 0) { Write-Output ("FAILED-ITEMS: " + ($script:failList -join ', ')) }
}
