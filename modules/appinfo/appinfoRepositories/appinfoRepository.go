package appinforepositories

import (
	"fmt"
	"strings"

	"github.com/WhoaWicked/board-game-shop/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepository interface {
	FindCategories(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
}

type appinfoRepository struct {
	db *sqlx.DB
}

func AppinfoRepository(db *sqlx.DB) IAppinfoRepository {
	return &appinfoRepository{db: db}
}

func (r *appinfoRepository) FindCategories(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	query := `SELECT id, title FROM categories`
	filterValues := make([]any, 0)
	if req.Title != "" {
		query += ` WHERE (LOWER(title) LIKE $1)`
		filterValues = append(filterValues, "%"+strings.ToLower(req.Title)+"%")
	}
	query += `;`
	categories := make([]*appinfo.Category, 0)
	if err := r.db.Select(&categories, query, filterValues...); err != nil {
		return nil, fmt.Errorf("select categories failed: %v", err)
	}
	return categories, nil
}
