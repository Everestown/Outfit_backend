package dto

type AddItemRequest struct {
	VariantID uint `json:"variant_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}
