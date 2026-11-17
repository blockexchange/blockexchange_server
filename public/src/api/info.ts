import { fetch_json } from "./util";
import type { Info } from "./types";

export const get_info = () => fetch_json<Info>(`/api/info`);
