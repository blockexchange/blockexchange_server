alter table tag add restricted bool not null default false;
insert into tag(name, description, restricted) values('featured', 'Featured schematics', true);