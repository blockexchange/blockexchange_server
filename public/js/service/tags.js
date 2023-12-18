import { get_tags as api_get_tags } from "../api/tags.js";

const store = Vue.reactive({
    list: []
});

export const fetch_tags = () => api_get_tags().then(l => store.list = l);
export const get_tags = () => store.list;