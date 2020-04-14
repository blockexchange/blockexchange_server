import Start from './views/Start.js';
import Search from './views/Search.js';
import UserSchemaList from './views/UserSchemaList.js';
import SchemaDetail from './views/SchemaDetail.js';

export default {
  "/": Start,
  "/search": Search,
  "/schema/:username": UserSchemaList,
  "/schema/:username/:schemaname": SchemaDetail
};
