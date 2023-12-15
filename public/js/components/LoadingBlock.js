
export default {
    props: ["fetch_data"],
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
    <div class="alert alert-primary" v-if="busy">
        <i class="fa fa-spinner fa-spin"></i> Loading
    </div>
    <slot :data="data" v-if="data"></slot>
    `
};