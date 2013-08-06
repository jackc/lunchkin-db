create table game_player(
  game_id int references game on delete cascade,
  player_id int references player,
  level smallint,
  effective_level smallint,
  winner boolean not null default false
);

---- CREATE above / DROP below ----

drop table game_player;
