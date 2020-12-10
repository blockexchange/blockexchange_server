export default {
	props: ["username", "schemaname"],
	template: /*html*/`
		<div>
			Schema view {{ username }} / {{ schemaname }}
		</div>
	`
};
