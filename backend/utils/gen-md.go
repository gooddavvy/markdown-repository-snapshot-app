package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
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

	tmpFile, err := ioutil.TempFile("", "*.zip")
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

func GenerateMarkdownFile(repoUrl string, ignoreList []string) (string, error) {
	var markdownFile string

	owner, repo, err := parseGitHubURL(repoUrl)
	if err != nil {
		return markdownFile, err
	}

	zipPath, err := downloadRepoAsZip(owner, repo)
	if err != nil {
		return markdownFile, err
	}
	defer os.Remove(zipPath)

	tempDir, err := os.MkdirTemp("", "repo")
	if err != nil {
		return markdownFile, err
	}
	defer os.RemoveAll(tempDir)

	err = unzip(zipPath, tempDir)
	if err != nil {
		return markdownFile, err
	}

	// Adjust the path to the unzipped content, which usually has a top-level directory
	files, err := os.ReadDir(tempDir)
	if err != nil || len(files) == 0 {
		return markdownFile, fmt.Errorf("failed to read unzipped contents")
	}

	// Assuming the first directory is the repo content directory
	repoContentDir := filepath.Join(tempDir, files[0].Name())

	var markdownContents []string
	err = filepath.Walk(repoContentDir, func(path string, info os.FileInfo, err error) error {
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
		return markdownFile, err
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

/* package utils



import (

 "bufio"

 "fmt"

 "os"

 "path/filepath"

 "strings"

)



// generateMarkdownSnapshot creates a Markdown file documenting the directory contents

func generateMarkdownFile(rootPath string, ignoreList []string) (string, error) {

 // Function to determine if the path should be ignored

 shouldIgnore := func(path string) bool {

  // Normalize the path to use forward slashes for consistent handling

  normalizedPath := filepath.ToSlash(path)



  for _, ignore := range ignoreList {

   // Normalize the ignore pattern to use forward slashes

   normalizedIgnore := filepath.ToSlash(ignore)



   // Match ignore patterns exactly from the root relative path

   trimmedPath := strings.TrimPrefix(normalizedPath, filepath.ToSlash(rootPath)+"/")

   if trimmedPath == normalizedIgnore || strings.HasPrefix(trimmedPath, normalizedIgnore+"/") {

    return true

   }

  }

  return false

 }



 // Walk the directory tree

 // func WalkDirAndWrite(dir string, info os.FileInfo, err error) error

 err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {

  if err != nil {

   return err

  }

  relativePath, err := filepath.Rel(rootPath, path)

  if err != nil {

   return err

  }

  if shouldIgnore(relativePath) {

   if info.IsDir() {

    return filepath.SkipDir

   }

   return nil

  }

  if !info.IsDir() {

   fileContent, err := os.ReadFile(path)

   if err != nil {

    return err

   }

   // Write the path and file content to the Markdown file

   fmt.Fprintf(writer, "### %s\n```\n%s\n```\n\n", relativePath, fileContent)

  }

  return nil

 })



 return err

}

*/
