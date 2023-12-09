import { login, logout, is_logged_in } from "../../service/login.js";
import { get_github_id, get_cdb_id, get_discord_id, get_mesehub_id, get_base_url } from "../../service/info.js";
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
        is_logged_in,
        get_github_id,
        get_cdb_id,
        get_discord_id,
        get_mesehub_id,
        get_base_url,
        validInput: function(){
            return this.username != "" && this.password != "";
        },
        github_href: function() {
			return `https://github.com/login/oauth/authorize?client_id=${get_github_id()}`;
		},
        cdb_href: function() {
            return `https://content.minetest.net/oauth/authorize/?response_type=code&client_id=${get_cdb_id()}&redirect_uri=${encodeURIComponent(get_base_url()+'/api/oauth_callback/cdb')}`;
        },
		discord_href: function() {
			return `https://discord.com/api/oauth2/authorize?client_id=${get_discord_id()}&redirect_uri=${encodeURIComponent(get_base_url()+'/api/oauth_callback/discord')}&response_type=code&scope=identify`;
		},
		mesehub_href: function() {
			return `https://git.minetest.land/login/oauth/authorize?client_id=${get_mesehub_id()}&redirect_uri=${encodeURIComponent(get_base_url()+'/api/oauth_callback/mesehub')}&response_type=code&state=STATE`;
		},
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
                    this.$router.push("/profile");
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
            <div class="col-md-2"></div>
            <div class="col-md-8 card" style="padding: 20px;">
                <div class="row">
                    <div class="col-md-6">
                        <h5>Login with username and password</h5>
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
                    <div class="col-md-6">
                        <h5>Login with external provider</h5>
                        <a :href="github_href" class="btn btn-secondary w-100" v-bind:class="{disabled:is_logged_in}" v-if="get_github_id">
                            <i class="fab fa-github"></i>
                            Login with Github
                        </a>
                        &nbsp;
                        <a :href="cdb_href" class="btn btn-secondary w-100" v-bind:class="{disabled:is_logged_in}" v-if="get_cdb_id">
                            <img :src="get_base_url + '/pics/contentdb.png'" height="24" width="24">
                            Login with ContentDB
                        </a>
                        &nbsp;
                        <a :href="discord_href" class="btn btn-secondary w-100" v-bind:class="{disabled:is_logged_in}" v-if="get_discord_id">
                            <i class="fab fa-discord"></i>
                            Login with Discord
                        </a>
                        &nbsp;
                        <a :href="mesehub_href" class="btn btn-secondary w-100" v-bind:class="{disabled:is_logged_in}" v-if="get_mesehub_id">
                            <img :src="get_base_url + '/pics/default_mese_crystal.png'">
                            Login with Mesehub
                        </a>
                    </div>
                </div>
            </div>
            <div class="col-md-4"></div>
        </div>
	`
};
