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
    mounted: function() {
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
    },
    template: /*html*/`
    <loading-spinner v-if="busy"/>
    <slot :data="data" v-if="data"></slot>
    `
};