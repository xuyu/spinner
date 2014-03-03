function trigger_ul_head_click(ev){
	var e = ev.target;
	var ul = $(e).parent();
	var head_height = $(e).outerHeight();
	var ul_height = ul.outerHeight();
	if (ul_height > head_height) {
		ul.height(head_height);
	} else {
		ul.height("");
	}
}

function group_tree_api(){
	$.getJSON("/spinner/webui/tree", function(data){
		g_datacenter_name.text(data.Name);
		g_datacenter_location.text(data.Location);
		var html = "";
		$.each(data.Groups, function(index, gp) {
			html = html + "<ul><h3 class='dt-group-head'>" + gp.Name;
			html = html + "<span>" + gp.Machines.length + "</span></h3>";
			$.each(gp.Machines, function(index, ma){
				var cla = "";
				if (ma.KeepAlive == undefined) {
					cla = "miss";
				} else if (ma.KeepAlive == 0) {
					cla == "miss";
				} else if (Math.abs(ma.KeepAlive - (new Date()).getTime()/1000) > 300) {
					cla = "miss";
				} else {
					cla = "alive";
				}
				html = html + "<li class='" + cla;
				html = html + "'><a href=\"javascript: task('dashboard', '" + ma.Hostname + "');\">";
				html = html + ma.IP + " @ " + ma.Hostname + "</a></li>"
			});
			html = html + "</ul>";
		});
		g_group_tree.html(html);
		g_group_tree.find("ul>h3.dt-group-head").click(trigger_ul_head_click);
		g_group_tree.find("ul:first-child").height(g_group_tree.find("ul:first-child>h3").outerHeight());
		$.each(g_group_tree.find("ul"), function(index, ul){
			if ($(ul).find("li.miss").length > 0) {
				$(ul).addClass("has-miss");
			}
		});
	});
}

function terminal_api(cmd){
	$.ajax(
		{
			type: "GET",
			url: "/spinner/webui/terminal",
			data: {h: g_cur_hostname, cmd: cmd},
			timeout: 60000,
			success: function(data){
				terminal_textarea_append(data);
			},
			complete: function(){
				g_terminal_input.val(terminal_prompt());
				g_terminal_input.removeAttr("readonly");
			}
		}
	).done(function(){
		g_terminal_input.focus();
	});
}

function open_file_api(file){
	$.ajax(
		{
			type: "GET",
			url: "/spinner/webui/open",
			data: {h: g_cur_hostname, file: file},
			timeout: 60000,
			success: function(data){
				var mode, m = /\.([a-z]+)$/.exec(file);
				if (m == null) {
					mode = "plain_text";
				} else {
					switch (m[1]) {
						case "c":
						case "cpp":
						case "h":
						case "hpp":
							mode = "c_cpp";
							break;
						case "go":
							mode = "golang";
							break;
						case "pb":
							mode = "protobuf";
							break
						case "json":
						case "lua":
						case "python":
						case "sh":
						case "xml":
						case "yaml":
							mode = m[1];
							break;
						default:
							mode = "plain_text";
					}
				}
				g_filesystem.find("h3").text(file);
				g_file_editor.getSession().setMode("ace/mode/" + mode);
				g_file_editor.setValue(data);
			}
		}
	);
}
