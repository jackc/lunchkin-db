package main

import (
	"errors"
	"net/http"
	"strings"
	// "bytes"
	"flag"
	"fmt"
	"github.com/jackc/pgx"
	qv "github.com/jackc/quo_vadis"
	"github.com/kylelemons/go-gypsy/yaml"
	// "io/ioutil"
	// "net/http"
	"os"
	"path/filepath"
)

var pool *pgx.ConnectionPool

type cliArgs struct {
	configPath string
}

func parseCliArgs() (args *cliArgs) {
	args = new(cliArgs)

	flag.StringVar(&args.configPath, "config", "config.yml", "path to config file")
	flag.Parse()

	return
}

func loadConfig(path string) (yf *yaml.File, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}

	yf, err = yaml.ReadFile(path)
	return
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

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Connection) (err error) {
	err = conn.Prepare("getPlayers", `
    select coalesce(array_to_json(array_agg(row_to_json(t))), '[]'::json)
    from (
      select player_id, name
      from player
      order by name
    ) t
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
      select game_id, date,
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
	args := parseCliArgs()

	var err error
	var yaml *yaml.File
	if yaml, err = loadConfig(args.configPath); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	var connectionParameters pgx.ConnectionParameters
	if connectionParameters, err = extractConnectionOptions(yaml); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	poolOptions := pgx.ConnectionPoolOptions{MaxConnections: 5, AfterConnect: afterConnect}
	pool, err = pgx.NewConnectionPool(connectionParameters, poolOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create database connection pool: %v\n", err)
		os.Exit(1)
	}

	router := qv.NewRouter()
	router.Get("/players", http.HandlerFunc(getPlayers))
	router.Post("/players", http.HandlerFunc(createPlayer))
	router.Delete("/players/:id", http.HandlerFunc(deletePlayer))
	router.Get("/games", http.HandlerFunc(getGames))
	router.Post("/games", http.HandlerFunc(createGame))
	router.Delete("/games/:id", http.HandlerFunc(deleteGame))
	http.Handle("/api/v1/", http.StripPrefix("/api/v1", router))
	http.Handle("/", NoDirListing(http.FileServer(http.Dir("./app/"))))

	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		os.Stderr.WriteString("Could not start web server!\n")
		os.Exit(1)
	}
}
