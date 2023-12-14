import { search_user } from "../api/user.js";
import format_time from "../util/format_time.js";

export default {
    props: ["name"],
    mounted: function() {
        if (this.name) {
            search_user({ name: this.name })
            .then(u => this.user = u[0]);
        }
    },
	data: function() {
		return {
			user: null
		};
	},
    methods: {
        format_time
    },
	template: /*html*/`
    <div v-if="user">
        <h5>
            User profile for
            <small class="text-body-secondary">{{user.name}}</small>
        </h5>
        <ul>
            <li>
                <b>Username:</b> {{user.name}}
            </li>
            <li>
                <b>Created:</b> {{format_time(user.created/1000)}}
            </li>
            <li>
                <b>ID:</b> <span class="badge bg-success">{{user.id}}</span>
            </li>
            <li>
                <b>Type:</b> <span class="badge bg-secondary">{{user.type}}</span>
            </li>
            <li>
                <b>Role:</b> <span class="badge bg-secondary">{{user.role}}</span>
            </li>
        </ul>
    </div>
	`
};
