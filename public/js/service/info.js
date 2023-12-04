
import { get_info } from "../api/info.js";

const store = Vue.reactive({});

export const fetch_info = () => get_info().then(i => Object.assign(store, i));

export const get_base_url = () => store.base_url;
export const get_github_id = () => store.oauth.github_id;
export const get_discord_id = () => store.oauth.discord_id;
export const get_mesehub_id = () => store.oauth.mesehub_id;
