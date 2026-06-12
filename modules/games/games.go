package games

import (
	"github.com/WhoaWicked/board-game-shop/modules/appinfo"
	"github.com/WhoaWicked/board-game-shop/modules/entities"
)

type Game struct {
	Id          string              `db:"id" json:"id"`
	Title       string              `db:"title" json:"title"`
	Description string              `db:"description" json:"description"`
	MinPlayer   int                 `db:"min_players" json:"min_players"`
	MaxPlayer   int                 `db:"max_players" json:"max_players"`
	PlayingTime int                 `db:"playing_time" json:"playing_time"`
	Status      string              `db:"status" json:"status"`
	Category    []*appinfo.Category `json:"categories"`
	Image       []*entities.Images  `json:"images"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

type GameFilter struct {
	Search      string `query:"search"` // title & description
	CountPlayer int    `query:"count_player"`
	MaxTime     int    `query:"max_time"`
	Status      string `query:"status"`
	Categories  []int  `query:"categories"`
	*entities.PaginationReq
	*entities.SortReq
}
