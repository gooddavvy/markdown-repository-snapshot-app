package utils

import (
	"context"

	"encoding/base64"

	"fmt"

	"net/url"

	"strings"

	"github.com/google/go-github/v61/github"

	"github.com/spf13/viper"

	"golang.org/x/oauth2"
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

func GenerateMarkdownFile(repoUrl string, ignoreList []string) (string, error) {
	var markdownFile string

	// Start working
	ctx := context.Background()
	token, ok := viper.Get("GITHUB_API_TOKEN").(string)

	if !ok || token == "" {

		return markdownFile, fmt.Errorf("GITHUB_API_TOKEN not set or empty")

	}

	ts := oauth2.StaticTokenSource(

		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	owner, repo, err := parseGitHubURL(repoUrl)

	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return markdownFile, err
	}

	// Recursively fetch all files and folders, excluding those in the ignore list

	var markdownContents []string

	err = fetchDirectoryContents(ctx, client, owner, repo, "", &markdownContents, ignoreList)

	if err != nil {

		return markdownFile, err

	}

	markdownFile = strings.Join(markdownContents, "\n\n")

	return markdownFile, nil

}

func fetchDirectoryContents(ctx context.Context, client *github.Client, owner, repo, path string, markdownContents *[]string, ignoreList []string) error {

	opt := &github.RepositoryContentGetOptions{}

	contents, dirContents, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opt)

	if err != nil {

		fmt.Println("Error fetching contents:", err)

		return err

	}

	// Check if we received a single file or a directory

	if contents != nil {

		// Single file

		if !isIgnored(contents.GetPath(), ignoreList) {

			if contents.GetType() == "file" {

				fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, contents.GetPath(), opt)

				if err != nil {

					return err

				}

				// Decode base64 content

				decodedContent, err := base64.StdEncoding.DecodeString(*fileContent.Content)

				if err != nil {

					return err

				}

				// Append file content in markdown format

				fileMarkdown := fmt.Sprintf("### %s\n```%s\n%s```", contents.GetPath(), getFileExtension(contents.GetPath()), string(decodedContent))

				*markdownContents = append(*markdownContents, fileMarkdown)

			}

		} else {

			// fmt.Println("Ignoring file:", contents.GetPath())

		}

	} else if dirContents != nil {

		// Directory

		for _, content := range dirContents {

			if !isIgnored(content.GetPath(), ignoreList) {

				if content.GetType() == "file" {

					fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, content.GetPath(), opt)

					if err != nil {

						return err

					}

					// Decode base64 content

					decodedContent, err := base64.StdEncoding.DecodeString(*fileContent.Content)

					if err != nil {

						return err

					}

					// Append file content in markdown format

					fileMarkdown := fmt.Sprintf("### %s\n\n```%s\n%s\n```", content.GetPath(), getFileExtension(content.GetPath()), string(decodedContent))

					*markdownContents = append(*markdownContents, fileMarkdown)

				} else if content.GetType() == "dir" {

					err := fetchDirectoryContents(ctx, client, owner, repo, content.GetPath(), markdownContents, ignoreList)

					if err != nil {

						return err

					}

				}

			} else {

				// fmt.Println("Ignoring directory or file:", content.GetPath())

			}

		}

	}

	return nil

}

func isIgnored(path string, ignoreList []string) bool {

	for _, ignore := range ignoreList {

		if strings.HasPrefix(path, ignore) {

			return true

		}

	}

	return false

}

func getFileExtension(filename string) string {

	parts := strings.Split(filename, ".")

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
