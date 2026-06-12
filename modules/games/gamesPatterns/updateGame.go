package gamespatterns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/WhoaWicked/board-game-shop/modules/games"
	"github.com/jmoiron/sqlx"
)

type IUpdateGameBuilder interface {
	initTransaction() error
	buildGameQuery()
	updateCategory() error
	updateGame() error
	commit() error
}

type updateGameBuilder struct {
	db     *sqlx.DB
	tx     *sqlx.Tx
	req    *games.Game
	query  string
	values []any
}

func UpdateGameBuilder(db *sqlx.DB, req *games.Game) IUpdateGameBuilder {
	return &updateGameBuilder{
		db:  db,
		req: req,
	}
}

type updateGameEngineer struct {
	builder IUpdateGameBuilder
}

func UpdateGameEngineer(b IUpdateGameBuilder) *updateGameEngineer {
	return &updateGameEngineer{builder: b}
}

func (b *updateGameBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}
func (b *updateGameBuilder) buildGameQuery() {
	b.query = `UPDATE games SET `
	field := make([]string, 0)
	if b.req.Title != "" {
		b.values = append(b.values, b.req.Title)
		field = append(field, fmt.Sprintf(`title = $%d`, len(b.values)))
	}
	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		field = append(field, fmt.Sprintf(`description = $%d`, len(b.values)))
	}
	if b.req.MinPlayer > 0 {
		b.values = append(b.values, b.req.MinPlayer)
		field = append(field, fmt.Sprintf(`min_players = $%d`, len(b.values)))
	}
	if b.req.MaxPlayer > 0 {
		b.values = append(b.values, b.req.MaxPlayer)
		field = append(field, fmt.Sprintf(`max_players = $%d`, len(b.values)))
	}
	if b.req.PlayingTime > 0 {
		b.values = append(b.values, b.req.PlayingTime)
		field = append(field, fmt.Sprintf(`playing_time = $%d`, len(b.values)))
	}
	if b.req.Status != "" {
		b.values = append(b.values, b.req.Status)
		field = append(field, fmt.Sprintf(`status = $%d`, len(b.values)))
	}
	b.query += strings.Join(field, ", ")
	b.values = append(b.values, b.req.Id)
	b.query += fmt.Sprintf(` WHERE id = $%d`, len(b.values))
}
func (b *updateGameBuilder) updateCategory() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if b.req.Category == nil {
		return nil
	}
	deleteQuery := `DELETE FROM games_categories WHERE game_id = $1;`
	if _, err := b.tx.ExecContext(ctx, deleteQuery, b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete old categories failed: %v", err)
	}
	if len(b.req.Category) == 0 {
		return nil
	}
	insertQuery := `INSERT INTO games_categories (game_id, category_id) VALUES `
	insertValues := make([]any, 0)
	valueString := make([]string, 0)
	for i, cat := range b.req.Category {
		valueString = append(valueString, fmt.Sprintf(`($%d, $%d)`, i*2+1, i*2+2))
		insertValues = append(insertValues, b.req.Id, cat.Id)
	}
	insertQuery += strings.Join(valueString, ", ") + ";"
	if _, err := b.tx.ExecContext(ctx, insertQuery, insertValues...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update categories failed: %v", err)
	}
	return nil
}

func (b *updateGameBuilder) updateGame() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if _, err := b.tx.ExecContext(ctx, b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update game failed: %v", err)
	}
	return nil
}
func (b *updateGameBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		b.tx.Rollback()
		return err
	}
	return nil
}

func (en *updateGameEngineer) UpdateGame() error {
	en.builder.initTransaction()
	en.builder.buildGameQuery()
	if err := en.builder.updateGame(); err != nil {
		return err
	}
	if err := en.builder.updateCategory(); err != nil {
		return err
	}
	en.builder.commit()
	return nil
}
