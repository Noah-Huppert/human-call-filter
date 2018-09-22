Vue.config.devtools = true;

/* API */
function makeAPIRequest(path, method) {
	return fetch(path, {
		method: method
	}).then(res => res.json());
}

const store = new Vuex.Store({
	state: {
		phoneNumbers: [],
		phoneCalls: [],
		challenges: []
	},
	mutations: {
		setPhoneNumbers: function(state, phoneNumbers) {
			state.phoneNumbers = phoneNumbers;
		},
		setPhoneCalls: function(state, phoneCalls) {
			state.phoneCalls = phoneCalls;
		},
		setChallenges: function(state, challenges) {
			state.challenges = challenges;
		}
	}
});

makeAPIRequest("/api/phone_numbers", "GET")
	.then(function(resp) {
		store.commit("setPhoneNumbers", resp.phone_numbers);
	});

makeAPIRequest("/api/phone_calls", "GET")
	.then(function(resp) {
		store.commit("setPhoneCalls", resp.phone_calls);
	});

makeAPIRequest("/api/challenges", "GET")
	.then(function(resp) {
		store.commit("setChallenges", resp.challenges);
	});

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
			<div class="navbar-item" v-on:click="closeNavbarMenu">
				<router-link to="/numbers">
					Numbers
				</router-link>
			</div>
			<div class="navbar-item" v-on:click="closeNavbarMenu">
				<router-link to="/calls">
					Calls
				</router-link>
			</div>
			<div class="navbar-item" v-on:click="closeNavbarMenu">
				<router-link to="/challenges">
					Challenges
				</router-link>
			</div>
		</div>
	</div>`,
	methods: {
		closeNavbarMenu: closeNavbarMenu
	}
});

/* Components */
Vue.component("data-table", {
	template: `<div class="container">
		<h1 class="title">{{ title }}</h1>

		<table class="table">
			<thead>
				<tr>
					<th v-for="name in headerNames">{{ name }}</th>
				</tr>
			</thead>
			<tbody>
				<component v-for="item in items"
					v-bind:is="rowComponent"
					v-bind:item="item"
					v-bind:class="[item.ID == selectedId ? 'selected' : '']">
				<component>
			</tbody>
		</table>
	</div>`,
	props: ["title", "items", "header-names", "row-component", "selected-id"]
});

/* Phone numbers page */
const phoneNumbersPage = Vue.component("phone-numbers-page", {
	template: `<div class="container">
		<data-table title="Phone Numbers"
			v-bind:items="phoneNumbers"
			v-bind:header-names="headerNames"
			v-bind:row-component="phoneNumberRow"
			v-bind:selected-id="selectedId">
		</data-table>
	</div>`,
	props: {
		selectedId: undefined
	},
	data: function() {
		return {
			headerNames: ["ID", "Number", "Name", "State", "City", "Zip Code"]
		};
	},
	computed: {
		phoneNumbers: function() {
			return store.state.phoneNumbers;
		}
	}
});

const phoneNumberRow = Vue.component("phone-number-row", {
	template: `<tr>
		<td>{{ item.ID }}</td>
		<td>{{ item.Number }}</td>
		<td>{{ item.Name }}</td>
		<td>{{ item.State }}</td>
		<td>{{ item.City }}</td>
		<td>{{ item.ZipCode }}</td>
	</tr>`,
	props: ["item"]
});

/* Phone calls page */
const phoneCallsPage = Vue.component("phone-calls-page", {
	template: `<div class="container">
		<data-table title="Phone Calls"
			v-bind:items="phoneCalls"
			v-bind:header-names="headerNames"
			v-bind:row-component="phoneCallRow"
			v-bind:selected-id="selectedId">
		</data-table>
	</div>`,
	props: ["selected-id"],
	data: function() {
		return {
			headerNames: ["ID", "Phone Number ID", "Twilio Call ID",
				"Date Received"]
		};
	},
	computed: {
		phoneCalls: function() {
			return store.state.phoneCalls;
		}
	}
});

const phoneCallRow = Vue.component("phone-call-row", {
	template: `<tr>
		<td>{{ item.ID }}</td>
		<td>
			<router-link v-bind:to="'/numbers?id=' + item.PhoneNumberID">
				{{ item.PhoneNumberID }}
			</router-link>
		</td>
		<td>{{ item.TwilioCallID }}</td>
		<td>{{ item.DateReceived }}</td>
	</tr>`,
	props: ["item"]
});

/* Challenges page */
const challengesPage = Vue.component("challenges-page", {
	template: `<div class="container">
		<data-table title="Challenges"
			v-bind:items="challenges"
			v-bind:header-names="headerNames"
			v-bind:row-component="challengeRow"
			v-bind:selected-id="selectedId">
		</data-table>
	</div>`,
	props: ["selected-id"],
	data: function() {
		return {
			headerNames: ["ID", "Phone Call ID", "Date Asked", "Operand A",
				"Operand B", "Solution", "Status"]
		};
	},
	computed: {
		challenges: function() {
			return store.state.challenges;
		}
	}
});

const challengeRow = Vue.component("challenge-row", {
	template: `<tr>
		<td>{{ item.ID }}</td>
		<td>
			<router-link v-bind:to="'/calls?id=' + item.PhoneCallID">
				{{ item.PhoneCallID }}
			</router-link>
		</td>
		<td>{{ item.DateAsked }}</td>
		<td>{{ item.OperandA }}</td>
		<td>{{ item.OperandB }}</td>
		<td>{{ item.Solution }}</td>
		<td>
			<div class="tag" v-bind:class="[statusClass]">
				{{ item.Status }}
			</div>
		</td>
	</tr>`,
	props: ["item"],
	data: function() {
		return {
			statusClass: ""
		};
	},
	watch: {
		item: function() {
			this.setStatusClass();
		}
	},
	methods: {
		setStatusClass: function() {
			switch (this.item.Status) {
				case "PASSED":
					this.statusClass = "is-success";
					break;

				case "FAILED":
					this.statusClass = "is-danger";
					break;

				case "ANSWERING":
					this.statusClass = "is-warning";
					break;
			}
		}
	},
	created: function() {
		this.setStatusClass()
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
					'selected-id': route.query.id
				};
			}
		},
		{
			path: "/calls",
			component: phoneCallsPage,
			props: function(route) {
				return {
					'selected-id': route.query.id
				};
			}
		},
		{
			path: "/challenges",
			component: challengesPage,
			props: function(route) {
				return {
					'selected-id': route.query.id
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
