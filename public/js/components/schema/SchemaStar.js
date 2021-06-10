import { count } from '../../api/star.js';

export default {
	props: ["schema"],
	data: function(){
		return {
			busy: true
		};
	},
	methods: {
		updateStars: function(){
			console.log("updateStars", this);
			count(this.schema.id)
			.then(stars => console.log(stars));
		}
	},
	watch: {
		schema: "updateStars"
	},
	created: function(){
		this.updateStars();
	},
	template: /*html*/`
		<span>
		</span>
	`
};
