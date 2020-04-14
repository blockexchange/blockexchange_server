import Start from './components/Start.js';
import Search from './components/Search.js';
import UserSchemaList from './components/UserSchemaList.js';
import SchemaDetail from './components/SchemaDetail.js';

export default {
  "/": Start,
  "/search": Search,
  "/schema/:username": UserSchemaList,
  "/schema/:username/:schemaname": SchemaDetail
};
