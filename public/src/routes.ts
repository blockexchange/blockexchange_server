import type { RouteRecordRaw } from 'vue-router';

import Import from './components/pages/Import.vue';
import Login from './components/pages/Login.vue';
import Mod from './components/pages/Mod.vue';
import Profile from './components/pages/Profile.vue';
import Register from './components/pages/Register.vue';
import Search from './components/pages/Search.vue';
import Start from './components/pages/Start.vue';
import Tags from './components/pages/Tags.vue';
import Users from './components/pages/Users.vue';
import UserDetail from './components/pages/UserDetail.vue';
import SchemaDetail from './components/pages/SchemaDetail.vue';
import UserSchemas from './components/pages/UserSchemas.vue';
import UserCollections from './components/pages/UserCollections.vue';
import UserCollection from './components/pages/UserCollection.vue';
import Mods from './components/pages/Mods.vue';

const routes: RouteRecordRaw[] = [{
	path: "/", component: Start
},{
	path: "/profile", component: Profile
},{
	path: "/user/:username", component: UserDetail, props: true
},{
	path: "/login", component: Login
},{
	path: "/register", component: Register
},{
	path: "/mod", component: Mod
},{
	path: "/mods", component: Mods
},{
	path: "/users", component: Users
},{
	path: "/search", component: Search,
},{
	path: "/schema/:username/:name", component: SchemaDetail, props: true
},{
	path: "/schema/:username", component: UserSchemas, props: true
},{
	path: "/collections/:username", component: UserCollections, props: true
},{
	path: "/collections/:username/:collection_name", component: UserCollection, props: true
},{
	path: "/import", component: Import
},{
	path: "/tags", component: Tags
}];

export default routes;
