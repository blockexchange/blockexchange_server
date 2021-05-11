
export default {
	props: ["schema"],
	template: /*html*/`
		<a v-if="schema.total_parts <= 50" class="btn btn-sm btn-primary" :href="'api/export_we/' + schema.id + '/' + schema.name + '.we'">
			<i class="fa fa-download"/> Export as WE Schema
		</a>
	`
};
