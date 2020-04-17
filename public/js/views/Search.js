
import SearchBar from '../components/SearchBar.js';
import SchemaList from '../components/SchemaList.js';

import store from '../store/search.js';

import debounce from '../util/debounce.js';

import { find_recent, find_by_keywords } from '../api/searchschema.js';
import { remove } from '../api/schema.js';


export default class {
  constructor(){
    this.debounced_search = debounce(this.search , 500);
    this.state = {
      result: []
    };
    this.search();
  }

  search() {
    if (store.keywords && store.keywords.length > 0)
      find_by_keywords(store.keywords)
      .then(l => this.state.result = l);
    else
      find_recent(20).then(l => this.state.result = l);
  }

  removeItem(schema) {
    remove(schema)
    .then(() => this.search());
  }

  changeKeywords(k){
    store.keywords = k;
    this.debounced_search();
  }

  view(){
    return("div", [
      m("div", m(SearchBar, {
        keywords: store.keywords,
        onChange: k => this.changeKeywords(k)
      })),
      m("div", m(SchemaList, {
        list: this.state.result,
        removeItem: schema => this.removeItem(schema)
      }))
    ]);
  }
}
