package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/super-l/machine-code/machine"
	"go.bug.st/serial"
)

func memStats(ctx *gin.Context) {
	stat, err := mem.VirtualMemory()
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, stat)
}

func cpuInfo(ctx *gin.Context) {
	info, err := cpu.Info()
	if err != nil {
		replyError(ctx, err)
		return
	}
	if len(info) == 0 {
		replyFail(ctx, "查询失败")
		return
	}
	replyOk(ctx, info[0])
}

func cpuStats(ctx *gin.Context) {
	times, err := cpu.Times(false)
	if err != nil {
		replyError(ctx, err)
		return
	}
	if len(times) == 0 {
		replyFail(ctx, "查询失败")
		return
	}
	replyOk(ctx, times[0])
}

func diskStats(ctx *gin.Context) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		replyError(ctx, err)
		return
	}

	usages := make([]*disk.UsageStat, 0)
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			replyError(ctx, err)
			return
		}
		usages = append(usages, usage)
	}
	replyOk(ctx, usages)
}

func protocolList(ctx *gin.Context) {
	//ps := protocols.Protocols()
	//replyOk(ctx, ps)
}

func protocolDetail(ctx *gin.Context) {
	//name := ctx.Param("name")
	//for _, p := range protocols.Protocols() {
	//	if p.Name == name {
	//		replyOk(ctx, p)
	//		return
	//	}
	//}
	replyFail(ctx, "找不到协议")
}

func serialPortList(ctx *gin.Context) {
	ports, err := serial.GetPortsList()
	if err != nil {
		replyError(ctx, err)
		return
	}
	replyOk(ctx, ports)
}

func machineInfo(ctx *gin.Context) {
	info := machine.GetMachineData()
	replyOk(ctx, gin.H{
		"sn":   info.SerialNumber,
		"mac":  info.Mac,
		"uuid": info.PlatformUUID,
		"cpu":  info.CpuId,
	})
}
