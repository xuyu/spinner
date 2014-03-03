// menu-nav
var g_pusher = $("#pusher");
var g_menu_nav = $("#menu-nav");
var g_menu_nav_open = true;

function open_menu_nav(){
	g_pusher.addClass("push-open");
	g_menu_nav.addClass("menu-nav-open");
	g_menu_nav_open = true;
}

function close_menu_nav(){
	g_pusher.removeClass("push-open");
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
var g_filesystem = $("#dt-filesystem");
var g_file_editor = ace.edit("file-editor");
g_file_editor.setTheme("ace/theme/xcode");
g_file_editor.setPrintMarginColumn(false);
var g_dashboard = $("#dt-dashboard");
var g_dashboard_btime = $("#dashboard-btime");
var g_dashboard_load = $("#dashboard-load");
var g_dashboard_cputimes = $("#dashboard-cputimes");
var g_dashboard_meminfo = $("#dashboard-meminfo");
var g_dashboard_diskusage = $("#dashboard-diskusage");
var g_dashboard_diskio = $("#dashboard-diskio");
var g_dashboard_netio = $("#dashboard-netio");

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

function task(name, hostname, options){
	if (name == g_cur_task) return
	close_last_task()
	change_cur_task(name);
	change_cur_hostname(hostname);
	switch (name) {
		case "group_tree":
			g_group_tree.addClass("cur-task");
			group_tree_api();
		break
		case "dashboard":
			g_dashboard.addClass("cur-task");
			dashboard_api();
		break
		case "file_system":
			g_filesystem.addClass("cur-task");
			if (options != null && options != undefined && options.file != undefined) {
				open_file_api(options.file);
			}
		break
		case "terminal":
			g_terminal.addClass("cur-task");
			g_terminal_textarea.empty();
			g_terminal_input.val(terminal_prompt());
		break
	}
	if (name != "dashboard") {
		g_dashboard_last_data = null;
		clearInterval(g_dashboard_runner);
		g_dashboard_runner = null;
	}
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
			var input = g_terminal_input.val().trim();
			g_terminal_input.val("");
			g_terminal_input.attr("readonly", "readonly");
			terminal_textarea_append(input);
			var cmd = input.replace(/^root@.+?# /, "").trim();
			var space = cmd.indexOf(" ");
			if (space == 4 && cmd.substr(0, space) == "open" && cmd.substr(space+1).length > 0
				&& cmd.substr(space+1)[0] == "/") {
				task("file_system", g_cur_hostname, {
					file: cmd.substr(space + 1).trim()
				});
			} else {
				terminal_api(cmd);
			}
			return false;
		}
	});
	g_terminal_input.keydown(function(e){
		if (e.keyCode == 8 || e.keyCode == 46) {
			if (g_terminal_input.val() == terminal_prompt()) {
				return false;
			}
		}
		return true;
	});
}

$(document).ready(function(){
	task('group_tree');
	terminal_input_bind_event();
})