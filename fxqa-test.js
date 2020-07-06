
(function (global) {
	var ua = navigator.userAgent.toLowerCase();
	isWindows = (ua.indexOf("windows") != -1 || ua.indexOf("win32") != -1);
	isMac = (ua.indexOf("macintosh") != -1 || ua.indexOf("mac os x") != -1);
	isAir = (ua.indexOf("adobeair") != -1);
	isLinux = (ua.indexOf("linux") != -1);
    isAndroid = (ua.indexOf("android") != -1 || ua.indexOf("adr") != -1);

	var _fxqa_port = 0;
	if (isWindows) {
		_fxqa_port = 9091;
	} else if (isMac) {
		_fxqa_port = 9092;
	} else {
		_fxqa_port = 9093;
	}

function fxqa_get_jscode() {
    var js_str = '';
    if(isWindows) {
        jQuery.support.cors = true;
        var nowTime = new Date().getTime();
        $.ajax({
            type: "GET",
            url: "http://127.0.0.1:" + _fxqa_port + "/code",
            //url: fxqa_cacheserver + '/string?key='+fxqa_tester+'_'+fxqa_platform_key+'_PREDEF_JSAPITEST_JSCODE&s='+nowTime,
            dataType: 'json',
            crossDomain: true,
            async: false,
            success: function (data, status, jqXHR) {
                if (data.ret == 0) {
                    js_str = data.code;
                } else {
                    js_str = '';
                }

            },
            error: function (xhr, status, error) {
                //alert(JSON.stringify(xhr));
                js_str = 'FXQAErr:' + error;
            }
        });
    }

    else {
        try {
            xhttp = new XMLHttpRequest();
            if (isAndroid) {
                xhttp.open("GET", "http:///10.0.2.2:" + _fxqa_port + "/code", false);
            }
            else {
                xhttp.open("GET", "http://127.0.0.1:" + _fxqa_port + "/code", false);
            }
            xhttp.send();
            //alert(xhttp.responseText);
            data = JSON.parse(xhttp.responseText);
            if (data.ret == 0) {
                js_str = data.code;
            } else {
                js_str = '';
            }
        } catch (e) {
            alert('get code err:' + e);
        }
    }
    return js_str;
	}


global.fxqa_log = function (err_code, log_str) {
    jQuery.support.cors = true;
    if(isWindows) {
        $.ajax({
            type: "POST",
            url: "http://127.0.0.1:" + _fxqa_port + "/log",
            data: {"str": log_str, "err": err_code},
            dataType: 'html',
            crossDomain: true,
            async: false,
            success: function (data, status, jqXHR) {
//            alert('Test Success End.');
            },
            error: function (xhr, status, error) {
                //alert("Log Print ERR:" + error);
            }
        });
    }
    else {
        var http = new XMLHttpRequest();
        if (isAndroid) {
            var url = "http://10.0.2.2:" + _fxqa_port + "/log";
        }
        else {
            var url = "http://127.0.0.1:" + _fxqa_port + "/log";
        }
        var params = "str=" + log_str + "&err=" + err_code;
        http.open("POST", url, false);
        http.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        http.send(params);
        if (http.status != 200) {
            alert('log send error');
        }
    }
	return;

}

function fxqa_finish() {
    if(isWindows){
        jQuery.support.cors = true;
        $.ajax({
            type: "DELETE",
            url: "http://127.0.0.1:"+_fxqa_port+"/code",
            dataType: 'html',
            crossDomain: true,
            async: false,
            success: function (data, status, jqXHR) {
//            alert('Test Success End.');
            },
            error: function (xhr, status, error) {
                //alert("Log Print ERR:" + error);
            }
        });
    }
    else {
        var http = new XMLHttpRequest();
        if (isAndroid) {
            var url = "http://10.0.2.2:" + _fxqa_port + "/clearcode";

        }
        else {
            var url = "http://127.0.0.1:" + _fxqa_port + "/clearcode";

        }
        http.open("GET", url, false);
        http.send();
        if (http.status != 200) {
            alert('delete code error:' + http.status);
        }
        return;
    }

}

function fxqa_status(action_s) {
	var js_str = '';
	//if(isWindows) {
        jQuery.support.cors = true;
        $.ajax({
            type: "POST",
            url: "http://127.0.0.1:" + _fxqa_port + "/status",
            dataType: 'json',
            data: {"action": action_s},
            crossDomain: true,
            async: false,
            success: function (data, status, jqXHR) {
//            alert('Test Success End.');
                if (action_s == 'check') {
                    js_str = data.status;
                }

            },
            error: function (xhr, status, error) {
                //alert("Log Print ERR:" + error);
            }
        });
   // }

	return js_str;
}

	function fxqa_update_teststatus(action_s) {
		var js_str = '';
		if(isWindows) {
            jQuery.support.cors = true;
            $.ajax({
                type: "POST",
                url: "http://127.0.0.1:" + _fxqa_port + "/test-status",
                dataType: 'json',
                data: {"action": action_s},
                crossDomain: true,
                async: false,
                success: function (data, status, jqXHR) {
//            alert('Test Success End.');
                    if (action_s == 'check') {
                        js_str = data.status;
                    }

                },
                error: function (xhr, status, error) {
                    //	alert("Log Print ERR:" + error);
                }
            });
        }
        else{
            var http = new XMLHttpRequest();
            if(isAndroid){
                var url = "http://10.0.2.2:" + _fxqa_port + "/test-status";
            }
            else {
                var url = "http://127.0.0.1:" + _fxqa_port + "/test-status";
            }
            var params = "action: " + action_s;
            http.open("POST", url, false);
            http.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            http.send(params);
            if (http.status != 200) {
                alert('status send error');
            }
        }
		return js_str;
	}

	function b64DecodeUnicode(str) {
		// Going backwards: from bytestream, to percent-encoding, to original string.
		return decodeURIComponent(atob(str).split('').map(function(c) {
			return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
		}).join(''));
	}

	function b64EncodeUnicode(str) {
		// first we use encodeURIComponent to get percent-encoded UTF-8,
		// then we convert the percent encodings into raw bytes which
		// can be fed into btoa.
		return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g,
			function toSolidBytes(match, p1) {
				return String.fromCharCode('0x' + p1);
			}));
	}

	global.fxqa_init = function () {
		//alert('xx autotest inside, before test start module shoule be init ok.');
		fxqa_status('start');
		setInterval(function () {
			var js_c = fxqa_get_jscode();
			if (js_c != '') {
				//alert(b64DecodeUnicode(js_c));
				//fxqa_log(0, js_c);
				try {
					var js_code = 'try{' + b64DecodeUnicode(js_c) + '}catch(error){fxqa_log(1, error);};';
					eval(js_code);
				} catch(error) {
					fxqa_log(1, error);
				}
				
				fxqa_finish();
			}
		}, 1000);
	}

})(this);