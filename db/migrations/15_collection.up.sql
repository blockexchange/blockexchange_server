-- collection table
create table collection(
    uid uuid primary key not null default gen_random_uuid(),
    user_uid uuid not null references public.user(uid) on delete cascade,
    name varchar not null,
    description text not null
);

-- optional foreign key
alter table schema add column collection_uid uuid references collection(uid) on delete set null;

-- search index
CREATE INDEX collection_search_idx ON collection USING GIN (to_tsvector('english', description || ' ' || name));
