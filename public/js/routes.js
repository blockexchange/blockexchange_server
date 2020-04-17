import Start from './views/Start.js';
import Search from './views/Search.js';
import UserSchemaList from './views/UserSchemaList.js';
import SchemaDetail from './views/SchemaDetail.js';
import UserList from './views/UserList.js';
import Login from './views/Login.js';
import Register from './views/Register.js';
import Mod from './views/Mod.js';

export default {
  "/": Start,
  "/users": UserList,
  "/login": Login,
  "/register": Register,
  "/search": Search,
  "/mod": Mod,
  "/schema/:username": UserSchemaList,
  "/schema/:username/:schemaname": SchemaDetail
};
