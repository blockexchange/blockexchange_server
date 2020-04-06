
import SchemaList from './SchemaList.js';
import { find_recent } from '../api/searchschema.js';


var list = [];
find_recent(50).then(l => list = l);

export default {
  view: () => SchemaList(list)
};
