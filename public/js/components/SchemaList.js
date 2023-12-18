import format_size from "../util/format_size.js";

export default {
    props: ["list"],
    methods: {
        format_size,
        preview_src: function(schema) {
            return BaseURL + '/api/schema/' + schema.id + '/screenshot?height=240&width=360';
        },
        schema_link: function(schema) {
            return '/schema/' + schema.username + '/' + schema.name;
        }
    },
    template: /*html*/`
    <div class="row">
        <div class="col-md-12" v-if="list && list.length == 0">
            <div class="alert alert-secondary">
                <i class="fa fa-circle-info"></i>
                No schematics found
            </div>
        </div>
        <div class="col-md-2" style="padding-bottom: 10px;" :key="schema.id" v-for="schema in list">
            <div class="card" style="min-height: 400px;">
                <router-link :to="schema_link(schema)">
                    <img
                        :src="preview_src(schema)"
                        class="card-img-top"
                        style="background-color: #3c3737; min-height: 200px;"/>
                </router-link>
                <div class="card-body">
                    <h5 class="card-title">
                        <p>
                        <router-link :to="schema_link(schema)">
                            <i class="fa fa-landmark"></i>
                            {{schema.name}}
                        </router-link>
                        <i class="fa fa-star" v-bind:style="{ color: schema.stars ? 'yellow' : '' }"></i>
                        <span class="badge bg-secondary rounded-pill">{{schema.stars}}</span>
                    </p>
                    <p>
                        <router-link :to="'/user/' + schema.username" class="text-muted">
                            <i class="fa fa-user"></i>
                            {{schema.username}}
                        </router-link>
                    </p>
                    </h5>
                    <p>
                        <span class="badge bg-success" v-for="tag in schema.tags" style="margin-right: 5px;">
                            <i class="fas fa-tag"></i>
                            {{tag}}
                        </span>                    
                    </p>
                    <p>
                        {{format_size(schema.total_size)}};
                        {{schema.size_x}} / {{schema.size_y}} / {{schema.size_z}} nodes
                    </p>
                </div>
            </div>
        </div>
    </div>
    `
};