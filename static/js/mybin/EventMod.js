var Ring={
    /**
     * 振铃
     */
    startRingTone: function () {
        console.log('开启振铃声音..>')
        let play = document.getElementById("ringtone").play();
        play.then(()=>{

        }).catch((err)=>{

        });

    },
    /**
     * 停止振铃
     */
    stopRingTone: function () {
        console.log('暂停振铃声音..>')
        document.getElementById("ringtone").pause();
    },

    startRingbackTone: function () {
        let play = document.getElementById("ringbacktone").play();
        play.then(()=>{

        }).catch((err)=>{


        });
    },
    stopRingbackTone: function () {
        document.getElementById("ringbacktone").pause();
    }


};