JsSIP.C.SESSION_EXPIRES = 120, JsSIP.C.MIN_SESSION_EXPIRES = 120;

// SIP用户
var _sipUserAgent = null;
// SIP对话
var _sipUserSession = null;
// 通话数据信息
var _newRTCSessionData = null;
// 远端视频流
var _sipRemoteStream = null;


// SIP通话-本地视频展示控件
//var _sipSelfView = null;

// SIP通话-远程视频展示控件
var _sipRemoteMedia = null;

//是否自动接听
var sipAutoAnswer = true;


/**
 * 创建SIP连接
 * @param wsUri            WebSocket连接地址(5066或7443):ws://10.10.10.123:5066
 * @param sipUri        Sip登录地址:sip:1003@10.10.10.123:5060
 * @param sipPwd        Sip登录密码:1234
 * @param callBack        回调函数
 * @returns
 */
function sipOnConn(wsUri, sipUri, sipPwd, callBack) {

    console.log('服务器wsUrl'+wsUri)

    if (!callBack) {
        callBack = new Object();
    }

    var socket = new JsSIP.WebSocketInterface(wsUri);

    var configuration = {
        sockets: [socket],
        outbound_proxy_set: wsUri,
        uri: sipUri,
        password: sipPwd,
        //contact_uri : sipUri + ";transport=wss",
        register: true,
        noAnswerTimeout: 120,
        //session_timers: false
    };

    // SIP服务器连接
    _sipUserAgent = new JsSIP.UA(configuration);

    // SIP服务器连接成功
    _sipUserAgent.on('connecting', function (data) {
        buttonDisabled('初始化')
        console.log("=============SIP服务器连接中: ", data);
    });

    // SIP服务器连接成功
    _sipUserAgent.on('connected', function (data) {
        buttonDisabled('成功')
        console.log("=============SIP服务器连接成功: ", data);

    });

    // SIP服务器断开连接
    _sipUserAgent.on('disconnected', function (data) {
        buttonDisabled('失败')
        console.log("=============SIP服务器断开连接, ", data);
        console.log("=============SIP服务器断开连接, Reason.>", data.error.reason);
        console.log("=============SIP服务器断开连接, code.>", data.error.code);
    });

    //1
    _sipUserAgent.on('registered', function (data) {
        buttonDisabled('成功')
        console.log("=============SIP成功注册成功: ", data);
        console.log("=============SIP成功注册成功: Response.>", data.Response);

    });

    //1
    _sipUserAgent.on('unregistered', function (data) {
        buttonDisabled('失败')
        console.log("=============SIP被解雇注册: ", data);
        console.log("=============SIP被解雇注册: Response.>", data.Response);
        console.log("=============SIP被解雇注册: Cause.>", data.Cause);

    });

    _sipUserAgent.on('registrationFailed', function (data) {
        buttonDisabled('失败')
        console.log("=============SIP由于注册失败而被解雇: ", data);
        console.log("=============SIP由于注册失败而被解雇: Response.>", data.Response);
        console.log("=============SIP由于注册失败而被解雇: Cause.>", data.Cause);

    });

    _sipUserAgent.on('registrationExpiring', function (data) {
        console.log("=============SIP..registrationExpiring: ", data);
        // 1.在注册到期之前发射几秒钟。如果应用程序没有为这个事件设置任何监听器，JsSIP将像往常一样重新注册。
        //
        // 2.如果应用程序订阅了这个事件，它负责ua.register()在registrationExpiring事件中调用（否则注册将过期）。
        //
        // 3.此事件使应用程序有机会在重新注册之前执行异步操作。对于那些在REGISTER请求中的自定义SIP头中使用外部获得的“令牌”的环境很有用。

    });

    _sipUserAgent.on('RTCDataChannel', function (data) {
        console.log("=============SIP..RTCDataChannel: ", data);
        // 1.在注册到期之前发射几秒钟。如果应用程序没有为这个事件设置任何监听器，JsSIP将像往常一样重新注册。
        //
        // 2.如果应用程序订阅了这个事件，它负责ua.register()在registrationExpiring事件中调用（否则注册将过期）。
        //
        // 3.此事件使应用程序有机会在重新注册之前执行异步操作。对于那些在REGISTER请求中的自定义SIP头中使用外部获得的“令牌”的环境很有用。

    });


    // SIP用户事件监听
    _sipUserAgent.on('newRTCSession', function (data) {
        console.log("=====================newRTCSession: ", data);

        _newRTCSessionData = data;
        _sipUserSession = data.session;

        /* 来电自动接听 */
        if (data.originator == "remote") {
            //接听
            if (sipAutoAnswer) {
                sipOnAnswer();
            }

            //来电话响铃中回调通知
            if (callBack.ringing) {

                callBack.ringing(data);
                //开启振铃
                Ring.startRingTone();
            }
        }

        data.session.on("accepted", function (e) {
            var pc = _sipUserSession.connection;
            // 获取远端视频流
            _sipRemoteStream = new MediaStream();
            pc.getReceivers().forEach(function (receiver) {
                _sipRemoteStream.addTrack(receiver.track);
            });

            _sipRemoteMedia = document.getElementsByName("sipRemoteView")[0];
            _sipRemoteMedia.srcObject = _sipRemoteStream;
            _sipRemoteMedia.play();

            //执行回调
            if (callBack.accepted) {
                callBack.accepted(data, e);
            }
        });

        data.session.on("sdp", function (e) {
            if (callBack.sdp) {
                callBack.sdp(data, e);
            }
        });

        data.session.on("connecting", function (e) {
            //执行回调
            if (callBack.connecting) {
                callBack.connecting(data, e);
            }
        });

        data.session.on("progress", function (e) {
            //来电或去电前期回调通知
            if (callBack.progress) {
                callBack.progress(data, e);
                buttonDisabled('来电或去电');

                //开启振铃
                Ring.startRingTone();
            }
        });

        data.session.on("confirmed", function (e) {
            //通话接听后回调通知
            if (callBack.confirmed) {
                callBack.confirmed(data, e);

                buttonDisabled('通话接听');

                //关闭振铃
                Ring.stopRingTone();
            }

        });

        data.session.on("failed", function (e) {
            //呼叫失败或未接通挂断回调通知
            if (callBack.failed) {
                callBack.failed(data, e);
                buttonDisabled('呼叫失败或未接通挂断');
                //关闭振铃
                Ring.stopRingTone();
            }
        });

        data.session.on("started", function (e) {
            //执行回调
            if (callBack.started) {
                callBack.started(data, e);
            }
        });

        data.session.on("ended", function (e) {
            //接通后正常挂断回调通知
            if (callBack.ended) {
                callBack.ended(data, e);

                buttonDisabled('接通后正常挂断');
            }
        });

        data.session.on("hold", function (e) {

            buttonDisabled('保持');

        })

        data.session.on("unhold", function (e) {

            buttonDisabled('取消保持');

        })


    });


    _sipUserAgent.start();

}


