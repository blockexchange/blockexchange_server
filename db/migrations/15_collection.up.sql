create table collection(
    uid uuid primary key not null default gen_random_uuid(),
    user_uid uuid not null references public.user(uid) on delete cascade,
    name varchar not null,
    description text not null
);

alter table schema add column collection_uid uuid;
