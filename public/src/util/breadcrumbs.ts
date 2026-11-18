export interface BreadcrumbItem {
	name: string;
	icon?: string;
	link?: string;
}

export const START: BreadcrumbItem = { name: "Home", icon: "home", link: "/" };
export const LOGIN: BreadcrumbItem = { name: "Log in", icon: "sign-in", link: "/login" };
export const PROFILE: BreadcrumbItem = { name: "Account", icon: "user", link: "/profile" };
export const REGISTER: BreadcrumbItem = { name: "Sign up", icon: "user-plus", link: "/register" };
export const MODS: BreadcrumbItem = { name: "Mods", icon: "box-archive", link: "/mods" };
export const MOD: BreadcrumbItem = { name: "Get the mod", icon: "download", link: "/mod" };
export const USERS: BreadcrumbItem = { name: "Users", icon: "users", link: "/users" };
export const SEARCH: BreadcrumbItem = { name: "Browse", icon: "search", link: "/search" };
export const SCHEMA_IMPORT: BreadcrumbItem = { name: "Import files", icon: "upload", link: "/import" };
export const TAGS: BreadcrumbItem = { name: "Tags", icon: "tags", link: "/tags" };

export const USER_DETAIL = (username: string): BreadcrumbItem => {
	return { name: username, icon: "user", link: `/user/${username}` };
};

export const USER_SCHEMAS = (username: string): BreadcrumbItem => {
	return { name: `Schematics by ${username}`, icon: "user", link: `/schema/${username}` };
};

export const USER_COLLECTIONS = (username: string): BreadcrumbItem => {
	return { name: `Collections by ${username}`, icon: "object-group", link: `/collections/${username}` };
};

export const USER_COLLECTION = (username: string, collection_name: string): BreadcrumbItem => {
	return { name: collection_name, icon: "object-group", link: `/collections/${username}/${collection_name}` };
};

export const SCHEMA_DETAIL = (username: string, name: string): BreadcrumbItem => {
	return { name: name, icon: "landmark", link: `/schema/${username}/${name}` };
};
