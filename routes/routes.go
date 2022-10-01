package routes

import (
	"camera_vnext_api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

var App *fiber.App

func init() {
	App = fiber.New(fiber.Config{
		StreamRequestBody: true,
		BodyLimit:         -1,
	})

	addApiRoutes()
}

func addApiRoutes() {
	api := App.Group("api")
	logger := log.Default()

	api.Get("", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"derp": "ERRYTHING IS AWESOME",
			"date": time.Now(),
		})
	})

	api.Get("guid", func(ctx *fiber.Ctx) error {
		uuid := utils.UUIDv4()
		return ctx.SendString(uuid)
	})

	api.Post("POST", func(ctx *fiber.Ctx) error {
		fileName := ctx.Query("fileName")
		if len(fileName) == 0 {
			logger.Println("FileName not given")
			return ctx.SendStatus(http.StatusBadRequest)
		}

		dirPath := path.Join(config.Instance.AssetPath)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				logger.Println(err.Error())
			}
		}

		if _, err := os.Stat(path.Join(dirPath, fileName)); os.IsExist(err) {
			fileName = "1." + fileName
		}

		reader := ctx.Context().RequestBodyStream()

		outFile, err := os.Create(filepath.Join(dirPath, fileName))
		if err != nil {
			logger.Println(err.Error())
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, reader)
		if err != nil {
			logger.Println(err.Error())
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.SendStatus(http.StatusOK)
	})
}
