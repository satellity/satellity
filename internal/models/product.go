package models

import (
	"context"
	"fmt"
	"satellity/internal/durable"
	"satellity/internal/session"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
)

// Product contains websites or apps
type Product struct {
	ProductID  string
	Name       string
	Body       string
	CoverURL   string
	Source     string
	Tags       pq.StringArray
	ViewsCount int64
	Score      int64
	UserID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	User *User
}

var productColumns = []string{"product_id", "name", "body", "cover_url", "source", "tags", "views_count", "score", "user_id", "created_at", "updated_at"}

func (p *Product) values() []interface{} {
	return []interface{}{p.ProductID, p.Name, p.Body, p.CoverURL, p.Source, p.Tags, p.ViewsCount, p.Score, p.UserID, p.CreatedAt, p.UpdatedAt}
}

func productFromRow(row durable.Row) (*Product, error) {
	var p Product
	err := row.Scan(&p.ProductID, &p.Name, &p.Body, &p.CoverURL, &p.Source, &p.Tags, &p.ViewsCount, &p.Score, &p.UserID, &p.CreatedAt, &p.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (user *User) CreateProduct(ctx context.Context, name, body, cover, source string, tags []string) (*Product, error) {
	name, body, cover, source = strings.TrimSpace(name), strings.TrimSpace(body), strings.TrimSpace(cover), strings.TrimSpace(source)
	if len(name) < 1 {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "name", "name is invalid", name)
	}
	if len(body) < 1 {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "body", "body is invalid", body)
	}
	if len(cover) < 1 {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "cover", "cover is invalid", cover)
	}
	if len(source) < 1 {
		return nil, session.BadDataErrorWithFieldAndData(ctx, "source", "source is invalid", source)
	}

	t := time.Now()
	p := &Product{
		ProductID: uuid.Must(uuid.NewV4()).String(),
		Name:      name,
		Body:      body,
		CoverURL:  cover,
		Source:    source,
		Tags:      pq.StringArray(tags),
		UserID:    user.UserID,
		CreatedAt: t,
		UpdatedAt: t,
	}
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		rows := [][]interface{}{
			p.values(),
		}
		_, err := tx.CopyFrom(ctx, pgx.Identifier{"products"}, productColumns, pgx.CopyFromRows(rows))
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	p.User = user
	return p, nil
}

func (user *User) UpdateProduct(ctx context.Context, productID, name, body, cover, source string, tags []string) (*Product, error) {
	name, body, cover, source = strings.TrimSpace(name), strings.TrimSpace(body), strings.TrimSpace(cover), strings.TrimSpace(source)
	var p *Product
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		existing, err := findProduct(ctx, tx, productID)
		if err != nil || existing == nil {
			return err
		}
		p = existing
		if p.UserID != user.UserID && !user.isAdmin() {
			return session.AuthorizationError(ctx)
		}
		p.User = user
		if p.UserID != user.UserID {
			u, err := findUserByID(ctx, tx, p.UserID)
			if err != nil {
				return err
			}
			p.User = u
		}
		if name != "" {
			p.Name = name
		}
		if body != "" {
			p.Body = body
		}
		if cover != "" {
			p.CoverURL = cover
		}
		if source != "" {
			p.Source = source
		}
		if len(tags) > 1 {
			p.Tags = pq.StringArray(tags)
		}
		p.UpdatedAt = time.Now()
		columns, placeholders := durable.PrepareColumnsWithParams([]string{"name", "body", "cover_url", "source", "tags", "updated_at"})
		values := []interface{}{p.Name, p.Body, p.CoverURL, p.Source, p.Tags, p.UpdatedAt}
		_, err = tx.Exec(ctx, fmt.Sprintf("UPDATE products SET (%s)=(%s) WHERE product_id='%s'", columns, placeholders, p.ProductID), values...)
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return p, nil
}

func FindProducts(ctx context.Context) ([]*Product, error) {
	var products []*Product
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM products ORDER BY score,created_at DESC LIMIT 99", strings.Join(productColumns, ","))
		rows, err := tx.Query(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		var userIds []string
		for rows.Next() {
			product, err := productFromRow(rows)
			if err != nil {
				return err
			}
			userIds = append(userIds, product.UserID)
			products = append(products, product)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		userSet, err := readUserSet(ctx, tx, userIds)
		if err != nil {
			return err
		}
		for i, product := range products {
			products[i].User = userSet[product.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return products, nil
}

func FindProduct(ctx context.Context, productID string) (*Product, error) {
	var p *Product
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		var err error
		p, err = findProduct(ctx, tx, productID)
		if err != nil || p == nil {
			return err
		}
		user, err := findUserByID(ctx, tx, p.UserID)
		if err != nil {
			return err
		}
		p.User = user
		return err
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return p, nil
}

func findProduct(ctx context.Context, tx pgx.Tx, productID string) (*Product, error) {
	if id, _ := uuid.FromString(productID); id.String() == uuid.Nil.String() {
		return nil, nil
	}
	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT %s FROM products WHERE product_id=$1", strings.Join(productColumns, ",")), productID)
	return productFromRow(row)
}

func RelatedProducts(ctx context.Context, id string) ([]*Product, error) {
	var products []*Product
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM products WHERE product_id>$1 LIMIT 3", strings.Join(productColumns, ","))
		rows, err := tx.Query(ctx, query, id)
		if err != nil {
			return err
		}
		defer rows.Close()

		var userIds []string
		for rows.Next() {
			product, err := productFromRow(rows)
			if err != nil {
				return err
			}
			userIds = append(userIds, product.UserID)
			products = append(products, product)
		}

		userSet, err := readUserSet(ctx, tx, userIds)
		if err != nil {
			return err
		}
		for i, product := range products {
			products[i].User = userSet[product.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return products, nil
}

func SearchProducts(ctx context.Context, query string) ([]*Product, error) {
	keys := strings.Split(strings.TrimSpace(query), ",")
	var products []*Product
	err := session.Database(ctx).RunInTransaction(ctx, func(tx pgx.Tx) error {
		query := fmt.Sprintf("SELECT %s FROM products WHERE $1 <@ tags LIMIT 50", strings.Join(productColumns, ","))
		rows, err := tx.Query(ctx, query, pq.Array(keys))
		if err != nil {
			return err
		}
		defer rows.Close()

		var userIds []string
		for rows.Next() {
			product, err := productFromRow(rows)
			if err != nil {
				return err
			}
			userIds = append(userIds, product.UserID)
			products = append(products, product)
		}

		userSet, err := readUserSet(ctx, tx, userIds)
		if err != nil {
			return err
		}
		for i, product := range products {
			products[i].User = userSet[product.UserID]
		}
		return nil
	})
	if err != nil {
		return nil, session.TransactionError(ctx, err)
	}
	return products, nil
}
