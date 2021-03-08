
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
