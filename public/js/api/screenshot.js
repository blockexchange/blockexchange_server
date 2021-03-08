
import memoize from '../util/memoize.js';

export const get_by_schemaid = memoize(schema_id => fetch(`api/schema/${schema_id}/screenshot`)
	.then(r => r.json())
);
