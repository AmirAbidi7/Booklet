package models

type Pdf struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Size     uint   `json:"size"`
	Pages    uint   `json:"pages"`
	Location string `json:"location"`
	Checksum string `json:"checksum"`
	Users    []User `gorm:"many2many:user_pdf;" json:"users"`
}

type UserPdf struct {
	PdfID       int  `gorm:"primaryKey" json:"pdf_id"`
	UserID      int  `gorm:"primaryKey" json:"user_id"`
	CurrentPage uint `json:"current_page"`
}
