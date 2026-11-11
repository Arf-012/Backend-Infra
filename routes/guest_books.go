package routes

import (
	"go-api-infra/handlers"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func GuestBookRoutes(router fiber.Router) {
	guest := router.Group("/guest-book")
	guest.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	guest.Post("/create", handlers.CreateGuest)
	guest.Post("/edit/:id", handlers.EditGuest)
	guest.Get("/showAll", handlers.GetAllGuestBook)
	guest.Get("/export", handlers.ExportGuestBook)
	guest.Get("/export-excel", handlers.ExportGuestBookExcel)
	guest.Get("/show/:id", handlers.GetGuestBook)
	guest.Put("/approve/:id", handlers.ApprovedGuest)
	guest.Put("/reject/:id", handlers.RejectGuestBook)
}