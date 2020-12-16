import StartPage from './components/pages/StartPage.js';
import LoginPage from './components/pages/LoginPage.js';
import SearchPage from './components/pages/SearchPage.js';
import RegisterPage from './components/pages/RegisterPage.js';
import OauthPage from './components/pages/OauthPage.js';
import UsersPage from './components/pages/UsersPage.js';
import SchemaPage from './components/pages/SchemaPage.js';
import ProfilePage from './components/pages/ProfilePage.js';

export default [{
  path: "/", component: StartPage
},{
  path: "/login", component: LoginPage
},{
  path: "/search", component: SearchPage
},{
	path: "/register", component: RegisterPage
},{
	path: "/oauth/:token", component: OauthPage
},{
	path: "/users", component: UsersPage
},{
	path: "/profile", component: ProfilePage
},{
	path: "/schema/:user_name/:schema_name", name: "schemapage", component: SchemaPage
}];
