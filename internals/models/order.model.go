package models

type OrderRequest struct {
	IsPaid          bool     `json:"is_paid"`
	TotalPrice      float64  `json:"total_price" binding:"required"`
	QRCode          string   `json:"qrcode" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	Email           string   `json:"email" binding:"required,email"`
	Phone           string   `json:"phone" binding:"required"`
	ScheduleID      int      `json:"id_schedule" binding:"required"`
	PaymentMethodID int      `json:"id_paymentmethod" binding:"required"`
	Seat            []string `json:"seat"`
}

type OrderResponse struct {
	ID     int      `json:"id" example:"101"`
	Name   string   `json:"name" example:"Rangga Saputra"`
	Email  string   `json:"email" example:"rangga@example.com"`
	Phone  string   `json:"phone" example:"+628123456789"`
	QRCode string   `json:"qrcode" example:"https://example.com/qrcode/101.png"`
	Seat   []string `json:"seat" example:"A1,A2,A3"`
}
