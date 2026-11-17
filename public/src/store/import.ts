import { reactive } from 'vue';
import type { ErrorResponse, Schema } from '../api/types';

export interface ImportResult {
	filename: string;
	busy: boolean;
	schema: Schema | null;
	error: ErrorResponse | null;
}

// import results, survive navigating away from the import page
export default reactive<{ results: ImportResult[] }>({
	results: []
});
