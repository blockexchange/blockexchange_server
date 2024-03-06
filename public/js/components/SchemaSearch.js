import { schema_search, schema_count } from "../api/schema.js";
import { get_tags } from "../service/tags.js";

import debounce from "../util/debounce.js";

import SchemaList from "./SchemaList.js";
import PagedContent from "./PagedContent.js";

const store = Vue.reactive({
    count: 0,
    busy: false,
    keywords: "",
    tag_uid: "",
    mod_name: "",
    tags: []
});

export default {
    components: {
        "schema-list": SchemaList,
        "paged-content": PagedContent
    },
    data: () => store,
    created: function() {
        this.tags = get_tags();

        // get data from query
        if (this.$route.query.q) {
            this.keywords = this.$route.query.q;
        }
        if (this.$route.query.tag_uid) {
            this.tag_uid = this.$route.query.tag_uid;
        }
        if (this.$route.query.mod_name) {
            this.mod_name = this.$route.query.mod_name;
        }

        // count and search on route load
        this.get_count();
        this.search();
    },
    watch: {
        "keywords": "get_count",
        "mod_name": "get_count",
        "tag_uid": "get_count"
    },
    methods: {
        search_body: function(limit, offset) {
            return {
                keywords: this.keywords ? this.keywords : null,
                tag_uid: this.tag_uid ? this.tag_uid : null,
                mod_name: this.mod_name ? this.mod_name : null,
                complete: true,
                limit: limit,
                offset: offset
            };
        },
        fetch_entries: function(limit, offset) {
            return schema_search(this.search_body(limit, offset));
        },
        count_entries: function() {
            return schema_count(this.search_body());
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
                    tag_uid: this.tag_uid ? this.tag_uid : undefined,
                    mod_name: this.mod_name ? this.mod_name : undefined,
                    page: this.$route.query.page
                }
            });
        }
    },
    template: /*html*/`
    <div class="row">
        <div class="col-md-6 col-xs-4">
            <input type="text" class="form-control" v-model="keywords" v-on:keyup.enter="search" placeholder="Keywords"/>
        </div>
        <div class="col-md-2 col-xs-4">
            <input type="text" class="form-control" v-model="mod_name" v-on:keyup.enter="search" placeholder="Modname"/>
        </div>
        <div class="col-md-2 col-xs-4">
            <select class="form-control" v-model="tag_uid">
                <option value="">All tags</option>
                <option v-for="tag in tags" :value="tag.uid">{{tag.name}}</option>
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
    <paged-content
        :fetch_entries="fetch_entries"
        :count_entries="count_entries"
        per_page="24">
        <template #body="{ list }">
            <schema-list :list="list"/>
        </template>
    </paged-content>
    `
};