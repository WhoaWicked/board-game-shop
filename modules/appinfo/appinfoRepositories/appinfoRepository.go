package appinforepositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/WhoaWicked/board-game-shop/modules/appinfo"
	"github.com/jmoiron/sqlx"
)

type IAppinfoRepository interface {
	FindCategories(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
	UpdateCategory(req *appinfo.Category) error
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

func (r *appinfoRepository) InsertCategory(req []*appinfo.Category) error {
	ctx := context.Background()
	query := `INSERT INTO categories (title) VALUES`
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	valuesStack := make([]any, 0, len(req))
	for i, cat := range req {
		valuesStack = append(valuesStack, cat.Title)
		if i != (len(req) - 1) {
			query += fmt.Sprintf(`
			($%d),`, i+1)
		} else {
			query += fmt.Sprintf(`
			($%d)`, i+1)
		}
	}
	query += `
	RETURNING id;`
	rows, err := tx.QueryxContext(ctx, query, valuesStack...)
	if err != nil {
		return fmt.Errorf("insert categories failed: %v", err)
	}
	defer rows.Close()
	var i int
	for rows.Next() {
		if err := rows.Scan(&req[i].Id); err != nil {
			return fmt.Errorf("scan categories id failed: %v", err)
		}
		i++
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *appinfoRepository) UpdateCategory(req *appinfo.Category) error {
	query := `UPDATE categories SET title = $1 WHERE id = $2;`
	result, err := r.db.ExecContext(context.Background(), query, req.Title, req.Id)
	if err != nil {
		return fmt.Errorf("update category failed: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected failed: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category id %d not found", req.Id)
	}
	return nil
}

func (r *appinfoRepository) DeleteCategory(categoryId int) error {
	ctx := context.Background()
	query := `DELETE FROM categories WHERE id = $1;`
	result, err := r.db.ExecContext(ctx, query, categoryId)
	if err != nil {
		return fmt.Errorf("delete category failed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected failed: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category id %d not found", categoryId)
	}

	return nil
}
