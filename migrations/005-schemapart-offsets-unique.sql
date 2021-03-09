
drop index schemapart_coords;
create unique index schemapart_coords on schemapart(schema_id, offset_x, offset_y, offset_z);

alter table schemapart
add constraint schemapart_unique_coords
unique using index schemapart_coords;
