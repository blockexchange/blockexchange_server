
create table user_schema_star (
  user_id serial references public.user(id) on delete cascade,
  schema_id serial references schema(id) on delete cascade,
  primary key (user_id, schema_id)
)
