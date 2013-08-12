package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func getPlayers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := pool.SelectValueTo(w, "getPlayers"); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func playerPath(player_id int32) string {
	return fmt.Sprintf("/api/v1/players/%d", player_id)
}

func createPlayer(w http.ResponseWriter, req *http.Request) {
	var player struct{ Name string }
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&player); err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, "Error decoding request: %v", err)
		return
	}

	if player.Name == "" {
		w.WriteHeader(422)
		fmt.Fprintln(w, `Request must include the attribute "name"`)
		return
	}

	if player_id, err := pool.SelectValue("createPlayer", player.Name); err == nil {
		w.Header().Add("Location", playerPath(player_id.(int32)))
		w.WriteHeader(http.StatusCreated)
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
		fmt.Fprintf(os.Stderr, "deletePlayer: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getGames(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := pool.SelectValueTo(w, "getGames"); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func gamePath(game_id int32) string {
	return fmt.Sprintf("/api/v1/games/%d", game_id)
}

func createGame(w http.ResponseWriter, req *http.Request) {
	type Player struct {
		Player_Id       int32
		Level           int16
		Effective_Level int16
		Winner          bool
	}
	var game struct {
		Date    string
		Length  int16
		Players []Player
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&game); err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, "Error decoding request: %v", err)
		return
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
	args = append(args, game.Date)
	args = append(args, game.Length)
	for i, p := range game.Players {
		fmt.Fprintf(&sql, "($%d, $%d, $%d, $%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4)
		args = append(args, p.Player_Id)
		args = append(args, p.Level)
		args = append(args, p.Effective_Level)
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

	if game_id, err := pool.SelectValue(sql.String(), args...); err == nil {
		w.Header().Add("Location", gamePath(game_id.(int32)))
		w.WriteHeader(http.StatusCreated)
	} else {
		fmt.Fprintf(os.Stderr, "createGame: %v", err)
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
		fmt.Fprintf(os.Stderr, "deleteGame: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getStandings(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := pool.SelectValueTo(w, "getStandings"); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
