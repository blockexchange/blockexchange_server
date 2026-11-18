<script setup lang="ts">
import { computed } from 'vue';
import { get_claims } from "../../service/login";

import PageHeader from "../PageHeader.vue";
import EmptyState from "../EmptyState.vue";
import UserRename from "../UserRename.vue";
import AccessTokens from "../AccessTokens.vue";
import ChangePassword from "../ChangePassword.vue";
import UnlinkOauth from "../UnlinkOauth.vue";
import { START, PROFILE } from "../../util/breadcrumbs";

const breadcrumb = [START, PROFILE];
const claims = computed(get_claims);
</script>

<template>
	<page-header title="Account settings" :breadcrumb="breadcrumb">
		<template #actions v-if="claims">
			<router-link :to="'/user/' + claims.username" class="btn btn-outline-secondary">
				<i class="fa fa-user"></i> View public profile
			</router-link>
		</template>
	</page-header>

	<div class="row justify-content-center" v-if="claims">
		<div class="col-xl-10 d-flex flex-column gap-4">
			<section class="card">
				<div class="card-header"><i class="fa fa-key"></i> Access tokens</div>
				<div class="card-body">
					<access-tokens :username="claims.username"/>
				</div>
			</section>

			<section class="card">
				<div class="card-header"><i class="fa fa-lock"></i> {{ claims.type == 'LOCAL' ? 'Password' : 'Login method' }}</div>
				<div class="card-body">
					<change-password v-if="claims.type == 'LOCAL'" :useruid="claims.user_uid"/>
					<unlink-oauth v-else :useruid="claims.user_uid" :provider="claims.type"/>
				</div>
			</section>

			<section class="card border-warning">
				<div class="card-header"><i class="fa fa-signature"></i> Username</div>
				<div class="card-body">
					<user-rename/>
				</div>
			</section>
		</div>
	</div>
	<empty-state icon="lock" title="You're not logged in" v-else>
		<router-link to="/login" class="btn btn-success mt-2">Log in</router-link>
	</empty-state>
</template>
