package repository

import (
	"backend/model"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetProducts(ctx context.Context) ([]model.Product, error) {
	rows, err := r.pool.Query(ctx, "select * from get_products;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := RowsToProductList(rows)

	return products, nil
}

func (r *Repository) GetProductById(ctx context.Context, id int) (model.Product, error) {
	var p model.Product

	row := r.pool.QueryRow(ctx, "select * from get_products where id=$1;", id)
	err := row.Scan(&p.Id, &p.Name, &p.Price, &p.Specifications, &p.Images, &p.CategoryName, &p.BrandName, &p.AvgRating)
	if err != nil {
		return model.Product{}, err
	}

	return p, nil
}

func (r *Repository) GetProductsByCategoryName(ctx context.Context, name string) ([]model.Product, error) {
	rows, err := r.pool.Query(ctx, "select * from get_products where category_name ILIKE $1", name)
	if err != nil {
		return nil, err
	}

	products := RowsToProductList(rows)

	return products, nil
}

func (r *Repository) GetProductsByName(ctx context.Context, searchQuery string) ([]model.Product, error) {
	rows, err := r.pool.Query(ctx, "select * from get_products where to_tsvector(name) @@ to_tsquery($1)", searchQuery)
	if err != nil {
		return nil, err
	}

	products := RowsToProductList(rows)

	return products, nil
}

func RowsToProductList(rows pgx.Rows) []model.Product {
	var products []model.Product
	for rows.Next() {
		var p model.Product

		err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Specifications, &p.Images, &p.CategoryName, &p.BrandName, &p.AvgRating)
		if err != nil {
			log.Println("scan error:", err)
		}

		products = append(products, p)
	}

	return products
}
