import Breadcrumb, { START, SCHEMA_IMPORT } from "../Breadcrumb.js";

import { import_schematic } from "../../api/import.js";
import { get_username } from "../../service/login.js";

const store = Vue.reactive({
	breadcrumb: [START, SCHEMA_IMPORT],
	results: []
});

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: () => store,
	methods: {
		import_file: function(file) {
			const res = Vue.reactive({
				filename: file.name,
				busy: true,
				schema: null,
				err: null
			});
			this.results.push(res);

			file.arrayBuffer()
			.then(buf => import_schematic(buf, file.name))
			.then(s => res.schema = s)
			.catch(e => res.error = e)
			.finally(() => res.busy = false);
		},
		on_upload: function() {
			this.results = [];
			const files = this.$refs.upload.files;
			for (let i=0; i<files.length; i++) {
				this.import_file(files[i]);	
			}
		}
	},
	computed: {
		username: get_username
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
			<div class="col-md-4"></div>
			<div class="col-md-4">
				<h5>
					Upload Worldedit or Blockexchange schematics
				</h5>
				<input ref="upload" type="file" class="form-control" v-on:change="on_upload" multiple accept=".we,.zip"/>
				<router-link to="/import-server">
					<i class="fa fa-server"></i>
					Import schematic from a public server
				</router-link>
				<hr>
				<ul>
					<li v-for="res in results">
						{{res.filename}}
						<i class="fa fa-spinner fa-spin" v-if="res.busy"></i>
						<i class="fa fa-check" v-if="res.schema"></i>
						<i class="fa fa-times" v-if="res.error"></i>
						<span class="badge bg-danger" v-if="res.error && res.error.message">
							{{res.error.message}}
						</span>
						<router-link :to="'/schema/' + username + '/' + res.schema.name" v-if="res.schema">
							{{res.schema.name}}
						</router-link>
					</li>
				</ul>
			</div>
			<div class="col-md-4"></div>
		</div>
	`
};
