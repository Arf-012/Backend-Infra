package mapper

import (
	"go-api-infra/dto"
	"go-api-infra/models"
	"time"
)

func GuestBookToDTO(guest *models.GuestBook) dto.GuestBookResponse {
	return dto.GuestBookResponse{
		ID: guest.ID,
		Title: guest.Title,
		InstanceName: guest.InstanceName,
		ContactPerson: guest.ContactPerson,
		Email: guest.Email,
		Phone: guest.Phone,
		Description: guest.Description,
		Status: guest.Status,
		RejectionReason: guest.RejectionReason,
		ApprovedAt: guest.ApprovedAt,
		RequestDate: guest.RequestDate,
		ShowInCalendar: guest.ShowInCalendar,
		CreatedAt: guest.CreatedAt,
	}
}

func GuestBookListToDTO(guests []models.GuestBook) []dto.GuestBookResponse {
	result := make([]dto.GuestBookResponse, 0)
	for _, guest := range guests {
		result = append(result, GuestBookToDTO(&guest))
	}
	return result
}

func UpdateGuestBookFromDTO(guest *models.GuestBook, dto dto.GuestBookRequest) {
	if dto.Title != "" {
		guest.Title = dto.Title
	}
	if dto.InstanceName != "" {
		guest.InstanceName = dto.InstanceName
	}
	if dto.ContactPerson != "" {
		guest.ContactPerson = dto.ContactPerson
	}
	if dto.Email != "" {
		guest.Email = dto.Email
	}
	if dto.Phone != "" {
		guest.Phone = dto.Phone
	}
	if dto.Description != "" {
		guest.Description = dto.Description
	}

	// Handle request_date (string → time.Time)
	if dto.RequestDate != "" {
		var parsed time.Time
		var err error

		// Coba parse dengan format RFC3339 dulu (misal: 2025-12-03T00:00:00Z)
		parsed, err = time.Parse(time.RFC3339, dto.RequestDate)
		if err != nil {
			// Kalau gagal, coba format tanggal saja (YYYY-MM-DD)
			parsed, err = time.Parse("2006-01-02", dto.RequestDate)
		}

		if err == nil {
			guest.RequestDate = parsed
		}
	}

	guest.UpdatedAt = time.Now()
}
