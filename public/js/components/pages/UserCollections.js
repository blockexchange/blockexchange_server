import Breadcrumb, { START, USER_COLLECTIONS } from "../Breadcrumb.js";
import { create_collection, get_collections_by_username } from "../../api/collection.js";
import { get_user_uid } from "../../service/login.js";
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
            <div class="row" v-if="!readonly">
                <div class="col-md-4"></div>
                <div class="col-md-4">
                    <div class="input-group has-validation">
                        <input type="text"
                            class="form-control"
                            v-bind:class="{'is-invalid': !valid_name || duplicate_name}"
                            placeholder="New collection name"
                            v-model="new_collection_name"/>
                        <button class="btn btn-secondary"
                            :disabled="!valid_name || duplicate_name || !new_collection_name"
                            v-on:click="create_collection">
                            <i class="fa fa-plus"></i> Create collection
                        </button>
                        <div class="invalid-feedback" v-if="new_collection_name && !valid_name">
                            Name is invalid, allowed chars: a to z, A to Z, 0 to 9 and -._
                        </div>   
                        <div class="invalid-feedback" v-if="new_collection_name && duplicate_name">
                            Name is already taken
                        </div>              
                    </div>
                </div>
                <div class="col-md-4"></div>
            </div>
            <hr>
            <div class="row">
                <div style="padding-bottom: 10px; width: 320px; min-height: 450px;" :key="collection.uid" v-for="collection in collections">
                    <div class="card" style="min-height: 400px;">
                        <div class="card-header">
                            {{collection.name}}
                        </div>
                        <div class="card-body">
                            <h5 class="card-title">Title</h5>
                            <p>
                                Stuff
                            </p>
                        </div>
                    </div>
                </div>
            </div>
            {{JSON.stringify(collections)}}
        </div>
	`
};
