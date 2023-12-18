export default {
    props: ["text"],
    data: function() {
        return {
            copied: false
        };
    },
    methods: {
        copy: function() {
            this.copied = true;
            navigator.clipboard.writeText(this.text);
        }
    },
    template: /*html*/`
    <span class="badge bg-secondary" v-on:click="copy">
        {{text}}
        &nbsp;
        <i class="fa fa-copy" v-bind:style="{'color':copied?'green':null}"></i>
    </span>
    `
};