import { reactive } from 'vue';
import { get_tags as api_get_tags } from "../api/tags";
import type { Tag } from "../api/types";

const store = reactive<{ list: Tag[] }>({
	list: []
});

export const fetch_tags = async () => {
	store.list = await api_get_tags();
};
export const get_tags = () => store.list;
