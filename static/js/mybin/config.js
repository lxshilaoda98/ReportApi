$(function () {
    //自动连接SIP服务器
    //let wsUri = "ws://192.168.0.109:5066"
//            let sipUrl = "1998@192.168.0.109:5060"
//            let sipPwd = "1234"
    let ctiServer = localStorage["CTIServerAddress"];
    let sipName = localStorage["DeviceName"];
    if (ctiServer == undefined || ctiServer == "" || sipName == "" || sipName == undefined) {
        console.log('配置为空.请检查')
        alert("配置为空.请检查配置");
    }else{
        let wsUri = `ws://${ctiServer}:5066`; //ws://xxx:5066   wss://xxx:7443
        let sipUrl = `${sipName}@${ctiServer}:5060`;
        let sipPwd = "1234";
        console.log('初始化配置')
        onConn(wsUri, sipUrl, sipPwd);
    }


    $('#Config').click(function () {

        layui.use(['layer', 'form', 'element'], function () {

            var $ = layui.jquery
                , layer = layui.layer
                , form = layui.form
                , element = layui.element;

            var index = layer.open({
                type: 1 //0（信息框，默认）1（页面层）2（iframe层）3（加载层）4（tips层）
                ,
                id: 'layer_1' //id
                ,
                title: '软电话配置' //标题
                ,
                shift: 2
                ,
                content: '<div class="layui-form layui-form-pane" lay-filter="form" action=""><div class="layui-form-item">' +

                '<label class="layui-form-label">Server地址:</label>' +
                '<div class="layui-input-block">' +
                '<input type="text" id="ctiserver" name="ctiserver" lay-verify="title" autocomplete="off" placeholder="ctiserver地址" class="layui-input" >' +
                '</div>' +

                '</div>' +
                '<div class="layui-form-item">' +

                '<label class="layui-form-label">分机号:</label>' +
                '<div class="layui-input-block">' +
                '<input type="text" id="devicename" name="devicename" lay-verify="title" autocomplete="off" placeholder="分机号" class="layui-input">' +
                '</div>' +

                '</div>' +

                '<div class="layui-form-item">' +

                '<label class="layui-form-label">自动就绪间隔(秒):</label>' +
                '<div class="layui-input-block">' +
                '<input type="text" id="AutoWrapEndTime" name="AutoWrapEndTime" lay-verify="title" autocomplete="off" placeholder="自动就绪间隔" class="layui-input">' +
                '</div>' +

                '</div>' +
                '<div class="layui-form-item">' +

                '<label class="layui-form-label">保存日志:</label>' +
                '<div class="layui-input-block">' +
                ' <input type="checkbox" id="savelog" name="savelog" lay-skin="switch" lay-text="ON|OFF">' +
                '</div>' +

                '</div>' +
                '</div>'
                //内容
                ,
                area: ['500px', '400px'] //宽高
                ,
                shade: [0.8, '#393D49'] //遮罩 不显示false
                ,
                shadeClose: false //是否点击遮罩关闭 默认false
                ,
                closeBtn: 1 //右上角关闭按钮样式 默认1 配置1和2 不显示0
                ,
                btn: ['确定', '取消', '下载日志', '清空日志'] //按钮
                ,
                btnAlign: 'c' //按钮排列 l（左对齐）c（居中对齐）r（右对齐。默认值）
                ,
                success: function () {
                    form.render();
                }
                ,
                btn1: function (index, layero) {
                    //按钮【按钮一】的回调
                    console.log('config...btnOK');

                    localStorage["CTIServerAddress"] = $('#ctiserver').val();
                    localStorage["DeviceName"] = $('#devicename').val();
                    localStorage["AutoWrapEndTime"] = $('#AutoWrapEndTime').val();

                    console.log('config...savelog->' + $('#savelog')[0].checked);
                    localStorage["SaveH5Log"] = $('#savelog')[0].checked;


                    layer.close(index);
                    window.location.reload()
                    return true;
                }
                ,
                btn2: function (index, layero) {
                    //按钮【按钮二】的回调
                    layer.close(index);

                    //return false 开启该代码可禁止点击该按钮关闭
                }
                ,
                btn3: function (index, layero) {
                    //layer.close();
                    l2i.download();
                }
                ,
                btn4: function (index, layero) {
                    //layer.close();
                    l2i.clear();

                    layer.msg('日志已经清空', {

                        time: 2000 //2秒关闭（如果不配置，默认是3秒）
                    });

                    return false;
                }
                ,
                cancel: function () {
                    //右上角关闭回调
                    //return false 开启该代码可禁止点击该按钮关闭
                    //return false;
                }
            });


            form.val('form', {
                "ctiserver": localStorage["CTIServerAddress"] // "name": "value"
                , "devicename": localStorage["DeviceName"] // "name": "value"
                , "AutoWrapEndTime": localStorage["AutoWrapEndTime"] // "name": "value"
                , "savelog": localStorage["SaveH5Log"] == 'false' ? false : true // "name": "value"

            });


        });
    });
});

function trans() {
    console.log('转接选择开始..>');


    $('#outPhone').textbox('setText', '');


    $('#transDiv').dialog({
        buttons: [{
            text: '确定',
            iconCls: 'icon-ok',
            handler: function () {
                transDiv();
            }
        },
            {
                text: '取消',
                iconCls: 'icon-cancel',
                handler: function () {
                    cancel('#transDiv');
                }
            }
        ]
    });
    $('#transDiv').dialog("open");
}

//关闭窗口
function cancel(type) {
    $(type).window('close');
}

//转接事件开始
function transDiv() {
    var outPhone = $('#outPhone').textbox('getText');
    TranNumber(outPhone);
    console.log('外线号码..'+outPhone);
    cancel('#transDiv');
}


