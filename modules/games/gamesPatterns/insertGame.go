package gamespatterns

import (
	"context"
	"fmt"
	"time"

	"github.com/WhoaWicked/board-game-shop/modules/games"
	"github.com/jmoiron/sqlx"
)

type IInsertGameBuilder interface {
	initTransaction() error
	insertGame() error
	insertCategory() error
	insertAttachment() error
	commit() error
	getGameId() string
}

type insertGameBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *games.Game
}

func InsertGameBuilder(db *sqlx.DB, req *games.Game) IInsertGameBuilder {
	return &insertGameBuilder{
		db:  db,
		req: req,
	}
}

type insertGameEngineer struct {
	builder IInsertGameBuilder
}

func InsertGameEngineer(b IInsertGameBuilder) *insertGameEngineer {
	return &insertGameEngineer{builder: b}
}

func (b *insertGameBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}
func (b *insertGameBuilder) insertGame() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	query := `
	INSERT INTO games (title, description, min_players, max_players, playing_time, status)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.Title,
		b.req.Description,
		b.req.MinPlayer,
		b.req.MaxPlayer,
		b.req.PlayingTime,
		b.req.Status).Scan(&b.req.Id); err != nil {
		return fmt.Errorf("insert game failed: %v", err)
	}
	return nil
}
func (b *insertGameBuilder) insertCategory() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	query := `INSERT INTO games_categories (game_id, category_id)
	VALUES`
	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Category {
		valueStack = append(valueStack, b.req.Id, b.req.Category[i].Id)
		if i != len(b.req.Category)-1 {
			query += fmt.Sprintf(`
			($%d, $%d),`, index+1, index+2)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d);`, index+1, index+2)
		}
		index += 2
	}
	if _, err := b.tx.ExecContext(ctx, query, valueStack...); err != nil {
		return fmt.Errorf("insert game categories failed: %v", err)
	}
	return nil
}
func (b *insertGameBuilder) insertAttachment() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	query := `INSERT INTO game_images (game_id, filename, url)
	VALUES`
	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Image {
		valueStack = append(valueStack,
			b.req.Id,
			b.req.Image[i].Filename,
			b.req.Image[i].Url)
		if i != len(b.req.Image)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
	}
	if _, err := b.tx.ExecContext(ctx, query, valueStack...); err != nil {
		return fmt.Errorf("insert game images failed: %v", err)
	}
	return nil
}
func (b *insertGameBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (b *insertGameBuilder) getGameId() string { return b.req.Id }

func (en *insertGameEngineer) InsertGame() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertGame(); err != nil {
		return "", err
	}
	if err := en.builder.insertCategory(); err != nil {
		return "", err
	}
	// if err := en.builder.insertAttachment(); err != nil {
	// 	return "", err
	// }
	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getGameId(), nil
}
