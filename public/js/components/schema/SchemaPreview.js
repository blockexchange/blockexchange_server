export default {
	props: ["screenshots", "schema", "version"],
	template: /*html*/`
	<div class="text-center" v-for="screenshot in screenshots">
		<img class="img-fluid" :src="'api/schema/' + schema.id + '/screenshot/' + screenshot.id + '?version=' + version"/>
	</div>
	`
};
