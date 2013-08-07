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

---- CREATE above / DROP below ----

drop view player_summary;

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
