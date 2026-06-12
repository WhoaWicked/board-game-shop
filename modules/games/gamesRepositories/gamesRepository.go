package gamesrepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/WhoaWicked/board-game-shop/modules/appinfo"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
	"github.com/WhoaWicked/board-game-shop/modules/games"
	gamespatterns "github.com/WhoaWicked/board-game-shop/modules/games/gamesPatterns"
	"github.com/jmoiron/sqlx"
)

type IGamesRepository interface {
	FindOneGame(gameId string) (*games.Game, error)
	FindGames(req *games.GameFilter) ([]*games.Game, int)
	InsertGame(req *games.Game) (*games.Game, error)
	DeleteGame(gameId string) error
	UpdateGame(req *games.Game) (*games.Game, error)
}

type gamesRepository struct {
	db *sqlx.DB
}

func GamesRepository(db *sqlx.DB) IGamesRepository {
	return &gamesRepository{db: db}
}

func (r *gamesRepository) FindOneGame(gameId string) (*games.Game, error) {
	query := `
	SELECT
	to_jsonb(t)
	FROM (
		SELECT
			g.id,
			g.title,
			g.description,
			g.status,
			g.min_players,
			g.max_players,
			g.playing_time,
			(
				SELECT coalesce(array_to_json(array_agg(it)), '[]'::json)
				FROM (
					SELECT
						i.id,
						i.filename,
						i.url
					FROM game_images i
					WHERE i.game_id = g.id
				) AS it
			) AS images,
			(
				SELECT coalesce(array_to_json(array_agg(cat)), '[]'::json)
				FROM (
					SELECT
						c.id,
						c.title
					FROM games_categories gc
					JOIN categories c on c.id = gc.category_id
					WHERE gc.game_id = g.id
				) AS cat
			) AS categories,
			g.created_at,
			g.updated_at
		FROM games g
		WHERE g.id = $1
	) AS t`
	gameBytes := make([]byte, 0)
	game := &games.Game{
		Image:    make([]*entities.Images, 0),
		Category: make([]*appinfo.Category, 0),
	}
	if err := r.db.Get(&gameBytes, query, gameId); err != nil {
		return nil, fmt.Errorf("find one game failed: %v", err)
	}
	if err := json.Unmarshal(gameBytes, game); err != nil {
		return nil, fmt.Errorf("unmarshal game failed: %v", err)
	}
	return game, nil
}

func (r *gamesRepository) FindGames(req *games.GameFilter) ([]*games.Game, int) {
	builder := gamespatterns.FindGameBuilder(r.db, req)
	engineer := gamespatterns.FindGameEngineer(builder)
	return engineer.FindGames(), engineer.CountGame()
}

func (r *gamesRepository) InsertGame(req *games.Game) (*games.Game, error) {
	builder := gamespatterns.InsertGameBuilder(r.db, req)
	engineer := gamespatterns.InsertGameEngineer(builder)
	gameId, err := engineer.InsertGame()
	if err != nil {
		return nil, err
	}
	game, err := r.FindOneGame(gameId)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (r *gamesRepository) DeleteGame(gameId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	query := `DELETE FROM games WHERE id = $1;`
	result, err := r.db.ExecContext(ctx, query, gameId)
	if err != nil {
		return fmt.Errorf("delete game failed: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected failed: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("game not found")
	}
	return nil
}

func (r *gamesRepository) UpdateGame(req *games.Game) (*games.Game, error) {
	builder := gamespatterns.UpdateGameBuilder(r.db, req)
	engineer := gamespatterns.UpdateGameEngineer(builder)
	if err := engineer.UpdateGame(); err != nil {
		return nil, err
	}
	game, err := r.FindOneGame(req.Id)
	if err != nil {
		return nil, err
	}
	return game, nil
}
