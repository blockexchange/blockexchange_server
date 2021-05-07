
alter table schema rename column size_x to size_x_plus;
alter table schema rename column size_y to size_y_plus;
alter table schema rename column size_z to size_z_plus;

alter table schema add column size_x_minus smallint default 0 not null;
alter table schema add column size_y_minus smallint default 0 not null;
alter table schema add column size_z_minus smallint default 0 not null;
