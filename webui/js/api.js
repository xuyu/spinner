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
		$("#datacenter-name").text(data.Name);
		$("#datacenter-location").text(data.Location);
		var groups = $("#dt-groups");
		$.each(data.Groups, function(index, gp) {
			var html = "<ul class='col-xs-6'><h3 class='dt-group-head'>" + gp.Name;
			html = html + "<span>" + gp.Machines.length + "</span></h3>";
			$.each(gp.Machines, function(index, ma){
				var cla = "";
				if (ma.KeepAlive == 0 || ma.KeepAlive == "" || Math.abs(ma.KeepAlive - (new Date()).getTime()/1000) > 300) {
					cla = "miss";
				} else {
					cla = "alive";
				}
				html = html + "<li class='" + cla + "'><a href='#'>" + ma.IP + "@" + ma.Hostname + "</a></li>"
			});
			html = html + "</ul";
			groups.html(html);
			$(".dt-group-head").click(trigger_ul_head_click);
		})
	})
}