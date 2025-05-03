
export function validate_pos_string(str) {
	const re = /^-?[0-9]+$/;
	const parts = str.split(",");
	if (parts.length != 3) {
		return false;
	}
	return (re.test(parts[0]) && re.test(parts[1]) && re.test(parts[2]));
}

export function parse_pos_string(str) {
    if (!validate_pos_string(str)) {
        return;
    }

	const parts = str.split(",");
    return {
        x: +parts[0],
        y: +parts[1],
        z: +parts[2]
    };
}

export function sort_pos(pos1, pos2) {
    return [{
        x: Math.min(pos1.x, pos2.x),
        y: Math.min(pos1.y, pos2.y),
        z: Math.min(pos1.z, pos2.z)
    }, {
        x: Math.max(pos1.x, pos2.x),
        y: Math.max(pos1.y, pos2.y),
        z: Math.max(pos1.z, pos2.z)
    }];
}