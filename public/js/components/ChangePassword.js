import { change_password } from "../api/user.js";

export default {
    props: ["useruid"],
    data: function() {
        return {
            old_password: "",
            new_password: "",
            password_wrong: false,
            success: false,
            busy: false
        };
    },
    methods: {
        change_password: function() {
            this.password_wrong = false;
            this.success = false;
            this.busy = true;
            change_password(this.useruid, {
                old_password: this.old_password,
                new_password: this.new_password
            }).then(success => {
                this.password_wrong = !success;
                this.success = success;
            }).finally(() => {
                this.busy = false;
            });
        }
    },
    template: /*html*/`
    <h5>Change password</h5>
    <div class="input-group has-validation">
        <input type="password" v-model="old_password" placeholder="Old password" class="form-control" v-bind:class="{'is-invalid': password_wrong}"/>
        <input type="password" v-model="new_password" placeholder="New password" class="form-control" v-bind:class="{'is-valid': success}"/>
        <button class="btn btn-primary" :disabled="!old_password || !new_password || busy" v-on:click="change_password">
            <i class="fa fa-edit" v-if="!busy"></i>
            <i class="fa fa-spinner fa-spin" v-if="busy"></i>
            Change password
        </button>
        <div class="invalid-feedback" v-if="password_wrong">
            Password invalid
        </div>
        <div class="valid-feedback" v-if="success">
            Password changed
        </div>
    </div>
    `
};