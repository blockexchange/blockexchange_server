import { MAPBLOCK_SIZE, MAPBLOCK_VOLUME, ParsedMapblock } from "@buckaroobanzay/luanti-viewer";
import type { MapblockPosition, NodeDefinition, ViewerOptions } from "@buckaroobanzay/luanti-viewer";

import { get_mediafiles_batch, get_nodedefinitions_batch } from "../api/media";
import { get_schemapart_chunk, SCHEMAPART_CHUNK_SIZE } from "../api/schemapart";
import type { Schema, SchemaPart } from "../api/types";
import { create_batch_loader } from "../util/batch_loader";

// max names per batch request, media files are bigger so fewer of them per request
const NODEDEF_BATCH_SIZE = 1000;
const MEDIA_BATCH_SIZE = 100;

// schematics bigger than this (in mapblocks) are only loaded around the camera
const VIEW_RANGE_THRESHOLD = 12;
const VIEW_RANGE = 12;

// schemapart node ids are stored as uint16 big endian with this offset
const NODE_ID_OFFSET = 32768;

// param1 of full sunlight: day light 15 (low nibble), night light 0 (high nibble)
const SUNLIGHT = 0x0f;

// built-in nodes without a stored definition, "ignore" would be rendered as an unknown node otherwise
const airlike_definition = (name: string): NodeDefinition => ({
	name: name,
	drawtype: "airlike",
	paramtype: "light",
	paramtype2: "none",
	light_source: 0,
	walkable: false,
	groups: null
});

const BUILTIN_NODES = new Map<string, NodeDefinition>(
	["air", "ignore"].map(name => [name, airlike_definition(name)])
);

interface SchemaPartMetadata {
	node_mapping: Record<string, number>;
	size: { x: number; y: number; z: number };
}

const base64_to_bytes = (b64: string) => Uint8Array.from(atob(b64), c => c.charCodeAt(0));

async function inflate(b64: string): Promise<Uint8Array> {
	const stream = new Blob([base64_to_bytes(b64)]).stream().pipeThrough(new DecompressionStream("deflate"));
	return new Uint8Array(await new Response(stream).arrayBuffer());
}

/**
 * converts a schematic part into a mapblock of the viewer
 * schemapart: x-major order (index = z + y * size_z + x * size_y * size_z), may be smaller than 16³ at the edges
 * mapblock: z-major order (index = z * 256 + y * 16 + x), always 16³, the missing nodes are filled with air
 * light: the schematics store no light (param1 = 0) for air, the faces next to it would be black, so built-in
 * airlike nodes without light get full sunlight
 */
export async function convert_schemapart(part: SchemaPart): Promise<ParsedMapblock> {
	const [raw_meta, data] = await Promise.all([inflate(part.metadata), inflate(part.data)]);
	const meta = JSON.parse(new TextDecoder().decode(raw_meta)) as SchemaPartMetadata;

	const sx = Math.min(meta.size.x, MAPBLOCK_SIZE);
	const sy = Math.min(meta.size.y, MAPBLOCK_SIZE);
	const sz = Math.min(meta.size.z, MAPBLOCK_SIZE);
	const count = meta.size.x * meta.size.y * meta.size.z;
	if (data.length < count * 4) {
		throw new Error(`schemapart data too short: ${data.length} bytes, expected ${count * 4}`);
	}

	// content ids are signed in the schemapart, the viewer expects unsigned ones
	const names = new Map<number, string>();
	for (const [name, id] of Object.entries(meta.node_mapping)) {
		names.set(id & 0xffff, name);
	}

	// id used for the padding outside of the schemapart
	let air_id = meta.node_mapping["air"];
	if (air_id === undefined) {
		air_id = 0;
		while (names.has(air_id)) {
			air_id++;
		}
		names.set(air_id, "air");
	}

	// ids of the built-in airlike nodes ("air", "ignore")
	const airlike_ids = new Set<number>();
	for (const [id, name] of names) {
		if (BUILTIN_NODES.has(name)) {
			airlike_ids.add(id);
		}
	}

	const node_ids = new Uint16Array(MAPBLOCK_VOLUME).fill(air_id & 0xffff);
	const param1 = new Uint8Array(MAPBLOCK_VOLUME).fill(SUNLIGHT);
	const param2 = new Uint8Array(MAPBLOCK_VOLUME);

	for (let x = 0; x < sx; x++) {
		for (let y = 0; y < sy; y++) {
			for (let z = 0; z < sz; z++) {
				const src = z + (y * meta.size.z) + (x * meta.size.y * meta.size.z);
				const dst = (z * MAPBLOCK_SIZE * MAPBLOCK_SIZE) + (y * MAPBLOCK_SIZE) + x;
				const id = (((data[src * 2]! << 8) | data[(src * 2) + 1]!) - NODE_ID_OFFSET) & 0xffff;
				const light = data[(count * 2) + src]!;
				node_ids[dst] = id;
				param1[dst] = (light == 0 && airlike_ids.has(id)) ? SUNLIGHT : light;
				param2[dst] = data[(count * 3) + src]!;
			}
		}
	}

	return new ParsedMapblock(node_ids, param1, param2, names);
}

