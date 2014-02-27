// menu-nav
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

// task
var g_cur_task = undefined;
var g_cur_hostname = undefined;
var g_cur_machine = $("#current-machine");

function change_cur_task(name){
	if (name != undefined && name != null && name != "") {
		g_cur_task = name;
	}
}

function change_cur_hostname(hostname){
	if (hostname != undefined && hostname != null && hostname != "") {
		g_cur_hostname = hostname;
		g_cur_machine.text("@" + hostname);
	}
}

function close_last_task(){
	$(".cur-task").removeClass("cur-task");
}

function task(name, hostname){
	if (name == g_cur_task) return
	close_last_task()
	change_cur_task(name);
	change_cur_hostname(hostname);
	switch (name) {
		case "group_tree":
			$("#dt-group-tree").addClass("cur-task");
			group_tree_api();
		break
		case "trust_central":
		break
		case "dashboard":
		break
		case "file_system":
		break
		case "terminal":
		break
	}
}

// datacenter label
var g_datacenter_name = $("#datacenter-name");
var g_datacenter_location = $("#datacenter-location");
var g_group_tree = $("#dt-group-tree");

$(document).ready(function(){
	open_menu_nav();
})