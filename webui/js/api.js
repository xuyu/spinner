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
				html = html + ma.IP + " @ " + ma.Hostname + "</a></li>";
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

function save_file(){
	var file = g_filesystem.find("h3>span").text();
	if (file == null || file == undefined || file.length == 0) {
		return;
	}
	var url = "/spinner/webui/save?h=" + encodeURIComponent(g_cur_hostname) + "&file=" + encodeURIComponent(file);
	$.ajax(
		{
			type: "POST",
			url: url,
			data: g_file_editor.getValue(),
			beforeSend: function(){
				g_filesystem.find("h3>a>img").addClass("moved");
			},
			complete: function(){
				setTimeout(function(){
					g_filesystem.find("h3>a>img").removeClass("moved");
				}, 500);
			}
		}
	);
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
				g_filesystem.find("h3>span").text(file);
				g_file_editor.setValue(data);
				g_file_editor.getSession().setMode("ace/mode/" + mode);
			}
		}
	);
}

function float2p(f){
	return Math.round(f * 100) / 100;
}

function readableFileSize(size) {
    var units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    var i = 0;
    while(size >= 1024) {
        size /= 1024;
        ++i;
    }
    return size.toFixed(1) + ' ' + units[i];
}

function cpu_usage(t1, t2){
	var u = Math.round((t2.user - t1.user) / (t2.total - t1.total) * 10000) / 100;
	var s = Math.round((t2.system - t1.system) / (t2.total - t1.total) * 10000) / 100;
	var i = Math.round((t2.iowait - t1.iowait) / (t2.total - t1.total) * 10000) / 100;
	return [u, s, i];
}

function diskio_speed(t1, t2, duration){
	var rincr = 0, wincr = 0;
	$.each(t2, function(name, v2){
		var v1 = t1[name];
		if (v1 == undefined) return;
		if (/^[0-9]+$/.test(name[name.length - 1])) {
			rincr = rincr + (v2.rbytes - v1.rbytes);
			wincr = wincr + (v2.wbytes - v1.wbytes);
		}
	});
	var rp = Math.round(rincr / duration);
	var wp = Math.round(wincr / duration);
	return [rp, wp];
}

function netio_speed(t1, t2, duration){
	var rincr = 0, sincr = 0;
	$.each(t2, function(name, v2){
		var v1 = t1[name];
		if (v1 == undefined) return;
		rincr = rincr + (v2.rbytes - v1.rbytes);
		sincr = sincr + (v2.sbytes - v1.sbytes);
	});
	var rp = Math.round(rincr / duration);
	var sp = Math.round(sincr / duration);
	return [rp, sp];
}

function dashboard_show_cputimes(times){
	var tds = g_dashboard_cputimes.find("table tr:last-child>td");
	$(tds[0]).text(times[0]);
	$(tds[1]).text(times[1]);
	$(tds[2]).text(times[2]);
}

function dashboard_show_meminfo(meminfo){
	var trs = g_dashboard_meminfo.find("table tr");
	$(trs[0]).find("td:last-child").text(readableFileSize(meminfo.free) + "/" + readableFileSize(meminfo.total));
	$(trs[1]).find("td:last-child").text(readableFileSize(meminfo.buffers) + "/" + readableFileSize(meminfo.cached));
	$(trs[2]).find("td:last-child").text(readableFileSize(meminfo.sfree) + "/" + readableFileSize(meminfo.stotal));
}

function dashboard_show_diskusage(usage){
	var table = g_dashboard_diskusage.find("table");
	var html = "<tr><td>mount</td><td>size</td><td>use%</td></tr>";
	$.each(usage, function(mount, value){
		var tr = "<tr><td>" + mount + "</td><td>" + readableFileSize(value[0]) + "</td><td>";
		tr = tr + Math.round((value[0] - value[1]) *100 / value[0]) + "%</td></tr>";
		html = html + tr;
	});
	table.html(html);
}

function dashboard_show_boottime(btime){
	g_dashboard_btime.find("span").text((new Date(btime * 1000)).toLocaleString());
}

function dashboard_show_load(c1, c2, load){
	var spans = g_dashboard_load.find("span");
	$(spans[0]).text(c1 + "/" + c2);
	$(spans[1]).text(float2p(load[0]));
	$(spans[2]).text(float2p(load[1]));
	$(spans[3]).text(float2p(load[2]));
}

function dashboard_show_diskio(speed){
	g_dashboard_diskio.find("table tr:first-child>td:last-child").text(readableFileSize(speed[0]) + "/s");
	g_dashboard_diskio.find("table tr:last-child>td:last-child").text(readableFileSize(speed[1]) + "/s");
}

function dashboard_show_netio(speed){
	g_dashboard_netio.find("table tr:first-child>td:last-child").text(readableFileSize(speed[0]) + "/s");
	g_dashboard_netio.find("table tr:last-child>td:last-child").text(readableFileSize(speed[1]) + "/s");
}

var g_dashboard_last_data = null;
var g_dashboard_last_data_ts = null;
function dashboard_show_next(data){
	if (g_dashboard_last_data == null) {
		g_dashboard_last_data = data;
		g_dashboard_last_data_ts = Math.round((new Date()).getTime() / 1000);
		return;
	}
	var duration = Math.round((new Date).getTime() / 1000) - g_dashboard_last_data_ts;
	dashboard_show_cputimes(cpu_usage(g_dashboard_last_data.cputimes, data.cputimes));
	dashboard_show_diskio(diskio_speed(g_dashboard_last_data.diskio, data.diskio, duration));
	dashboard_show_netio(netio_speed(g_dashboard_last_data.netio, data.netio, duration));
	g_dashboard_last_data = data;
}

var g_dashboard_runner = null;
function dashboard_api(){
	$.getJSON("/spinner/webui/dashboard", {h: g_cur_hostname}, function(data){
		dashboard_show_boottime(data.btime);
		dashboard_show_load(data.physicpu, data.logicpu, data.load);
		dashboard_show_meminfo(data.meminfo);
		dashboard_show_diskusage(data.diskusage);
		dashboard_show_next(data);
		if (g_dashboard_runner == null) {
			g_dashboard_runner = setInterval(function(){
				dashboard_api();
			}, 2000);
		}
	});
}
