export default {
	props: ["schema"],
	template: /*html*/`
	<ul>
		<li>
			<b>Changed: </b>{{ new Date(+schema.created).toDateString() }}
		</li>
		<li>
			<b>Size: </b>{{ schema.total_size | prettysize }}
		</li>
		<li>
			<b>Dimensions: </b>{{ schema.max_x }} / {{ schema.max_y }} / {{ schema.max_z }} nodes
		</li>
		<li>
			<b>Parts: </b>{{ schema.total_parts }}
		</li>
		<li>
			<b>Downloads: </b>{{ schema.downloads }}
		</li>
		<li>
			<b>License: </b><license-badge style="display: inline;" :license="schema.license"/>
		</li>
	</ul>
	`
};
