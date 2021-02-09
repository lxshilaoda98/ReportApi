
function errKeys(errKey)
    session:consoleLog('info',"go..errkey.>"..errKey)
    session:streamFile("/usr/local/freeswitch/sounds/ivrwav/errKey.wav");
    if errKey =="menu2" then
        menu2()
    elseif errKey =="key1" then
        key1()
    elseif errKey=="key2" then
        key2()
    elseif errKey=="key3" then
        key3()
    elseif errKey=="key4" then
        key4()
    end
end


function menu()
    if not session:ready() then return end
     local getWeek = os.date("%w");
     local getHour = os.date("%H");
     local getMinute = os.date("%M");--分钟
     session:consoleLog("info","week.>"..getWeek.."..>hour..>"..getHour.."..>getMinute..>"..getMinute);
    if tonumber(getHour) >= 9 and tonumber(getHour) <=17 then
        session:consoleLog("info","workTime.>")
        session:playAndGetDigits(1,1,3,3000,"#","/usr/local/freeswitch/sounds/ivrwav/workTime.wav","","^[0-9]$")
        menu2()
    else
        session:consoleLog("info","NoworkTime.>")
        noWorkTime();
    end
end


--播放二级菜单
function menu2()
    if not session:ready() then return end

    digit = session:playAndGetDigits(1,1,1,5000,"","/usr/local/freeswitch/sounds/ivrwav/menu2.wav","","^[0-4]|\\*$")

    if(digit=="1") then
        key1()
    elseif (digit=="2") then
        key2()
    elseif (digit=="3") then
        key3()
    elseif (digit =="4") then
        key4()
    elseif (digit=="0") then
        key0()
    elseif (digit =="*") then
        menu2()
    else
        errKeys("menu2")
    end
end


function key1()
    if not session:ready() then return end
    digit = session:playAndGetDigits(1,1,1,5000,"","/usr/local/freeswitch/sounds/ivrwav/key1.wav","","^0|\\*|#$")
    if (digit=="*") then
        key1()
    elseif(digit=="#") then
        menu2()
    elseif (digit=="0") then
        skill = "support@default";
        toAgent(skill)
    else
        session:consoleLog('info',"errorkey"..digit)
        errKeys('key1')--应该进入按键错误节点
    end
end

function key2()
    if not session:ready() then return end
    digit = session:playAndGetDigits(1,1,1,5000,"","/usr/local/freeswitch/sounds/ivrwav/key2.wav","","^0|\\*|#$")
    if (digit=="*") then
        key2()
    elseif(digit=="#") then
        menu2()
    elseif (digit=="0") then
        skill = "support@default";
        toAgent(skill)
    else
        errKeys("key2")--应该进入按键错误节点
    end
end

function key3()
    if not session:ready() then return end
    digit = session:playAndGetDigits(1,1,1,5000,"","/usr/local/freeswitch/sounds/ivrwav/key3.wav","","^0|\\*|#$")
    if (digit=="*") then
        key3()
    elseif(digit=="#") then
        menu2()
    elseif (digit=="0") then
        skill = "support@default";
        toAgent(skill)
    else
        errKeys("key3")--应该进入按键错误节点
    end
end

function key4()
    if not session:ready() then return end
    digit = session:playAndGetDigits(1,1,1,5000,"","/usr/local/freeswitch/sounds/ivrwav/key4.wav","","^0|\\*|#$")
    if (digit=="*") then
        key4()
    elseif(digit=="#") then
        menu2()
    elseif (digit=="0") then
        skill = "support@default";
        toAgent(skill)
    else
        errKeys("key4")--应该进入按键错误节点
    end
end

function key0()
    if not session:ready() then return end

    skill = "support@default";
    toAgent(skill)

end

function toAgent(skill)
    if not session:ready() then return end
    session:streamFile("/usr/local/freeswitch/sounds/ivrwav/key0.wav");

    session:execute("callcenter",skill)
    return "break";

end

function goodbye()
    if not session:ready() then return end
    session:hangup()
end

function noWorkTime()
    if not session:ready() then return end
    --session:speak("say","see you.")
    session:streamFile("/usr/local/freeswitch/sounds/ivrwav/noWorkTime.wav");
    session:hangup()
end

 --应答本次呼叫
 session:answer()
 session:setAutoHangup(false)

 menu()