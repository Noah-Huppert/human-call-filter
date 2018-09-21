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
					<li v-on:click="closeNavbarMenu">
						<router-link to="/calls">
							Calls
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
				<phone-number-row v-for="number in phoneNumbers"
					v-bind:number="number"
					v-bind:selected-id="id">
				</phone-number-row>
			</tbody>
		</table>
	</div>`,
	props: {
		id: undefined
	},
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

Vue.component("phone-number-row", {
	template: `<tr v-bind:class="{ selected: isSelected }">
		<td>{{ number.ID }}</td>
		<td>{{ number.Number }}</td>
		<td>{{ number.Name }}</td>
		<td>{{ number.State }}</td>
		<td>{{ number.City }}</td>
		<td>{{ number.ZipCode }}</td>
	</tr>`,
	props: ["number", "selected-id"],
	data: function() {
		return {
			isSelected: this.getIsSelected()
		};
	},
	watch: {
		selectedId: function() {
			this.isSelected = this.getIsSelected();
		}
	},
	methods: {
		getIsSelected: function() {
			return this.number.ID == this.selectedId;
		}
	}
});

const phoneCallsPage = Vue.component("phone-calls-page", {
	template: `<div class="container">
		<h1 class="title">Phone Calls</h1>

		<table class="table">
			<thead>
				<tr>
					<th>ID</th>
					<th>Phone Number ID</th>
					<th>Twilio Call ID</th>
					<th>Date Received</th>
				</tr>
			</thead>
			<tbody>
				<phone-call-row v-for="call in phoneCalls"
					v-bind:call="call"
					v-bind:selected-id="id">
				</phone-call-row>
			</tbody>
		</table>
	</div>`,
	props: {
		id: undefined
	},
	data: function() {
		return {
			phoneCalls: this.phoneCalls
		};
	},
	created: function() {
		this.phoneCalls = [];
		var self = this;

		makeAPIRequest("/api/phone_calls", "GET")
			.then(function(resp) {
				self.phoneCalls = resp.phone_calls;
			});
	}
});

Vue.component("phone-call-row", {
	template: `<tr v-bind:class="{ selected: isSelected }">
		<td>{{ call.ID }}</td>
		<td>
			<router-link v-bind:to="'/numbers?id=' + call.PhoneNumberID">
				{{ call.PhoneNumberID }}
			</router-link>
		</td>
		<td>{{ call.TwilioCallID }}</td>
		<td>{{ call.DateReceived }}</td>
	</tr>`,
	props: ["call", "selected-id"],
	data: function() {
		return {
			isSelected: this.getIsSelected()
		};
	},
	watch: {
		selectedId: function() {
			this.isSelected = this.getIsSelected();
		}
	},
	methods: {
		getIsSelected: function() {
			return this.call.ID == this.selectedId;
		}
	}
});

/* Router */
const router = new VueRouter({
	routes: [
		{
			path: "/",
			redirect: "/numbers"
		},
		{
			path: "/numbers",
			component: phoneNumbersPage,
			props: function(route) {
				return {
					id: route.query.id
				};
			}
		},
		{
			path: "/calls",
			component: phoneCallsPage,
			props: function(route) {
				return {
					id: route.query.id
				};
			}
		}
	]
});

/* Root */
var app = new Vue({
	el: "#app",
	router: router
});
