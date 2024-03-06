
import { get_schemamod_count } from "../../api/schemamod.js";
import Breadcrumb, { START, MODS } from "../Breadcrumb.js";
import LoadingBlock from "../LoadingBlock.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
        "loading-block": LoadingBlock
	},
	data: function() {
		return {
			breadcrumb: [START, MODS]
		};
	},
    methods: {
        fetch_data: function() {
            return {
                mods: get_schemamod_count()
            };
        }
    },
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
        <loading-block :fetch_data="fetch_data" v-slot="{ data }">
            <div class="container">
                <button class="btn btn-secondary m-1" :disabled="true" v-for="mod in data.mods">
                    {{mod.mod_name}}
                    <span class="badge bg-primary rounded-pill">{{mod.count}}</span>
                </button>
            </div>
        </loading-block>
		`
};
