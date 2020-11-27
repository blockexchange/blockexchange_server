Vue.component('start-page', {
	template: `
    <div style="text-align: center;">
      <img src="pics/blockexchange.png"/>
      <hr/>
      <h4>Exchange your schemas across worlds with ease</h4>
      <hr/>

      <div>
        <a href="#!/search" class="btn btn-primary">
          <i class="fa fa-search"></i>
          Search
        </a>
        <a href="#!/users" class="btn btn-primary">
          <i class="fa fa-users"></i>
          Users
        </a>
        <a href="#!/mod" class="btn btn-primary">
          <i class="fa fa-download"></i>
          Mod/Installation
        </a>
        <a href="https://github.com/blockexchange" class="btn btn-secondary">
          <i class="fab fa-github"></i>
          Source
        </a>
      </div>
    </div>
  `
});
