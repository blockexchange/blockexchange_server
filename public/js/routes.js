import StartPage from './components/start-page.js';
import LoginPage from './components/login-page.js';
import SearchPage from './components/search-page.js';

export default [{
  path: "/", component: StartPage
},{
  path: "/login", component: LoginPage
},{
  path: "/search", component: SearchPage
}];
