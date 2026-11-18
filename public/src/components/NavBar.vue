<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { get_claims, has_permission, logout } from "../service/login";
import ThemeSwitcher from "./ThemeSwitcher.vue";

const route = useRoute();
const router = useRouter();

const claims = computed(get_claims);
const is_admin = computed(() => has_permission('ADMIN'));

const menu_open = ref(false);
const user_menu_open = ref(false);
const user_menu = ref<HTMLElement | null>(null);

// close all menus after navigating
watch(() => route.fullPath, () => {
	menu_open.value = false;
	user_menu_open.value = false;
});

function on_document_click(e: MouseEvent) {
	if (user_menu.value && !user_menu.value.contains(e.target as Node)) {
		user_menu_open.value = false;
	}
}
onMounted(() => document.addEventListener("click", on_document_click));
onBeforeUnmount(() => document.removeEventListener("click", on_document_click));

async function do_logout() {
	await logout();
	user_menu_open.value = false;
	await router.push("/");
}
</script>

<template>
	<nav class="navbar navbar-expand-lg bg-body-tertiary bx-navbar">
		<div class="container-xl">
			<router-link to="/" class="navbar-brand d-flex align-items-center gap-2">
				<img src="/pics/blockexchange_logo_small.png" alt="">
				<span>Block exchange</span>
			</router-link>
			<button class="navbar-toggler" type="button"
				:aria-expanded="menu_open"
				aria-label="Toggle navigation"
				@click="menu_open = !menu_open">
				<span class="navbar-toggler-icon"></span>
			</button>
			<div class="collapse navbar-collapse" :class="{show: menu_open}">
				<ul class="navbar-nav me-auto mb-2 mb-lg-0 gap-lg-1">
					<li class="nav-item">
						<router-link to="/search" class="nav-link">
							<i class="fa fa-magnifying-glass fa-fw"></i> Browse
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/mods" class="nav-link">
							<i class="fa fa-box-archive fa-fw"></i> Mods
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/users" class="nav-link">
							<i class="fa fa-users fa-fw"></i> Users
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/mod" class="nav-link">
							<i class="fa fa-download fa-fw"></i> Get the mod
						</router-link>
					</li>
				</ul>
				<div class="d-flex flex-column flex-lg-row gap-2" v-if="!claims">
					<router-link to="/login" class="btn btn-outline-secondary btn-sm">
						Log in
					</router-link>
					<router-link to="/register" class="btn btn-success btn-sm">
						Sign up
					</router-link>
					<theme-switcher/>
				</div>
				<div class="d-flex flex-column flex-lg-row gap-2 align-items-lg-center" v-else>
					<router-link to="/import" class="btn btn-success btn-sm">
						<i class="fa fa-upload"></i> Upload
					</router-link>
					<div class="dropdown" ref="user_menu">
						<button class="btn btn-outline-secondary btn-sm dropdown-toggle w-100"
							type="button"
							:aria-expanded="user_menu_open"
							@click="user_menu_open = !user_menu_open">
							<i class="fa fa-user"></i>
							{{ claims.username }}
						</button>
						<ul class="dropdown-menu" :class="{show: user_menu_open}">
							<li>
								<router-link :to="'/user/' + claims.username" class="dropdown-item">
									<i class="fa fa-cubes fa-fw"></i> My schematics
								</router-link>
							</li>
							<li v-if="is_admin">
								<router-link :to="'/collections/' + claims.username" class="dropdown-item">
									<i class="fa fa-object-group fa-fw"></i> My collections
								</router-link>
							</li>
							<li>
								<router-link to="/import" class="dropdown-item">
									<i class="fa fa-upload fa-fw"></i> Import files
								</router-link>
							</li>
							<li>
								<router-link to="/profile" class="dropdown-item">
									<i class="fa fa-gear fa-fw"></i> Account &amp; access tokens
								</router-link>
							</li>
							<template v-if="is_admin">
								<li><hr class="dropdown-divider"></li>
								<li>
									<router-link to="/tags" class="dropdown-item">
										<i class="fa fa-tags fa-fw"></i> Tags
									</router-link>
								</li>
							</template>
							<li><hr class="dropdown-divider"></li>
							<li>
								<button class="dropdown-item" type="button" @click="do_logout">
									<i class="fa fa-right-from-bracket fa-fw"></i> Log out
								</button>
							</li>
						</ul>
					</div>
					<theme-switcher/>
				</div>
			</div>
		</div>
	</nav>
</template>
