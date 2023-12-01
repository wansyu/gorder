import requests
import json

apiKey = "apikey"
programName = "getip"

url = "http://localhost:8080/call-program"

data = {
    "program_name": programName,
    "key": apiKey
}

headers = {
    "Content-Type": "application/json"
}

response = requests.post(url, data=json.dumps(data), headers=headers)

print(response.json())