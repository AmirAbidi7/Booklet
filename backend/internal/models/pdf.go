package models

import "time"

type Pdf struct {
	ID        uint   `json:"id"`
	Size      uint   `json:"size"`
	Pages     uint   `json:"pages"`
	Checksum  string `json:"-"`
	CreatedAt time.Time
	Users     []User `gorm:"many2many:user_pdf;" json:"users"`
}

type UserPdf struct {
	PdfID       uint   `gorm:"primaryKey" json:"pdf_id"`
	UserID      uint   `gorm:"primaryKey" json:"user_id"`
	Title       string `json:"title"`
	CurrentPage uint   `json:"current_page" gorm:"default:1"`
	CreatedAt   time.Time
}
