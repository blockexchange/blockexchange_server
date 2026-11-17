import { fetch_json } from "./util";
import type { Tag } from "./types";

export const get_tags = () => fetch_json<Tag[]>(`/api/tags`);
