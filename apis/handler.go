package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/n1n1n1_owner/ReportApi/models"
	Crmhelper "github.com/n1n1n1_owner/ReportApi/models/Helper"
	"github.com/n1n1n1_owner/ReportApi/models/ivrConter/jys"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"net/http"
	"strconv"
	"time"
	//jys "github.com/n1n1n1_owner/ReportApi/models/ivrConter/jys"
)

func GetDialplan(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")

	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	p := models.Dialplan{}
	dialplans, err := p.GetDialplan(start, end)
	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": dialplans,
		"count":  len(dialplans),
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}

//IVR呼叫量统计呼入参数
func GetIvrStatis(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	types := c.DefaultQuery("type", "日")
	sTime_epoch := c.DefaultQuery("sTime_epoch", string(time.Now().Unix()))
	eTime_epoch := c.DefaultQuery("eTime_epoch", string(time.Now().Unix()))

	timeLayout := "2006-01-02 15:04:05"
	//sTime_epoch string 转换成int64
	sTime, err := strconv.ParseInt(sTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. sTime_epoch .>Err.>", err)
	}
	var data_sTime = time.Unix(sTime, 0).Format(timeLayout)

	//eTime_epoch string 转换成int64
	eTime, err := strconv.ParseInt(eTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. eTime_epoch .>Err.>", err)
	}
	var data_eTime = time.Unix(eTime, 0).Format(timeLayout)

	//分页相关 转换
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	p := models.IvrStatis{}
	ivrStatis, err := p.GetIvrStatis(types, data_sTime, data_eTime, start, end)
	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": ivrStatis,
		"count":  len(ivrStatis),
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}

