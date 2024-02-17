import format_size from "../../util/format_size.js";
import format_time from "../../util/format_time.js";

import ModalPrompt from "../ModalPrompt.js";
import ClipboardCopy from "../ClipboardCopy.js";

import {
    schema_update,
    schema_set_tags,
    schema_update_screenshot,
    schema_delete,
    schema_update_mods,
    schema_update_info
} from "../../api/schema.js";
import { get_schema_star, star_schema, unstar_schema, count_schema_stars } from "../../api/schema_star.js";
import { get_tags } from "../../service/tags.js";
import { is_logged_in, has_permission } from "../../service/login.js";

const max_placement_tool_size = 100;

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
            cache_counter: Date.now(),
            edit_mode: false,
            error_response: null,
            screenshot_busy: false,
            mods_busy: false,
            delete_prompt: false,
            schema_star: null
        };
    },
    methods: {
        format_size,
        format_time,
        has_permission,
        get_tags: function() {
            return get_tags().filter(t => !t.restricted || has_permission("ADMIN"));
        },
        update_star: function() {
            if (is_logged_in()) {
                get_schema_star(this.schema.uid)
                .then(s => this.schema_star = s)
                .then(() => count_schema_stars(this.schema.uid))
                .then(count => this.schema.stars = count);
            }
        },
        star: function() {
            star_schema(this.schema.uid).then(() => this.update_star());
        },
        unstar: function() {
            unstar_schema(this.schema.uid).then(() => this.update_star());
        },
        save: function() {
            this.error_response = null;
            schema_update(this.schema)
            .then(() => {
                // set tags
                schema_set_tags(this.schema.uid, this.schema.tags);

                this.edit_mode = false;
                this.$router.push(`/schema/${this.username}/${this.schema.name}`);
            })
            .catch(e => {
                this.error_response = e;
            });
        },
        update_mods: function() {
            this.mods_busy = true;
            schema_update_mods(this.schema.uid)
            .then(m => {
                this.schema.mods = m;
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
        logged_in: is_logged_in,
        can_use_placement_tool: function() {
            return (this.schema.size_x <= max_placement_tool_size &&
                    this.schema.size_y <= max_placement_tool_size &&
                    this.schema.size_z <= max_placement_tool_size);
        }
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
                        Used mods
                    </div>
                    <div class="card-body">
                        <a v-if="!edit_mode && schema.cdb_collection"
                            class="btn btn-xs btn-outline-success"
                            target="new"
                            :href="'https://content.minetest.net/collections/' + schema.cdb_collection + '/'">
                            <i class="fa fa-box-archive"></i>
                            Open CDB Mod collection
                        </a>
                        <hr v-if="schema.cdb_collection">
                        <input v-if="edit_mode" type="text" class="form-control" v-model="schema.cdb_collection" placeholder="CDB Collection in the 'username/collectionname' form"/>
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
                            <img :src="BaseURL + '/api/schema/' + schema.uid + '/screenshot?cachebust=' + cache_counter" class="img-fluid">
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
                                <u>
                                    <h4>Online</h4>
                                </u>
                                <h5>Single placement</h5>
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
                                <h5 v-if="can_use_placement_tool">Multiple placements with the placement tool</h5>
                                <ul v-if="can_use_placement_tool">
                                    <li>
                                        Download and install the <router-link to="/mod">blockexchange</router-link>-mod
                                    </li>
                                    <li>
                                        Create a placement tool with <clipboard-copy :text="'/bx_placer ' + username + ' ' + schema.name"></clipboard-copy>
                                    </li>
                                    <li>
                                        Point and click to place the schematic
                                    </li>
                                </ul>
                            </div>
                            <div class="col-md-6">
                                <u>
                                    <h4>Offline</h4>
                                </u>
                                <div class="btn-group">
                                    <a class="btn btn-outline-success"
                                        :href="BaseURL + '/api/export_bx/' + schema.uid + '/' + schema.name + '.zip'">
                                        <i class="fa fa-download"></i>
                                        Download as BX schematic
                                    </a>
                                    <a class="btn btn-outline-success"
                                        v-bind:class="{disabled: schema.total_parts >= 200}"
                                        :href="BaseURL + '/api/export_we/' + schema.uid + '/' + schema.name + '.we'">
                                        <i class="fa fa-download"></i>
                                        Download as WE schematic
                                        <i class="fa fa-triangle-exclamation"
                                            v-if="schema.total_parts >= 200"
                                            style="color: red;"
                                            title="WE-Import disable due to schematic size"></i>
                                        <i class="fa fa-triangle-exclamation"
                                            v-else-if="schema.total_parts > 50"
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
                                        Check if the placement fits with <clipboard-copy :text="'/bx_allocate_local ' + username + ' ' + schema.name"></clipboard-copy>
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