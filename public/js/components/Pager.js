
export default {
    props: ["current", "pages"],
    methods: {
        selectPage: function(num){
            this.$emit("switch", num);
        },
        next: function(){
            this.selectPage(this.current + 1);
        },
        previous: function(){
            this.selectPage(this.current - 1);
        }
    },
    template: /*html*/`
        <nav aria-label="Page navigation example" v-if="pages > 1">
            <ul class="pagination">
                <li class="page-item" v-bind:class="{ 'disabled': current == 1 }">
                    <a class="page-link" v-on:click="previous()">Previous</a>
                </li>
                <li class="page-item" v-for="index in pages" v-bind:class="{ 'active': index == current }">
                    <a class="page-link" v-on:click="selectPage(index)">{{ index }}</a>
                </li>
                <li class="page-item" v-bind:class="{ 'disabled': current == pages }">
                    <a class="page-link" v-on:click="next()">Next</a>
                </li>
            </ul>
        </nav>
    `
};