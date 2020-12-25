const vue =new Vue({

    el: '#app',
    data: function () {
        return {
            CAni:"", //弹屏来电的号码
            dialogFormTransfer: false, //盲转dialog
            TransferAni: '',//盲转相关
            StateValue: false,

            zxAniText: '',//咨询相关
            zxDialog: false,
            zxAgentDialog: false, //咨询坐席
            ExtensionVisible: false, //配置分机的dialog
            isAble: true,
            gridClickOid: "",
            dialogVisible: false,
            visible: false,
            dialogFormVisible: false,
            dialogTableVisible: false,
            formLabelWidth: '90px',
            FromTask: '380px',

            form: {
                extension: localStorage["extension"],
                name: '',
                region: '',
                date1: '',
                date2: '',
                delivery: false,
                type: [],
                resource: '',
                desc: ''
            },
            gridData: null
        }
    },
    created() {
        window.addEventListener('beforeunload', e => this.beforeunloadFn(e))
    },
    destroyed() {
        window.removeEventListener('beforeunload', e => this.beforeunloadFn(e))
    },
    methods: {
        //弹屏挂机
        AlertAniHangup(){
            this.dialogVisible = false
            console.log("拒接电话..开始")
            var SendJson = {
                MSG: "HangupCall",
                Data: {
                    "CallId": CallId
                }
            }
            var str = JSON.stringify(SendJson);
            ws.send(str);
        },
        //弹屏接听
        AlertTalk(){
            this.dialogVisible = false
            console.log("摘机开始..发送EvetMessage"+CallId)
            var SendJson = {
                MSG: "TalkCall",
                Data: {
                    "CallId": CallId
                }
            }
            var str = JSON.stringify(SendJson);
            ws.send(str);
        },
        //改变坐席的状态/ 空闲和小休的状态切换
        ChangeState() {//切换空闲或者小休
            var state;
            if (this.StateValue) { //表示要空闲
                state = "Available"
            } else { //表示要小休
                state = "On Break"
            }
            var SendJson = {
                MSG: "AgentChange",
                Data: {
                    "AgentId": extension,
                    "State": state
                }
            }
            var str = JSON.stringify(SendJson);
            console.log("Request Message..>" + str)
            ws.send(str);
        },
        //关闭来电弹屏框
        handleClose(done) {
            this.$confirm('确认关闭？')
                .then(_ => {
                    done();
                    console.log("拒接电话..开始")
                    var SendJson = {
                        MSG: "HangupCall",
                        Data: {
                            "CallId": CallId
                        }
                    }
                    var str = JSON.stringify(SendJson);
                    ws.send(str);
                })
                .catch(_ => {
                });
        },
        //外呼号码
        MakeCallAni() {
            CallDirection = "out"; //设置为外呼的情况
            var SendJson = {
                MSG: "MakeCall",
                Data: {
                    "extension": extension,
                    "CalledNumber": this.form.name,
                }
            }
            var str = JSON.stringify(SendJson);
            ws.send(str);
            console.log(str)
            this.dialogFormVisible = false
            ButtonVisForStatus("回铃状态")
        },
        //开始要外呼到坐席,给socket发送要外呼的信息
        TranAgentClick() {
            if (this.gridClickOid == "") {
                alert("请选择要转给的坐席!")
                return
            }
            var SendJson = {
                MSG: "MakeCallAgent",
                Data: {
                    "extension": extension,
                    "CalledNumber": this.gridClickOid,
                }
            }
            var str = JSON.stringify(SendJson);
            ws.send(str);
            console.log(str)
            this.dialogTableVisible = false;
            ButtonVisForStatus("回铃状态")
        },
        //点击坐席列表后，给一个字段赋值。确认列表选中的值
        handleCurrentChange(val) {
            //console.log(val.name)
            this.gridClickOid = val.name
        },
        //获取空闲的坐席信息..
        async CallAgentList() {
            this.gridClickOid = '';
            this.dialogTableVisible = true;
            var AgentList;
            await axios.get(HttpUrl + "/callcenter/GetAvailableUser?AgentId=" + extension).then(function (response) {
                if (response.data != "无可用坐席") {
                    console.log('找到可用坐席.')
                    AgentList = response.data;
                } else {
                    console.log('无可用空闲坐席.')
                }
            }, function (err) {
                console.log(err)
                this.gridData = "未找到坐席";
            });
            this.gridData = AgentList
        },
        //保持
        HoldCall() {
            var SendJson = {
                MSG: "HoldCall",
                Data: {
                    "CallId": CallId
                }
            }
            var str = JSON.stringify(SendJson);
            ws.send(str);
        },
        //取消保持
        UnHoldCall() {
            var SendJson = {
                MSG: "UnHoldCall",
                Data: {
                    "CallId": CallId
                }
            }
            var str = JSON.stringify(SendJson);
            ws.send(str);
        },
        //处理挂机流程
        CallHangup() {
            console.log("准备挂机")
            var SendJson = {
                MSG: "HangupCall",
                Data: {
                    "CallId": CallId
                }
            }
            var str = JSON.stringify(SendJson);
            ws.send(str);
        },
        //读取分机
        ConfigClick() {
            this.ExtensionVisible = true;
            var GetSessionExtension = localStorage["extension"];
            console.log("GetSessionExtension.." + GetSessionExtension)
            if (GetSessionExtension != undefined) {
                console.log("获取到session的数据为：" + GetSessionExtension)
                $("#ConfigExtension").val(GetSessionExtension);
            }
        },
        //处理分机配置
        ExtensionConfig() {
            var Extension = $("#ConfigExtension").val();
            localStorage["extension"] = Extension;
            this.ExtensionVisible = false;
        },
        //满意度按钮
        async ToSatisfaction() {
            console.log("要转接的callid为：" + CallId)
            var AgentId = extension;
            if (CallId == undefined) {
                alert("无法转接到满意度.")
                return;
            }
            await axios.get(HttpUrl + `/bgapi/fs/fs_transfer_satis?uuid=${CallId}&FileName=Satis.lua&AgentId=${AgentId}&CallDirection=in`)
                .then(function (response) {
                    if (response.data.code == "0") {
                        console.log('成功转入满意度');
                    }
                }, function (err) {
                    console.log('转入失败：' + err)
                })
        },
        //盲转
        async TransferPone() {
            this.dialogFormTransfer = false;
            console.log("CalledId.>" + zxCalledId + "..>CallId..>" + zxCallId + "..>CallDirection.>" + CallDirection)
            var TranAni = this.TransferAni;
            var TranCallId = zxCalledId;

            if (TranCallId == undefined) {
                this.$message.error('未在通话中,无法使用!');
                return;
            }
            await axios.get(HttpUrl + `/mod/Blind_turn?uuid=${TranCallId}&phone=${TranAni}`)
                .then(function (response) {
                    console.log("reponse." + response.data.code)
                    if (response.data.code == "0") {
                        ButtonVisForStatus("话后状态")
                        console.log("已经成功转给：" + TranAni)
                    } else {
                        this.$message.error('失败：' + response.data.val);
                    }
                }, function (err) {
                    console.log(err)
                    this.$message.error('失败：' + err);
                })
        },
        //咨询后操作
        async AfterOperation(command) {
            console.log("按键为："+command)
            var recv_dtmf=command;
            var TranCallId = zxCalledId;
            if (TranCallId == undefined) {
                this.$message.error('未在通话中,无法使用!');
                return;
            }
            await axios.get(HttpUrl + `/mod/recv_dtmf?uuid=${TranCallId}&key=${recv_dtmf}`)
                .then(function (response) {
                    console.log("reponse." + response.data.code)
                    if (response.data.code == "0") {
                        if(recv_dtmf=="1"){
                            ButtonVisForStatus("服务状态")
                        }else if(recv_dtmf=="2"){
                            ButtonVisForStatus("话后状态")
                        }else if(recv_dtmf=="3"){
                            ButtonVisForStatus("会议状态")
                        }
                    } else {
                        this.$message.error('失败：' + response.data.val);
                    }
                }, function (err) {
                    console.log(err)
                    this.$message.error('失败：' + err);
                })
        },
        //咨询外线
        async zxAniClick() {
            var TranAni = this.zxAniText;
            if (TranAni != '' && TranAni != undefined) {
                console.log("CalledId.>" + zxCalledId + "..>CallId..>" + zxCallId + "..>CallDirection.>" + CallDirection)
                var TranAni = this.zxAniText;
                var TranCallId = zxCalledId;
                if (CallDirection == "out") {
                    TranCallId = zxCalledId;
                }
                if (TranCallId == undefined) {
                    this.$message.error('未在通话中,无法使用!');
                    return;
                }
                //准备开始发起咨询的请求..
                console.log("Get Http ..>");
                await axios.get(HttpUrl + `/mod/att_xferByuuid?uuid=${TranCallId}&phone=${TranAni}`)
                    .then(function (response) {
                        console.log("转入xfer.>>");
                        console.log("reponse." + response.data.code)
                        if (response.data.code == "0") {
                            ButtonVisForStatus("咨询中状态")

                        } else {
                            this.$message.error('失败：' + response.data.val);
                        }
                    }, function (err) {
                        console.log(err)
                        this.$message.error('失败：' + err);
                })
                this.zxDialog = false
            }else{
                this.$message.error('号码无效！');
            }

        },
        //咨询坐席列表
        async zxAgentList() {
            this.zxAgentDialog = true
            this.gridClickOid = '';
            var AgentList;
            await axios.get(HttpUrl + "/callcenter/GetAvailableUser?AgentId=" + extension).then(function (response) {
                if (response.data != "无可用坐席") {
                    console.log('找到可用坐席.')
                    AgentList = response.data;
                } else {
                    console.log('无可用空闲坐席.')
                }
            }, function (err) {
                console.log(err)
                this.gridData = "未找到坐席";
            });
            this.gridData = AgentList
        },
        //页面关闭或者刷新
        beforeunloadFn(e) {

            var SendJson = {
                MSG: "AgentChange",
                Data: {
                    "AgentId": extension,
                    "State": "Logged Out"
                }
            }
            var str = JSON.stringify(SendJson);
            console.log("Request Message..>" + str)
            ws.send(str);
            // ...
        },
        //用来显示隐藏一个dialog的一个方法 外层js 调用vue的方法来给vue赋值
        DlVisible(e){
            if (e=="振铃状态"){
                this.CAni=CAni;
                this.dialogVisible = true; //弹屏框
            }else{
                this.dialogVisible = false;//弹屏框
            }

        }
    }
})

