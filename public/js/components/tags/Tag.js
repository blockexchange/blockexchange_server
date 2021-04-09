
import { remove } from '../../api/schematag.js';
import TagStore from '../../store/tag.js';
import loginStore from '../../store/login.js';

export default {
    props: ["tag", "tag_id", "schema_id", "user_id"],
    created: function(){
        if (!this.tag && this.tag_id){
            //search for the actual tag in the store
            this.tag = TagStore.tags.find(t => t.id == this.tag_id);
        }

        // set name
        this.name = this.tag.name;

        if (loginStore.claims && loginStore.claims.user_id == this.user_id){
            // the user owns the schema
            this.can_delete = true;
        }
    },
    methods: {
        remove: function(){
            //TODO
            console.log("delete tag: ", this.tag);
            remove(this.schema_id, this.tag_id);
            this.$emit("removed", this.id);
        }
    },
    data: function() {
        return {
            can_delete: false,
            name: ""
        };
    },
    template: `
        <span class="badge badge-success">
            <i class="fas fa-tag"></i>
            {{ name }}
            <i v-if="can_delete" class="fa fa-times" v-on:click="remove"></i>
        </span>
    `
};