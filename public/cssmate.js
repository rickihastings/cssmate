(function() {
	/**
	 * cssmate
	 * 
	 * Ricki Hastings (c) 2014 - https://github.com/rickihastings
	 * https://github.com/rickihastings/cssmate
	 */

	var connection;
	var files = [];
	var tags = {};
	var config = {
		disabled: false,
		hostname: '10.0.33.34',
		port: 3000,
	}

	function init() {
		getCSSDeps();
		websocketConnect();
	}

	function getCSSDeps() {
		var links = document.getElementsByTagName('link');

		for (var i = 0, len = links && links.length; i < len; ++i) {
			var link = links[i],
				filename = link.href.replace(link.baseURI, '');
			
			tags[filename] = link;
			files.push(filename);
		}
		// find out css dependencies
	}

	function websocketConnect() {
		connection = new WebSocket('ws://' + config.hostname + ':' + config.port + '/websocket');

		connection.onopen = function() {
			connection.send(JSON.stringify({command: 'watching', data: files.join(',')}));
		}
		// when the connection is open, send some data to the server

		connection.onerror = function (error) {
			websocketConnect();
		}
		// log errors

		connection.onmessage = function (e) {
			var json;

			try {
				json = JSON.parse(e.data);
			} catch (e) {
				console.log('Failed to parse output, your cssmate binary is botched!');
			}

			if (!json || !json.command || !json.data) {
				return;
			}

			switch (json.command) {
				case 'watching':
					if (!json.data) {
						config.disabled = true;
					}
					break;
				case 'changed':
					reloadStylesheet(json.data);
					break;
			}
		}
		// log messages from the server
	}

	function reloadStylesheet(filename) {
		if (!tags[filename]) {
			console.log('Apparently', filename, 'has been changed but we don\'t have any record of having it?')
			return;
		}

		tags[filename].href = filename + '?id=' + new Date().getMilliseconds();
		// we whack a timestamp on the end so the browser knows its different and reloads it
	}

	init();
})();