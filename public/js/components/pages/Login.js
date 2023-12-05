import { login, logout, is_logged_in } from "../../service/login.js";
import Breadcrumb, { LOGIN, START } from "../Breadcrumb.js";

export default {
	data: function() {
        return {
			breadcrumb: [START, LOGIN],
            username: "",
            password: "",
            busy: false,
            error_message: ""
        };
    },
    components: {
        "bread-crumb": Breadcrumb
    },
    computed: {
        is_logged_in: is_logged_in,
        validInput: function(){
            return this.username != "" && this.password != "";
        }
    },
    methods: {
        login: function() {
            this.busy = true;
            this.error_message = "";
            login(this.username, this.password)
            .then(success => {
                this.busy = false;
                if (!success) {
                    // no luck
                    this.error_message = "Login failed!";
                } else {
                    // go to base page
                    this.$router.push("/");
                }
            });
        },
        logout: function() {
            this.busy = true;
            logout()
            .then(() => this.busy = false);
        }
    },
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
                <h4>Login</h4>
                <form @submit.prevent="login">
                    <input type="text"
                        class="form-control"
                        placeholder="Username"
                        :disabled="is_logged_in"
                        v-model="username"/>
                    <input type="password"
                        class="form-control"
                        placeholder="Password"
                        :disabled="is_logged_in"
                        v-model="password"/>
                    <button class="btn btn-primary w-100" v-if="!is_logged_in" type="submit" :disabled="!validInput">
                        <i class="fa-solid fa-right-to-bracket"></i>
                        Login
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                        <span class="badge bg-danger">{{error_message}}</span>
                    </button>
                    <a class="btn btn-secondary w-100" v-if="is_logged_in" v-on:click="logout">
                        <i class="fa-solid fa-right-from-bracket"></i>
                        Logout
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                    </a>
                </form>
            </div>
            <div class="col-md-4"></div>
        </div>
	`
};
