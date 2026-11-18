<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { theme, type ThemePreference } from "../service/theme";

const options: { value: ThemePreference; label: string; icon: string }[] = [
	{ value: "auto", label: "System", icon: "fa-circle-half-stroke" },
	{ value: "light", label: "Light", icon: "fa-sun" },
	{ value: "dark", label: "Dark", icon: "fa-moon" }
];

const open = ref(false);
const root = ref<HTMLElement | null>(null);
const current = computed(() => options.find(o => o.value == theme.value)!);

function select(v: ThemePreference) {
	theme.value = v;
	open.value = false;
}

function on_document_click(e: MouseEvent) {
	if (root.value && !root.value.contains(e.target as Node)) {
		open.value = false;
	}
}
onMounted(() => document.addEventListener("click", on_document_click));
onBeforeUnmount(() => document.removeEventListener("click", on_document_click));
</script>

<template>
	<div class="dropdown" ref="root">
		<button class="btn btn-sm btn-outline-secondary dropdown-toggle w-100"
			type="button"
			:aria-expanded="open"
			:title="'Theme: ' + current.label"
			@click="open = !open">
			<i class="fa-solid fa-fw" :class="current.icon"></i>
			<span class="d-lg-none ms-1">Theme: {{ current.label }}</span>
			<span class="visually-hidden">Toggle theme</span>
		</button>
		<ul class="dropdown-menu dropdown-menu-end" :class="{show: open}">
			<li v-for="o in options" :key="o.value">
				<button type="button" class="dropdown-item d-flex align-items-center gap-2"
					:class="{active: o.value == theme}"
					:aria-pressed="o.value == theme"
					@click="select(o.value)">
					<i class="fa-solid fa-fw" :class="o.icon"></i>
					{{ o.label }}
					<i class="fa fa-check ms-auto" v-if="o.value == theme"></i>
				</button>
			</li>
		</ul>
	</div>
</template>
