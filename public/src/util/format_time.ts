function pad(n: number) {
	return n<10 ? `0${n}` : `${n}`;
}

export default function format_time(ts: number): string {
	const d = new Date(ts);
	return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}` +
		` ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
}
