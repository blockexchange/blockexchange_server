
export default {
    props: ["title"],
    emits: ["close"],
    template: /*html*/`
    <div class="modal-backdrop show"></div>
    <div class="modal show" style="display: block;" tabindex="-1">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h1 class="modal-title fs-5">{{title}}</h1>
                    <button type="button" class="btn-close" v-on:click="$emit('close')"></button>
                </div>
                <div class="modal-body">
                    <slot name="body"></slot>
                </div>
                <div class="modal-footer">
                    <slot name="footer"></slot>
                </div>
            </div>
        </div>
    </div>
    `
};