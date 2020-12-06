export default {
	template: /*html*/`
    <div style="text-align: center;">
      <img src="pics/blockexchange.png"/>
			<div class="alert alert-danger" role="alert">
			  Blockexchange is currently undergoing major refactoring and may be unusable for now
				<br>
				Stay tuned for the next official release announcement
			</div>
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
          <i class="fa fa-github"></i>
          Source
        </a>
				<a href="https://github.com/login/oauth/authorize?client_id=68c2728e22f3a4b02dc0" class="btn btn-secondary">
          Login test
        </a>
      </div>
    </div>
  `
};
