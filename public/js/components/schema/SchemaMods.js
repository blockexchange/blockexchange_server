export default {
	props: ["schema"],
	template: /*html*/`
	<ul>
		<li v-for="mod in schema.mods">
			<span class="badge bg-primary">
				{{ mod }}
			</span>
		</li>
	</ul>
	`
};
