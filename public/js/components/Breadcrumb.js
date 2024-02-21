
export default {
    props: ["items"],
    methods: {
        get_icon_class: function(item) {
            const cl = {};
            if (item.icon){
                cl.fa = true;
                cl["fa-" + item.icon] = true;
            }
            return cl;
        }
    },
    template: /*html*/`
    <nav>
        <ol class="breadcrumb text-bg-secondary" style="padding: 5px; border-radius: 5px;">
            <li class="breadcrumb-item" v-for="item in items">
                <router-link class="link-light" :to="item.link" v-if="item.link">
                    <i v-bind:class="get_icon_class(item)" v-if="item.icon"></i>
                    {{item.name}}
                </router-link>
                <span v-else>
                    <i v-bind:class="get_icon_class(item)" v-if="item.icon"></i>
                    {{item.name}}
                </span>
            </li>
        </ol>
    </nav>
    `
};

export const START = { name: "Start", icon: "home", link: "/" };
export const LOGIN = { name: "Login", icon: "sign-in", link: "/login" };
export const PROFILE = { name: "Profile", icon: "user", link: "/profile" };
export const REGISTER = { name: "Register", icon: "user-plus", link: "/register" };
export const MOD = { name: "Mod/Installation", icon: "download", link: "/mod" };
export const USERS = { name: "Users", icon: "users", link: "/users" };
export const SEARCH = { name: "Search", icon: "search", link: "/search" };
export const MY_SCHEMAS = { name: "My schemas", icon: "home", link: "/schemas" };
export const SCHEMA_IMPORT = { name: "Schematic import", icon: "upload", link: "/import" };
export const TAGS = { name: "Tags", icon: "tags", link: "/tags" };

export const USER_DETAIL = username => {
    return { name: `Profile for ${username}`, icon: "user", link: `/user/${username}` };
};

export const USER_SCHEMAS = username => {
    return { name: `Schematics for ${username}`, icon: "user", link: `/schema/${username}` };
};

export const USER_COLLECTIONS = username => {
    return { name: `Collections for ${username}`, icon: "object-group", link: `/collections/${username}` };
};

export const USER_COLLECTION = (username, collection_name) => {
    return { name: `Collection ${collection_name}`, icon: "object-group", link: `/collections/${username}/${collection_name}` };
};

export const SCHEMA_DETAIL = (username, name) => {
    return { name: `Schematic ${name}`, icon: "bookmark", link: `/schema/${username}/${name}` };
};