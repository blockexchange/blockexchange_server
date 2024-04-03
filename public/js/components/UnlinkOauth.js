import { unlink_oauth } from "../api/user.js";
import { check_login } from "../service/login.js";

export default {
    props: ["useruid"],
    data: function() {
        return {
            password_1: "",
            password_2: "",
            busy: false
        };
    },
    methods: {
        unlink: function() {
            this.busy = true;
            unlink_oauth(this.useruid, {
                new_password: this.password_1
            })
            .then(() => check_login(true))
            .then(() => location.reload());
        }
    },
    template: /*html*/`
    <h5>Unlink OAuth provider</h5>
    <div class="input-group">
        <input type="password" v-model="password_1" placeholder="New password" class="form-control"/>
        <input type="password" v-model="password_2" placeholder="New password (repeat)" class="form-control"/>
        <button class="btn btn-warning" :disabled="!password_1 || !password_2 || password_1 != password_2 || busy" v-on:click="unlink">
            <i class="fa fa-link-slash" v-if="!busy"></i>
            <i class="fa fa-spinner fa-spin" v-if="busy"></i>
            Unlink and set password
        </button>
    </div>
    <div class="alert alert-warning">
        <i class="fa fa-triangle-exclamation"></i>
        <b>Warning:</b> This unlinks the account from the OAuth provider, this action is irreversible!
    </div>
    `
};