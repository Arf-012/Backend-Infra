package handlers

import (
	"encoding/csv"
	"fmt"
	"go-api-infra/database"
	"go-api-infra/dto"
	"go-api-infra/mapper"
	"go-api-infra/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
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
	return ExportGuestBookExcel(c)
}

func ExportGuestBookExcel(c *fiber.Ctx) error {
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

	templatePath := "Rekap Buku Tamu Digital.xlsx"
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		// Jika template tidak ada, buat file baru saja sebagai fallback
		f = excelize.NewFile()
	}
	defer f.Close()

	sheetName := "Sheet1"
	// Memastikan kita menulis di sheet yang benar, atau ambil sheet pertama
	index, _ := f.GetSheetIndex(sheetName)
	if index == -1 {
		sheetName = f.GetSheetName(0)
	}

	// Mulai mengisi data dari baris tertentu, misalnya baris ke-2 (asumsi baris 1 adalah header di template)
	startRow := 2

	// Style untuk header: Background Biru, Teks Putih, Font Bold
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#0000FF"}, Pattern: 1},
		Font: &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// Header sesuai permintaan
	headers := []string{"NO", "Tanggal", "Waktu", "Judul", "Nama/Instansi", "Nama Yang Dituju", "Email", "Nomor Whatsapp", "Deskripsi"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	for i, g := range guests {
		row := startRow + i
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), g.CreatedAt.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), g.CreatedAt.Format("15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), g.Title)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), g.InstanceName)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), g.ContactPerson)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), g.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), g.Phone)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), g.Description)
	}

	fileName := fmt.Sprintf("rekap_guest_books_%s.xlsx", time.Now().Format("20060102"))

	// Simpan ke file temporary untuk dikirim
	tmpPath := "/tmp/" + fileName
	if err := f.SaveAs(tmpPath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate excel file",
		})
	}
	defer os.Remove(tmpPath)

	return c.Download(tmpPath, fileName)
}
