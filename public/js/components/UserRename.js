import { get_claims } from "../service/login.js";
import { search_user, save_user, get_user } from "../api/user.js";
import { check_login } from "../service/login.js";

export default {
    data: function() {
        return {
            name: "",
            username_taken: false
        };
    },
    mounted: function() {
        this.name = get_claims().username;
    },
    methods: {
        rename: function() {
            search_user({ name: this.name })
            .then(list => {
                this.username_taken = list.length > 0;
                if (this.username_taken) {
                    return;
                }
                return get_user(get_claims().user_id)
                .then(user => {
                    user.name = this.name;
                    return save_user(user);
                })
                .then(() => check_login(true));
            });
        }
    },
    watch: {
        name: function() {
            this.username_taken = false;
        }
    },
    computed: {
        can_rename: function() {
            return (this.name && this.name != get_claims().username);
        },
        username_valid: function() {
            return /^[a-zA-Z0-9_.-]*$/.test(this.name);
        },
        is_valid: function() {
            return this.username_valid && !this.username_taken;
        }
    },
    template: /*html*/`
        <h5>Change username</h5>
        <div class="input-group has-validation">
            <input type="text" v-model="name" class="form-control" v-bind:class="{'is-invalid':!is_valid}"/>
            <button class="btn btn-warning" :disabled="!can_rename || !is_valid" v-on:click="rename">
                <i class="fa fa-edit"></i>
                Rename
            </button>
            <div class="invalid-feedback" v-if="!username_valid">
                Username is invalid, allowed chars: a to z, A to Z, 0 to 9 and -, _
            </div>
            <div class="invalid-feedback" v-if="username_taken">
                Username is already taken
            </div>
        </div>
    `
};
