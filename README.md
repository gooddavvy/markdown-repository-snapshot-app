# markdown-repository-snapshot-app

Description: An app that you input a GitHub repository's URL into, alongside ignore files&folders, and it outputs a downloadable `output.md` file

Release Date: May 17, 2024

Phase: Completed

**Need to work on:**

- [x] Making it work via GitHub API and GitHub API token
- [x] Doing it via getting and unzipping zip file for repository

**Feel free to look at the code--the app is all open-source!**

# How to use

**Setting it up:**

First, ensure you have installed [GoLang](https://go.dev/dl/) and [Python](https://python.org/downloads/).

After this, you can open your terminal, navigate to the desired directory, and run:

```
git clone https://github.com/gooddavvy/markdown-repository-snapshot-app
cd markdown-repository-snapshot-app
code .
```

**Running it:**

Then, open a VS-Code terminal in the workspace, and run:

```
cd backend
go mod tidy
go run main.go
```

After that, open another one, and run:

```
cd frontend
pip install -r requirements.txt
python app.py
```

If everything works as expected, it should automatically open `http://localhost:8501` in your default browser.

# Note

If you encounter any issues, please feel free to talk in the [Issues Section](https://github.com/gooddavvy/markdown-repository-snapshot-app/issues)!
