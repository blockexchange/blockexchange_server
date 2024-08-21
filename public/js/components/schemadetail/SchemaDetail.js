import format_size from "../../util/format_size.js";
import format_time from "../../util/format_time.js";

import ModalPrompt from "../ModalPrompt.js";
import SchemaStar from "./SchemaStar.js";
import SchemaDownload from "./SchemaDownload.js";
import SchemaCollection from "./SchemaCollection.js";
import ModBadge from "../ModBadge.js";
import TagBadge from "../TagBadge.js";

import {
    schema_update,
    schema_set_tags,
    schema_update_screenshot,
    schema_delete,
    schema_update_mods,
    schema_update_info
} from "../../api/schema.js";
import { get_tags } from "../../service/tags.js";
import { is_logged_in, has_permission } from "../../service/login.js";

export default {
    components: {
        "modal-prompt": ModalPrompt,
        "schema-star": SchemaStar,
        "schema-download": SchemaDownload,
        "schema-collection": SchemaCollection,
        "mod-badge": ModBadge,
        "tag-badge": TagBadge
    },
    props: {
        search_result: { type: Object, required: true },
        allow_edit: { type: Boolean, default: false }
    },
    data: function() {
        return {
            BaseURL,
            cache_counter: Date.now(),
            edit_mode: false,
            error_response: null,
            screenshot_busy: false,
            mods_busy: false,
            delete_prompt: false,
            username: this.search_result.username,
            schema: this.search_result.schema,
            mods: this.search_result.mods,
            tags: this.search_result.tags
        };
    },
    methods: {
        format_size,
        format_time,
        has_permission,
        get_tags: function() {
            return get_tags().filter(t => !t.restricted || has_permission("ADMIN"));
        },
        save: function() {
            this.error_response = null;
            schema_update(this.schema)
            .then(() => {
                // set tags
                schema_set_tags(this.schema.uid, this.tags);

                this.edit_mode = false;
                this.$router.push(`/schema/${this.username}/${this.schema.name}`);
                this.$emit("save");
            })
            .catch(e => {
                this.error_response = e;
            });
        },
        update_mods: function() {
            this.mods_busy = true;
            schema_update_mods(this.schema.uid)
            .then(m => {
                this.mods = m;
                this.mods_busy = false;
            });
        },
        update_info: function() {
            schema_update_info(this.schema.uid)
            .then(s => {
                this.schema.total_size = s.total_size;
                this.schema.total_parts = s.total_parts;
            });
        },
        update_screenshot: function() {
            this.screenshot_busy = true;
            schema_update_screenshot(this.schema.uid)
            .then(() => {
                this.screenshot_busy = false;
                this.cache_counter = Date.now();
            });
        },
        delete_schema: function() {
            schema_delete(this.schema.uid)
            .then(() => this.$router.push(`/schema/${this.username}`));
        },
        markdown: function(txt) {
            return DOMPurify.sanitize(marked.parse(txt));
        }
        
    },
    computed: {
        logged_in: is_logged_in
    },
    template: /*html*/`
    <div>
        <modal-prompt v-if="delete_prompt" title="Confirm deletion" v-on:close="delete_prompt = false">
            <template #body>
                <p>Confirm deletion of schematic <b>{{schema.name}}</b></p>
            </template>
            <template #footer>
                <a class="btn btn-secondary" v-on:click="delete_prompt = false">
                    <i class="fa fa-times"></i> Abort
                </a>
                <a class="btn btn-danger" v-on:click="delete_schema">
                    <i class="fa fa-trash"></i> Delete
                </a>
            </template>
        </modal-prompt>
        <div class="row">
            <div class="col-md-6">
                <h3 v-if="!edit_mode">
                    {{schema.name}}
                    <small class="text-muted">by {{username}}</small>
                    &nbsp;
                    <schema-star :schema="schema"/>
                </h3>
                <div v-else class="input-group has-validation">
                    <input type="text" class="form-control" v-model="schema.name" v-bind:class="{'is-invalid': error_response}">
                    <div class="invalid-feedback" v-if="error_response && error_response.name_taken">
                        Schematic name already taken
                    </div>
                    <div class="invalid-feedback" v-if="error_response && error_response.name_invalid">
                        Schematic name invalid, allowed chars: a to z, A to Z, 0 to 9 and -, _
                    </div>
                </div>
            </div>
            <div class="col-md-6">
                <div class="btn-group float-end" v-if="allow_edit && !edit_mode">
                    <a class="btn btn-sm btn-secondary" v-on:click="edit_mode = true">
                        <i class="fa fa-edit"></i> Edit
                    </a>
                    <a class="btn btn-sm btn-secondary" v-bind:class="{'disabled': screenshot_busy}" v-on:click="update_screenshot">
                        <i v-if="screenshot_busy" class="fa fa-spinner fa-spin"></i>
                        <i v-else class="fa fa-image"></i>
                        Update screenshot
                    </a>
                    <a class="btn btn-sm btn-secondary" v-bind:class="{'disabled': mods_busy}" v-on:click="update_mods">
                        <i v-if="mods_busy" class="fa fa-spinner fa-spin"></i>
                        <i v-else class="fa fa-box-archive"></i>
                        Update modnames
                    </a>
                    <a class="btn btn-sm btn-secondary" v-on:click="update_info">
                        <i class="fa fa-chart-simple"></i>
                        Update stats
                    </a>
                    <a class="btn btn-sm btn-danger" v-on:click="delete_prompt = true">
                        <i class="fa fa-trash"></i> Delete
                    </a>
                </div>
                <div class="btn-group float-end" v-if="edit_mode">
                    <a class="btn btn-sm btn-success" v-on:click="save">
                        <i class="fa fa-save"></i> Save
                    </a>
                </div>
            </div>
        </div>
        <hr>
        <div class="row">
            <div class="col-md-4">
                <div class="card">
                    <div class="card-header">
                        <i class="fa fa-magnifying-glass"></i>
                        Details
                    </div>
                    <div class="card-body">
                        <ul>
                            <li><b>Created:</b> {{format_time(schema.created)}}</li>
                            <li><b>Modified:</b> {{format_time(schema.mtime)}}</li>
                            <li><b>Size:</b> {{format_size(schema.total_size)}}</li>
                            <li><b>Dimensions:</b> {{schema.size_x}} / {{schema.size_y}} / {{schema.size_z}} nodes</li>
                            <li><b>Parts:</b> {{schema.total_parts}}</li>
                            <li><b>Downloads:</b> {{schema.downloads}}</li>
                            <li>
                                <b>License:</b>
                                <span v-if="!edit_mode">
                                    <img v-if="schema.license == 'CC0'" :src="BaseURL + '/pics/license_cc0.png'">
                                    <img v-else-if="schema.license == 'CC-BY-SA'" :src="BaseURL + '/pics/license_cc-by-sa.png'">
                                    <span v-else class="badge bg-secondary">{{schema.license}}</span>
                                </span>
                                <select v-else class="form-control" v-model="schema.license">
                                    <option value="CC0">CC0</option>
                                    <option value="CC-BY-SA">CC-BY-SA</option>
                                    <option value="MIT">MIT</option>
                                    <option value="Proprietary">Proprietary</option>
                                </select>
                            </li>
                        </ul>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-header">
                        <i class="fa fa-file-lines"></i>
                        Description
                    </div>
                    <div class="card-body">
                        <p v-if="!edit_mode">{{schema.short_description}}</p>
                        <input v-else type="text" class="form-control" v-model="schema.short_description" placeholder="Short description"/>
                        <hr>
                        <div v-if="!edit_mode" v-html="markdown(schema.description)"></div>
                        <textarea v-else v-model="schema.description" rows="10" class="form-control"></textarea>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-header">
                        <i class="fa fa-box-archive"></i>
                        Used mods
                    </div>
                    <div class="card-body">
                        <a v-if="!edit_mode && schema.cdb_collection"
                            class="btn btn-xs btn-outline-success w-100"
                            target="new"
                            :href="'https://content.minetest.net/collections/' + schema.cdb_collection + '/'">
                            <i class="fa fa-box-archive"></i>
                            Open CDB Mod collection
                            <i class="fa fa-up-right-from-square"></i>
                        </a>
                        <hr v-if="schema.cdb_collection">
                        <input v-if="edit_mode" type="text" class="form-control" v-model="schema.cdb_collection" placeholder="CDB Collection in the 'username/collectionname' form"/>
                        <mod-badge v-for="mod in mods" :name="mod"/>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-header">
                        <i class="fa fa-tags"></i>
                        Tags
                    </div>
                    <div class="card-body">
                        <tag-badge v-if="!edit_mode" v-for="tag in tags" :name="tag"/>
                        <div v-else>
                            <ul>
                                <li v-for="tag in get_tags()">
                                    <input type="checkbox" class="form-check-input" :value="tag.name" v-model="tags"/>
                                    <span class="badge bg-success">
                                        <i class="fa fa-tag"></i>
                                        {{tag.name}}
                                    </span>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
                <br>
                <div class="card" v-if="edit_mode || search_result.collection_name">
                    <div class="card-header">
                        <i class="fa fa-object-group"></i>
                        Collection
                    </div>
                    <div class="card-body">
                        <schema-collection
                            :schema="schema"
                            :username="username"
                            :edit_mode="edit_mode"
                            :collection_name="search_result.collection_name"/>
                    </div>
                </div>
            </div>
            <div class="col-md-8">
                <div class="card">
                    <div class="card-header">
                        <i class="fa fa-image"></i>
                        Preview
                    </div>
                    <div class="card-body">
                        <div class="text-center" style="min-height: 600px;">
                            <img :src="BaseURL + '/api/schema/' + schema.uid + '/screenshot?cachebust=' + cache_counter" class="img-fluid">
                        </div>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-header">
                        <i class="fa fa-download"></i>
                        Download
                    </div>
                    <div class="card-body">
                        <schema-download :schema="schema" :username="username"/>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
};