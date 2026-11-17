<script setup lang="ts">
import BreadCrumb from "../Breadcrumb.vue";
import FeaturedSchemas from "../FeaturedSchemas.vue";
import Stats from "../Stats.vue";
import { START } from "../../util/breadcrumbs";

const breadcrumb = [START];
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<div class="text-center container">
		<h3>
			Blockexchange
			<small class="text-muted">
				- exchange your schematics across worlds with ease
			</small>
		</h3>
		<hr/>
		<stats/>
		<br>
		<br>
		<router-link to="/search" class="btn btn-primary">
			<i class="fa fa-search"></i> Search
		</router-link>
		&nbsp;
		<router-link to="/users" class="btn btn-primary">
			<i class="fa fa-users"></i> Users
		</router-link>
		&nbsp;
		<router-link to="/mod" class="btn btn-primary">
			<i class="fa fa-download"></i> Mod/Installation
		</router-link>
		&nbsp;
		<a class="btn btn-secondary" href="https://github.com/blockexchange/blockexchange_server" target="new">
			<i class="fa-brands fa-github"></i> Browse the source
		</a>
		&nbsp;
		<a class="btn btn-secondary" href="https://discord.gg/ye9aCUJPdB" target="new">
			<i class="fa-brands fa-discord"></i> Join the discord community
		</a>
		<hr/>
		<h5>Featured schematics</h5>
		<featured-schemas/>
	</div>
</template>
