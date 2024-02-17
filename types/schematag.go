package types

type SchemaTag struct {
	UID       string `json:"uid" ksql:"uid"`
	TagUID    string `json:"tag_uid" ksql:"tag_uid"`
	SchemaUID string `json:"schema_uid" ksql:"schema_uid"`
}
