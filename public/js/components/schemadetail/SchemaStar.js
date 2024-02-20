import { is_logged_in } from "../../service/login.js";
import { get_schema_star, star_schema, unstar_schema, count_schema_stars } from "../../api/schema_star.js";

export default {
    props: ["schema"],
    data: function() {
        return {
            schema_star: null
        };
    },
    mounted: function() {
        this.update_star();
    },
    methods: {
        update_star: function() {
            if (is_logged_in()) {
                get_schema_star(this.schema.uid)
                .then(s => this.schema_star = s)
                .then(() => count_schema_stars(this.schema.uid))
                .then(count => this.schema.stars = count);
            }
        },
        star: function() {
            star_schema(this.schema.uid).then(() => this.update_star());
        },
        unstar: function() {
            unstar_schema(this.schema.uid).then(() => this.update_star());
        }
    },
    computed: {
        logged_in: is_logged_in
    },
    template: /*html*/`
    <button class="btn btn-outline-primary" :disabled="!logged_in" v-on:click="schema_star ? unstar() : star()">
        <i class="fa fa-star" v-bind:style="{color: schema_star ? 'yellow' : ''}"></i>
        <span class="badge bg-secondary rouded-pill">{{schema.stars}}</span>
        {{ schema_star ? 'Unstar' : 'Star' }}
    </button>
    `
};