import Breadcrumb, { START, USER_SCHEMAS, SCHEMA_DETAIL, THREE_VIEW } from "../Breadcrumb.js";

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
    template: /*html*/`
        <bread-crumb :items="breadcrumb"/>
    `
};
