

export default {
    props: ["fetch_entries", "count_entries"],
    data: function() {
        return {
            list: [],
            busy: false,
            per_page: 20,
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
            this.current_page = this.$router.currentRoute.value.query.page || 1;
            this.count_entries()
            .then(c => this.total_count = c)
            .then(() => this.fetch_entries(this.per_page, (this.current_page-1)*this.per_page))
            .then(list => this.list = list)
            .then(() => this.busy = false);
        }
    },
    computed: {
        pages: function() {
            return Math.ceil(this.total_count / this.per_page);
        }
    },
    template: /*html*/`
    <table>
        <thead>
            <slot name="header"></slot>
        </thead>
        <tbody>
            <slot name="body" :list="list"></slot>
        </tbody>
        <div class="btn-group">
            <router-link
                :to="$router.currentRoute.value.path + '?page=' + i"
                v-for="i in pages"
                class="btn btn-xs btn-secondary">
                {{i}}
            </router-link>
        </div>
    </table>
    `
};
