Vue.config.devtools = true;

/* API */
function makeAPIRequest(path, method) {
	return fetch(path, {
		method: method
	}).then(res => res.json());
}

/* Navigation */
function toggleNavbarMenu() {
	document.getElementById("nav-menu").classList.toggle("is-active");
}

function closeNavbarMenu() {
	console.log("close");
	document.getElementById("nav-menu").classList.remove("is-active");
}

Vue.component("navbar-brand", {
	template: `<div class="navbar-brand">
		<div class="navbar-item">Human Call Filter</div>	
		<a role="button" class="navbar-burger" v-on:click="toggleNavbarMenu">
			<span></span>
			<span></span>
			<span></span>
		</a>		
	</div>`,
	methods: {
		toggleNavbarMenu: toggleNavbarMenu,
	}
});

Vue.component("navbar-menu", {
	template: `<div id="nav-menu" class="navbar-menu">
		<div class="navbar-end">
			<div class="navbar-item">
				<ul>
					<li v-on:click="closeNavbarMenu">
						<router-link to="/numbers">
							Numbers
						</router-link>
					</li>
				</ul>
			</div>
		</div>
	</div>`,
	methods: {
		closeNavbarMenu: closeNavbarMenu
	}
});

/* Pages */
const phoneNumbersPage = Vue.component("phone-numbers-page", {
	template: `<div class="container">
		<h1 class="title">Phone Numbers</h1>

		<table class="table">
			<thead>
				<tr>
					<th>ID</th>
					<th>Number</th>
					<th>Name</th>
					<th>State</th>
					<th>City</th>
					<th>Zip Code</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="number in phoneNumbers">
					<td>{{ number.ID }}</td>
					<td>{{ number.Number }}</td>
					<td>{{ number.Name }}</td>
					<td>{{ number.State }}</td>
					<td>{{ number.City }}</td>
					<td>{{ number.ZipCode }}</td>
				</tr>
			</tbody>
		</table>
	</div>`,
	data: function() {
		return {
			phoneNumbers: this.phoneNumbers
		};
	},
	created: function() {
		this.phoneNumbers = [];
		var self = this;

		makeAPIRequest("/api/phone_numbers", "GET")
			.then(function(resp) {
				self.phoneNumbers = resp.phone_numbers;
			});
	}
});

/* Router */
const router = new VueRouter({
	routes: [
		{ path: "/", redirect: "/numbers" },
		{ path: "/numbers", component: phoneNumbersPage }
	]
});

/* Root */
var app = new Vue({
	el: "#app",
	router: router
});
