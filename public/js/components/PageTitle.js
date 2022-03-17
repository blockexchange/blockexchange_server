
export default {
    props: ["major", "minor"],
    template: /*html*/`
        <h3 v-if="major || minor">
            <span v-if="major">{{ major }}</span>
			&nbsp;
            <small v-if="minor" class="text-muted">{{ minor }}</small>
        </h3>
    `
}