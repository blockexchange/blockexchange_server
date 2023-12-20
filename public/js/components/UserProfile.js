import format_time from "../util/format_time.js";

export default {
    props: ["user"],
    methods: {
        format_time
    },
	template: /*html*/`
    <div v-if="user">
        <div class="row">
            <div class="col-md-6">
                <h5>
                    User profile for
                    <small class="text-body-secondary">{{user.name}}</small>
                </h5>
                <ul>
                    <li>
                        <b>Username:</b> {{user.name}}
                    </li>
                    <li>
                        <b>Created:</b> {{format_time(user.created)}}
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
            <div class="col-md-6 text-end">
                <img v-if="user.avatar_url" :src="user.avatar_url" height="64" width="64"/>
            </div>
        </div>
        <router-link :to="'/schema/' + user.name" class="btn btn-xs btn-outline-success">
            Show all schematics of {{user.name}}
        </router-link>
    </div>
	`
};
