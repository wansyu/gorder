$apiKey = "apikey"
$programName = "getip"

$url = "http://localhost:8080/call-program"

$jsonData  = @{
    "program_name" = $programName
    "key" = $apiKey
} 
$jsonData  = $jsonData  | ConvertTo-Json 
$headers = @{ "Content-type" = "application/json" }

$response = Invoke-RestMethod  -Uri $url -Method Post -Headers $headers -Body $jsonData  

$response
