export default {
	props: ["screenshots", "schema"],
	template: /*html*/`
	<div v-for="screenshot in screenshots">
		<img class="img-fluid" :src="'api/schema/' + schema.id + '/screenshot/' + screenshot.id"/>
	</div>
	`
};
