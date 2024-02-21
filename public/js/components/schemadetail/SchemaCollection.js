import { get_collections_by_username } from "../../api/collection.js";

export default {
    props: ["schema", "username", "edit_mode", "collection_name"],
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
            <router-link :to="'/collections/' + username + '/' + collection_name" class="btn btn-success w-100" v-if="collection_name">
                <i class="fa fa-object-group"></i>
                Part of the <span class="badge bg-primary">{{collection_name}}</span> collection
            </router-link>
        </div>
        <div v-else>
            <select class="form-control" v-model="schema.collection_uid">
                <option :value="null"></option>
                <option v-for="c in collections" :value="c.uid">{{c.name}}</option>
            </select>
        </div>
    `
};