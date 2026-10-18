import Breadcrumb, { START, USER_SCHEMAS, SCHEMA_DETAIL, THREE_VIEW } from "../Breadcrumb.js";
import { schema_search } from "../../api/schema.js";
import ThreePreview from "../ThreePreview.js";

export default {
    props: ["username", "name"],
    components: {
        "bread-crumb": Breadcrumb,
        "three-preview": ThreePreview
    },
    data: function() {
        return {
            breadcrumb: [
                START,
                USER_SCHEMAS(this.username),
                SCHEMA_DETAIL(this.username, this.name),
                THREE_VIEW(this.username, this.name)
            ],
            schema_uid: null
        };
    },
    mounted: async function() {
        const list = await schema_search({ schema_name: this.name, user_name: this.username })
        const schema = list[0].schema;
        this.schema_uid = schema.uid;
    },
    template: /*html*/`
        <bread-crumb :items="breadcrumb"/>
        <three-preview style="height: 480px; width: 100%" v-if="schema_uid" :schema_uid="schema_uid"/>
    `
};
