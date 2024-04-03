import { get_stats } from "../service/info.js";
import format_size from "../util/format_size.js";
import format_count from "../util/format_count.js";

export default {
    methods: {
        format_size,
        format_count
    },
    computed: {
        stats: get_stats
    },
    template: /*html*/`
    <span v-if="stats">
        <b>Database Statistics:</b>
        Schematics: <span class="badge bg-secondary">{{format_count(stats.schema_count)}}</span>
        Schematic parts: <span class="badge bg-secondary">{{format_count(stats.schemapart_count)}}</span>
        Users: <span class="badge bg-secondary">{{format_count(stats.user_count)}}</span>
        Total size: <span class="badge bg-secondary">{{format_size(stats.total_size)}}</span>
    </span>
    `
};