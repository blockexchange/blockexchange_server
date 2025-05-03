import ClipboardCopy from "../ClipboardCopy.js";

const max_placement_tool_size = 100;

export default {
    props: ["schema", "username"],
    components: {
        "clipboard-copy": ClipboardCopy
    },
    data: function() {
        return {
            BaseURL
        };
    },
    computed: {
        can_use_placement_tool: function() {
            return (this.schema.size_x <= max_placement_tool_size &&
                    this.schema.size_y <= max_placement_tool_size &&
                    this.schema.size_z <= max_placement_tool_size);
        }
    },
    template: /*html*/`
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
                        v-if="schema.total_parts >= 500"
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
                    Check if the placement fits with <clipboard-copy :text="'/bx_allocate_local ' + schema.name"></clipboard-copy>
                </li>
                <li>
                    Place the schematic with <clipboard-copy :text="'/bx_load_local ' + schema.name"></clipboard-copy>
                </li>
            </ul>
        </div>
    </div>
    `
};