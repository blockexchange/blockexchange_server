import SchemaView from '../schema/SchemaView.js';

export default {
	components: {
		"schema-view": SchemaView
	},
	template: /*html*/`
		<div>
			<schema-view :username="$route.params.username" :schemaname="$route.params.schemaname"/>
		</div>
	`
};
