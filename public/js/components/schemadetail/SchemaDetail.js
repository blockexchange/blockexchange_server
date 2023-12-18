import format_size from "../../util/format_size.js";
import format_time from "../../util/format_time.js";

import { schema_update } from "../../api/schema.js";

export default {
    props: {
        schema: { type: Object, required: true },
        username: { type: String, required: true },
        allow_edit: { type: Boolean, default: false }
    },
    data: function() {
        return {
            BaseURL,
            edit_mode: false
        };
    },
    methods: {
        format_size,
        format_time,
        save: function() {
            schema_update(this.schema)
            .then(() => {
                // ok
            })
            .catch(e => console.error(e))
            .finally(() => this.edit_mode = false);
        }
    },
    template: /*html*/`
    <div>
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
                <input v-else type="text" class="form-control" v-model="schema.name">
            </div>
            <div class="col-md-6">
                <div class="btn-group float-end" v-if="allow_edit && !edit_mode">
                    <a class="btn btn-sm btn-secondary" v-on:click="edit_mode = true">
                        <i class="fa fa-edit"></i> Edit
                    </a>
                    <a class="btn btn-sm btn-secondary">
                        <i class="fa fa-image"></i> Update screenshot
                    </a>
                    <a class="btn btn-sm btn-danger">
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
                        <span class="badge bg-success" v-for="tag in schema.tags" style="margin-right: 5px;">
                            <i class="fas fa-tag"></i>
                            {{tag}}
                        </span>
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