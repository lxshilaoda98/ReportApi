function WorkTime()
digit_1612245502403 = session:playAndGetDigits(1,1,3,5000,"","./static/res/uploadFile/ivrWav/ring.wav","","^[0-4]|#$")
if(digit_1612245502403=="4") then 
digit_1612694094486 = session:playAndGetDigits(5,6,6,6000,"","./static/res/uploadFile/ivrWav/满意.wav","","^[0-3]$")
 
end 

end

function OffTime()
digit_1612365273926 = session:playAndGetDigits(1,1,1,5000,"","./static/res/uploadFile/ivrWav/ring.wav","","^[0-0]$")

end

--funcMusic
--funcAgent
--funcGroup
function Menu() 
	if not session:ready() then return end
		local yy = os.date("%Y%m%d");
		local getWeek = os.date("%w");
		local getHour = os.date("%H");
		local getMinute = os.date("%M");
		local getSecond = os.date("%S");
		local week= string.find("4,3,5,1,2", getWeek,1)
		local ifWorkTime=0;
		if week ~=nil
		then
		session:consoleLog("info","今天为工作日,继续查找验证.>\n");
		ifWorkTime =1
		-- body
		--上班--看下今天是否是特殊非工作日
		local xx= string.find("20210210,20210211", yy,1)
		if xx ~=nil then
		--找到休息日了
		ifWorkTime=0
		session:consoleLog("info","找到特殊非工作日，今天休息.>\n");
		else
		session:consoleLog("info","没有找到特殊非工作日！！！今天上班.>\n");
		end
		else
		--不上班，看下是否今天是上班日
		local td= string.find("20210201", yy,1)
		if td ~=nil then
		--今天上班
		ifWorkTime=1
		session:consoleLog("info","找到特殊工作日，今天上班.>\n");
		else
		session:consoleLog("info","没有找到特殊工作日！！！今天休息.>\n");
		end
		end
		if ifWorkTime ==1 then
		local ttTime = getHour..getMinute..getSecond
		session:consoleLog("info","今天上班.>查找是否在上班时间内"..ttTime.."\n");
		local intTTime = tonumber(ttTime)
		if intTTime >=090000 and intTTime < 120000 then
		-- body
		session:consoleLog("info","早上上班中\n");
		else if intTTime >=120000 and intTTime <203000 then
		session:consoleLog("info","下午上班中\n");
		else
		session:consoleLog("info","下班了\n");
		ifWorkTime = 0
		end
		end
		end
		if ifWorkTime ==1
		then
		WorkTime()
		else
		OffTime()
		end
		session:consoleLog("info",ifWorkTime.."\n");
		session:consoleLog("info","week.>"..getWeek.."..>hour..>"..getHour.."..>getMinute..>"..getMinute.."..getSecond..>"..getSecond);
		end

session:answer()
session:setAutoHangup(false)
Menu()
