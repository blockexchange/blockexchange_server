import Start from './views/Start.js';
import Search from './views/Search.js';
import UserSchemaList from './views/UserSchemaList.js';
import SchemaDetail from './views/SchemaDetail.js';
import Login from './views/Login.js';
import Register from './views/Register.js';

export default {
  "/": Start,
  "/login": Login,
  "/register": Register,
  "/search": Search,
  "/schema/:username": UserSchemaList,
  "/schema/:username/:schemaname": SchemaDetail
};
