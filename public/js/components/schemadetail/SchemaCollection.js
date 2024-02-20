import { get_collections_by_username } from "../../api/collection.js";

export default {
    props: ["schema", "username", "edit_mode"],
    data: function() {
        return {
            collections: []
        };
    },
    mounted: function() {
        get_collections_by_username(this.username).then(c => this.collections = c);
    },
    template: /*html*/`
        <div v-if="!edit_mode">
            <span class="badge bg-success" v-if="schema.collection_uid" style="margin-right: 5px;">
                <i class="fa fa-object-group"></i>
                {{schema.collection_uid}}
            </span>
        </div>
        <div v-else>
            <select class="form-control" v-model="schema.collection_uid">
                <option value=""></option>
                <option v-for="c in collections" :value="c.uid">{{c.name}}</option>
            </select>
        </div>
    `
};