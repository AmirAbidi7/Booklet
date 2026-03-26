package models

type Pdf struct {
	ID       uint   `json:"id"`
	Size     uint   `json:"size"`
	Pages    uint   `json:"pages"`
	Checksum string `json:"-"`
	Users    []User `gorm:"many2many:user_pdf;" json:"users"`
}

type UserPdf struct {
	PdfID       uint   `gorm:"primaryKey" json:"pdf_id"`
	Title       string `json:"title"`
	UserID      uint   `gorm:"primaryKey" json:"user_id"`
	CurrentPage uint   `json:"current_page" gorm:"default:1"`
}
