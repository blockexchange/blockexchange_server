export interface BreadcrumbItem {
	name: string;
	icon?: string;
	link?: string;
}

export const START: BreadcrumbItem = { name: "Start", icon: "home", link: "/" };
export const LOGIN: BreadcrumbItem = { name: "Login", icon: "sign-in", link: "/login" };
export const PROFILE: BreadcrumbItem = { name: "Profile", icon: "user", link: "/profile" };
export const REGISTER: BreadcrumbItem = { name: "Register", icon: "user-plus", link: "/register" };
export const MODS: BreadcrumbItem = { name: "Mods", icon: "box-archive", link: "/mods" };
export const MOD: BreadcrumbItem = { name: "Mod/Installation", icon: "download", link: "/mod" };
export const USERS: BreadcrumbItem = { name: "Users", icon: "users", link: "/users" };
export const SEARCH: BreadcrumbItem = { name: "Search", icon: "search", link: "/search" };
export const MY_SCHEMAS: BreadcrumbItem = { name: "My schemas", icon: "home", link: "/schemas" };
export const SCHEMA_IMPORT: BreadcrumbItem = { name: "Schematic import", icon: "upload", link: "/import" };
export const TAGS: BreadcrumbItem = { name: "Tags", icon: "tags", link: "/tags" };

export const USER_DETAIL = (username: string): BreadcrumbItem => {
	return { name: `Profile for ${username}`, icon: "user", link: `/user/${username}` };
};

export const USER_SCHEMAS = (username: string): BreadcrumbItem => {
	return { name: `Schematics for ${username}`, icon: "user", link: `/schema/${username}` };
};

export const USER_COLLECTIONS = (username: string): BreadcrumbItem => {
	return { name: `Collections for ${username}`, icon: "object-group", link: `/collections/${username}` };
};

export const USER_COLLECTION = (username: string, collection_name: string): BreadcrumbItem => {
	return { name: `Collection ${collection_name}`, icon: "object-group", link: `/collections/${username}/${collection_name}` };
};

export const SCHEMA_DETAIL = (username: string, name: string): BreadcrumbItem => {
	return { name: `Schematic ${name}`, icon: "landmark", link: `/schema/${username}/${name}` };
};
