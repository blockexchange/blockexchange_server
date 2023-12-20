import format_size from "../../util/format_size.js";
import format_time from "../../util/format_time.js";

import ModalPrompt from "../ModalPrompt.js";
import ClipboardCopy from "../ClipboardCopy.js";

import { schema_update, schema_set_tags, schema_update_screenshot, schema_delete } from "../../api/schema.js";
import { get_schema_star, star_schema, unstar_schema, count_schema_stars } from "../../api/schema_star.js";
import { get_tags } from "../../service/tags.js";
import { is_logged_in } from "../../service/login.js";

export default {
    components: {
        "modal-prompt": ModalPrompt,
        "clipboard-copy": ClipboardCopy
    },
    props: {
        schema: { type: Object, required: true },
        username: { type: String, required: true },
        allow_edit: { type: Boolean, default: false }
    },
    mounted: function() {
        this.update_star();
    },
    data: function() {
        return {
            BaseURL,
            edit_mode: false,
            error_response: null,
            screenshot_busy: false,
            delete_prompt: false,
            schema_star: null
        };
    },
    methods: {
        format_size,
        format_time,
        get_tags,
        update_star: function() {
            if (is_logged_in()) {
                get_schema_star(this.schema.id)
                .then(s => this.schema_star = s)
                .then(() => count_schema_stars(this.schema.id))
                .then(count => this.schema.stars = count);
            }
        },
        star: function() {
            star_schema(this.schema.id).then(() => this.update_star());
        },
        unstar: function() {
            unstar_schema(this.schema.id).then(() => this.update_star());
        },
        save: function() {
            this.error_response = null;
            schema_update(this.schema)
            .then(() => {
                // set tags
                schema_set_tags(this.schema.id, this.schema.tags);

                this.edit_mode = false;
                this.$router.push(`/schema/${this.username}/${this.schema.name}`);
            })
            .catch(e => {
                this.error_response = e;
            });
        },
        update_screenshot: function() {
            this.screenshot_busy = true;
            schema_update_screenshot(this.schema.id)
            .then(() => this.screenshot_busy = false);
        },
        delete_schema: function() {
            schema_delete(this.schema.id)
            .then(() => this.$router.push(`/schema/${this.username}`));
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
                    <button class="btn btn-outline-primary" :disabled="!logged_in" v-on:click="schema_star ? unstar() : star()">
                        <i class="fa fa-star" v-bind:style="{color: schema_star ? 'yellow' : ''}"></i>
                        <span class="badge bg-secondary rouded-pill">{{schema.stars}}</span>
                        {{ schema_star ? 'Unstar' : 'Star' }}
                    </button>
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
                        Description
                    </div>
                    <div class="card-body">
                        <pre v-if="!edit_mode">{{schema.description}}</pre>
                        <textarea v-else v-model="schema.description" rows="10" class="form-control"></textarea>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-header">
                        Used mods
                    </div>
                    <div class="card-body">
                        <span v-for="mod in schema.mods" class="badge bg-primary" style="margin-right: 5px;">
                            <i class="fa fa-box-archive"></i>
                            {{mod}}
                        </span>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-header">
                        Tags
                    </div>
                    <div class="card-body">
                        <span v-if="!edit_mode" class="badge bg-success" v-for="tag in schema.tags" style="margin-right: 5px;">
                            <i class="fas fa-tag"></i>
                            {{tag}}
                        </span>
                        <div v-else>
                            <ul>
                                <li v-for="tag in get_tags()">
                                    <input type="checkbox" class="form-check-input" :value="tag.name" v-model="schema.tags"/>
                                    <span class="badge bg-success">
                                        <i class="fas fa-tag"></i>
                                        {{tag.name}}
                                    </span>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            <div class="col-md-8">
                <div class="card">
                    <div class="card-header">
                        Preview
                    </div>
                    <div class="card-body">
                        <div class="text-center" style="min-height: 600px;">
                            <img :src="BaseURL + '/api/schema/' + schema.id + '/screenshot'" class="img-fluid">
                        </div>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-header">
                        Download
                    </div>
                    <div class="card-body">
                        <div class="row">
                            <div class="col-md-6">
                                <h5>Online</h5>
                                <ul>
                                    <li>
                                        Download and install the <router-link to="/mod">blockexchange</router-link>-mod
                                    </li>
                                    <li>
                                        Select the origin position with <clipboard-copy :text="'/bx_pos1'"></clipboard-copy>
                                    </li>
                                    <li>
                                        Check if the placement fits with <clipboard-copy :text="'/bx_allocate ' + username + ' ' + schema.name"></clipboard-copy>
                                    </li>
                                    <li>
                                        Place the schematic with <clipboard-copy :text="'/bx_load ' + username + ' ' + schema.name"></clipboard-copy>
                                    </li>
                                </ul>
                            </div>
                            <div class="col-md-6">
                                <h5>Offline</h5>
                                <div class="btn-group">
                                    <a class="btn btn-outline-primary"
                                        :href="BaseURL + '/api/export_bx/' + schema.id + '/' + schema.name + '.zip'">
                                        <i class="fa fa-download"></i>
                                        Download as BX schematic
                                    </a>
                                    <a class="btn btn-outline-primary" v-if="schema.total_parts < 200"
                                        :href="BaseURL + '/api/export_we/' + schema.id + '/' + schema.name + '.we'">
                                        <i class="fa fa-download"></i>
                                        Download as WE schematic
                                        <i class="fa fa-triangle-exclamation"
                                            v-if="schema.total_parts > 50"
                                            style="color: yellow;"
                                            title="WE-Import might be slow due to schematic size"></i>
                                    </a>
                                </div>
                                <ul>
                                    <li>
                                        Download and install the <router-link to="/mod">blockexchange</router-link>-mod
                                    </li>
                                    <li>
                                        Download the Blockexchange schematic with the above button and place it in the
                                        <b class="text-muted">{world-folder}/bxschems/</b> folder
                                    </li>
                                    <li>
                                        Select the origin position with <clipboard-copy :text="'/bx_pos1'"></clipboard-copy>
                                    </li>
                                    <li>
                                        Check if the placement fits with <clipboard-copy :text="'/bx_allocate_locale ' + username + ' ' + schema.name"></clipboard-copy>
                                    </li>
                                    <li>
                                        Place the schematic with <clipboard-copy :text="'/bx_load_local ' + username + ' ' + schema.name"></clipboard-copy>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
};