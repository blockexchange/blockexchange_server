import SchemaList from "./SchemaList.js";

import { schema_search } from "../api/schema.js";
import { get_tags } from "../service/tags.js";

export default {
    components: {
        "schema-list": SchemaList
    },
    data: function() {
        return {
            list: []
        };
    },
    mounted: function() {
        const tags = get_tags();
        const featured = tags.find(t => t.name == "featured");
        if (!featured) {
            return;
        }

        schema_search({
            tag_uid: featured.uid,
            limit: 6,
            complete: true
        })
        .then(l => this.list = l);
    },
    template: /*html*/`
        <schema-list :list="list" :show_not_found="false">
        </schema-list>
    `
};