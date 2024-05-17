# markdown-directory-snapshot-app

Description: An app that you input a GitHub repository's link into, alongside ignore files&folders, and it outputs a downloadable `output.md` file
Release Date: May 16, 2024 - Evening
Phase: BETA

**Feel free to look at the code--the app is all open-source!**

# How to use

**Setting it up:**

First, ensure you have installed [GoLang](https://go.dev/dl/) and [Python](https://python.org/downloads/).

Then, generate a [new GitHub API key](https://github.com/settings/tokens/new) if you haven't already.

After this, you can open your terminal, navigate to the desired directory, and run:

```
git clone https://github.com/gooddavvy/markdown-repository-snapshot-app
cd markdown-repository-snapshot-app
code .
```

Then, in the `backend` folder, add a `.env` file:

```env
GITHUB_API_TOKEN = your_github_api_token

```

Be sure to replace `your_github_api_token` with your actual GitHub API token.

**Running it:**

Then, open a VS-Code terminal in the workspace, and run:

```
cd backend
go run main.go
```

After that, open another one, and run:

```
pip install -r requirements.txt
python -m streamlit run app.py
```

If it works as expected, it should automatically open `http://localhost:8501` in your default browser.

# Important considerations

This application is still in its BETA phase, so you might not want to clone this repository just yet. Please feel free to talk in the [Issues Section](https://github.com/gooddavvy/markdown-repository-snapshot-app/issues)!
