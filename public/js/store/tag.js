import { get_all } from '../api/tag.js';

const store = Vue.observable({
    tags: []
});

get_all().then(tags => store.tags = tags);

export default store;