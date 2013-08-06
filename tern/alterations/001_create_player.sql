create table player(
  player_id serial primary key,
  name varchar not null unique
);

---- CREATE above / DROP below ----

drop table player;
