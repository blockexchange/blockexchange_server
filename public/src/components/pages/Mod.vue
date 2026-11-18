<script setup lang="ts">
import { computed } from 'vue';
import PageHeader from "../PageHeader.vue";
import ClipboardCopy from "../ClipboardCopy.vue";
import { START, MOD } from "../../util/breadcrumbs";
import { is_logged_in } from "../../service/login";

const breadcrumb = [START, MOD];
const logged_in = computed(is_logged_in);
</script>

<template>
	<page-header
		title="Get the blockexchange mod"
		subtitle="Install the mod once, then load and save schematics with chat commands in any world."
		:breadcrumb="breadcrumb">
		<template #actions>
			<a class="btn btn-success" href="https://content.luanti.org/packages/BuckarooBanzay/blockexchange/" target="_blank" rel="noopener">
				<img src="/pics/contentdb.png" width="20" height="20" alt=""> Open on ContentDB
			</a>
		</template>
	</page-header>

	<div class="row g-4">
		<div class="col-lg-8">
			<section class="card mb-4">
				<div class="card-header">Setup</div>
				<div class="card-body d-flex flex-column gap-4">
					<div class="bx-step">
						<span class="bx-step-number">1</span>
						<div class="flex-grow-1">
							<h3 class="h6 fw-bold">Install from the content browser</h3>
							<p class="text-muted-bx">
								Open Luanti, go to the <b>Content</b> tab, search for <b>blockexchange</b> and install it.
							</p>
							<img src="/pics/help_install_mod.png" class="img-fluid rounded border border-secondary" alt="The blockexchange mod in the Luanti content browser" loading="lazy"/>
						</div>
					</div>
					<div class="bx-step">
						<span class="bx-step-number">2</span>
						<div class="flex-grow-1">
							<h3 class="h6 fw-bold">Allow web access for the mod</h3>
							<p class="text-muted-bx">
								The mod talks to this server over HTTP. Add it to the trusted HTTP mods in your settings
								(search for "http" in the settings tab) or add this line to <code>minetest.conf</code>:
							</p>
							<p><clipboard-copy text="secure.http_mods = blockexchange"/></p>
							<img src="/pics/help_install_secure_mods.png" class="img-fluid rounded border border-secondary" alt="The secure HTTP mods setting" loading="lazy"/>
							<div class="bx-note info mt-3">
								<i class="fa fa-circle-info"></i>
								Without this setting only the offline commands (local files) work.
							</div>
						</div>
					</div>
					<div class="bx-step">
						<span class="bx-step-number">3</span>
						<div class="flex-grow-1">
							<h3 class="h6 fw-bold">Enable the mod and grant the privilege</h3>
							<p class="text-muted-bx mb-2">
								Enable <b>blockexchange</b> for your world, join it and grant yourself the privilege:
							</p>
							<clipboard-copy text="/grantme blockexchange"/>
						</div>
					</div>
				</div>
			</section>

			<section class="card mb-4">
				<div class="card-header"><i class="fa fa-download"></i> Place a schematic</div>
				<div class="card-body">
					<ol class="mb-0 d-flex flex-column gap-2">
						<li>Find a build on the <router-link to="/search">browse page</router-link>.</li>
						<li>Stand where it should go and mark the origin with <clipboard-copy text="/bx_pos1"/></li>
						<li>Copy the <b>/bx_load</b> command from the schematic page and run it.</li>
					</ol>
					<div class="bx-note mt-3">
						<i class="fa fa-triangle-exclamation"></i>
						Placing a schematic replaces whatever is already there. Use <b>/bx_allocate</b> first to preview the area.
					</div>
				</div>
			</section>

			<section class="card mb-4">
				<div class="card-header"><i class="fa fa-upload"></i> Upload a build</div>
				<div class="card-body">
					<p class="text-muted-bx" v-if="!logged_in">
						You need an account to upload.
						<router-link to="/register">Sign up</router-link> or <router-link to="/login">log in</router-link> first.
					</p>
					<ol class="d-flex flex-column gap-2">
						<li>Create an access token on your <router-link to="/profile">account page</router-link> and run the <b>/bx_login</b> command shown there.</li>
						<li>Mark two opposite corners of your build with <clipboard-copy text="/bx_pos1"/> and <clipboard-copy text="/bx_pos2"/></li>
						<li>Upload it with <clipboard-copy text="/bx_save my_build"/></li>
					</ol>
					<p class="small text-muted-bx mb-0">
						Uploading needs the <b>blockexchange</b> or <b>blockexchange_protected_upload</b> privilege on that server.
						Have a WorldEdit file instead? Use the <router-link to="/import">file import</router-link>.
					</p>
				</div>
			</section>
		</div>

		<aside class="col-lg-4">
			<section class="card mb-4">
				<div class="card-header"><i class="fa-brands fa-youtube"></i> Video walkthrough</div>
				<div class="ratio ratio-16x9">
					<iframe src="https://www.youtube-nocookie.com/embed/x0geU_EyO-0" title="Blockexchange installation video" loading="lazy" allowfullscreen></iframe>
				</div>
			</section>
			<section class="card">
				<div class="card-header"><i class="fa fa-terminal"></i> Common commands</div>
				<ul class="list-group list-group-flush small">
					<li class="list-group-item"><b>/bx</b><br><span class="text-muted-bx">Info and settings menu</span></li>
					<li class="list-group-item"><b>/bx_load</b> user name<br><span class="text-muted-bx">Place a schematic at pos1</span></li>
					<li class="list-group-item"><b>/bx_allocate</b> user name<br><span class="text-muted-bx">Preview the area a schematic needs</span></li>
					<li class="list-group-item"><b>/bx_placer</b> user name<br><span class="text-muted-bx">Get a tool that places it on click</span></li>
					<li class="list-group-item"><b>/bx_save</b> name<br><span class="text-muted-bx">Upload the marked area</span></li>
					<li class="list-group-item"><b>/bx_cancel</b><br><span class="text-muted-bx">Stop a running job</span></li>
				</ul>
				<div class="card-body">
					<a href="https://github.com/blockexchange/blockexchange/blob/master/readme.md#chat-commands" target="_blank" rel="noopener" class="small">
						Full command reference <i class="fa fa-up-right-from-square"></i>
					</a>
				</div>
			</section>
		</aside>
	</div>
</template>
