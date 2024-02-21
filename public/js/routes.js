import Import from './components/pages/Import.js';
import Login from './components/pages/Login.js';
import Mod from './components/pages/Mod.js';
import Profile from './components/pages/Profile.js';
import Register from './components/pages/Register.js';
import Search from './components/pages/Search.js';
import Start from './components/pages/Start.js';
import Tags from './components/pages/Tags.js';
import Users from './components/pages/Users.js';
import UserDetail from './components/pages/UserDetail.js';
import SchemaDetail from './components/pages/SchemaDetail.js';
import UserSchemas from './components/pages/UserSchemas.js';
import UserCollections from './components/pages/UserCollections.js';
import UserCollection from './components/pages/UserCollection.js';

export default [{
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
	path: "/users", component: Users
},{
	path: "/search", component: Search
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
