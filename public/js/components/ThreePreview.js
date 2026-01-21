import { count_schemaparts_by_mtime, get_next_schemapart_by_mtime, get_next_schemapart_by_offset, get_schemapart_by_offset } from "../api/schema_part.js";
import { get_colormapping } from "../api/colormapping.js";
import { decompressSync } from 'fflate';

import { Color, PerspectiveCamera, Scene, WebGLRenderer, AxesHelper, AmbientLight } from "three";
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import Stats from 'three/examples/jsm/libs/stats.module.js'

function base64ToArrayBuffer(base64) {
    var binaryString = atob(base64);
    var bytes = new Uint8Array(binaryString.length);
    for (var i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
}

function ab2str(buf) {
    return String.fromCharCode.apply(null, new Uint16Array(buf));
}

export default {
    props: ["schema_uid"],
    data: function() {
        return {
            active: true
        };
    },
    mounted: async function() {
        const part = await get_next_schemapart_by_mtime(this.schema_uid, 0);
        console.log(part);

        const buf = base64ToArrayBuffer(part.metadata);
        const ab = decompressSync(new Uint8Array(buf));
        console.log(ab);

        console.log(JSON.parse(ab2str(ab)));

        const mapbuf = base64ToArrayBuffer(part.data);
        const mapdata = decompressSync(new Uint8Array(mapbuf));

        console.log(mapdata);

        const cm = await get_colormapping();
        console.log(cm["default:stone"]);

        //

        const el = this.$refs.target;

        this.scene = new Scene();
        this.camera = new PerspectiveCamera(75, el.innerWidth / el.innerHeight, 0.1, 1000)
        this.stats = new Stats()

        el.parentElement.appendChild(this.stats.dom)
        this.scene.background = new Color(1,1,1)

        this.renderer = new WebGLRenderer({ canvas: el })
        this.renderer.setSize(el.innerWidth, el.innerHeight)

        this.controls = new OrbitControls(this.camera, el)
        this.controls.listenToKeyEvents(document.body)
        this.controls.minDistance = 5
        this.controls.maxDistance = 500

        this.controls.target.x = 30
        this.controls.target.y = -30
        this.controls.target.z = -30

        const axesHelper = new AxesHelper( 5 );
        this.scene.add( axesHelper );

        this.scene.add(new AmbientLight(0xffffff)); //TODO: dark/light mode

        this.controls.addEventListener('change', () => this.render());

        this.animate()

        const part_count = await count_schemaparts_by_mtime(this.schema_uid, 0);
        console.log("Parts: " + part_count);

        var pos_x = 0, pos_y = 0, pos_z = 0;
        var schema_part = await get_schemapart_by_offset(this.schema_uid, pos_x, pos_y, pos_z);

        do {
            var schema_part = await get_next_schemapart_by_offset(this.schema_uid, pos_x, pos_y, pos_z);
            if (!schema_part) {
                break;
            }

            // TODO: render stuff

            pos_x = schema_part.offset_x;
            pos_y = schema_part.offset_y;
            pos_z = schema_part.offset_z;
        } while (true);
    },
    methods: {
        animate: function() {
            if (!this.active) {
                return;
            }
            this.controls.update()
            window.requestAnimationFrame(() => this.animate())
        },
        render: function() {
            this.stats.begin()
            this.renderer.render(this.scene, this.camera)
            this.stats.end()
            console.log(`Calls: ${this.renderer.info.render.calls}, Triangles: ${this.renderer.info.render.triangles}`)
        }
    },
    beforeUnmount: function() {
        this.scene.remove();
        this.controls.remove();
        this.active = false;
    },
    template: /*html*/`
        <canvas ref="target"></canvas>
    `
};
