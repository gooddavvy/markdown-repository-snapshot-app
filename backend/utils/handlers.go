package utils

import (
	"encoding/json"
	_ "errors"

	"github.com/gofiber/fiber/v2"
)

func GenMdHandler(ctx *fiber.Ctx) error {
	reqBody := make(map[string]interface{})
	if err := json.Unmarshal(ctx.Body(), &reqBody); err != nil {
		ctx.Context().SetContentType("application/json")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	repoUrl, linkOk := reqBody["RepositoryURL"].(string)
	if !linkOk {
		ctx.Context().SetContentType("application/json")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid RepositoryURL",
		})
	}

	directoryName, dirnameOk := reqBody["DirectoryName"].(string)
	if !dirnameOk {
		directoryName = ""
	}

	ignoreListInterface, listOk := reqBody["IgnoreList"].([]interface{})
	if !listOk {
		ctx.Context().SetContentType("application/json")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid IgnoreList",
		})
	}

	ignoreList := make([]string, len(ignoreListInterface))
	for i, v := range ignoreListInterface {
		ignoreList[i], listOk = v.(string)
		if !listOk {
			ctx.Context().SetContentType("application/json")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "IgnoreList contains non-string value",
			})
		}
	}

	snapshot, err := GenerateMarkdownFile(repoUrl, ignoreList, directoryName)
	if err != nil {
		ctx.Context().SetContentType("application/json")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating Markdown file: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).SendString(snapshot)
}
