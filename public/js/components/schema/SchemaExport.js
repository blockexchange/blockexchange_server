
export default {
	props: ["schema"],
	template: /*html*/`
		<a class="btn btn-sm btn-primary" :href="'api/export/' + schema.id + '/' + schema.name + '.we'">
			<i class="fa fa-download"/> Export as WE Schema
		</a>
	`
};
