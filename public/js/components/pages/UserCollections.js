import { create_collection, get_collections_by_username } from "../../api/collection.js";
import { get_user_uid } from "../../service/login.js";

import Breadcrumb, { START, USER_COLLECTIONS } from "../Breadcrumb.js";
import LoadingSpinner from "../LoadingSpinner.js";

export default {
    props: ["username", "readonly"],
	components: {
        "bread-crumb": Breadcrumb,
        "loading-spinner": LoadingSpinner
	},
	data: function() {
		return {
			breadcrumb: [START, USER_COLLECTIONS(this.username)],
            collections: [],
            busy: false,
            new_collection_name: ""
		};
	},
    mounted: function() {
        this.update();
    },
    methods: {
        update: function() {
            this.busy = true;
            get_collections_by_username(this.username)
            .then(c => {
                c.sort((a,b) => a.name > b.name);
                this.collections = c;
                this.busy = false;
            });
        },
        create_collection: function() {
            create_collection({
                name: this.new_collection_name,
                user_uid: get_user_uid()
            })
            .then(() => {
                this.new_collection_name = "";
                this.update();
            });
        }
    },
    computed: {
        valid_name: function() {
            return /^[a-zA-Z0-9_.-]*$/.test(this.new_collection_name);
        },
        duplicate_name: function() {
            return this.collections.some(c => c.name == this.new_collection_name);
        }
    },
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
        <loading-spinner v-if="busy"/>
        <div v-else>
            <table class="table table-dark table-striped table-condensed">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr :key="collection.uid" v-for="collection in collections">
                        <td>
                            <router-link
                                class="link-info link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover"
                                :to="'/collections/' + username + '/' + collection.name">
                                <i class="fa fa-object-group"></i>
                                {{collection.name}}
                            </router-link>
                        </td>
                        <td>
                            <div class="btn-group">
                                <button class="btn btn-danger">
                                    <i class="fa fa-trash"></i>
                                    Remove
                                </button>
                                <button class="btn btn-primary">
                                    <i class="fa fa-edit"></i>
                                    Edit
                                </button>
                            </div>
                        </td>
                    </tr>
                    <tr v-if="!readonly">
                        <td>
                            <input type="text"
                                class="form-control"
                                v-bind:class="{'is-invalid': !valid_name || duplicate_name}"
                                placeholder="New collection name"
                                v-model="new_collection_name"/>
                            <div class="invalid-feedback" v-if="new_collection_name && !valid_name">
                                Name is invalid, allowed chars: a to z, A to Z, 0 to 9 and -._
                            </div>   
                            <div class="invalid-feedback" v-if="new_collection_name && duplicate_name">
                                Name is already taken
                            </div>
                        </td>
                        <td>
                            <button class="btn btn-secondary"
                                :disabled="!valid_name || duplicate_name || !new_collection_name"
                                v-on:click="create_collection">
                                <i class="fa fa-plus"></i> Create collection
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
	`
};
