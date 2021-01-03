

-- COLLECTION

create table collection (
	id serial primary key not null,
  user_id serial references public.user(id) on delete cascade,
  name varchar not null
);

create table collection_schema (
	collection_id bigint not null references collection(id) on delete cascade,
	schema_id bigint not null references schema(id) on delete cascade,
	primary key (collection_id, schema_id)
);
