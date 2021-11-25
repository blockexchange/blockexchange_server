export default {
	props: ["schema"],
	template: /*html*/`
	<div class="row">
		<div class="col-md-6">
			<div class="card">
				<div class="card-header">
					Online
				</div>
				<div class="card-body">
					<pre>
/bx_pos1
/bx_load {{ schema.user.name }} {{ schema.name }}
</pre>
				</div>
			</div>
		</div>
		<div class="col-md-6">
			<div class="card">
				<div class="card-header">
					Offline
				</div>
				<div class="card-body">
					<a v-if="schema.total_parts <= 50" class="btn btn-sm btn-primary"
						:href="'api/export_we/' + schema.user.name + '/' + schema.name + '/' + schema.name + '.we'">
						<i class="fa fa-download"/> Export as WE Schema
					</a>
					<a class="btn btn-sm btn-primary"
						:href="'api/export_bx/' + schema.user.name + '/' + schema.name + '/' + schema.name + '.zip'">
						<i class="fa fa-download"/> Export as BX Schema
					</a>
				</div>
			</div>
		</div>
	</div>
	`
};
