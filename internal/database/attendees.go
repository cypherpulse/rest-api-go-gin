package database

import "database/sql"

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	UserId  int `json:"userId" binding:"required"`
	EventId int `json:"eventId" binding:"required"`
}
