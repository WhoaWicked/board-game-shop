package gamespatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/WhoaWicked/board-game-shop/modules/games"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type IFindGameBuilder interface {
	initQuery()
	initCountQuery()
	buildWhereSearch()
	buildWhereStatus()
	buildWhereCountPlayer()
	buildWhereMaxTime()
	buildWhereCategories()
	buildSort()
	buildPaginate()
	closeQuery()
	getQuery() string
	setQuery(query string)
	getValues() []any
	setValues(data []any)
	setLastIndex(n int)
	getDb() *sqlx.DB
	reset()
}

type findGameBuilder struct {
	db        *sqlx.DB
	req       *games.GameFilter
	query     string
	values    []any
	lastIndex int
}

func FindGameBuilder(db *sqlx.DB, req *games.GameFilter) IFindGameBuilder {
	return &findGameBuilder{
		db:     db,
		req:    req,
		values: make([]any, 0),
	}
}

type findGameEngineer struct {
	builder IFindGameBuilder
}

func FindGameEngineer(b IFindGameBuilder) *findGameEngineer {
	return &findGameEngineer{builder: b}
}

func (b *findGameBuilder) initQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg(at))
	FROM (
		SELECT
			g.id,
			g.title,
			g.description,
			g.min_players,
			g.max_players,
			g.playing_time,
			g.status,
			(
				SELECT coalesce(array_to_json(array_agg(it)), '[]'::json)
				FROM (
					SELECT i.id, i.filename, i.url
					FROM game_images i
					WHERE i.game_id = g.id
				) AS it
			) AS images,
			(
			SELECT coalesce(array_to_json(array_agg(cat)), '[]'::json)
			FROM (
				SELECT c.id, c.title
				FROM games_categories gc_sub
				JOIN categories c ON c.id = gc_sub.category_id
				WHERE gc_sub.game_id = g.id
			) AS cat
		) AS categories,
		g.created_at,
		g.updated_at
	FROM games g
	WHERE 1 = 1`
}
func (b *findGameBuilder) initCountQuery() {
	b.query += `
	SELECT
	COUNT(*) AS count
	FROM games g
	WHERE 1 = 1`
}
func (b *findGameBuilder) buildWhereSearch() {
	if b.req.Search != "" {
		b.values = append(b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%")
		query := fmt.Sprintf(`
		AND (g.title ILIKE $%d OR g.description ILIKE $%d)`, b.lastIndex+1, b.lastIndex+2)
		temp := b.getQuery()
		temp += query
		b.setQuery(temp)
		b.lastIndex = len(b.values)
	}
}
func (b *findGameBuilder) buildWhereStatus() {
	if b.req.Status != "" {
		b.values = append(b.values, strings.ToLower(b.req.Status))
		query := fmt.Sprintf(`
		AND g.status = $%d`, b.lastIndex+1)
		temp := b.getQuery()
		temp += query
		b.setQuery(temp)
		b.lastIndex = len(b.values)
	}
}

func (b *findGameBuilder) buildWhereCountPlayer() {
	if b.req.CountPlayer > 0 {
		b.values = append(b.values,
			b.req.CountPlayer,
			b.req.CountPlayer)
		query := fmt.Sprintf(`
		AND (g.min_players <= $%d AND g.max_players >= $%d)`, b.lastIndex+1, b.lastIndex+2)
		temp := b.getQuery()
		temp += query
		b.setQuery(temp)
		b.lastIndex = len(b.values)
	}
}

func (b *findGameBuilder) buildWhereMaxTime() {
	if b.req.MaxTime > 0 {
		b.values = append(b.values, b.req.MaxTime)
		query := fmt.Sprintf(`
		AND g.playing_time <= $%d`, b.lastIndex+1)
		temp := b.getQuery()
		temp += query
		b.setQuery(temp)
		b.lastIndex = len(b.values)
	}
}
func (b *findGameBuilder) buildWhereCategories() {
	if len(b.req.Categories) > 0 {
		b.values = append(b.values, pq.Array(b.req.Categories))
		query := fmt.Sprintf(`
		AND EXISTS (
		SELECT 1
		FROM games_categories gc
		WHERE gc.game_id = g.id
		AND gc.category_id = ANY($%d)
		)`, b.lastIndex+1)
		temp := b.getQuery()
		temp += query
		b.setQuery(temp)
		b.lastIndex = len(b.values)
	}
}
func (b *findGameBuilder) buildSort() {
	gameByMap := map[string]string{
		"id":         "g.id",
		"title":      "g.title",
		"created_at": "g.created_at",
	}
	if gameByMap[b.req.OrderBy] == "" {
		b.req.OrderBy = gameByMap["created_at"]
	} else {
		b.req.OrderBy = gameByMap[b.req.OrderBy]
	}
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[b.req.Sort] == "" {
		b.req.Sort = sortMap["DESC"]
	} else {
		b.req.Sort = sortMap[strings.ToUpper(b.req.Sort)]
	}
	b.query += fmt.Sprintf(`
	ORDER BY %s %s`, b.req.OrderBy, b.req.Sort)
	b.lastIndex = len(b.values)

}
func (b *findGameBuilder) buildPaginate() {
	b.values = append(b.values,
		(b.req.Page-1)*b.req.Limit,
		b.req.Limit)
	b.query += fmt.Sprintf(`
	OFFSET $%d LIMIT $%d`, b.lastIndex+1, b.lastIndex+2)
	b.lastIndex = len(b.values)
}
func (b *findGameBuilder) closeQuery() {
	b.query += `
	) AS at;`
}
func (b *findGameBuilder) getQuery() string      { return b.query }
func (b *findGameBuilder) setQuery(query string) { b.query = query }
func (b *findGameBuilder) getValues() []any      { return b.values }
func (b *findGameBuilder) setValues(data []any)  { b.values = data }
func (b *findGameBuilder) setLastIndex(n int)    { b.lastIndex = n }
func (b *findGameBuilder) getDb() *sqlx.DB       { return b.db }
func (b *findGameBuilder) reset() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastIndex = 0
}

func (en *findGameEngineer) FindGames() []*games.Game {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer en.builder.reset()
	en.builder.initQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	en.builder.buildWhereCountPlayer()
	en.builder.buildWhereMaxTime()
	en.builder.buildWhereCategories()
	en.builder.buildSort()
	en.builder.buildPaginate()
	en.builder.closeQuery()
	gamesRaw := make([]byte, 0)
	if err := en.builder.getDb().GetContext(ctx, &gamesRaw, en.builder.getQuery(), en.builder.getValues()...); err != nil {
		log.Printf("get games failed: %v", err)
		return make([]*games.Game, 0)
	}
	games := make([]*games.Game, 0)
	if err := json.Unmarshal(gamesRaw, &games); err != nil {
		log.Printf("unmarshal games failed: %v", err)
	}
	return games
}

func (en *findGameEngineer) CountGame() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer en.builder.reset()
	en.builder.initCountQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	en.builder.buildWhereCountPlayer()
	en.builder.buildWhereMaxTime()
	en.builder.buildWhereCategories()
	var count int
	if err := en.builder.getDb().GetContext(ctx, &count, en.builder.getQuery(), en.builder.getValues()...); err != nil {
		log.Printf("count games failed: %v\n", err)
		return 0
	}
	return count
}
