package sys_environment

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func GetHostInfo() *host.InfoStat {
	info, _ := host.Info()
	return info
}

func GetPrivateKey() ([]byte, error) {
	//获取项目根目录
	dir, err2 := os.Getwd()
	if err2 != nil {
		return nil, err2
	}
	clog.Infof("项目路径: " + dir)
	con, err := os.ReadFile(dir + GetOperatingSystem() + "key" + GetOperatingSystem() + "private.key")
	return con, err
}

func GetPublicKey() ([]byte, error) {
	//获取项目根目录
	dir, err2 := os.Getwd()
	if err2 != nil {
		return nil, err2
	}
	con, err := os.ReadFile(dir + GetOperatingSystem() + "key" + GetOperatingSystem() + "public.txt")

	return con, err
}

// GetOperatingSystem 获取操作系统并返回对应文件夹路径符号
func GetOperatingSystem() string {
	switch runtime.GOOS {
	case "windows":
		return "\\"
	case "linux":
		return "/"
	}
	return ""
}

func GetCpuPercent() []cpu.InfoStat {
	info, _ := cpu.Info()
	return info
}

func GetMemPercent() *mem.VirtualMemoryStat {
	memInfo, _ := mem.VirtualMemory()
	return memInfo
}

func GetSwapMemory() *mem.SwapMemoryStat {
	memory, _ := mem.SwapMemory()
	return memory
}

func GetSwapDevice() []*mem.SwapDevice {
	devices, _ := mem.SwapDevices()
	return devices
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

func GetAllDiskPartition() []disk.PartitionStat {
	info, _ := disk.Partitions(true)
	return info
}

func GetSpecifyDiskPartition(driveLetter string) *disk.UsageStat {
	info, _ := disk.Usage(driveLetter)
	return info
}

func GetIOCounters() map[string]disk.IOCountersStat {
	info, _ := disk.IOCounters()
	return info
}

// GetInternalIP
/** 获取内网ip
 * @param netMode 网络模式(tcp、udp、tcp4、udp4等等)
 * @return []net.ConnectionStat 网络信息数组
 */
func GetInternalIP() (localIp []string) {
	info, _ := net.Interfaces()
	for _, stat := range info {
		for _, addr := range stat.Addrs {
			if strings.Contains(addr.Addr, "24") {
				localIp = append(localIp, addr.Addr)
			}
		}
	}
	return localIp
}

// GetExternal 获取本机公网ip
func GetExternal() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//s := buf.String()
	return string(content)
}
