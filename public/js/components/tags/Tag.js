
export default {
    props: ["tag"],
    template: `
        <span class="badge badge-success">
            <i class="fas fa-tag"></i>
            {{ tag.name }}
        </span>
    `
};