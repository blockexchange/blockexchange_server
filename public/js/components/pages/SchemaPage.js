import SchemaView from '../schema/SchemaView.js';

export default {
	components: {
		"schema-view": SchemaView
	},
	template: /*html*/`
		<div>
			<schema-view :user_name="$route.params.user_name" :schema_name="$route.params.schema_name"/>
		</div>
	`
};
