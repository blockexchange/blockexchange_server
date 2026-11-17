import { reactive } from 'vue';

// search form state, survives navigating away from the search page
export default reactive({
	count: 0,
	busy: false,
	keywords: "",
	tag_name: "",
	mod_name: ""
});
