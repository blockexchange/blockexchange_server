// collects single key lookups and resolves them with as few batch requests as possible,
// results (also missing ones) are cached per key
export function create_batch_loader<T>(
	fetch_batch: (keys: string[]) => Promise<Map<string, T>>,
	max_batch_size: number,
	delay_ms = 10
): (key: string) => Promise<T | null> {
	const cache = new Map<string, Promise<T | null>>();
	let pending = new Map<string, { resolve: (v: T | null) => void; reject: (e: unknown) => void }>();
	let timer: ReturnType<typeof setTimeout> | null = null;

	const flush = () => {
		timer = null;
		const batch = pending;
		pending = new Map();

		const keys = Array.from(batch.keys());
		for (let i = 0; i < keys.length; i += max_batch_size) {
			const slice = keys.slice(i, i + max_batch_size);
			fetch_batch(slice)
				.then(result => slice.forEach(k => batch.get(k)!.resolve(result.get(k) ?? null)))
				.catch(e => slice.forEach(k => {
					// allow a retry on the next lookup
					cache.delete(k);
					batch.get(k)!.reject(e);
				}));
		}
	};

	return (key: string) => {
		let p = cache.get(key);
		if (!p) {
			p = new Promise((resolve, reject) => pending.set(key, { resolve, reject }));
			cache.set(key, p);
			if (!timer) {
				timer = setTimeout(flush, delay_ms);
			}
		}
		return p;
	};
}
