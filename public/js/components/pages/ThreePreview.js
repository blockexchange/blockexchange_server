import Breadcrumb, { START, USER_SCHEMAS, SCHEMA_DETAIL, THREE_VIEW } from "../Breadcrumb.js";
import { schema_search } from "../../api/schema.js";
import { get_schemapart_by_ofset } from "../../api/schema_part.js";
import { decompressSync } from 'fflate';

function base64ToArrayBuffer(base64) {
    var binaryString = atob(base64);
    var bytes = new Uint8Array(binaryString.length);
    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
}

function ab2str(buf) {
  return String.fromCharCode.apply(null, new Uint16Array(buf));
}

export default {
    props: ["username", "name"],
    components: {
        "bread-crumb": Breadcrumb
    },
    data: function() {
        return {
            breadcrumb: [
                START,
                USER_SCHEMAS(this.username),
                SCHEMA_DETAIL(this.username, this.name),
                THREE_VIEW(this.username, this.name)
            ]
        };
    },
    mounted: async function() {
        const list = await schema_search({ schema_name: this.name, user_name: this.username })
        const schema = list[0].schema;
        console.log(schema);

        const part = await get_schemapart_by_ofset(schema.uid, 0, 0, 0);
        console.log(part);

        const buf = base64ToArrayBuffer(part.metadata);
        const ab = decompressSync(new Uint8Array(buf));
        console.log(ab);

        console.log(JSON.parse(ab2str(ab)));

        const mapbuf = base64ToArrayBuffer(part.data);
        const mapdata = decompressSync(new Uint8Array(mapbuf));

        console.log(mapdata);
    },
    template: /*html*/`
        <bread-crumb :items="breadcrumb"/>
    `
};
