param(
  [string]$ServiceId = "",
  [string]$ServiceName = "portfolio-backend",
  [string]$EnvFile = ".env.render",
  [switch]$Deploy
)

$ErrorActionPreference = "Stop"

function Get-RenderHeaders {
  $apiKey = $env:RENDER_API_KEY
  if (-not $apiKey) { $apiKey = $env:RENDER_API_TOKEN }
  if (-not $apiKey) {
    throw "RENDER_API_KEY (or RENDER_API_TOKEN) must be set in your environment."
  }
  return @{
    Authorization = "Bearer $apiKey"
    Accept        = "application/json"
    "Content-Type"= "application/json"
  }
}

function Get-RenderServiceIdByName([string]$name, [hashtable]$headers) {
  $services = @()
  $cursor = $null
  while ($true) {
    $uri = "https://api.render.com/v1/services?limit=100"
    if ($cursor) { $uri = "$uri&cursor=$cursor" }
    $page = Invoke-RestMethod -Method Get -Headers $headers -Uri $uri
    if (-not $page -or $page.Count -eq 0) { break }

    foreach ($item in $page) {
      if ($item.service) { $services += $item.service }
    }

    $next = $page[-1].cursor
    if (-not $next -or $next -eq $cursor) { break }
    $cursor = $next
  }

  $matches = $services | Where-Object { $_.name -eq $name }
  if (-not $matches -or $matches.Count -eq 0) {
    $known = ($services | Select-Object -ExpandProperty name | Sort-Object -Unique) -join ", "
    throw "No Render service named '$name' found. Known services: $known"
  }
  if ($matches.Count -gt 1) {
    $ids = ($matches | Select-Object -ExpandProperty id) -join ", "
    throw "Multiple services named '$name' found (ids: $ids). Re-run with -ServiceId."
  }
  return $matches[0].id
}

function Read-DotEnv([string]$path) {
  if (-not (Test-Path $path)) { throw "Env file not found: $path" }
  $vars = @{}
  foreach ($line in Get-Content $path) {
    $trim = $line.Trim()
    if (-not $trim) { continue }
    if ($trim.StartsWith("#")) { continue }
    $idx = $trim.IndexOf("=")
    if ($idx -lt 1) { continue }
    $key = $trim.Substring(0, $idx).Trim()
    $val = $trim.Substring($idx + 1).Trim()
    if (($val.StartsWith("\"") -and $val.EndsWith("\"")) -or ($val.StartsWith("'") -and $val.EndsWith("'"))) {
      $val = $val.Substring(1, $val.Length - 2)
    }
    if ($key) { $vars[$key] = $val }
  }
  return $vars
}

$headers = Get-RenderHeaders
if (-not $ServiceId) {
  $ServiceId = Get-RenderServiceIdByName -name $ServiceName -headers $headers
}

$envVars = Read-DotEnv -path $EnvFile
if ($envVars.Keys.Count -eq 0) { throw "No variables found in $EnvFile" }

Write-Host "Applying $($envVars.Keys.Count) env vars to service $ServiceId..."
foreach ($key in ($envVars.Keys | Sort-Object)) {
  $escapedKey = [uri]::EscapeDataString($key)
  $uri = "https://api.render.com/v1/services/$ServiceId/env-vars/$escapedKey"
  $body = @{ value = $envVars[$key] } | ConvertTo-Json
  Invoke-RestMethod -Method Put -Headers $headers -Uri $uri -Body $body | Out-Null
  Write-Host "Set $key"
}

if ($Deploy) {
  $deployUri = "https://api.render.com/v1/services/$ServiceId/deploys"
  Invoke-RestMethod -Method Post -Headers $headers -Uri $deployUri -Body "{}" | Out-Null
  Write-Host "Triggered deploy."
}

