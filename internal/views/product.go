package views

import (
	"net/http"
	"satellity/internal/models"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/gofrs/uuid"
)

// ProductView is the response body of product
type ProductView struct {
	Type       string    `json:"type"`
	ProductID  string    `json:"product_id"`
	ShortID    string    `json:"short_id"`
	Name       string    `json:"name"`
	Body       string    `json:"body"`
	CoverURL   string    `json:"cover_url"`
	Source     string    `json:"source"`
	Tags       []string  `json:"tags"`
	ViewsCount int64     `json:"views_count"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       UserView  `json:"user"`
}

func buildProduct(p *models.Product) ProductView {
	view := ProductView{
		Type:       "product",
		ProductID:  p.ProductID,
		Name:       p.Name,
		Body:       p.Body,
		CoverURL:   p.CoverURL,
		Source:     p.Source,
		Tags:       p.Tags,
		ViewsCount: p.ViewsCount,
		UserID:     p.UserID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
	id, _ := uuid.FromString(p.ProductID)
	view.ShortID = base58.Encode(id.Bytes())
	if p.User != nil {
		view.User = buildUser(p.User)
	}
	return view
}

// RenderProduct response a product
func RenderProduct(w http.ResponseWriter, r *http.Request, product *models.Product) {
	RenderResponse(w, r, buildProduct(product))
}

// RenderProducts response a bundle of products
func RenderProducts(w http.ResponseWriter, r *http.Request, products []*models.Product) {
	productViews := make([]ProductView, len(products))
	for i, product := range products {
		productViews[i] = buildProduct(product)
	}
	RenderResponse(w, r, productViews)
}
