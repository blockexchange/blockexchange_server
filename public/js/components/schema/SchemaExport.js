
export default {
	props: ["schema"],
	template: /*html*/`
	<div>
		<a v-if="schema.total_parts <= 50" class="btn btn-sm btn-primary" :href="'api/export_we/' + schema.id + '/' + schema.name + '.we'">
			<i class="fa fa-download"/> Export as WE Schema
		</a>
		<a class="btn btn-sm btn-primary" :href="'api/export_bx/' + schema.id + '/' + schema.name + '.zip'">
			<i class="fa fa-download"/> Export as BX Schema
		</a>
	</div>
	`
};
