package dto

import "time"

type GuestBookRequest struct {
	Title 			string 		`form:"title"`
	InstanceName 	string 		`form:"instance_name"`
	ContactPerson 	string 		`form:"contact_person"`
	Email 			string 		`form:"email"`
	Phone 			string 		`form:"phone"`
	Description 	string 		`form:"description"`
	RequestDate 	string		`form:"request_date"`
}

type GuestBookResponse struct {
	ID				uint		`json:"id"`
	Title 			string 		`json:"title"`
	InstanceName 	string 		`json:"instance_name"`
	ContactPerson 	string 		`json:"contact_person"`
	Email 			string 		`json:"email"`
	Phone 			string 		`json:"phone"`
	Description 	string 		`json:"description"`
	Status 			string 		`json:"status"`
	RejectionReason string 		`json:"rejection_reason"`
	ApprovedAt 		time.Time 	`json:"approved_at"`
	RequestDate 	time.Time 	`json:"request_date"`
	ShowInCalendar 	bool 		`json:"show_in_calendar"`
	CreatedAt 		time.Time 	`json:"created_at"`
}