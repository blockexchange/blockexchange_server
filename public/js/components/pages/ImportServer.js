
import Breadcrumb, { START, SERVER_IMPORT } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, SERVER_IMPORT],
			host: "",
			port: 30000,
			pos1: "",
			pos2: ""
		};
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
							<td>Server-host</td>
							<td>
								<input type="text" class="form-control" v-model="host" placeholder="Server-host"/>
							</td>
						</tr>
						<tr>
							<td>Server-port</td>
							<td>
								<input type="number" class="form-control" v-model="port" placeholder="Server-port"/>
							</td>
						</tr>
						<tr>
							<td>Pos 1</td>
							<td>
								<input type="text" class="form-control" v-model="pos1" placeholder="x,y,z"/>
							</td>
						</tr>
						<tr>
							<td>Pos 2</td>
							<td>
								<input type="text" class="form-control" v-model="pos2" placeholder="x,y,z"/>
							</td>
						</tr>
						<tr>
							<td colspan="2">
								<button class="btn btn-success w-100">
									<i class="fa fa-plus"></i>
									Import
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
