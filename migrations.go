package main

import (
	"fmt"
	"github.com/JackC/pgx"
	mig "github.com/JackC/pgx/migrate"
)

func migrate(connectionParameters pgx.ConnectionParameters) (err error) {
	var conn *pgx.Connection
	conn, err = pgx.Connect(connectionParameters)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := conn.Close()
		if err == nil {
			err = closeErr
		}
	}()

	var m *mig.Migrator
	m, err = mig.NewMigrator(conn, "schema_version")
	if err != nil {
		return
	}

	m.OnStart = func(migration *mig.Migration) {
		fmt.Printf("Migrating %d: %s\n", migration.Sequence, migration.Name)
	}

	m.AppendMigration("Create player", `
		create table player(
		  player_id serial primary key,
		  name varchar not null unique
		);
	`)

	m.AppendMigration("Create game", `
		create table game(
		  game_id serial primary key,
		  date date not null
		);
	`)

	m.AppendMigration("Create game_player", `
		create table game_player(
		  game_id int references game on delete cascade,
		  player_id int references player,
		  level smallint,
		  effective_level smallint,
		  winner boolean not null default false,
		  primary key (game_id, player_id)
		);

		create index on game_player (player_id);
	`)

	m.AppendMigration("Create player_summary", `
		create view player_summary as
		  select player_id, name,
		  count(*) as num_games,
		  sum(winner::int) as num_wins,
		  sum(case when winner then ceil(num_players::float8/num_winners::float8) else 0 end)::bigint as num_points
		from (
		select *,
		  count(*) over (partition by game_id) as num_players,
		  sum(winner::int) over (partition by game_id) as num_winners
		from player
		  join game_player using(player_id)
		  join game using(game_id)
		) t
		group by player_id, name;
	`)

	m.AppendMigration("Add rating to player_summary", `
		drop view player_summary;

		create view player_summary as
		select player_id,
		  name,
		  num_games,
		  num_wins,
		  num_points::bigint,
		  round((num_points / num_games)::numeric, 3) as rating
		from (
		  select player_id,
		    name,
		    count(*) as num_games,
		    sum(winner::int) as num_wins,
		    sum(case when winner then ceil(num_players::float8/num_winners::float8) else 0 end) as num_points
		  from (
		  select *,
		    count(*) over (partition by game_id) as num_players,
		    sum(winner::int) over (partition by game_id) as num_winners
		  from player
		    join game_player using(player_id)
		    join game using(game_id)
		  ) t
		  group by player_id, name
		) t;
	`)

	m.AppendMigration("Add length to game", `
		alter table game add column length smallint;
		comment on column game.length is 'Number of rounds played';
	`)

	return m.Migrate()
}
