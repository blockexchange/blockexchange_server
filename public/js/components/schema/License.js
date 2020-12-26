export default {
	props: ["license"],
	computed: {
		imgsrc: function(){
			switch (this.license){
				case "CC0": return "pics/license_cc0.png";
			}
		}
	},
	template: /*html*/`
		<div>
			<img v-if="imgsrc" :src="imgsrc"/>
			<div v-else class="badge badge-primary">
				{{ license }}
			</div>
		</div>
	`
};
