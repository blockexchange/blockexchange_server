
import SearchBar from './SearchBar.js';
import SchemaList from './SchemaList.js';

import debounce from '../util/debounce.js';

import { find_recent, find_by_keywords } from '../api/searchschema.js';

var keywords = "";
var list = [];

const debounced_search = debounce(() => {
  if (keywords && keywords.length > 0)
    find_by_keywords(keywords).then(l => list = l);
  else
    find_recent(20).then(l => list = l);
}, 500);

find_recent(20).then(l => list = l);

function changeKeywords(k){
  keywords = k;
  debounced_search();
}

export default {
  view(){
    return("div", [
      m("div", m(SearchBar, { keywords: keywords, onChange: changeKeywords })),
      m("div", m(SchemaList, { list: list }))
    ]);
  }
};