//坐席呼叫量统计
func GetAgent_CallStatis(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	sTime_epoch := c.DefaultQuery("sTime_epoch", string(time.Now().Unix()))
	eTime_epoch := c.DefaultQuery("eTime_epoch", string(time.Now().Unix()))

	timeLayout := "2006-01-02 15:04:05"
	//sTime_epoch string 转换成int64
	sTime, err := strconv.ParseInt(sTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. sTime_epoch .>Err.>", err)
	}
	var data_sTime = time.Unix(sTime, 0).Format(timeLayout)

	//eTime_epoch string 转换成int64
	eTime, err := strconv.ParseInt(eTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. eTime_epoch .>Err.>", err)
	}
	var data_eTime = time.Unix(eTime, 0).Format(timeLayout)

	//分页相关 转换
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	p := models.Agent_CallStatis{}
	ivrStatis, err := p.GetAgent_CallStatis(data_sTime, data_eTime, start, end)
	if err != nil {
		fmt.Printf("[GetAgent_CallStatis] sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": ivrStatis,
		"count":  len(ivrStatis),
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}

//综合呼叫统计
func GetCallCountStatis(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	sTime_epoch := c.DefaultQuery("sTime_epoch", string(time.Now().Unix()))
	eTime_epoch := c.DefaultQuery("eTime_epoch", string(time.Now().Unix()))

	selectType := c.DefaultQuery("selectType", "日")

	timeLayout := "2006-01-02 15:04:05"
	//sTime_epoch string 转换成int64
	sTime, err := strconv.ParseInt(sTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. sTime_epoch .>Err.>", err)
	}
	var data_sTime = time.Unix(sTime, 0).Format(timeLayout)

	//eTime_epoch string 转换成int64
	eTime, err := strconv.ParseInt(eTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. eTime_epoch .>Err.>", err)
	}
	var data_eTime = time.Unix(eTime, 0).Format(timeLayout)

	//分页相关 转换
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	p := models.CallCountStatis{}
	ivrStatis, err := p.GetCallCountStatis(data_sTime, data_eTime, start, end, selectType)
	if err != nil {
		fmt.Printf("[GetCallCountStatis] sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": ivrStatis,
		"count":  len(ivrStatis),
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)

}

func GetJysAgent(c *gin.Context) {
	Ani := c.DefaultQuery("Ani", "17600082595")
	a := jys.Agent{}
	agent, err := a.GetAgentIDForJys(Ani)
	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": agent,
	})
}

//#####系统监控方法

//系统信息 返回json
func GetComputeInfoJson(c *gin.Context) {

	info, _ := host.Info() //获取计算机信息
	fmt.Println(info)

	c.JSON(http.StatusOK, gin.H{
		"result": info,
	})
}

//系统信息，映射到页面
func GetComputeInfoMonitor(c *gin.Context) {
	info, _ := host.Info() //获取计算机信息
	day, hour, minute := Crmhelper.ResolveTime(int(info.Uptime))
	upTimeDay := strconv.Itoa(day) + "天"
	bootTime := time.Unix(int64(info.BootTime), 0).Format("2006-01-02 15:04:05")
	fmt.Println(day, hour, minute, upTimeDay, bootTime)

	//stats, e := host.Users()
	//if e!=nil{
	//	fmt.Println("get host users failed .>",e)
	//}else{
	//	fmt.Println(stats)
	//}

	c.HTML(200, "index.html", gin.H{
		"hostname":             info.Hostname,
		"uptime":               upTimeDay,
		"bootTime":             bootTime,
		"procs":                info.Procs,
		"os":                   info.OS,
		"platform":             info.Platform,
		"platformFamily":       info.PlatformFamily,
		"platformVersion":      info.PlatformVersion,
		"kernelVersion":        info.KernelVersion,
		"kernelArch":           info.KernelArch,
		"virtualizationSystem": info.VirtualizationSystem,
		"virtualizationRole":   info.VirtualizationRole,
		"hostID":               info.HostID,
	})
}

//获取cpu信息，映射到页面

func GetCpuInfoMonitor(c *gin.Context) {
	InfoStat, err := cpu.Info()
	if err != nil {
		fmt.Println("Get cpu info Err.>", err)
		return
	}
	iCpu, _ := cpu.Counts(true) //cpu逻辑数量
	fmt.Println("cpu逻辑核心数量", iCpu)
	wCpu, _ := cpu.Counts(false)
	fmt.Println("cpu物理核心数量", wCpu)
	fmt.Println(InfoStat)
	fmt.Println(len(InfoStat))

	//执行所消耗的时间
	//TimesStat包含CPU执行各种工作所花费的时间。 时间单位以秒为单位。 它基于linux / proc / stat文件。
	//用户时间（User  Time）即us所对应的列，表示CPU执行用户进程所占用的时间，通常情况下希望us的占比越高越好
	//系统时间（System Time）即sy所对应该的列，表示CPU自内核态所花费的时间，sy占比比较高通常意味着系统在某些方面设计得不合理，比如频繁的系统调用导致的用户态和内核态的频繁切换
	//Nice时间（Nice Time）即ni所对应的列，表示系统在调整进程优先级的时候所花费的时间
	//空闲时间（Idle Time）即id所对应的列，表示系统处于空闲期，等待进程运行，这个过程所占用的时间。当然，我们希望id的占比越低越好
	//等待时间（Waiting Time）即wa所对应的列，表示CPU在等待I/O操作所花费的时间，系统不应该花费大量的时间来进行等待，否则便表示可能有某个地方设计不合理
	//硬件中断处理时间（Hard Irq Time）即hi对应的列，表示系统处理硬件中断所占用的时间
	//软件中断处理时间（Soft Irq Time）即si对应的列，表示系统处理软件中断所占用的时间
	//丢失时间（Steal Time）即st对应的列，实在硬件虚拟化开始流行后操作系统新增的一列，表示被强制等待虚拟CPU的时间
	tsInfo, _ := cpu.Times(false)
	fmt.Println(tsInfo)
	//user User表示：CPU一共花了多少比例的时间运行在用户态空间或者说是用户进程(running user space processes)。
	// 典型的用户态空间程序有：Shells、数据库、web服务器……

	//Nice表示：可理解为，用户空间进程的CPU的调度优先级，范围为[-20,19]

	//System System的含义与User相似。System表示：CPU花了多少比例的时间在内核空间运行。分配内存、IO操作、创建子进程……都是内核操作。这也表明，当IO操作频繁时，System参数会很高。

	//Wait  在计算机中，读写磁盘的操作远比CPU运行的速度要慢，CPU负载处理数据，而数据一般在磁盘上需要读到内存中才能处理。
	//当CPU发起读写操作后，需要等着磁盘驱动器将数据读入内存(可参考：JAVA IO 以及 NIO 理解)，
	// 从而导致CPU 在等待的这一段时间内无事可做。CPU处于这种等待状态的时间由Wait参数来衡量。

	//Idle Idel表示：CPU处于空闲状态时间比例。一般而言，idel + user + nice 约等于100%

	//iowait iowait其实是一种特殊形式的CPU空闲

	//irq 中断请求？

	//softirq 软中断？

	//steal 你的虚拟机（VM）会与虚拟环境的宿主机上的多个虚拟机实例共享物理资源。
	// 其中之一共享的就是CPU时间切片。如果你的VM的物理机虚拟比是1/4，
	// 那么它的CPU使用率不会限制于25%的CPU时间切片－它能够使用超过它设置的虚拟比。（有别于内存的使用，内存大小是严格控制的）。

	//guest 当客户CPU执行hlt指令，它自愿产生CPU时间给其他客户。


	//百分比计算每个CPU使用或组合使用的cpu的百分比。 如果给定的间隔为0，它将把当前的CPU时间与上次调用进行比较。
	// 每个cpu返回一个值，如果percpu设置为false，则返回一个值。
	syinfo, _ := cpu.Percent(time.Duration(time.Second), false)
	fmt.Println(syinfo)



}



//获取物理内存和交换区内存信息
func GetMemInfoMonitor(c *gin.Context){
	memInfo, _ := mem.SwapMemory()
	fmt.Println(memInfo)
	fmt.Println(Crmhelper.FormatFileSize(int64(memInfo.Total)))
	fmt.Println(Crmhelper.FormatFileSize(int64(memInfo.Used)))
	fmt.Println(Crmhelper.FormatFileSize(int64(memInfo.Free)))
	info2, _ := mem.SwapMemory()
	fmt.Println(info2)
}

//可以通过psutil获取磁盘分区、磁盘使用率和磁盘IO信息
func GetDiskInfoMonitor(c *gin.Context){
	query := c.DefaultQuery("id", "0")
	info, _ := disk.Partitions(true) //所有分区
	//fmt.Println(info)
	numbers:=[]string{} //磁盘信息
	SumNumber:=[]string{}//磁盘总容量
	UseNumber:=[]string{}//磁盘使用量
	DisBL :=[]string{} //磁盘使用情况，用了多少+剩余多少
	infoStr:=[]*disk.UsageStat{}//磁盘列表详细信息
	infoStrForDiskName :=[] disk.IOCountersStat{}

	RWIoNumber := make(map[string]uint64)
	for k,_ :=range info {
		numbers = append(numbers, info[k].Device)
		info2, _ := disk.Usage(info[k].Device) //指定某路径的硬盘使用情况

		i, e := strconv.Atoi(Crmhelper.GetGbFileSize(info2.Total))
		if e != nil {
			fmt.Println("convert String to Int Err..>",e)
		}
		fi, e := strconv.Atoi(Crmhelper.GetGbFileSize(info2.Free))
		if e != nil {
			fmt.Println("convert String to Int Err..>",e)
		}
		ui, _ := strconv.Atoi(Crmhelper.GetGbFileSize(info2.Used))

		infoStr = append(infoStr,info2)

		infoStr[k].Total = uint64(i)
		infoStr[k].Free = uint64(fi)
		infoStr[k].Used = uint64(ui)
		upc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", info2.UsedPercent), 64)
		infoStr[k].UsedPercent = upc
		ipc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", info2.InodesUsedPercent), 64)
		infoStr[k].InodesUsedPercent = ipc
		SumStr :=Crmhelper.GetGbFileSize(info2.Total)
		SumNumber=append(SumNumber,SumStr)

		UseStr :=Crmhelper.GetGbFileSize(info2.Used)
		UseNumber=append(UseNumber,UseStr)

		FeeStr :=Crmhelper.GetGbFileSize(info2.Free)
		//添加指定的磁盘使用情况
		if query ==info[k].Device {
			DisBL = append(DisBL, UseStr)
			DisBL = append(DisBL, FeeStr)
		}
	}
	info3, _ := disk.IOCounters("D:") //所有硬盘的io信息
	for m,v := range info3{
		if (query != "0"){
			if query == m{
				infoStrForDiskName = append(infoStrForDiskName, v)

				RWIoNumber["ReadCount"]=v.ReadCount
				RWIoNumber["WriteCount"]=v.WriteCount
				RWIoNumber["IopsInProgress"]=v.IopsInProgress
				//RWIoNumber=append(RWIoNumber,v.ReadCount,v.WriteCount,v.IopsInProgress )
			}
			//fmt.Println("找到值..>",query)
		}

		//fmt.Println("v..>",v,"m..>",m)
	}
	c.HTML(200,"DiskInfo.html",gin.H{
			"RWIoNumber":RWIoNumber,
			"DisBL":DisBL,
			"IoDiskInfo":infoStrForDiskName, //渲染磁盘io
			"cMapDiskForIo":info3,
			"DiskInfo":infoStr,
			"DiskLen":numbers,
			"SumNumber":SumNumber,
			"UseNumber":UseNumber,
	})
}

//获取当前网络连接信息
func GetNetInfoMonitor(c *gin.Context){
	info, _ := net.Connections("all") //可填入tcp、udp、tcp4、udp4等等 all 查询全部的连接
	//fmt.Println(info)

	//获取网络读写字节／包的个数
	info2, _ := net.IOCounters(false)
	//fmt.Println(info2)
	//c.HTML(200,"index.html",gin.H{
	//	"data":info2,
	//})
	c.HTML(http.StatusOK,"index.html",gin.H{
		"code":0,
		"res":info2,
		"nres":info,
	})
}

func GetProcessByIdInfoMonitor(c *gin.Context){
	pid:=c.DefaultQuery("id","")
	pidInt, _ := strconv.Atoi(pid)

	fmt.Println("进入byID..>",c.DefaultQuery("id",""))
	newProcess, _ := process.NewProcess(int32(pidInt))

	fmt.Println(newProcess)

	//info2,_ := process.GetWin32Proc(1120) //对应pid的进程信息
	//fmt.Println(info2)
	//
	//fmt.Println(info2[0].ParentProcessID) //获取父进程的pid

}



func GetFreeSwitchInfoMonitor(c *gin.Context){
	c.HTML(http.StatusOK,"FreeswitchMonitor.html",gin.H{

	})
}