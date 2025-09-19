import { schema_search } from "../api/schema.js";
import { get_schemapart_by_ofset } from "../api/schema_part.js";
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
        const list = await schema_search({ schema_name: this.name, user_name: this.username })
        const schema = list[0].schema;
        console.log(schema);

        const part = await get_schemapart_by_ofset(schema.uid, 0, 0, 0);
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

        this.scene.add(new AmbientLight(0xffffff));

        this.controls.addEventListener( 'change', () => this.render() );

        this.animate()
    },
    methods: {
        animate: function() {
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
        console.log("beforeUnmount")
        this.active = false;
    },
    template: /*html*/`
        <canvas ref="target"></canvas>
    `
};
