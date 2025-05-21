package api

import (
	// "github.com/Himneesh-Kalra/go-vector-db/db"
	"github.com/Himneesh-Kalra/go-vector-db/db"
	"github.com/Himneesh-Kalra/go-vector-db/models"
	"github.com/Himneesh-Kalra/go-vector-db/storage"
	"github.com/gofiber/fiber/v2"
)

type SearchRequest struct {
	Query []float32
	K     int
}

func SetupRoutes() *fiber.App {
	app := fiber.New()

	//INSERT a vector
	app.Post("/insert/:table", func(c *fiber.Ctx) error {
		var vec models.Vector
		table := c.Params("table")
		if err := c.BodyParser(&vec); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
		}
		storage.InsertVector(table, vec)
		return c.JSON(fiber.Map{"status": "vector inserted", "id": vec.ID})
	})

	//DELETE a vector
	app.Delete("/delete/:table/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		table := c.Params("table")
		success := storage.DeleteVector(table, id)
		if !success {
			return c.Status(404).JSON(fiber.Map{"error": "Vector not found"})
		}
		return c.JSON(fiber.Map{"Status": "Vector deleted", "id": id})
	})

	//GETALL Vectors
	app.Get("/vectors/:table", func(c *fiber.Ctx) error {
		table := c.Params("table")
		data, ok := storage.GetTable(table)
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "table not found"})
		}
		return c.JSON(data)
	})
	
	//TOPK Search
	app.Post("/search/:table", func(c *fiber.Ctx) error {
		var req SearchRequest
		table := c.Params("table")

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid"})
		}

		results := db.SearchTopK(req.Query, req.K, table)
		return c.JSON(results)
	})

	return app
}
