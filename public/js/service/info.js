
import { get_info } from "../api/info.js";

const store = Vue.reactive({});

export const fetch_info = () => get_info().then(i => Object.assign(store, i));

export const get_base_url = () => store.base_url;
export const get_github_login = () => store.oauth_login.github;
export const get_cdb_login = () => store.oauth_login.cdb;
export const get_discord_login = () => store.oauth_login.discord;
export const get_mesehub_login = () => store.oauth_login.mesehub;
export const get_codeberg_login = () => store.oauth_login.codeberg;
export const get_stats = () => store.stats;