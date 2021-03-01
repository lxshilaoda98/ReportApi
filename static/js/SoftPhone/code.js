var ws;//socket
var extension;//分机号
var CallId;//主叫id
var CalledId;//被叫id
var HttpUrl="http://127.0.0.1:3333"; //访问http接口使用的地址
var socketUrl = "ws://127.0.0.1:4444"; //一般情况ws 的地址就是http接口的地址
var CallDirection;//呼叫方向，用来判断主叫和被叫
var zxCallId;
var zxCalledId;

var CAni; //来电的号码

$(document).ready(function () {
    extension=localStorage["extension"];
    console.log("connection Esl Node Socket ..>");
    console.log("socket URl ...>" + socketUrl)
    ws = new WebSocket(socketUrl);
//申请一个WebSocket对象，参数是服务端地址，同http协议使用http://开头一样，WebSocket协议的url使用ws://开头，另外安全的WebSocket协议使用wss://开头
    ws.onopen = function () {
        //当WebSocket创建成功时，触发onopen事件
        console.log("open.socket.>Extension=>"+extension);
        var SendJson = {
            MSG:"Login",
            Data:{
                "extension":extension,
                "AgentId":extension,//"1999",
            }
        }
        var str = JSON.stringify(SendJson);
        ws.send(str); //将消息发送到服务端
        //初始化页面按钮信息
        ButtonInit();
    }
    ws.onmessage = function (e) {
        console.log("Response Server Message ..e>"+e);
        //当客户端收到服务端发来的消息时，触发onmessage事件，参数e.data包含server传递过来的数据
        console.log("Response Server Message ..e.data>"+e.data);
        let eventName = JSON.parse(e.data);
        switch (eventName["Response"]) {
            case "Login":
                if(eventName["success"]=="true"){
                    console.log(eventName["message"])
                }else{
                    alert("登录失败了:"+eventName["message"]);
                    ButtonVisForStatus("登录失败")
                }
                break;
            case "MakeCall":
                if(eventName["success"]=="true"){
                    CallId=eventName["CallId"]
                    console.log("MakeCall返回成功,回去callid..>"+eventName["CallId"])
                    if(CallId==undefined){
                        console.log("回铃未摘")
                        ButtonVisForStatus("回铃未摘");
                        return;
                    }
                    ButtonVisForStatus("外呼电话");

                }else{
                    CallId=null
                }
                break;
            case "MakeCallAgent":
                if(eventName["success"]=="true"){
                    CallId=eventName["CallId"]
                    console.log("MakeCallAgent返回成功,回去callid..>"+eventName["CallId"])
                    if(CallId==undefined){
                        console.log("回铃未摘")
                        ButtonVisForStatus("回铃未摘");
                        return;
                    }
                    ButtonVisForStatus("外呼坐席");


                }else{
                    CallId=null
                }
                break;
            case "HoldCall":
                if(eventName["success"]=="true"){
                    $("#HoldVis").hide(); //隐藏保持按钮
                    $("#UnHoldVis").show();
                    ButtonVisForStatus("保持状态");
                }
                break;
            case "UnHoldCall":
                if(eventName["success"]=="true"){
                    $("#UnHoldVis").hide(); //隐藏取消保持按钮
                    $("#HoldVis").show();
                    ButtonVisForStatus("服务状态")
                }
                break;
            case "CallHangup":
                if(eventName["success"]=="true"){
                     console.log("DisPhone..clear.>.CallId>"+ CallId);
                     CallId=null
                }
                break;
            case "AgentChange":
                if (eventName["message"]=="+OK\n" && eventName["success"]=="true"){
                    // if (eventName["status"]=="Available"){
                    //     ButtonVisForStatus("空闲状态")
                    // }else if (eventName["status"]=="On Break"){
                    //     ButtonVisForStatus("小休状态")
                    // }
                    console.log("服务端返回，切换状态成功！")
                }

                break;
            case "EVENT":
                console.log("收到ESL EVENT 事件,事件:"+eventName["EventName"]);
                if (eventName["EventName"]=="CHANNEL_CALLSTATE"){
                    CAni=eventName["CallerANI"];

                    CallDirection="out";
                    //来电话了，可能不是callcenter过来的数据.
                    CallId =eventName["CallerUniqueID"]
                    CalledId = eventName["ChannelCallUUID"]
                }else if(eventName["EventName"]=="bridge-agent-start"){


                    CallId =eventName["CCMemberSessionUUID"]
                    CalledId =eventName["CCAgentUUID"]

                    zxCallId=eventName["CCMemberSessionUUID"];
                    zxCalledId=eventName["CCAgentUUID"];
                }else if(eventName["EventName"] =="CHANNEL_ANSWER"){
                    //如果呼叫方向为out. 这样的话可能是坐席外呼。这样保持.满意度.转接.等。都需要将被叫转出去。
                    // （设置callid=CallerUniqueID）calledid =ChannelCallUUID
                    if (CallDirection=="out"){
                        CallId = eventName["CallerUniqueID"]
                        CalledId =eventName["ChannelCallUUID"]
                        console.log("本次可能是外呼应答..CallId:"+CallId+"..>>CalledId:"+CalledId)
                    }
                }else if (eventName["EventName"] =="CeeCallUpdate"){
                    console.log("呼叫中心接话，切换callID")
                    //CalledId =eventName["CalledId"]
                }else if (eventName["EventName"] =="agent-offering"){
                    CAni=eventName["CCMemberCIDNumber"];
                }
                EventHandle(eventName["EventName"]);
                break;
            default:
                console.log(e.data)
                break;
        }

    }
    ws.onclose = function (e) {
        //当客户端收到服务端发送的关闭连接请求时，触发onclose事件
        console.log("close");
    }
    ws.onerror = function (e) {
        //如果出现连接、处理、接收、发送数据失败的时候触发onerror事件
        console.log(e);
    }

});
function EventHandle(EventName){
    //进行event的事件处理方法..
    switch (EventName) {
        case "HEARTBEAT":
            console.log("心跳事件");
            break;
        case "CHANNEL_ANSWER":
            console.log("[maybe]被叫应答了.!");
            ButtonVisForStatus("服务状态")
            break;
        case "CHANNEL_HANGUP_COMPLETE":
            console.log("通话挂断了!");
            //清理本次callid和calledid
            CallId = "";
            CalledId="";
            ButtonVisForStatus("话后状态")
            break;
        case "CHANNEL_CALLSTATE":

            this.CallVLVisible = true;
            console.log("来电话了!"+this.CallVLVisible);
            ButtonVisForStatus("振铃状态")
            break;
        case "agent-offering":
            console.log("CCMemberUUID.>>")
            ButtonVisForStatus("振铃状态")
            break;
        case "bridge-agent-start":
            console.log("CCMemberUUID.>>")
            ButtonVisForStatus("服务状态")
            break;

        default :
            console.log("未知处理事件:"+EventName)
            break;

    }
}


