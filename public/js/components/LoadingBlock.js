import LoadingSpinner from "./LoadingSpinner.js";

export default {
    props: ["fetch_data"],
    components: {
        "loading-spinner": LoadingSpinner
    },
    data: function() {
        return {
            busy: false,
            data: null
        };
    },
    methods: {
        update_data: function() {
            this.busy = true;
            const promise_map = this.fetch_data();
            const result_map = {};
            const promise_list = [];
            Object.keys(promise_map).forEach(key => {
                const p = promise_map[key];
                promise_list.push(p);
                p.then(result => result_map[key] = result);
            });
            Promise.all(promise_list).then(() => {
                this.data = result_map;
                this.busy = false;
            });
        }
    },
    mounted: function() {
        this.update_data();
    },
    template: /*html*/`
    <loading-spinner v-if="busy"/>
    <slot :data="data" :update_data="update_data" v-if="data"></slot>
    `
};