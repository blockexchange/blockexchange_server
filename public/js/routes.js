import Start from './components/Start.js';
import Search from './components/Search.js';
import GroupSchemaList from './components/GroupSchemaList.js';
import SchemaDetail from './components/schemadetail/SchemaDetail.js';
import SchemaEdit from './components/schemaedit/SchemaEdit.js';
import UserList from './components/UserList.js';
import Login from './components/login/Login.js';
import Register from './components/Register.js';
import Mod from './components/Mod.js';
import Import from './components/import/Import.js';

export default {
  "/": Start,
  "/users": UserList,
  "/login": Login,
  "/register": Register,
  "/search": Search,
  "/mod": Mod,
  "/import": Import,
  "/schema/:username": GroupSchemaList,
  "/schema/:username/:schemaname": SchemaDetail,
  "/schema/:username/:schemaname/edit": SchemaEdit
};
