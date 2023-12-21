-- additional search index
create index schema_mtime on schema(mtime);

-- stars column
alter table schema add stars int not null default 0;
create index schema_stars on schema(stars);
update schema s set stars = (select count(*) from user_schema_star where schema_id = s.id);