const pos_key = (x: number, y: number, z: number) => `${x},${y},${z}`;

const media_type = (filename: string) => {
	const ext = filename.split(".").pop()?.toLowerCase();
	switch (ext) {
	case "png": return "image/png";
	case "jpg":
	case "jpeg": return "image/jpeg";
	case "tga": return "image/x-tga";
	default: return "application/octet-stream";
	}
};

// node definitions are shared between all views
const get_nodedefinition = create_batch_loader<NodeDefinition>(async names => {
	const result = new Map<string, NodeDefinition>();
	for (const nd of await get_nodedefinitions_batch(names)) {
		try {
			const def = JSON.parse(nd.definition) as NodeDefinition;
			result.set(nd.name, { ...def, name: def.name ?? nd.name });
		} catch (e) {
			console.warn(`invalid node definition '${nd.name}'`, e);
		}
	}
	return result;
}, NODEDEF_BATCH_SIZE);

// media files are shared between all views
const get_media = create_batch_loader<Blob>(async filenames => {
	const result = new Map<string, Blob>();
	for (const mf of await get_mediafiles_batch(filenames)) {
		result.set(mf.name, new Blob([base64_to_bytes(mf.data)], { type: media_type(mf.name) }));
	}
	return result;
}, MEDIA_BATCH_SIZE);

// mapblock bounds of the schematic, the schematic origin is at mapblock (0,0,0)
export const get_schema_bounds = (schema: Schema) => ({
	min: { x: 0, y: 0, z: 0 },
	max: {
		x: Math.max(Math.ceil(schema.size_x / MAPBLOCK_SIZE) - 1, 0),
		y: Math.max(Math.ceil(schema.size_y / MAPBLOCK_SIZE) - 1, 0),
		z: Math.max(Math.ceil(schema.size_z / MAPBLOCK_SIZE) - 1, 0)
	}
});

// options for "createView" that load the schematic parts in chunks and the node definitions and media in batches
export function create_schema_viewer_options(schema: Schema): ViewerOptions {
	// raw parts by node offset, kept to reload mapblocks that came into range again
	const parts = new Map<string, SchemaPart>();
	// pending or finished chunk requests by chunk origin
	const chunks = new Map<string, Promise<void>>();

	const load_chunk = (x: number, y: number, z: number) => {
		const key = pos_key(x, y, z);
		let p = chunks.get(key);
		if (!p) {
			p = get_schemapart_chunk(schema.uid, x, y, z)
				.then(list => list.forEach(part => parts.set(pos_key(part.offset_x, part.offset_y, part.offset_z), part)))
				.catch(e => {
					// retry on the next request
					chunks.delete(key);
					throw e;
				});
			chunks.set(key, p);
		}
		return p;
	};

	const get_mapblock = async (pos: MapblockPosition) => {
		const ox = pos.x * MAPBLOCK_SIZE, oy = pos.y * MAPBLOCK_SIZE, oz = pos.z * MAPBLOCK_SIZE;
		const chunk_origin = (v: number) => Math.floor(v / SCHEMAPART_CHUNK_SIZE) * SCHEMAPART_CHUNK_SIZE;
		await load_chunk(chunk_origin(ox), chunk_origin(oy), chunk_origin(oz));

		const part = parts.get(pos_key(ox, oy, oz));
		return part ? convert_schemapart(part) : null;
	};

	const bounds = get_schema_bounds(schema);
	const max_extent = Math.max(bounds.max.x, bounds.max.y, bounds.max.z) + 1;

	return {
		getNodeDefinitions: async name => BUILTIN_NODES.get(name) ?? get_nodedefinition(name),
		getMapblock: get_mapblock,
		getMedia: get_media,
		minPos: bounds.min,
		maxPos: bounds.max,
		viewRange: max_extent > VIEW_RANGE_THRESHOLD ? VIEW_RANGE : null
	};
}
