
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
        <div class="container">
            <h5>Used mods</h5>
            <loading-block :fetch_data="fetch_data" v-slot="{ data }">
                <router-link :to="{path: '/search', query: {mod_name: mod.mod_name}}" class="btn btn-secondary m-1" v-for="mod in data.mods">
                    <i class="fa fa-box-archive"></i>
                    {{mod.mod_name}}
                    <span class="badge bg-primary rounded-pill">{{mod.count}}</span>
                </router-link>
            </loading-block>
        </div>
		`
};
