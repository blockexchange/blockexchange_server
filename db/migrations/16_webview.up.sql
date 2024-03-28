
create table mod(
    name varchar(256) primary key not null, -- name of the mod
    source varchar(256) not null, -- source repository/url
    media_license varchar(64) not null -- license of media assets ("name" field from https://content.minetest.net/api/licenses/)
);

create table nodedefinition(
    name varchar(256) primary key not null, -- name of the node ("mynode:mymod")
    mod_name varchar(256) not null references mod(name) on delete cascade, -- mod-reference
    definition text not null -- node-definition in json
);

create table mediafile(
    name varchar(256) primary key not null, -- name of the media file ("default_stone.png", "moreblocks_slope.obj")
    mod_name varchar(256) not null references mod(name) on delete cascade, -- mod-reference
    data bytea not null -- content
);