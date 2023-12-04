import Login from './components/pages/Login.js';
import Profile from './components/pages/Profile.js';
import Start from './components/pages/Start.js';

export default [{
	path: "/", component: Start
},{
	path: "/profile", component: Profile
},{
	path: "/login", component: Login
}];
