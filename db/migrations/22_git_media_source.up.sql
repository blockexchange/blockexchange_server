-- git repositories that textures and models are imported from
create table git_media_source(
    id uuid primary key not null default gen_random_uuid(),
    git_url varchar(256) not null unique,
    -- contentdb package ("author/name"), used to look up the licenses, matched by the git url if empty
    cdb_package varchar(128) not null default '',
    -- license names from https://content.luanti.org/api/licenses/, fetched from the contentdb if empty
    media_license varchar(64) not null default '',
    license varchar(64) not null default '',
    last_updated_timestamp bigint not null default 0,
    last_updated_success boolean not null default false,
    last_updated_message varchar not null default ''
);

-- media files now belong to a git source (or none, for the existing files, the luanti connections and uploads)
alter table mediafile drop column mod_name;
alter table mediafile add column git_media_source_id uuid references git_media_source(id) on delete cascade;
create index mediafile_git_media_source_id on mediafile(git_media_source_id);
