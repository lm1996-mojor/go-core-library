package sys_environment

import (
	"fmt"
	"io"
	sdkNet "net"
	"net/http"
	"os"
	"runtime"
	"strings"

	clog "github.com/lm1996-mojor/go-core-library/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
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
	interfaces, err := sdkNet.Interfaces()
	if err != nil {
		fmt.Println("获取网卡信息出错：", err)
		return nil
	}
	ipHeader := ""
	switch runtime.GOOS {
	case "windows":
		ipHeader = "192"
	case "linux":
		ipHeader = "172"
	}
	for _, iface := range interfaces {
		addrs, err1 := iface.Addrs()
		if err1 != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*sdkNet.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					if strings.Split(ipNet.IP.String(), ".")[0] == ipHeader {
						localIp = append(localIp, ipNet.IP.String())
					}
				}
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

func GetIp() (ip string) {
	interfaces, err := sdkNet.Interfaces()
	if err != nil {
		fmt.Println("获取网卡信息出错：", err)
		return ""
	}
	ipHeader := ""
	switch runtime.GOOS {
	case "windows":
		ipHeader = "192"
	case "linux":
		ipHeader = "172"
	}
	endFlag := false
	for _, iface := range interfaces {
		addrs, err1 := iface.Addrs()
		if err1 != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*sdkNet.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					if strings.Split(ipNet.IP.String(), ".")[0] == ipHeader {
						ip = ipNet.IP.String()
						endFlag = true
						break
					}
				}
			}
		}
		if endFlag {
			break
		}
	}
	return ip
}
