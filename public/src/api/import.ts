import { fetch_json_checked } from "./util";
import type { ErrorResponse, Schema } from "./types";

export const import_schematic = (buf: ArrayBuffer, filename: string) => fetch_json_checked<Schema, ErrorResponse>(`/api/import/${filename}`, {
	method: "POST",
	body: buf
});
