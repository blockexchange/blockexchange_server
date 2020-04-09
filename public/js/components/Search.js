
import SearchBar from './SearchBar.js';
import SchemaList from './SchemaList.js';

import store from '../store/search.js';

import debounce from '../util/debounce.js';

import { find_recent, find_by_keywords } from '../api/searchschema.js';

const debounced_search = debounce(() => {
  if (store.keywords && store.keywords.length > 0)
    find_by_keywords(store.keywords).then(l => store.result = l);
  else
    find_recent(20).then(l => store.result = l);
}, 500);

find_recent(20).then(l => store.result = l);

function changeKeywords(k){
  store.keywords = k;
  debounced_search();
}

export default {
  view(){
    return("div", [
      m("div", m(SearchBar, { keywords: store.keywords, onChange: changeKeywords })),
      m("div", m(SchemaList, { list: store.result }))
    ]);
  }
};
