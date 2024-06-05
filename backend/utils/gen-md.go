package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Parse the GitHub repository URL and return the owner and repository name
func parseGitHubURL(repoURL string) (owner, repo string, err error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) >= 2 {
		return parts[0], parts[1], nil
	}
	return "", "", fmt.Errorf("invalid GitHub repository URL")
}

func downloadRepoAsZip(owner, repo string) (string, error) {
	urls := []string{
		fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/main.zip", owner, repo),
		fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/master.zip", owner, repo),
	}

	var resp *http.Response
	var err error
	for _, url := range urls {
		resp, err = http.Get(url)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		resp.Body.Close()
		resp = nil
	}

	if resp == nil {
		return "", fmt.Errorf("failed to download repo: 404 Not Found")
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateMarkdownFile creates a Markdown file documenting the directory contents
func GenerateMarkdownFile(repoUrl string, ignoreList []string, folder string) (string, error) {
	var markdownFile string
	var dirname string

	dirname = folder

	owner, repo, err := parseGitHubURL(repoUrl)
	if err != nil {
		return "", err
	}

	zipPath, err := downloadRepoAsZip(owner, repo)
	if err != nil {
		return "", err
	}
	defer os.Remove(zipPath)

	tempDir, err := os.MkdirTemp("", "repo")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	err = unzip(zipPath, tempDir)
	if err != nil {
		return "", err
	}

	// Adjust the path to the unzipped content, which usually has a top-level directory
	files, err := os.ReadDir(tempDir)
	if err != nil || len(files) == 0 {
		return "", fmt.Errorf("failed to read unzipped contents")
	}

	// Assuming the first directory is the repo content directory
	repoContentDir := filepath.Join(tempDir, files[0].Name())

	var markdownContents []string
	var walkDir string

	if dirname == "" {
		walkDir = repoContentDir
	} else {
		walkDir = filepath.Join(repoContentDir, dirname)
	}

	err = filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath := strings.TrimPrefix(path, repoContentDir+string(os.PathSeparator))
		if info.IsDir() || isIgnored(relativePath, ignoreList) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fileMarkdown := fmt.Sprintf("### %s\n\n```%s\n%s\n```", relativePath, getFileExtension(relativePath), string(content))
		markdownContents = append(markdownContents, fileMarkdown)

		return nil
	})
	if err != nil {
		return "", err
	}

	markdownFile = strings.Join(markdownContents, "\n\n")
	return markdownFile, nil
}

func isIgnored(path string, ignoreList []string) bool {
	for _, ignore := range ignoreList {
		if strings.HasPrefix(path, ignore) {
			return true
		}
	}
	return false
}

func getFileExtension(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

/*

 */
