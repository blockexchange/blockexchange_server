
export default {
    props: ["total", "route", "limit"],
    data: function() {
        return {
            page: 1
        };
    },
    methods: {
        update: function(){
            const r = this.route;
			this.page = r.query.page || 1;
            this.$emit("fetchData", this.limit, (this.page-1) * this.limit);
        },
        selectPage: function(num){
			this.$router.push({
				path: this.route.path,
				query: {
					page: num
				}
			});
        },
        next: function(){
            this.selectPage(+this.page + 1);
        },
        previous: function(){
            this.selectPage(+this.page - 1);
        }
    },
    computed: {
        pages: function() {
            return this.total ? Math.ceil(this.total / this.limit) : 0;
        }
    },
    watch: {
		"route": "update"
    },
    created: function(){
        this.update();
    },
    template: /*html*/`
        <nav aria-label="Page navigation example" v-if="pages && pages > 1">
            <ul class="pagination">
                <li class="page-item" v-bind:class="{ 'disabled': page == 1 }">
                    <a class="page-link" v-on:click="previous()">Previous</a>
                </li>
                <li class="page-item" v-for="index in pages" v-bind:class="{ 'active': index == page }">
                    <a class="page-link" v-on:click="selectPage(index)">{{ index }}</a>
                </li>
                <li class="page-item" v-bind:class="{ 'disabled': page == pages }">
                    <a class="page-link" v-on:click="next()">Next</a>
                </li>
            </ul>
        </nav>
    `
};