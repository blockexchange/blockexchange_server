import html from './html.js';

export default {
  view: () => html`
    <div class="row">
      <div class="col-md-4">
      </div>
      <div class="col-md-4">
        <form onsubmit=${e => e.preventDefault()}>
          <input class="form-control" placeholder="Username"/>
          <input class="form-control" placeholder="Password"/>
          <input class="form-control" placeholder="Password"/>
          <button class="btn btn-primary btn-block">
            Register
          </button>
        </form>
      </div>
      <div class="col-md-4">
      </div>
    </div>
  `
};
