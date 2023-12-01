$apiKey = "apikey"
$programName = "getip"

$url = "http://localhost:8080/call-program"

$body = @{
    "program_name" = $programName
    "key" = $apiKey
} | ConvertTo-Json

$response = Invoke-RestMethod -Method Post -Uri $url -Body $body -ContentType "application/json"

Write-Host $response
