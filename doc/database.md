

## user

* **id** serial
* **role** varchar
* **name** varchar
* **hash** varchar
* **created** bigint

## schema

* **id** serial
* **name** varchar(64)
* **description** text
* **complete** boolean
* **size_x** smallint
* **size_y** smallint
* **size_z** smallint
* **created** bigint
* **part_length** smallint
* **total_size** int
* **total_parts** int

## user_schema_star

* **user_id** serial (user.id)
* **schema_id** serial (schema.id)

## schmemapart

* **id** serial
* **schema_id** serial (schema.id)
* **offset_x** smallint
* **offset_y** smallint
* **offset_z** smallint
* **data** blob
* **metadata** blob

## schemamod

* **id** serial
* **schema_id** serial (schema.id)
* **mod_name** varchar(32)

## schematag

* **id** serial
* **schema_id** serial (schema.id)
* **tag_name** varchar(32)
