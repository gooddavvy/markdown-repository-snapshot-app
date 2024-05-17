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

    # Input for repository URL
    repository_url = st.text_input(label="Repository URL",
                                   placeholder="https://github.com/spf13/viper")

    if "ignore_list" not in st.session_state:
        st.session_state["ignore_list"] = []

    if "generate_md" not in st.session_state:
        st.session_state["generate_md"] = False

    def add_ignore_item():
        ignore_item = st.session_state.get("ignore_item_input", "")
        if ignore_item:
            st.session_state.ignore_list.append(ignore_item)
            st.session_state["ignore_item_input"] = ""

    def set_generate_md():
        st.session_state["generate_md"] = True

    # Input for ignore item and add button
    ignore_item = st.text_input(label="New Ignore Item",
                                key="ignore_item_input",
                                placeholder=random.choice([".github", "README.md", "remote"]))

    st.button("Add Ignore Item", on_click=add_ignore_item)

    # Display the ignore list items
    for item in st.session_state.ignore_list:
        st.write(item)

    st.button("Generate Markdown", on_click=set_generate_md)

    if st.session_state["generate_md"]:
        # Call the genMd function
        ignore_list = st.session_state.ignore_list
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
            st.download_button(
                label="Download Markdown",
                data=res,
                file_name="output.md",
                mime="text/markdown"
            )

if __name__ == "__main__":
    main()