var g_menu_nav = $("#menu-nav");
var g_menu_nav_open = false;

function open_menu_nav(){
	g_menu_nav.addClass("menu-nav-open");
	g_menu_nav_open = true;
}

function close_menu_nav(){
	g_menu_nav.removeClass("menu-nav-open");
	g_menu_nav_open = false
}

function trigger_menu_nav(){
	if (g_menu_nav_open) {
		close_menu_nav()
	} else {
		open_menu_nav()
	}
}