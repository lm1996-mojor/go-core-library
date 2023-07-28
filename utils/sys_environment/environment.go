package sys_environment

import (
	"os"
	"runtime"

	clog "mojor/go-core-library/log"
)

//var Asset embed.FS

//func GetApplicationYaml() []byte {
//	con, err := Asset.ReadFile("application.yaml")
//	if err != nil {
//		panic(err)
//	}
//
//	return con
//}
//func main() {
//	key, err := GetPrivateKey()
//	if err != nil {
//		return
//	}
//	fmt.Println(key)
//}

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

//
//func GetAllLocaleFiles() []string {
//	var ret = make([]string, 0)
//
//	dirs, err := Asset.ReadDir("locales")
//	if err == nil {
//		for _, dir := range dirs {
//			if dir.IsDir() {
//				dirName := dir.Name()
//				files, err := Asset.ReadDir("locales/" + dirName)
//				if err == nil {
//					for _, file := range files {
//						if !file.IsDir() {
//							ret = append(ret, dirName+"/"+file.Name())
//						}
//					}
//				}
//			}
//		}
//	}
//
//	return ret
//}

//func GetLocaleFile(fileName string) ([]byte, error) {
//	return Asset.ReadFile("locales/" + fileName)
//}

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

//GetCurrentAbPath 获取当前文件绝对路径
//func GetCurrentAbPath() string {
//	dir := getCurrentAbPathByExecutable()
//	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
//	if strings.Contains(dir, tmpDir) {
//		return getCurrentAbPathByCaller()
//	}
//	return dir
//}
//
//// 获取当前执行文件绝对路径
//func getCurrentAbPathByExecutable() string {
//	exePath, err := os.Executable()
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
//	return res
//}
//
//// 获取当前执行文件绝对路径（go run）
//func getCurrentAbPathByCaller() string {
//	var abPath string
//	_, filename, _, ok := runtime.Caller(0)
//	if ok {
//		abPath = path.Dir(filename)
//	}
//	return abPath
//}
