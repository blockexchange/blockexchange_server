<script setup lang="ts">
import { computed } from 'vue';
import { get_claims } from "../../service/login";

import BreadCrumb from "../Breadcrumb.vue";
import UserProfile from "../UserProfile.vue";
import UserRename from "../UserRename.vue";
import AccessTokens from "../AccessTokens.vue";
import ChangePassword from "../ChangePassword.vue";
import UnlinkOauth from "../UnlinkOauth.vue";
import { START, PROFILE } from "../../util/breadcrumbs";

const breadcrumb = [START, PROFILE];
const claims = computed(get_claims);
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<div class="row">
		<div class="col-md-2"></div>
		<div class="col-md-8 card" style="padding: 20px;" v-if="claims">
			<user-profile :username="claims.username"/>
			<hr>
			<access-tokens :username="claims.username"/>
			<hr>
			<change-password v-if="claims.type == 'LOCAL'" :useruid="claims.user_uid"/>
			<unlink-oauth v-if="claims.type != 'LOCAL'" :useruid="claims.user_uid"/>
			<user-rename/>
		</div>
		<div class="col-md-8" v-else>
			<div class="alert alert-secondary">
				<i class="fa fa-circle-info"></i>
				Please <router-link to="/login">log in</router-link> to view your profile
			</div>
		</div>
		<div class="col-md-2"></div>
	</div>
</template>
