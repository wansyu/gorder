import requests
import json

apiKey = "somepassword"
programName = "getip"

url = "http://localhost:8080/call-program"

data = {
    "program_name": programName,
    "key": apiKey
}

headers = {
    "Content-Type": "application/json; charset=utf-8"
}

response = requests.post(url, data=json.dumps(data), headers=headers)

print(response.json())
