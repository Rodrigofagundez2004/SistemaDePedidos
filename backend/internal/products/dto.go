package products

type CreateProductRequest struct 
{
	CategoryID  string `json:"category_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	PriceCents  int    `json:"price_cents" binding:"required,min=0"`
	Currency    string `json:"currency"`
	Stock       int    `json:"stock" binding:"min=0"`
	ImageURL    string `json:"image_url"`

}
type UpdateProductRequest struct{
	Name        *string `json:"name"`
	Description *string `json:"description"`
	PriceCents  *int    `json:"price_cents"`
	Stock       *int    `json:"stock"`
	ImageURL    *string `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}