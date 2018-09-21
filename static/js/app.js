var app = new Vue({
	el: "#app",
	data: {
		hello: "hello world"
	},
	methods: {
		toggleNavbarMenu: toggleNavbarMenu
	}
});

function toggleNavbarMenu() {
	document.getElementById("nav-menu").classList.toggle("is-active");
}
