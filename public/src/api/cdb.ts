import { fetch_json_checked } from "./util";
import type { ErrorResponse } from "./types";

export interface CdbVcsPackage {
	author: string;
	name: string;
	title: string;
	type: "mod" | "game" | "txp";
	repo: string | null;
}

let packages: Promise<CdbVcsPackage[]> | null = null;

// the server caches the list for 24 hours, fetch it only once per page load
export const get_cdb_vcs_packages = () => {
	if (!packages) {
		packages = fetch_json_checked<CdbVcsPackage[], ErrorResponse>("/api/cdb/vcs-packages").catch(e => {
			packages = null; // allow a retry on the next focus
			throw e;
		});
	}
	return packages;
};