/*SIP呼叫回调函数示例*/
var __sipCallBack = {
    'ringing': function (data) {
        console.log('=============来电话响铃中回调通知-ringing: ', data);
    },
    'accepted': function (data, e) {
        console.log('=============accepted: ', data, e);
    },
    'connecting': function (data, e) {
        console.log('=============connecting: ', data, e);
    },
    'sending': function (data, e) {
        console.log('=============sending: ', data, e);
    },
    'progress': function (data, e) {
        console.log('=============来电或去电前期回调通知-progress: ', data, e);
    },
    'confirmed': function (data, e) {

        console.log('=============通话接听后回调通知-confirmed: ', data, e);

    },
    'failed': function (data, e) {
        console.log('=============呼叫失败或未接通挂断回调通知-failed: ', data, e);
    },
    'started': function (data, e) {
        console.log('=============started: ', data, e);
    },
    'hold': function (data, e) {
        console.log('=============hold.', data, e);
    },
    'ended': function (data, e) {
        console.log('=============接通后正常挂断回调通知-ended: ', data, e);
    }

};


var _sipEventHandlers = {
	'progress' : function(e) {
		console.log('=============正在呼叫....');
	},
	'failed' : function(e) {
		console.log('=============呼叫失败: ', e);
	},
	'ended' : function(e) {
		console.log('=============呼叫结束: ', e);
	},
	'confirmed' : function(e) {
		console.log('=============confirmed接通: ', e);
	},
	'addstream' : function(e) {
		console.log('=============addstream');
	},
	'succeeded' : function(e) {
		console.log('=============succeeded');
	},
	'failed' : function(e) {
		console.log('=============failed: ', e);
	}
};

/* SIP拨打电话 */
function sipOnCall(user) {
    if (!_sipUserAgent)
        return;

    var options = {
        'eventHandlers' : _sipEventHandlers,
        'mediaConstraints': {'audio': true, 'video': false}
    };
    _sipUserSession = _sipUserAgent.call(user, options);
}

/* SIP电话接听 */
function sipOnAnswer() {
    if (_sipUserSession) {
        var options = {
            //'eventHandlers' : _sipEventHandlers,
            'mediaConstraints': {'audio': true, 'video': false, 'mandatory': {'maxWidth': 640, 'maxHeight': 360}}
        };
        _sipUserSession.answer(options);
    }
}

/**
 * 拒接听
 */
function sipHangUp() {

    _sipUserSession.terminate();

}


/* SIP电话挂断 */
function sipOnHangUp() {
    if (_sipUserSession) {
        _sipUserSession.terminate();
    }
}

/**
 * 发送DTMF
 * @param num
 */
function sipSendDTMF(num) {
    if (!_sipUserSession)
        return;

    var extraHeaders = [ 'X-Foo: foo', 'X-Bar: bar' ];
    var options = {
        'mediaConstraints': {'duration': 160, 'interToneGap': 1200,'extraHeaders': extraHeaders}
    };

    _sipUserSession.sendDTMF(num,options);
}

/**
 * 保持
 */
function sipOnHole() {
    if (!_sipUserSession)
        return;

    _sipUserSession.hold();
}
/**
 * 取消保持
 */
function sipUnHole() {
    if (!_sipUserSession)
        return;

    _sipUserSession.unhold();
}

