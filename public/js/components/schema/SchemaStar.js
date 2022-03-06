import { get, remove, create } from '../../api/star.js';
import loginStore from '../../store/login.js';

export default {
	props: ["schema"],
	data: function(){
		return {
			stars: null
		};
	},
	methods: {
		updateStars: function(){
			let user_id;
			if (loginStore.claims){
				user_id = loginStore.claims.user_id;
			}

			get(this.schema.id, user_id)
			.then(stars => this.stars = stars);
		},
		removeStar: function(){
			if (loginStore.claims) {
				remove(this.schema.id).then(() => this.updateStars());
			}
		},
		addStar: function(){
			if (loginStore.claims) {
				create(this.schema.id).then(() => this.updateStars());
			}
		}
	},
	watch: {
		schema: "updateStars"
	},
	created: function(){
		this.updateStars();
	},
	template: /*html*/`
		<span v-if="stars">
			<i v-if="stars.starred" class="fa-solid fa-star" style="color: yellow;" v-on:click="removeStar()"></i>
			<i v-else class="fa-regular fa-star" v-on:click="addStar()"></i>
			<span v-if="stars.count > 0" class="badge bg-secondary rounded-pill">{{ stars.count }}</span>
		</span>
	`
};
