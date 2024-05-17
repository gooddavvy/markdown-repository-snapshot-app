import streamlit as st
import requests
import json
import random

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


def main():
    st.title("Markdown Repository Snapshot Application")
    use_json = True

    # Input for repository URL and declaration ignore list
    repository_url = st.text_input(label="Repository URL",
                                   placeholder="https://github.com/spf13/viper")

    if "ignore_list" not in st.session_state:
        st.session_state["ignore_list"] = []

    if "ignore_item" not in st.session_state:
        st.session_state["ignore_item"] = ""

    def append_ignore_item(ignore_item):
            st.session_state["ignore_list"].append(ignore_item)

    st.write("Current Ignore List:", st.session_state["ignore_list"])
    # for ignore_item in st.session_state["ignore_list"]:
    #     st.write(ignore_item)

    ignore_item = st.text_input(label="New Ignore Item",
                                key="ignore_item",
                                placeholder=random.choice([".github", "README.md", "remote"]))

    global x
    x = False
    def add_ignore_item():
        global x
        if x is True:
            if ignore_item != "":
                append_ignore_item(ignore_item)
                append_ignore_item(ignore_item)

                x = False
                st.rerun()
        else:
            x = True
            add_ignore_item()



    st.button(label="Add Ignore Item", on_click=add_ignore_item)

    if st.button("Generate Markdown"):
        ignore_list = st.session_state["ignore_list"]
        response = genMd(repository_url, ignore_list)
        res = response.text

        # Convert response text to dictionary
        try:
            res = response.json()
        except json.JSONDecodeError:
            use_json = False
            res = response.text

        # Display the response
        st.write("Response Status Code:", response.status_code)
        st.write("**See response content below...**\n")
        if use_json:
            st.json(res)
        else:
            # Create a downloadable link for the raw text response
            st.download_button(
                label="Download Markdown",
                data=res,
                file_name="output.md",
                mime="text/markdown"
            )

if __name__ == "__main__":
    main()