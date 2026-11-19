-- luanti media connection table
create table luanti_media_connection(
    id uuid primary key not null default gen_random_uuid(),
    hostname varchar not null,
    port int not null,
    username varchar not null,
    password varchar not null,
    last_updated_timestamp bigint not null default 0,
    last_updated_success boolean not null default false,
    last_updated_message varchar not null default ''
);
