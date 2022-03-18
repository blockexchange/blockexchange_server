
export default {
    props: ["link"],
    data: function() {
        return {
            success: false,
            failure: false
        };
    },
    methods: {
		copyLink: function(){
            navigator.clipboard.writeText(this.link)
            .then(() => this.success = true, () => this.failure = true);
		}
    },
    computed: {
        style: function() {
            if (this.success) {
                return {color: "green"};
            } else if (this.failure) {
                return {color: "red"};
            }
        }
    },
    template: /*html*/`
        <i class="fa-solid fa-copy" v-bind:style="style" v-on:click="copyLink()"></i>
    `
};