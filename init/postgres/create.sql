create sequence if not exists token_id_seq;

create table if not exists t_token
(
  id bigint default nextval('token_id_seq'::regclass) constraint token_pkey primary key,
  token text not null unique,
  status smallint default 0,
  origin_lat numeric not null,
  origin_long numeric not null,
  destinations text not null,
  create_time timestamp default now() not null
);

create table if not exists t_route
(
  token_id bigint not null constraint route_pkey primary key,
  result text,
  create_time timestamp default now() not null
);