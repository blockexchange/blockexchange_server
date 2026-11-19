import { ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import debounce from './debounce';

// two-way binds a text filter to the "q" query param, resets paging on change
export function use_query_filter(wait = 300) {
	const route = useRoute();
	const router = useRouter();

	const get_query = () => typeof route.query.q == "string" ? route.query.q : "";
	const filter = ref(get_query());

	const apply = debounce((value: string) => {
		const query = { ...route.query };
		delete query.page;
		if (value) {
			query.q = value;
		} else {
			delete query.q;
		}
		router.replace({ path: route.path, query });
	}, wait);

	watch(filter, v => {
		if (v.trim() != get_query()) {
			apply(v.trim());
		}
	});
	// back/forward navigation
	watch(() => route.query.q, () => {
		if (filter.value.trim() != get_query()) {
			filter.value = get_query();
		}
	});

	return { filter, get_query };
}
