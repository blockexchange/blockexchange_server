//import { schema_search } from "../api/schema.js";

const store = Vue.reactive({

});

export default {
    data: () => store,
    template: /*html*/`
    <div class="row">
        <div class="col-md-4">
            <input type="text" class="form-control" placeholder="Keywords"/>
        </div>
        <div class="col-md-4">
            <select class="form-control">
                <option>All tags</option>
            </select>
        </div>
        <div class="col-md-4">
            <button class="btn btn-success w-100">
                <i class="fa-solid fa-magnifying-glass"></i>
                Search
            </button>
        </div>
    </div>
    `
};