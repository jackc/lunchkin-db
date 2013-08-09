alter table game add column length smallint;
comment on column game.length is 'Number of rounds played';

---- CREATE above / DROP below ----

alter table game drop column length;
