$(function() {
	function toggle(event) {
		var self = event.target
		$('#sider a.current').removeClass('current')
		$(self).addClass('current')
		// 如果是目录，展开/或关闭其下的文件
		if($(self).attr("data-IsDir") === "true") {
			if ($(self).siblings('ul').is(":visible")) {
				$(self).siblings('ul').hide().end()
				return
			}
			$.ajax({
				type: "GET",
				url: "/nav/sub",
				data: {"dir": $(self).attr("data-Path")},
				success: appendSub(self),
				dataType: 'json'
			});
		} else {
			$.ajax({
				type: "GET",
				url: "/main/view",
				data: {"file": $(self).attr("data-Path")},
				success: viewFile(self),
				dataType: 'text'
			});
		}
	}

	function appendSub(self) {
		return function(data) {
			var ul = '<ul>'
			for(var i = 0; i < data.length; i++) {
				ul += '<li><a href="#" data-IsDir="' + data[i].IsDir
				+ '" data-Path="'+ data[i].Path + '">' + data[i].Name + '</a></li>'
			}
			ul += "</ul>"
			$(self).siblings('ul').remove()
			$(self).parent().append(ul).siblings().children("ul").hide()
			$(self).parent().find('ul>li>a').click(toggle)
		}
	}

	function viewFile(self) {
		return function(data) {
			//console.log(typeof(data))
			$('#content').val(data)
		}
	}

	$('#sider a').click(toggle)

	function saveFileAjax() {
		if($(self).attr("data-IsDir") === "true") {
			alert("selected file is directory")
			return
		}
		$.ajax({
			type: "POST",
			url: "/main/save?file=" + $('#sider a.current').attr('data-path'),
			processData:false,
			data: $('#content').val(),
			success: saveFile,
			dataType: "text",
			contentType: "text/plain"
		});
	}

	function saveFile(data) {
		alert(data)
	}
	$('#content').keydown(function(e) {
		// tab is tab
		if (e.which == 9) {
			e.preventDefault();
			var start = this.selectionStart, end = this.selectionEnd;
			var text = this.value;
			var tab = '\t';
			text = text.substr(0, start) + tab + text.substr(start);
			this.value = text;
			this.selectionStart = start + tab.length;
			this.selectionEnd = end + tab.length;
		}
	})
	// 捕获键盘事件
	$(document).keydown(function(event) {
		// console.log(event.which)
		// Ctrl + S
		if (!(String.fromCharCode(event.which).toLowerCase() == 's' && event.ctrlKey)) return true;
		event.preventDefault();

		saveFileAjax()

		return false;
	})

})
