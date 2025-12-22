package dto

type CreateOrderRequest struct {
	TotalAmount float64 `json:"total_amount" binding:"required"`
	Postcode    string  `json:"postcode" binding:"required"`
	Country     string  `json:"country" binding:"required"`
	City        string  `json:"city" binding:"required"`
	Street      string  `json:"street" binding:"required"`
	House       string  `json:"house" binding:"required"`
	Apartment   string  `json:"apartment"`
	Phone       string  `json:"phone" binding:"required"`
	Surname     string  `json:"surname" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Patronymic  string  `json:"patronymic"`
}
