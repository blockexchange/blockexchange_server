export default {
	props: ["schema"],
	template: /*html*/`
	<ul>
		<li v-for="mod in Object.keys(schema.mods)">
			<span class="badge badge-primary">
				{{ mod }}
			</span>
			<span class="badge badge-success">
				{{ schema.mods[mod] }}
			</span>
		</li>
	</ul>
	`
};
