package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/JackC/pgx"
	qv "github.com/JackC/quo_vadis"
	"github.com/kylelemons/go-gypsy/yaml"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var pool *pgx.ConnectionPool

var config struct {
	assetPath     string
	configPath    string
	listenAddress string
	listenPort    string
}

func init() {
	var err error
	var yf *yaml.File

	flag.StringVar(&config.listenAddress, "address", "127.0.0.1", "address to listen on")
	flag.StringVar(&config.listenPort, "port", "8080", "port to listen on")
	flag.StringVar(&config.assetPath, "assetpath", "assets", "path to assets")
	flag.StringVar(&config.configPath, "config", "config.yml", "path to config file")
	flag.Parse()

	givenCliArgs := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		givenCliArgs[f.Name] = true
	})

	if config.configPath, err = filepath.Abs(config.configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid config path: %v\n", err)
		os.Exit(1)
	}

	if config.assetPath, err = filepath.Abs(config.assetPath); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid asset path: %v\n", err)
		os.Exit(1)
	}

	if yf, err = yaml.ReadFile(config.configPath); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if !givenCliArgs["address"] {
		if address, err := yf.Get("address"); err == nil {
			config.listenAddress = address
		}
	}

	if !givenCliArgs["assetpath"] {
		if assetpath, err := yf.Get("assetpath"); err == nil {
			config.assetPath = assetpath
		}
	}

	if !givenCliArgs["port"] {
		if port, err := yf.Get("port"); err == nil {
			config.listenPort = port
		}
	}

	var connectionParameters pgx.ConnectionParameters
	if connectionParameters, err = extractConnectionOptions(yf); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err = migrate(connectionParameters); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	poolOptions := pgx.ConnectionPoolOptions{MaxConnections: 5, AfterConnect: afterConnect}
	pool, err = pgx.NewConnectionPool(connectionParameters, poolOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create database connection pool: %v\n", err)
		os.Exit(1)
	}
}

func extractConnectionOptions(config *yaml.File) (connectionOptions pgx.ConnectionParameters, err error) {
	connectionOptions.Host, _ = config.Get("database.host")
	connectionOptions.Socket, _ = config.Get("database.socket")
	if connectionOptions.Host == "" && connectionOptions.Socket == "" {
		err = errors.New("Config must contain database.host or database.socket but it does not")
		return
	}
	port, _ := config.GetInt("database.port")
	connectionOptions.Port = uint16(port)
	if connectionOptions.Database, err = config.Get("database.database"); err != nil {
		err = errors.New("Config must contain database.database but it does not")
		return
	}
	if connectionOptions.User, err = config.Get("database.user"); err != nil {
		err = errors.New("Config must contain database.user but it does not")
		return
	}
	connectionOptions.Password, _ = config.Get("database.password")
	return
}

type player struct {
	player_id int32
	name      string
}

type GamePlayer struct {
	nulls          int64
	GameId         int32
	PlayerId       int32
	Name           string
	Level          int16
	EffectiveLevel int16
	Winner         bool
}

type Game struct {
	nulls   int64
	Players []GamePlayer
	Date    time.Time
	GameId  int32
	Length  int16
}

func SelectAllGamesWithDetails() (games []Game, err error) {
	gameIdToIndexMap := make(map[int32]int)
	games = make([]Game, 0, 8)
	err = pool.SelectFunc("select game_id, date, length from game", func(r *pgx.DataRowReader) error {
		var g Game
		g.Players = make([]GamePlayer, 0, 4)
		g.GameId = r.ReadValue().(int32)
		g.Date = r.ReadValue().(time.Time)
		g.Length = r.ReadValue().(int16)
		gameIdToIndexMap[g.GameId] = len(games)
		games = append(games, g)
		return nil
	})
	if err != nil {
		return
	}

	err = pool.SelectFunc("select game_id, player_id, name, level, effective_level, winner from game_player join player using(player_id)", func(r *pgx.DataRowReader) error {
		var gp GamePlayer
		gp.GameId = r.ReadValue().(int32)
		gp.PlayerId = r.ReadValue().(int32)
		gp.Name = r.ReadValue().(string)
		gp.Level = r.ReadValue().(int16)
		gp.EffectiveLevel = r.ReadValue().(int16)
		gp.Winner = r.ReadValue().(bool)

		gameIdx := gameIdToIndexMap[gp.GameId]
		games[gameIdx].Players = append(games[gameIdx].Players, gp)
		return nil
	})
	if err != nil {
		return
	}

	return
}

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Connection) (err error) {
	err = conn.Prepare("getPlayers", `
    select player_id, name
    from player
    order by name
  `)
	if err != nil {
		return
	}

	err = conn.Prepare("createPlayer", `
    insert into player(name) values($1) returning player_id
  `)
	if err != nil {
		return
	}

	err = conn.Prepare("deletePlayer", `
    delete from player where player_id=$1
  `)
	if err != nil {
		return
	}

	err = conn.Prepare("getGames", `
    select coalesce(array_to_json(array_agg(row_to_json(g))), '[]'::json)
    from (
      select game_id, date, length,
        (
          select coalesce(array_to_json(array_agg(row_to_json(t))), '[]'::json)
          from (
            select player_id, name, level, effective_level, winner
            from game_player
              join player using(player_id)
            where game.game_id=game_player.game_id
          ) t
        ) players
      from game
    ) g
  `)
	if err != nil {
		return
	}

	err = conn.Prepare("deleteGame", `
    delete from game where game_id=$1
  `)
	if err != nil {
		return
	}

	return
}

func NoDirListing(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func main() {
	router := qv.NewRouter()
	router.Get("/standings", http.HandlerFunc(getStandings))
	router.Get("/players", http.HandlerFunc(getPlayers))
	router.Post("/players", http.HandlerFunc(createPlayer))
	router.Post("players/:id/delete", http.HandlerFunc(deletePlayer))
	router.Get("/games", http.HandlerFunc(getGames))
	router.Post("games/:id/delete", http.HandlerFunc(deleteGame))
	http.Handle("/", router)
	http.Handle("/assets/", NoDirListing(http.StripPrefix("/assets/", http.FileServer(http.Dir(config.assetPath)))))

	listenAt := fmt.Sprintf("%s:%s", config.listenAddress, config.listenPort)
	fmt.Printf("Starting to listen on: %s\n", listenAt)

	if err := http.ListenAndServe(listenAt, nil); err != nil {
		os.Stderr.WriteString("Could not start web server!\n")
		os.Exit(1)
	}
}
