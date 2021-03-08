
module.exports = class {
	constructor(){
		// array of numbers/bytes
		this.parts = [];
	}

	add(o){
		if (typeof(o) == "string"){
			// string type
			const buf = Buffer.from(o);
			for (let i=0; i<buf.length; i++){
				if (buf[i] == 0x0A){
					// insert backslash before
					this.parts.push(0x5C);
				} else if (buf[i] == 0x1B){
					// translation marker, next bytes: 0x28 (start) or 0x45 (end)
				}
				this.parts.push(buf[i]);
			}
		} else if (typeof(o) == "number"){
			this.parts.push(o);
		}
	}

	toBuffer(){
		return Buffer.from(this.parts);
	}
};
