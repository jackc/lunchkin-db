package main

import (
	"bytes"
	"fmt"
	"github.com/JackC/pgx"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getPlayers(w http.ResponseWriter, req *http.Request) {
	players := make([]*player, 0, 16)

	err := pool.SelectFunc("getPlayers", func(r *pgx.DataRowReader) error {
		var p player
		p.player_id = r.ReadValue().(int32)
		p.name = r.ReadValue().(string)
		players = append(players, &p)
		return nil
	})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	RenderPlayersIndex(w, players)
}

func playerPath(player_id int32) string {
	return fmt.Sprintf("/players/%d", player_id)
}

func deletePlayerPath(player_id int32) string {
	return fmt.Sprintf("/players/%d/delete", player_id)
}

func createPlayer(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("player_name")

	if name == "" {
		w.WriteHeader(422)
		fmt.Fprintln(w, `Request must include the attribute "name"`)
		return
	}

	if _, err := pool.SelectValue("createPlayer", name); err == nil {
		http.Redirect(w, req, "/players", http.StatusSeeOther)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func deletePlayer(w http.ResponseWriter, req *http.Request) {
	var err error
	var playerId int64
	if playerId, err = strconv.ParseInt(req.FormValue("id"), 10, 32); err != nil {
		http.NotFound(w, req)
		return
	}

	if _, err := pool.Execute("deletePlayer", int32(playerId)); err != nil {
		fmt.Fprintf(os.Stderr, "deletePlayer: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/players", http.StatusSeeOther)
}

func getGames(w http.ResponseWriter, req *http.Request) {
	games, err := SelectAllGamesWithDetails()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	RenderGamesIndex(w, games)
}

func gamePath(game_id int32) string {
	return fmt.Sprintf("/api/v1/games/%d", game_id)
}

type DateField struct {
	Raw    string
	Errors []error
	Value  time.Time
}

func (f *DateField) Parse(raw string) {
	f.Raw = raw
	f.Errors = make([]error, 0)

	var err error
	f.Value, err = time.Parse("2006-01-02", raw)
	if err != nil {
		f.Errors = append(f.Errors, err)
	}
}

type Int16Field struct {
	Raw    string
	Errors []error
	Value  int16
}

func (f *Int16Field) Parse(raw string) {
	f.Raw = raw
	f.Errors = make([]error, 0)

	var err error
	var i64 int64
	i64, err = strconv.ParseInt(raw, 10, 16)
	if err != nil {
		f.Errors = append(f.Errors, err)
	}
	f.Value = int16(i64)
}

type Int32Field struct {
	Raw    string
	Errors []error
	Value  int32
}

func (f *Int32Field) Parse(raw string) {
	f.Raw = raw
	f.Errors = make([]error, 0)

	var err error
	var i64 int64
	i64, err = strconv.ParseInt(raw, 10, 32)
	if err != nil {
		f.Errors = append(f.Errors, err)
	}
	f.Value = int32(i64)
}

type GameFormPlayer struct {
	PlayerId       Int32Field
	Level          Int16Field
	EffectiveLevel Int16Field
	Winner         bool
}
type GameForm struct {
	Date    DateField
	Length  Int16Field
	Players []GameFormPlayer
}

func newGame(w http.ResponseWriter, req *http.Request) {
	players, err := SelectAllPlayers()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	RenderGamesNew(w, players)
}

func deleteGamePath(gameId int32) string {
	return fmt.Sprintf("/games/%d/delete", gameId)
}

func createGame(w http.ResponseWriter, req *http.Request) {

	var err error
	var game GameForm

	err = req.ParseForm()
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, "Unable to parse form data: %v", err)
		return
	}

	game.Date.Parse(req.FormValue("Game.Date"))
	game.Length.Parse(req.FormValue("Game.Length"))

	playerIds := req.Form["Game.Players.Ids"]
	game.Players = make([]GameFormPlayer, len(playerIds))
	for i, playerId := range playerIds {
		game.Players[i].PlayerId.Parse(playerId)
		game.Players[i].Level.Parse(req.FormValue("Game.Players." + playerId + ".Level"))
		game.Players[i].EffectiveLevel.Parse(req.FormValue("Game.Players." + playerId + ".EffectiveLevel"))

		if req.FormValue("Game.Players."+playerId+".Winner") != "" {
			game.Players[i].Winner = true
		}
	}

	args := make([]interface{}, 0, 1+len(game.Players)*4)
	var sql bytes.Buffer
	sql.WriteString(`
		    with g as (
		      insert into game(date, length) values($1, $2) returning game_id
		    ), p as (
		    insert into game_player(game_id, player_id, level, effective_level, winner)
		    select game_id, player_id, level, effective_level, winner
		    from g
		      cross join (values`)
	args = append(args, game.Date.Value)
	args = append(args, game.Length.Value)
	for i, p := range game.Players {
		fmt.Fprintf(&sql, "($%d, $%d, $%d, $%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4)
		args = append(args, p.PlayerId.Value)
		args = append(args, p.Level.Value)
		args = append(args, p.EffectiveLevel.Value)
		args = append(args, p.Winner)
		if i < (len(game.Players) - 1) {
			sql.WriteString(", ")
		}
	}
	sql.WriteString(`) as gp(player_id, level, effective_level, winner)`)
	sql.WriteString(`
		    )
		    select game_id
		    from g
		  `)

	if _, err := pool.SelectValue(sql.String(), args...); err == nil {
		http.Redirect(w, req, "/standings", http.StatusSeeOther)
	} else {
		fmt.Fprintf(os.Stderr, "createGame: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}

func deleteGame(w http.ResponseWriter, req *http.Request) {
	var err error
	var gameId int64
	if gameId, err = strconv.ParseInt(req.FormValue("id"), 10, 32); err != nil {
		http.NotFound(w, req)
		return
	}

	if _, err := pool.Execute("deleteGame", int32(gameId)); err != nil {
		fmt.Fprintf(os.Stderr, "deleteGame: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/games", http.StatusSeeOther)
}

func sortLink(label, column, defaultSortDir, sortCol, sortDir string) string {
	var newSortDir string
	if column == sortCol {
		if sortDir == "asc" {
			newSortDir = "desc"
		} else {
			newSortDir = "asc"
		}
	} else {
		newSortDir = defaultSortDir
	}
	return fmt.Sprintf(`<a href="?sortCol=%s&sortDir=%s">%s</a>`, column, newSortDir, label)
}

func getStandings(w http.ResponseWriter, req *http.Request) {
	sortCol := req.FormValue("sortCol")
	switch sortCol {
	case "name", "num_games", "num_wins", "num_points", "rating":
	default:
		sortCol = "rating"
	}

	sortDir := req.FormValue("sortDir")
	switch sortDir {
	case "asc", "desc":
	default:
		sortDir = "desc"
	}

	rows, err := pool.SelectRows(fmt.Sprintf(`
    select player_id, name, num_games, num_wins, num_points, round(rating, 3)::varchar as rating
    from player_summary
    order by %s %s, name asc
  `, sortCol, sortDir))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	RenderStandings(w, rows, sortCol, sortDir)
}
