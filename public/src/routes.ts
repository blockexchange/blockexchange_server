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
import Nodedefinitions from './components/pages/Nodedefinitions.vue';
import Mediafiles from './components/pages/Mediafiles.vue';
import LuantiConnections from './components/pages/LuantiConnections.vue';
import NotFound from './components/pages/NotFound.vue';

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
	// lazy loaded, keeps three.js out of the main bundle
	path: "/schema/:username/:name/3d", component: () => import('./components/pages/SchemaView3d.vue'), props: true
},{
	// luanti wasm client, runs in the cross-origin isolated document under /wasm/ (see main.ts)
	path: "/schema/:username/:name/online", component: () => import('./components/pages/SchemaEditOnline.vue'), props: true,
	meta: { isolated: true }
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
},{
	path: "/nodedefinitions", component: Nodedefinitions
},{
	path: "/mediafiles", component: Mediafiles
},{
	path: "/luanti-connections", component: LuantiConnections
},{
	path: "/:pathMatch(.*)*", component: NotFound
}];

export default routes;
