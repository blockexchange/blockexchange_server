

## user

* **id** bigint
* **name** varchar

## schema

* **id** bigint
* **description** text
* **complete** boolean
* **size_x** int
* **size_y** int
* **size_z** int
* **created** datetime
* **part_length** int
* **total_size** int
* **total_parts** int

## schmemapart

* **id** bigint
* **schema_id** bigint (schema.id)
* **offset_x** int
* **offset_y** int
* **offset_z** int
* **data** blob

## schemamod

* **id** bigint
* **schema_id** bigint (schema.id)
* **mod_name** varchar

## schematag

* **id** bigint
* **schema_id** bigint (schema.id)
* **tag_name** varchar
