import format_size from "../../util/format_size.js";
import format_time from "../../util/format_time.js";

import ModalPrompt from "../ModalPrompt.js";

import { schema_update, schema_set_tags, schema_update_screenshot, schema_delete } from "../../api/schema.js";
import { get_tags } from "../../service/tags.js";

export default {
    components: {
        "modal-prompt": ModalPrompt
    },
    props: {
        schema: { type: Object, required: true },
        username: { type: String, required: true },
        allow_edit: { type: Boolean, default: false }
    },
    data: function() {
        return {
            BaseURL,
            edit_mode: false,
            error_response: null,
            screenshot_busy: false,
            delete_prompt: false
        };
    },
    methods: {
        format_size,
        format_time,
        get_tags,
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
            .then(() => this.$router.push(`/`));
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
                    <button class="btn btn-outline-primary">
                        <i class="fa fa-star" v-bind:style="{color: 'yellow'}"></i>
                        <span class="badge bg-secondary rouded-pill">{{schema.stars}}</span>
                        Star
                    </button>
                </h3>
                <div v-else class="input-group has-validation">
                    <input type="text" class="form-control" v-model="schema.name" v-bind:class="{'is-invalid': error_response}">
                    <div class="invalid-feedback" v-if="error_response && error_response.name_taken">
                        Schema name already taken
                    </div>
                    <div class="invalid-feedback" v-if="error_response && error_response.name_invalid">
                        Schema name invalid, allowed chars: a to z, A to Z, 0 to 9 and -, _
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
                    <div class="card-body">
                        <h5 class="card-title">Details</h5>
                        <ul>
                            <li><b>Created:</b> {{format_time(schema.created / 1000)}}</li>
                            <li><b>Modified:</b> {{format_time(schema.mtime / 1000)}}</li>
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
                                    <option value="Proprietary">Proprietary</option>
                                </select>
                            </li>
                        </ul>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">Description</h5>
                        <pre v-if="!edit_mode">{{schema.description}}</pre>
                        <textarea v-else v-model="schema.description" rows="10" class="form-control"></textarea>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">Used mods</h5>
                        <details>
                            <summary>Click to open</summary>
                            <ul>
                                <li v-for="mod in schema.mods">
                                    <span class="badge bg-primary">{{mod}}</span>
                                </li>
                            </ul>
                        </details>
                    </div>
                </div>
                <br>
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">Tags</h5>
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
                    <div class="card-body">
                        <h5 class="card-title">Preview</h5>
                        <div class="text-center">
                            <img :src="BaseURL + '/api/schema/' + schema.id + '/screenshot'" class="img-fluid">
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
};