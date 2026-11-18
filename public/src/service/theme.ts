import { ref, watch } from 'vue';

export type ThemePreference = "auto" | "light" | "dark";

const STORAGE_KEY = "bx-theme";
const media = window.matchMedia("(prefers-color-scheme: dark)");

function load(): ThemePreference {
	try {
		const v = localStorage.getItem(STORAGE_KEY);
		if (v == "light" || v == "dark") {
			return v;
		}
	} catch {
		// storage blocked, use the system theme
	}
	return "auto";
}

// the selected preference, "auto" follows the operating system
export const theme = ref<ThemePreference>(load());

const resolve = (t: ThemePreference) => t == "auto" ? (media.matches ? "dark" : "light") : t;

function apply() {
	document.documentElement.setAttribute("data-bs-theme", resolve(theme.value));
}

watch(theme, t => {
	try {
		if (t == "auto") {
			localStorage.removeItem(STORAGE_KEY);
		} else {
			localStorage.setItem(STORAGE_KEY, t);
		}
	} catch {
		// not persisted, still applied for this page view
	}
	apply();
});

// follow system changes while in "auto" mode
media.addEventListener("change", () => {
	if (theme.value == "auto") {
		apply();
	}
});

apply();
