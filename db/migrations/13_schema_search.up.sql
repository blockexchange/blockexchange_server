alter table schema drop column search_tokens;
CREATE INDEX schema_search_idx ON schema USING GIN (to_tsvector('english', description || ' ' || name));