function ButtonInit() {
    console.log("Html Button Init..")
    $('#UnHoldVis').hide(); //隐藏取消保持按钮

    ButtonVisForStatus("小休状态");

}

//根据状态和操作来判断按钮的显示隐藏
function ButtonVisForStatus(Status) {
    switch (Status) {
        case "小休状态":
            $("#TaskVis").attr("class", "el-button el-button--default");
            $("#TaskVis").removeAttr("disabled");//启用外呼
            $("#TaskAgentVis").attr("class", "el-button el-button--default");
            $("#TaskAgentVis").removeAttr("disabled");//启用外呼坐席
            $("#HangupVis").attr("class", "el-button el-button--default is-disabled");//禁用el的按钮（挂机按钮）
            $("#HangupVis").attr("disabled", "disabled"); //禁用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");

            //启用的时候class 为：el-button el-button--default

            $("#StatusInput").val('小休状态');

            break;
        case "回铃状态":
            $("#StatusInput").val('回铃状态');
            //除了挂机按钮，其他都禁用
            $("#TaskVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskVis").attr("disabled", "disabled");
            $("#TaskAgentVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskAgentVis").attr("disabled", "disabled");
            $("#HangupVis").attr("class", "el-button el-button--default");//启用el的按钮（挂机按钮）
            $("#HangupVis").removeAttr("disabled");//启用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");

            break;
        case "外呼电话":
            $("#StatusInput").val('呼叫中状态');
            //除了挂机按钮，其他都禁用
            $("#TaskVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskVis").attr("disabled", "disabled");
            $("#TaskAgentVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskAgentVis").attr("disabled", "disabled");
            $("#HangupVis").attr("class", "el-button el-button--default");//启用el的按钮（挂机按钮）
            $("#HangupVis").removeAttr("disabled");//启用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");
            break;
        case "外呼坐席":
            $("#StatusInput").val('呼叫中状态');
            //除了挂机按钮，其他都禁用
            $("#TaskVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskVis").attr("disabled", "disabled");
            $("#TaskAgentVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskAgentVis").attr("disabled", "disabled");
            $("#HangupVis").attr("class", "el-button el-button--default");//启用el的按钮（挂机按钮）
            $("#HangupVis").removeAttr("disabled");//启用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");
            break;
        case "服务状态":
            $("#StatusInput").val('服务状态');
            vue.DlVisible('服务状态');
            //禁用后操作和外呼操作
            $("#TaskVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskVis").attr("disabled", "disabled");
            $("#TaskAgentVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskAgentVis").attr("disabled", "disabled");
            $("#HangupVis").attr("class", "el-button el-button--default");//启用el的按钮（挂机按钮）
            $("#HangupVis").removeAttr("disabled");//启用el的按钮 （挂机按钮）
            $("#HoldVis").removeAttr("disabled");//用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default");
            $("#TransferPoneVis").removeAttr("disabled");//用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default");
            $("#zxAniVis").removeAttr("disabled");//用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default");
            $("#zxAgentVis").removeAttr("disabled");//用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").removeAttr("disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default");

            break;
        case "话后状态":
            vue.DlVisible('话后状态');
            $("#StatusInput").val('话后状态');
            $("#TaskVis").attr("class", "el-button el-button--default");
            $("#TaskVis").removeAttr("disabled");//启用外呼
            $("#TaskAgentVis").attr("class", "el-button el-button--default");
            $("#TaskAgentVis").removeAttr("disabled");//启用外呼坐席
            $("#HangupVis").attr("class", "el-button el-button--default is-disabled");//禁用el的按钮（挂机按钮）
            $("#HangupVis").attr("disabled", "disabled"); //禁用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");
            break;
        case "振铃状态":
            $("#StatusInput").val('振铃状态');

            vue.DlVisible('振铃状态');

            //除了挂机按钮，其他都禁用
            $("#TaskVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskVis").attr("disabled", "disabled");
            $("#TaskAgentVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskAgentVis").attr("disabled", "disabled");
            $("#HangupVis").attr("class", "el-button el-button--default");//启用el的按钮（挂机按钮）
            $("#HangupVis").removeAttr("disabled");//启用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");
            break;
        case "回铃未摘":
            $("#StatusInput").val('小休状态');
            //this.StateValue = false;
            $("#TaskVis").attr("class", "el-button el-button--default");
            $("#TaskVis").removeAttr("disabled");//启用外呼
            $("#TaskAgentVis").attr("class", "el-button el-button--default");
            $("#TaskAgentVis").removeAttr("disabled");//启用外呼坐席
            $("#HangupVis").attr("class", "el-button el-button--default is-disabled");//禁用el的按钮（挂机按钮）
            $("#HangupVis").attr("disabled", "disabled"); //禁用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");
            break;
        case "空闲状态":
            $("#StatusInput").val('空闲状态');
            //this.StateValue = true;
            $("#TaskVis").attr("class", "el-button el-button--default");
            $("#TaskVis").removeAttr("disabled");//启用外呼
            $("#TaskAgentVis").attr("class", "el-button el-button--default");
            $("#TaskAgentVis").removeAttr("disabled");//启用外呼坐席
            $("#HangupVis").attr("class", "el-button el-button--default is-disabled");//禁用el的按钮（挂机按钮）
            $("#HangupVis").attr("disabled", "disabled"); //禁用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");
            break;
        case "登录失败":
            $("#StatusInput").val('未登录状态');

            $("#TaskVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskVis").attr("disabled", "disabled");
            $("#TaskAgentVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskAgentVis").attr("disabled", "disabled");
            $("#HangupVis").attr("class", "el-button el-button--default is-disabled");//禁用el的按钮（挂机按钮）
            $("#HangupVis").attr("disabled", "disabled"); //禁用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled", "disabled");//禁用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled", "disabled");//禁用盲转按钮
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled", "disabled");//禁用咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled", "disabled");//禁用咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").attr("disabled", "disabled");//禁用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default is-disabled");
            $("#SatisfactionVis").attr("disabled", "disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default is-disabled");
            break;
        case "咨询中状态":
            $("#StatusInput").val('咨询中状态');
            //禁用后操作和外呼操作
            $("#TaskVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskVis").attr("disabled", "disabled");
            $("#TaskAgentVis").attr("class", "el-button el-button--default  is-disabled");
            $("#TaskAgentVis").attr("disabled", "disabled");
            $("#HangupVis").attr("class", "el-button el-button--default");//启用el的按钮（挂机按钮）
            $("#HangupVis").removeAttr("disabled");//启用el的按钮 （挂机按钮）
            $("#HoldVis").attr("disabled","disabled");//用保持按钮
            $("#HoldVis").attr("class", "el-button el-button--default is-disabled");
            $("#TransferPoneVis").attr("disabled","disabled");
            $("#TransferPoneVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAniVis").attr("disabled","disabled");//咨询外线按钮
            $("#zxAniVis").attr("class", "el-button el-button--default is-disabled");
            $("#zxAgentVis").attr("disabled","disabled");//咨询坐席按钮
            $("#zxAgentVis").attr("class", "el-button el-button--default is-disabled");
            $("#AfterOperationVis").removeAttr("disabled");//启用后操作按钮
            $("#AfterOperationVis").attr("class", "el-button el-button--default");
            $("#SatisfactionVis").removeAttr("disabled");//禁用满意度按钮
            $("#SatisfactionVis").attr("class", "el-button el-button--default");
            break;
        case "会议状态":
            $("#StatusInput").val('会议状态');
            break;
        default:
            break;
    }
}
