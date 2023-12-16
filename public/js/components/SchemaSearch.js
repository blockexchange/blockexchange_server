import { schema_search, schema_count } from "../api/schema.js";
import { get_tags } from "../api/tags.js";

import debounce from "../util/debounce.js";

import SchemaList from "./SchemaList.js";

const store = Vue.reactive({
    list: [],
    count: 0,
    busy: false,
    keywords: "",
    tag_id: -1,
    tags: [],
    tags_loaded: false
});

// load tags once
get_tags().then(tags => {
    store.tags = tags;
    store.tags_loaded = true;
});

export default {
    components: {
        "schema-list": SchemaList
    },
    data: () => store,
    mounted: function() {
        // get data from query
        if (!this.keywords && this.$route.query.q) {
            this.keywords = this.$route.query.q;
        }
        if (this.tag_id < 0 && this.$route.query.tag_id) {
            this.tag_id = +this.$route.query.tag_id;
        }

        // count and search on route load
        this.get_count();
        this.search();
    },
    watch: {
        "keywords": "get_count",
        "tag_id": "get_count"
    },
    methods: {
        search_body: function() {
            return {
                keywords: this.keywords ? this.keywords : null,
                tag_id: this.tag_id >= 0 ? this.tag_id : null,
                complete: true,
                limit: 24
            };
        },
        get_count: debounce(function() {
            this.busy = true;
            schema_count(this.search_body())
            .then(c => this.count = c)
            .then(() => this.busy = false);
        }, 250),
        search: function() {
            // set data on query
            this.$router.push({
                path: this.$route.path,
                query: {
                    q: this.keywords ? this.keywords : undefined,
                    tag_id: this.tag_id >= 0 ? this.tag_id : undefined
                }
            });

            schema_search(this.search_body())
            .then(list => {
                this.list = list;
                this.busy = false;
            });
        }
    },
    template: /*html*/`
    <div class="row" v-if="tags_loaded">
        <div class="col-md-8 col-xs-4">
            <input type="text" class="form-control" v-model="keywords" v-on:keyup.enter="search" placeholder="Keywords"/>
        </div>
        <div class="col-md-2 col-xs-4">
            <select class="form-control" v-model="tag_id">
                <option value="-1">All tags</option>
                <option v-for="tag in tags" :value="tag.id">{{tag.name}}</option>
            </select>
        </div>
        <div class="col-md-2 col-xs-4">
            <button class="btn btn-success w-100" v-on:click="search">
                <i class="fa fa-spinner fa-spin" v-if="busy"></i>
                <i class="fa-solid fa-magnifying-glass" v-else></i>
                Search
                <span class="badge rounded-pill bg-secondary">{{count}}</span>
            </button>
        </div>
    </div>
    <hr>
    <schema-list :list="list"/>
    `
};