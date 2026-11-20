import type { SchemaPart } from "./types";

// size of the cube covered by a chunk request in nodes
export const SCHEMAPART_CHUNK_SIZE = 64;

// returns all parts with an offset between (x,y,z) (inclusive) and (x,y,z) + 64 (exclusive)
export const get_schemapart_chunk = async (schema_uid: string, x: number, y: number, z: number): Promise<SchemaPart[]> => {
	const r = await fetch(`/api/schemapart_chunk/${schema_uid}/${x}/${y}/${z}`);
	if (r.status == 204) {
		return [];
	}
	if (r.status != 200) {
		throw new Error(`schemapart chunk request failed with status ${r.status}`);
	}
	return (await r.json() as SchemaPart[] | null) ?? [];
};
