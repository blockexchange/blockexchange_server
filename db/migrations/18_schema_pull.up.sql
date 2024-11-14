
create table schema_pull(
    schema_uid uuid primary key not null references public.schema(uid) on delete cascade,
    pos_x int not null,
    pos_y int not null,
    pos_z int not null,
    interval int not null,
    next_run bigint not null default 0,
    hostname varchar(128) not null,
    port int not null
);

create table schema_pull_client(
    uid uuid primary key not null default gen_random_uuid(),
    schema_pull_uid uuid not null references public.schema_pull(schema_uid) on delete cascade,
    enabled boolean not null,
    username varchar(64) not null,
    password varchar(64) not null,
    last_run bigint not null default 0,
    last_error boolean not null,
    last_message varchar(256) not null
);