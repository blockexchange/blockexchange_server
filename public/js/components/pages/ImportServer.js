
import Breadcrumb, { START, SERVER_IMPORT } from "../Breadcrumb.js";
import { schema_create } from "../../api/schema.js";
import { validate_pos_string, parse_pos_string, sort_pos } from "../../util/pos.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, SERVER_IMPORT],
			name: "",
			host: "",
			port: 30000,
			pos1: "0,0,0",
			pos2: "10,10,10",
			error_message: "",
		};
	},
	methods: {
		create: async function() {
			this.error_message = "";
			try {
				let pos1 = parse_pos_string(this.pos1);
				let pos2 = parse_pos_string(this.pos2);
				[pos1, pos2] = sort_pos(pos1, pos2);

				const result = await schema_create({
					name: this.name,
					size_x: pos2.x - pos1.x + 1,
					size_y: pos2.y - pos1.y + 1,
					size_z: pos2.z - pos1.z + 1
				});
			} catch (e) {
				this.error_message = e.message;
			}
		}
	},
	computed: {
		ready: function() {
			return (this.name && this.host && this.port && this.pos1 && this.pos2);
		},
		error_pos1: function() {
			return !validate_pos_string(this.pos1);
		},
		error_pos2: function() {
			return !validate_pos_string(this.pos2);
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
        <div class="row">
			<div class="col-md-2"></div>
			<div class="col-md-8">
				<h5>Import schematic from a public server</h5>
				<div class="alert alert-info">
					<i class="fa fa-info"></i>
					<b>Note:</b> TODO
				</div>
				<table class="table table-dark table-striped table-condensed">
					<tbody>
						<tr>
							<td>Schematic name</td>
							<td>
								<input type="text" class="form-control" v-model="name" placeholder="Schematic name"/>
							</td>
						</tr>
						<tr>
							<td>Server host</td>
							<td>
								<input type="text" class="form-control" v-model="host" placeholder="Server-host"/>
							</td>
						</tr>
						<tr>
							<td>Server port</td>
							<td>
								<input type="number" class="form-control" v-model="port" placeholder="Server-port"/>
							</td>
						</tr>
						<tr>
							<td>Pos 1</td>
							<td>
								<div class="input-group has-validation">
									<input type="text" class="form-control" v-model="pos1" placeholder="x,y,z" v-bind:class="{'is-invalid': error_pos1}"/>
									<div class="invalid-feedback" v-if="error_pos1">
										Invalid position, expected format: "x,y,z"
									</div>
								</div>
							</td>
						</tr>
						<tr>
							<td>Pos 2</td>
							<td>
								<div class="input-group has-validation">
									<input type="text" class="form-control" v-model="pos2" placeholder="x,y,z" v-bind:class="{'is-invalid': error_pos2}"/>
									<div class="invalid-feedback" v-if="error_pos2">
										Invalid position, expected format: "x,y,z"
									</div>
								</div>
							</td>
						</tr>
						<tr>
							<td colspan="2">
								<button class="btn btn-success w-100" v-on:click="create" :disabled="!ready">
									<i class="fa fa-plus"></i>
									Import
									<span class="badge bg-danger">
										{{error_message}}
									</span>
								</button>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
			<div class="col-md-2"></div>
        </div>
		`
};
