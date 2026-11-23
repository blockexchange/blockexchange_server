// JSON types mirroring the Go structs in types/ and api/

// "LOCAL" or the name of the oauth provider
export type UserType = string;
export type UserRole = "DEFAULT" | "ADMIN";
export type Permission = "UPLOAD" | "OVERWRITE" | "MANAGEMENT" | "MEDIA" | "ADMIN";

export interface Claims {
	username: string;
	user_uid: string;
	type: UserType;
	permissions: Permission[];
}

export interface OauthLogin {
	github: string;
	discord: string;
	mesehub: string;
	cdb: string;
	codeberg: string;
}

export interface Stats {
	schema_count: number;
	schemapart_count: number;
	user_count: number;
	total_size: number;
}

export interface Info {
	api_version_major: number;
	api_version_minor: number;
	name: string;
	owner: string;
	base_url: string;
	oauth_login: OauthLogin;
	stats: Stats;
}

export interface ErrorResponse {
	message: string;
	status: number;
}

export interface AccessToken {
	uid: string;
	name: string;
	token: string;
	user_uid: string;
	created: number;
	expires: number;
	usecount: number;
}

export interface Collection {
	uid?: string;
	user_uid: string;
	name: string;
	description?: string;
}

export interface Tag {
	uid: string;
	name: string;
	description: string;
	restricted: boolean;
}

// editable fields of a tag
export type TagRequest = Omit<Tag, "uid">;

export interface User {
	uid: string;
	created: number;
	name: string;
	type: UserType;
	role: UserRole;
	external_id: string | null;
	avatar_url: string;
}

export interface UserSearch {
	name?: string;
	limit?: number;
	offset?: number;
}

export interface ChangePassword {
	old_password?: string;
	new_password: string;
}

export interface RegisterRequest {
	name: string;
	password: string;
	captcha_id: string;
	captcha_answer: string;
}

export interface RegisterResponse {
	success: boolean;
	error_password_too_short: boolean;
	error_captcha: boolean;
	error_invalid_username: boolean;
	error_username_taken: boolean;
}

export interface Schema {
	uid: string;
	type: number;
	created: number;
	mtime: number;
	user_uid: string;
	collection_uid: string | null;
	name: string;
	description: string;
	short_description: string;
	cdb_collection: string;
	complete: boolean;
	size_x: number;
	size_y: number;
	size_z: number;
	total_size: number;
	total_parts: number;
	downloads: number;
	views: number;
	license: string;
	stars: number;
}

export interface SchemaUpdateError {
	name_taken: boolean;
	name_invalid: boolean;
}

export interface SchemaSearchRequest {
	user_uid?: string | null;
	type?: number | null;
	schema_uid?: string | null;
	tag_uid?: string | null;
	tag_name?: string | null;
	schema_name?: string | null;
	mod_name?: string | null;
	user_name?: string | null;
	keywords?: string | null;
	complete?: boolean | null;
	collection_uid?: string | null;
	collection_name?: string | null;
	with_collection?: boolean | null;
	order_direction?: string | null;
	order_column?: string | null;
	from_mtime?: number | null;
	until_mtime?: number | null;
	limit?: number | null;
	offset?: number | null;
}

export interface SchemaSearchResponse {
	schema: Schema;
	collection_name: string | null;
	username: string;
	tags: string[];
	mods: string[];
}

export interface SchemaModCount {
	mod_name: string;
	count: number;
}

export type SchemaModDependencyType = "mod" | "game";

export interface SchemaModDependency {
	uid: string;
	schema_uid: string;
	mod_type: SchemaModDependencyType;
	git_url: string;
	// defaults to "latest" if empty
	version: string;
}

// editable fields of a schema mod dependency
export type SchemaModDependencyRequest = Omit<SchemaModDependency, "uid" | "schema_uid">;

export interface Nodedefinition {
	name: string;
	mod_name: string;
	definition: string;
}

export interface MediafileInfo {
	name: string;
	mod_name: string;
	size: number;
}

export interface LuantiMediaConnection {
	id: string;
	hostname: string;
	port: number;
	username: string;
	last_updated_timestamp: number;
	last_updated_success: boolean;
	last_updated_message: string;
}

export interface LuantiMediaConnectionRequest {
	hostname: string;
	port: number;
	username: string;
	// empty on update keeps the stored password
	password: string;
}

// media file including the base64 encoded data
export interface Mediafile {
	name: string;
	mod_name: string;
	data: string;
}

// schematic part, offsets are in nodes relative to the schematic origin (multiples of 16)
export interface SchemaPart {
	schema_uid: string;
	offset_x: number;
	offset_y: number;
	offset_z: number;
	mtime: number;
	// base64 (unpadded) of the zlib compressed node ids, param1 and param2
	data: string;
	// base64 (unpadded) of the zlib compressed json metadata (node_mapping, size, ...)
	metadata: string;
}
