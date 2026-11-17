import { fetch_json } from "./util";
import type { SchemaModCount } from "./types";

export const get_schemamod_count = () => fetch_json<SchemaModCount[]>(`/api/schemamod/count`);
