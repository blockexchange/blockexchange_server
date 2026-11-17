const ranges = [{
	size: 1000*1000,
	suffix: "M"
},{
	size: 1000,
	suffix: "K"
}];

export default function format_count(value: number | string): string {
	const count = +value;
	for (const range of ranges) {
		if (count > range.size){
			return (Math.floor(count / range.size * 100) / 100) + " " + range.suffix;
		}
	}

	return `${count}`;
}
