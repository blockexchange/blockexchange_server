import { get_access_tokens, create_access_token, delete_access_token } from "../api/access_token.js";
import format_time from "../util/format_time.js";

export default {
    props: ["username"],
    data: function() {
        return {
            list: [],
            new_name: "",
            new_days: 365
        };
    },
    mounted: function() {
        this.update();
    },
    methods: {
        format_time,
        is_expired: function(at) {
            return Date.now() > at.expires;
        },
        update: function() {
            get_access_tokens().then(l => this.list = l);
        },
        create: function() {
            create_access_token({
                name: this.new_name,
                expires: Date.now() + (1000 * 3600 * 24 * this.new_days)
            })
            .then(() => {
                this.new_name = "";
                this.update();
            });
        },
        remove: function(at) {
            delete_access_token(at.id)
            .then(() => this.update());
        }
    },
    template: /*html*/`
    <div>
        <h5>Manage access-tokens</h5>
        <table class="table table-striped table-condensed table-dark">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Token/Login</th>
                    <th>Created</th>
                    <th>Expires</th>
                    <th>Use-count</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="at in list" :key="at.id" v-bind:class="{'table-warning': is_expired(at)}">
                    <td>{{at.name}}</td>
                    <td>
                        <pre>/bx_login {{username}} {{at.token}}</pre>
                    </td>
                    <td>{{format_time(at.created)}}</td>
                    <td>{{format_time(at.expires)}}</td>
                    <td>{{at.usecount}}</td>
                    <td>
                        <a class="btn btn-danger" v-on:click="remove(at)">
                            <i class="fa fa-trash"></i>
                            Delete
                        </a>
                    </td>
                </tr>
                <tr>
                    <td colspan="2">
                        <label>Token-Name (usually the server-name you are using it on)</label>
                        <input class="form-control" placeholder="Name" v-model="new_name">
                    </td>
                    <td></td>
                    <td>
                        <label>Expiration (in days)</label>
                        <input class="form-control" placeholder="Expiration (days)" v-model="new_days">
                    </td>
                    <td></td>
                    <td>
                        <a class="btn btn-primary" v-bind:class="{disabled: !new_name || !new_days || new_days <= 0}" v-on:click="create">
                            <i class="fa fa-plus"></i>
                            Create
                        </a>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    `
};