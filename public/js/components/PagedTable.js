

export default {
    props: ["fetch_entries", "count_entries"],
    data: function() {
        return {
            list: [],
            busy: false,
            total_count: 0,
            current_page: this.$router.currentRoute.value.query.page || 1
        };
    },
    mounted: function() {
        this.busy = true;
        this.count_entries()
        .then(c => this.total_count = c)
        .then(() => this.fetch_entries(10, 0))
        .then(list => this.list = list)
        .then(() => this.busy = false);
    },
    methods: {

    },
    template: /*html*/`
    <table>
        <slot name="header"></slot>
        <slot name="body" :list="list"></slot>
    </table>
    `
};