/**
 * 转接盲转
 * @param Number
 * @constructor
 */
function TranNumber(Number) {
    if (!_sipUserSession)
        return;

    _sipUserSession.refer(Number)
}

/* SIP发送消息 */
function sipOnSendMsg(user, msg) {
    if (!_sipUserAgent)
        return;

    var options = {
        //"eventHandlers" : _sipEventHandlers
    };
    _sipUserAgent.sendMessage(user, msg, options);
}


/**
 * 获取iFrame中所有的视频控件
 * @returns
 */
function getAllSipRemoteView() {
    var allSipRemoteViews = [];
    var iframes = document.getElementsByTagName("iframe");
    for (var i = 0; i < iframes.length; i++) {
        var sipRemoteViews = iframes[i].contentDocument.getElementsByName("sipRemoteView");
        if (sipRemoteViews.length > 0) {
            allSipRemoteViews.push.apply(allSipRemoteViews, sipRemoteViews);
        }
    }
    return allSipRemoteViews;
}

/**
 * 判断是否是视频
 * @returns
 */
function isSipVideo() {
    if (_newRTCSessionData && _newRTCSessionData.request && _newRTCSessionData.request.sdp && _newRTCSessionData.request.sdp.media) {
        if (_newRTCSessionData.request.sdp.media.length > 1) {
            return true;
        }
    }
    return false;
}


/**
 * 判断按钮显示隐藏
 * @param StatusName
 */
function buttonDisabled(StatusName) {
    console.log('StatusName..>' + StatusName)
    switch (StatusName) {
        case '初始化' :
            butColor('callButton', '禁用');//拨打按钮禁用
            butColor('aswButton', '禁用');//接听按钮禁用
            butColor('hanguButton', '禁用');//挂断按钮禁用
            butColor('holdButton', '禁用');//保持 禁用
            butColor('oNHoleButton', '禁用');//取消保持 禁用
            butColor('transButton', '禁用');//转接 禁用
            butColor('BhpButton', '禁用');//拨号盘按钮禁用


            break;
        case '成功' :
            butColor('callButton', '启用');//拨打按钮启用
            break;
        case '失败' :
            butColor('callButton', '禁用'); //拨打按钮禁用
            butColor('aswButton', '禁用');//接听按钮禁用
            break;
        case '通话接听' :
            butColor('callButton', '禁用'); //拨打按钮禁用
            butColor('hanguButton', '启用');//挂断按钮启用
            butColor('holdButton', '启用');//保持按钮启用
            butColor('transButton', '启用');//转接按钮启用
            butColor('BhpButton', '启用');//拨号盘按钮启用
            break;
        case '接通后正常挂断':
            butColor('callButton', '启用'); //拨打按钮启用
            butColor('hanguButton', '禁用');//挂断按钮禁用
            butColor('holdButton', '禁用');//保持按钮禁用
            butColor('transButton', '禁用');//转接按钮禁用
            butColor('BhpButton', '禁用');//拨号盘按钮禁用
            $("#holdButton").text('保持');
            break;
        case '呼叫失败或未接通挂断':
            butColor('callButton', '启用'); //拨打按钮启用
            butColor('hanguButton', '禁用');//挂断按钮禁用
            butColor('holdButton', '禁用');//保持按钮禁用
            butColor('transButton', '禁用');//转接按钮禁用
            butColor('BhpButton', '禁用');//拨号盘按钮禁用
            $("#holdButton").text('保持');
            break;

        case '保持':
            butColor('callButton', '禁用'); //拨打按钮禁用
            butColor('hanguButton', '启用');//挂断按钮禁用
            butColor('transButton', '禁用');//转接按钮禁用
            butColor('BhpButton', '禁用');//拨号盘按钮禁用
            $("#holdButton").text('取消保持'); // 只支持修改文本
            break;
        case '取消保持':
            butColor('callButton', '禁用'); //拨打按钮禁用
            butColor('hanguButton', '启用');//挂断按钮禁用
            butColor('transButton', '启用');//转接按钮禁用
            butColor('BhpButton', '启用');//拨号盘按钮启用
            $("#holdButton").text('保持'); // 只支持修改文本
            break;
        case '来电或去电' :
            butColor('callButton', '禁用');//拨打按钮禁用
            butColor('aswButton', '禁用');//接听按钮禁用
            butColor('hanguButton', '启用');//挂断按钮禁用
            butColor('holdButton', '禁用');//保持 禁用
            butColor('oNHoleButton', '禁用');//取消保持 禁用
            butColor('transButton', '禁用');//转接 禁用
            butColor('BhpButton', '禁用');//拨号盘按钮禁用


            break;
        default:
            break;

    }

}

/**
 * 控制按钮显示隐藏+样式
 * @param butOid
 * @param butType
 */
function butColor(butOid, butType) {
    switch (butType) {
        case '启用':
            $('#' + butOid + '').attr("disabled", false);
            $('#' + butOid + '').css("background-color", "#4CAF50");
            break;
        case '禁用':
            $('#' + butOid + '').attr("disabled", true);
            $('#' + butOid + '').css("background-color", "#DCDCDC");
            break;
    }

}


