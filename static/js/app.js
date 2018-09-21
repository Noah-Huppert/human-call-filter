Vue.config.devtools = true;

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
const numbersPage = Vue.component("numbers-page", {
	template: `<div class="container">
		<h1 class="title">Numbers</h1>
	</div>`
});

/* Router */
const router = new VueRouter({
	routes: [
		{ path: "/numbers", component: numbersPage }
	]
});

/* Root */
var app = new Vue({
	el: "#app",
	router: router,
	data: {
	},
});
