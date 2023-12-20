import format_time from "../util/format_time.js";
import { count_user_schema_stars, search_users } from "../api/user.js";
import { schema_count } from "../api/schema.js";

import LoadingBlock from "./LoadingBlock.js";

export default {
    props: ["username"],
    components: {
        "loading-block": LoadingBlock
    },
    methods: {
        format_time,
        fetch_data: function() {
            var user;
            return {
				user: search_users({ name: this.username })
                    .then(l => l[0])
                    .then(u => {
                        user = u;
                        return count_user_schema_stars(u.id);
                    })
                    .then(stars => {
                        return Object.assign(user, {
                            stars
                        });
                    }),
                schema_count: schema_count({ user_name: this.username })
            };
        }
    },
	template: /*html*/`
    <loading-block :fetch_data="fetch_data" v-slot="{ data }">
        <div class="row">
            <div class="col-md-6">
                <h5>
                    User profile for
                    <small class="text-body-secondary">{{data.user.name}}</small>
                </h5>
                <ul>
                    <li>
                        <b>Username:</b> {{data.user.name}}
                    </li>
                    <li>
                        <b>Schematics:</b> <span class="badge bg-secondary rouded-pill">{{data.schema_count}}</span>
                    </li>
                    <li>
                        <b>Stars:</b>
                        <i class="fa fa-star" v-bind:style="{color: data.user.stars ? 'yellow' : ''}"></i>
                        <span class="badge bg-secondary rouded-pill">{{data.user.stars}}</span>
                    </li>
                    <li>
                        <b>Created:</b> {{format_time(data.user.created)}}
                    </li>
                    <li>
                        <b>ID:</b> <span class="badge bg-success">{{data.user.id}}</span>
                    </li>
                    <li>
                        <b>Type:</b> <span class="badge bg-secondary">{{data.user.type}}</span>
                    </li>
                    <li>
                        <b>Role:</b> <span class="badge bg-secondary">{{data.user.role}}</span>
                    </li>
                </ul>
            </div>
            <div class="col-md-6 text-end">
                <img v-if="data.user.avatar_url" :src="data.user.avatar_url" height="64" width="64"/>
            </div>
        </div>
        <router-link :to="'/schema/' + data.user.name" class="btn btn-xs btn-outline-success">
            Show all schematics of {{data.user.name}}
        </router-link>
    </loading-block>
	`
};
