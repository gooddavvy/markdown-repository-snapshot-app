###### This file is just for testing.

import requests
import json

# Define the dictionary with complex structures
params = {
    "RepositoryURL": "https://github.com/spf13/viper",
    "IgnoreList": [
        ".github",
        "README.md"
    ]
}

# Convert the dictionary to a JSON string
json_params = json.dumps(params)
print("JSON Parameters:", json_params)

# URL for the POST request
url = "http://localhost:8081/gen-md"

# Send the POST request with the JSON string as the body
res = requests.post(url, data=json_params, headers={"Content-Type": "application/json"})

# Print detailed request information for debugging
print("Response Status Code:", res.status_code)
print("Response Content:", res.text)
