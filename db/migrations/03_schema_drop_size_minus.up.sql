alter table schema drop column size_x_minus;
alter table schema drop column size_y_minus;
alter table schema drop column size_z_minus;

alter table schema rename column size_x_plus to size_x;
alter table schema rename column size_y_plus to size_y;
alter table schema rename column size_z_plus to size_z;
