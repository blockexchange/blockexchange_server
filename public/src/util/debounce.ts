// https://davidwalsh.name/javascript-debounce-function
export default function debounce<A extends unknown[]>(func: (...args: A) => void, wait: number, immediate?: boolean) {
	let timeout: ReturnType<typeof setTimeout> | undefined;
	return function(...args: A) {
		const later = () => {
			timeout = undefined;
			if (!immediate) func(...args);
		};
		const callNow = immediate && !timeout;
		clearTimeout(timeout);
		timeout = setTimeout(later, wait);
		if (callNow) func(...args);
	};
}
