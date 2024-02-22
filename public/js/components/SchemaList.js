import format_size from "../util/format_size.js";

export default {
    props: {
        list: { type: Array },
        show_not_found: { type: Boolean, default: true }
    },
    methods: {
        format_size,
        preview_src: function(schema) {
            return BaseURL + '/api/schema/' + schema.uid + '/screenshot?height=240&width=360';
        },
        schema_link: function(entry) {
            return '/schema/' + entry.username + '/' + entry.schema.name;
        }
    },
    template: /*html*/`
    <div class="row">
        <div class="col-md-12" v-if="list && list.length == 0 && show_not_found">
            <div class="alert alert-secondary">
                <i class="fa fa-circle-info"></i>
                No schematics found
            </div>
        </div>
        <div style="padding-bottom: 10px; width: 320px; min-height: 450px;" :key="entry.schema.uid" v-for="entry in list">
            <div class="card" style="min-height: 400px;">
                <router-link :to="schema_link(entry)">
                    <img
                        :src="preview_src(entry.schema)"
                        class="card-img-top"
                        style="background-color: #3c3737; min-height: 200px;"/>
                </router-link>
                <div class="card-body">
                    <h5 class="card-title">
                        <p>
                        <router-link :to="schema_link(entry)" class="link-success link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover">
                            <i class="fa fa-landmark"></i>
                            {{entry.schema.name}}
                        </router-link>
                        <i class="fa fa-star" v-if="entry.schema.stars" v-bind:style="{ color: entry.schema.stars ? 'yellow' : '' }"></i>
                        <span class="badge bg-secondary rounded-pill" v-if="entry.schema.stars">{{entry.schema.stars}}</span>
                    </p>
                    <p>
                        <router-link :to="'/user/' + entry.username" class="link-light link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover">
                            <i class="fa fa-user"></i>
                            {{entry.username}}
                        </router-link>
                    </p>
                    </h5>
                    <p>
                        <span class="badge bg-success" v-for="tag in entry.tags" style="margin-right: 5px;">
                            <i class="fas fa-tag"></i>
                            {{tag}}
                        </span>                    
                    </p>
                    <p>
                        <router-link
                            :to="'/collections/' + entry.username + '/' + entry.collection_name"
                            class="link-info link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover"
                            v-if="entry.collection_name">
                            <i class="fa fa-object-group"></i>
                            {{entry.collection_name}}
                        </router-link>
                    </p>
                    <p>
                        {{format_size(entry.schema.total_size)}};
                        {{entry.schema.size_x}} / {{entry.schema.size_y}} / {{entry.schema.size_z}} nodes
                    </p>
                </div>
            </div>
        </div>
    </div>
    `
};