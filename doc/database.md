

## user

* **id** bigint
* **name** varchar

## schema

* **id** bigint
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

## schmemapart

* **id** bigint
* **schema_id** bigint (schema.id)
* **offset_x** smallint
* **offset_y** smallint
* **offset_z** smallint
* **data** blob
* **metadata** blob

## schemamod

* **id** bigint
* **schema_id** bigint (schema.id)
* **mod_name** varchar(32)

## schematag

* **id** bigint
* **schema_id** bigint (schema.id)
* **tag_name** varchar(32)
