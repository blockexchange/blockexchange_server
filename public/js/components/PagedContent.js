export default {
    props: ["fetch_entries", "count_entries", "per_page"],
    data: function() {
        return {
            list: [],
            busy: false,
            total_count: 0,
            current_page: 1
        };
    },
    mounted: function() {
        this.update();
    },
    watch: {
        "$route": "update"
    },
    methods: {
        update: function() {
            this.busy = true;
            this.list = [];
            this.current_page = this.$route.query.page || 1;
            this.count_entries()
            .then(c => this.total_count = c)
            .then(() => this.fetch_entries(+this.per_page, (this.current_page-1)*this.per_page))
            .then(list => this.list = list)
            .then(() => this.busy = false);
        },
        get_route: function(page) {
            return {
                path: this.$route.path,
                query: Object.assign({}, this.$route.query, {
                    page: page
                })
            };
        }
    },
    computed: {
        pages: function() {
            return Math.ceil(this.total_count / this.per_page);
        }
    },
    template: /*html*/`
        <slot name="body" :list="list"></slot>
        <slot name="busy" v-if="busy">
            <div class="alert alert-primary">
                <i class="fa fa-spinner fa-spin"></i> Loading...
            </div>
        </slot>
        <div class="btn-group" v-if="!busy">
            <router-link
                :to="get_route(i)"
                v-for="i in pages"
                class="btn btn-xs btn-secondary"
                v-bind:class="{disabled: this.current_page == i}">
                {{i}}
            </router-link>
        </div>
    `
};
