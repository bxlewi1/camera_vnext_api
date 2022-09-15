package routes

import (
	"camera_vnext_api/config"
	"github.com/gofiber/fiber/v2"
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
		BodyLimit: 100 * 1024 * 1024 * 1024,
	})

	App.Server().StreamRequestBody = true

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

	api.Post("POST", func(ctx *fiber.Ctx) error {
		fileName := ctx.Query("FileName")
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

	api.Post("", func(ctx *fiber.Ctx) error {
		fileName := ctx.FormValue("FileName")
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

		upload, err := ctx.FormFile("Upload")
		if err != nil {
			logger.Println(err.Error())
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		err = ctx.SaveFile(upload, path.Join(dirPath, fileName))
		if err != nil {
			logger.Println(err.Error())
			return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.SendStatus(http.StatusOK)
	})
}
