function CallingAni(ani) {
    let option = {
        title: '来电信息',
        popupTime: false,
        hook: {
            initStart: function () {
                // 检查之前老旧实例如果存在则销毁
                if (document.querySelector('#modal-layer-container'))
                    ModalLayer.removeAll();
            }
        },
        mask:false,
        contentFullContainer: false,
        displayProgressBar: true,
        displayProgressBarPos: 'left',
        content: '收到'+ani+'来电,是否接听?',
        okText: '接听',
        noText: '拒绝',
        callback: {
            ok: function () {
                this.hide();
                sipOnAnswer();
            },
            no: function () {
                this.hide();
                sipHangUp();
            }
        }
    }

    ModalLayer.confirm(option);
}