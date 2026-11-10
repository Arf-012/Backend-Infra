package models

import (
	"time"

	"gorm.io/gorm"
)

type GuestBook struct {
	gorm.Model
	Title 	string `json:"title"`
	InstanceName string `json:"instance_name"`
	ContactPerson string `json:"contact_person"`
	Email 	string 	`json:"email"`
	Phone 	string `json:"phone"`
	Description string `json:"description"`
	Status string `json:"status" gorm:"type:enum('pending','approved','rejected','hidden');default:'pending'"`
	RejectionReason string `json:"rejection_reason" gorm:"type:text"`
	ApprovedAt time.Time 	`json:"approved_at"`
	RequestDate time.Time 	`json:"request_date" gorm:"type:date"`
	ShowInCalendar 	bool `json:"show_in_calendar"`
}
