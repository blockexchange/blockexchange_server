<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import FeaturedSchemas from "../FeaturedSchemas.vue";
import Stats from "../Stats.vue";
import { is_logged_in } from "../../service/login";

const router = useRouter();
const keywords = ref("");
const logged_in = computed(is_logged_in);

function search() {
	router.push({ path: "/search", query: { q: keywords.value || undefined } });
}
</script>

<template>
	<section class="bx-hero">
		<h1>Share your builds <span class="accent">across worlds</span></h1>
		<p class="lead">
			Block exchange is a free schematic library for Luanti (formerly Minetest).
			Save a build on one server, then place it on any other with a single chat command.
		</p>
		<form class="bx-hero-search" @submit.prevent="search" role="search">
			<div class="input-group input-group-lg">
				<input type="search" class="form-control" v-model="keywords" placeholder="Search schematics, e.g. castle" aria-label="Search schematics"/>
				<button class="btn btn-success" type="submit">
					<i class="fa fa-magnifying-glass"></i>
					<span class="d-none d-sm-inline"> Search</span>
				</button>
			</div>
		</form>
		<div class="d-flex flex-wrap justify-content-center gap-2">
			<router-link to="/search" class="btn btn-outline-light">
				Browse all schematics
			</router-link>
			<router-link to="/mod" class="btn btn-outline-light">
				<i class="fa fa-download"></i> Get the mod
			</router-link>
		</div>
	</section>

	<stats class="mb-5"/>

	<section class="mb-5">
		<div class="d-flex justify-content-between align-items-end mb-3">
			<h2 class="bx-section-title mb-0">Featured schematics</h2>
			<router-link to="/search" class="small">See all <i class="fa fa-arrow-right"></i></router-link>
		</div>
		<featured-schemas/>
	</section>

	<section class="mb-4">
		<h2 class="bx-section-title">How it works</h2>
		<div class="row g-3">
			<div class="col-md-4">
				<div class="card h-100"><div class="card-body bx-step">
					<span class="bx-step-number">1</span>
					<div>
						<h3 class="h6 fw-bold">Install the mod</h3>
						<p class="mb-0 small text-muted-bx">
							Add the <b>blockexchange</b> mod to your world from the in-game content browser.
							<router-link to="/mod">Setup guide</router-link>
						</p>
					</div>
				</div></div>
			</div>
			<div class="col-md-4">
				<div class="card h-100"><div class="card-body bx-step">
					<span class="bx-step-number">2</span>
					<div>
						<h3 class="h6 fw-bold">Find a build</h3>
						<p class="mb-0 small text-muted-bx">
							<router-link to="/search">Browse</router-link> schematics by keyword, tag or the
							<router-link to="/mods">mods they use</router-link>.
						</p>
					</div>
				</div></div>
			</div>
			<div class="col-md-4">
				<div class="card h-100"><div class="card-body bx-step">
					<span class="bx-step-number">3</span>
					<div>
						<h3 class="h6 fw-bold">Place or share it</h3>
						<p class="mb-0 small text-muted-bx">
							Copy the chat command from a schematic page to place it.
							<template v-if="logged_in">
								<router-link to="/import">Upload</router-link> your own builds to share them.
							</template>
							<template v-else>
								<router-link to="/register">Create an account</router-link> to share your own builds.
							</template>
						</p>
					</div>
				</div></div>
			</div>
		</div>
	</section>
</template>
