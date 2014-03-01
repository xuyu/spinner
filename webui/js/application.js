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

// datacenter label
var g_datacenter_name = $("#datacenter-name");
var g_datacenter_location = $("#datacenter-location");
var g_group_tree = $("#dt-group-tree");
var g_terminal = $("#dt-terminal");
var g_terminal_textarea = $("#dt-terminal pre");
var g_terminal_input = $("#dt-terminal input");

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
			g_group_tree.addClass("cur-task");
			group_tree_api();
		break
		case "trust_central":
		break
		case "dashboard":
		break
		case "file_system":
		break
		case "terminal":
			g_terminal.addClass("cur-task");
			g_terminal_textarea.empty();
			g_terminal_input.val(terminal_prompt());
		break
	}
	close_menu_nav();
}

function terminal_prompt(){
	return "root@" + g_cur_hostname + "# ";
}

function terminal_textarea_append(text){
	var original = g_terminal_textarea.text();
	if (original != "") {
		original = original + "\n";
	}
	g_terminal_textarea.text(original + text);
	g_terminal_textarea.scrollTop(g_terminal_textarea[0].scrollHeight - g_terminal_textarea.height());
}

function terminal_input_bind_event(){
	g_terminal_input.keyup(function(e){
		if (e.keyCode == 13) {
			var cmd = g_terminal_input.val();
			g_terminal_input.val("");
			g_terminal_input.attr("readonly", "readonly");
			terminal_textarea_append(cmd);
			terminal_api(cmd.replace(/^root@.+?# /, ""));
		}
	});
}

$(document).ready(function(){
	task('group_tree');
	terminal_input_bind_event();
})