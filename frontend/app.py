from flask import Flask, request, jsonify, render_template
import requests
import webbrowser
import json

app = Flask(__name__)

# Function to send request and get response from server
def genMd(repository_url, ignore_list):
    params = {
        "RepositoryURL": repository_url,
        "IgnoreList": ignore_list
    }

    # Convert the dictionary to a JSON string
    json_params = json.dumps(params)

    # URL for the POST request
    url = "http://localhost:8081/gen-md"

    # Send the POST request with the JSON string as the body
    res = requests.post(url, data=json_params, headers={"Content-Type": "application/json"})

    # Return the response
    return res


@app.route("/")
def index():
    return render_template("index.html")

@app.route("/generate", methods=["POST"])
def generate():
    data = request.json
    repository_url = data.get("repository_url")
    ignore_list = data.get("ignore_list", [])
    response = genMd(repository_url, ignore_list)
    try:
        res = response.json()
    except json.JSONDecodeError:
        res = response.text
    return jsonify({"status_code": response.status_code, "response": res})


if __name__ == "__main__":
    webbrowser.open("http://localhost:8501/")
    app.run(debug=True, host="0.0.0.0", port=8501)