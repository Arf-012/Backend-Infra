package handlers

import (
	"encoding/csv"
	"fmt"
	"go-api-infra/database"
	"go-api-infra/dto"
	"go-api-infra/mapper"
	"go-api-infra/models"
	"time"

	"github.com/gofiber/fiber/v2"
)


func CreateGuest(c *fiber.Ctx) error {
	var input dto.GuestBookRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "invalid form data",
		})
	}

	guestbook := models.GuestBook{
		Title: input.Title,
		InstanceName: input.InstanceName,
		ContactPerson: input.ContactPerson,
		Email: input.Email,
		Phone: input.Phone,
		Description: input.Description,
		Status: "pending",
	}

	if input.RequestDate != "" {
		parsed, err := time.Parse("2006-01-02", input.RequestDate)
		if err == nil {
			guestbook.RequestDate = parsed
		}
	}


	if err := database.DB.Create(&guestbook).Error; err != nil {	
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to create guest book",
		})
	}

	response := mapper.GuestBookToDTO(&guestbook) 

	return c.Status(200).JSON(response)
}

func EditGuest(c *fiber.Ctx) error {
	id := c.Params("id")
	var guest models.GuestBook

	if err := database.DB.First(&guest, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "cannot find guest",
		})
	}

	var input dto.GuestBookRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "invalid form data",
		})
	}

	mapper.UpdateGuestBookFromDTO(&guest, input)

	if err := database.DB.Save(&guest).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"erros": "failed to save guest",
		})
	}

	response := mapper.GuestBookToDTO(&guest)
	
	return c.Status(200).JSON(response)
}

func GetAllGuestBook(c *fiber.Ctx) error {
	var guests []models.GuestBook	

	if err := database.DB.Order("created_at DESC").Find(&guests).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch guestbook",
		})
	}

	response := mapper.GuestBookListToDTO(guests)

	return c.Status(200).JSON(response)
}

func GetGuestBook(c *fiber.Ctx) error {
	id := c.Params("id")
	var guest models.GuestBook

	if err := database.DB.First(&guest, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch guestbook",
		})
	}

	response := mapper.GuestBookToDTO(&guest)

	return c.Status(200).JSON(response)
}

func DeleteGuest(c *fiber.Ctx) error {
	id := c.Params("id")
	var guest models.GuestBook

	if err := database.DB.First(&guest, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch guest",
		})
	}

	if err := database.DB.Delete(&guest).Error; err != nil{
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to delete guest",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "successfully delete guestbooks",
	})
}

func ApprovedGuest(c *fiber.Ctx) error {
	id := c.Params("id")
	var guest models.GuestBook

	if err := database.DB.First(&guest, id).Error; err != nil {
		return c.Status(200).JSON(fiber.Map{
			"error": "failed to fetch",
		})	
	}


	guest.Status = "approved"
	guest.ApprovedAt = time.Now()

	if err := database.DB.Save(&guest).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to save guest",
		})
	}

	response := mapper.GuestBookToDTO(&guest)

	return c.Status(200).JSON(response)
}

func RejectGuestBook(c *fiber.Ctx) error {
	id := c.Params("id")
	rejection_reason := c.FormValue("rejection_reason")
	var guest models.GuestBook

	if err := database.DB.First(&guest, id).Error; err != nil {
		return c.Status(200).JSON(fiber.Map{
			"error": "failed to fetch",
		})
	}

	guest.RejectionReason = rejection_reason
	guest.Status = "rejected"

	if err := database.DB.Save(&guest).Error; err != nil {
		return c.Status(200).JSON(fiber.Map{
			"error": "failed to saved",
		})
	}

	response := mapper.GuestBookToDTO(&guest)

	return c.Status(200).JSON(response)

}

func ExportGuestBook(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")

	var guests []models.GuestBook
	query := database.DB.Model(&models.GuestBook{})

	if from != "" && to != "" {
		query = query.Where("created_at BETWEEN ? AND ?", from+" 00:00:00", to+" 23:59:59")
	}

	if err := query.Order("created_at DESC").Find(&guests).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch data for export",
		})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=guest_books_%s.csv", time.Now().Format("20060102")))

	writer := csv.NewWriter(c)
	defer writer.Flush()

	// Header
	writer.Write([]string{"ID", "Title", "Instance Name", "Contact Person", "Email", "Phone", "Status", "Request Date", "Created At"})

	for _, g := range guests {
		writer.Write([]string{
			fmt.Sprintf("%d", g.ID),
			g.Title,
			g.InstanceName,
			g.ContactPerson,
			g.Email,
			g.Phone,
			g.Status,
			g.RequestDate.Format("2006-01-02"),
			g.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return nil
}