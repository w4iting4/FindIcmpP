package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	// "strings"
)

func execCommand(str string) (string, error) {
	cmd := exec.Command("CMD", "/C", str)
	//stdout, err := cmd.StdoutPipe() //通道通过wait关闭 获取输出对象，可以从该对象中读取输出结果
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := cmd.Start(); err != nil { // 运行命令
	//	log.Fatal(err)
	//}
	//if opBytes, err := ioutil.ReadAll(stdout); err != nil { // 读取输出结果
	//	log.Fatal(err)
	//} else {
	//	if err := cmd.Wait(); err != nil {
	//		return "", err
	//	}
	//	file := string(opBytes)
	//	return file, nil
	//}
	//return "", err
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println("[Error]=> " + fmt.Sprint(err))
		return "", err
	}
	return out.String(), nil
}

func Chcp65001() {
	log.Println("[Info]-> 初始化当前终端编码为UTF-8")
	str, err := execCommand("chcp 65001")
	if err != nil {
		log.Fatal("[Error]=> 初始化终端失败")
	}
	fmt.Println(str)
}

func StartGetFile() string {
	log.Println("[Info]-> 开始抓包")
	file, err := execCommand("netsh trace start capture=yes Protocol=1")
	if err != nil {
		log.Fatal("[Error]=> 请使用管理员权限")
	}
	log.Printf("[Info]: %s", file)
	var fileSilce = strings.Split(file, "\r\n")
	for i, v := range fileSilce { //遍历切片确认文件路径
		if strings.Contains(fileSilce[i], "Trace File") {
			var traceFilePath = strings.Split(v, " ")
			for index, value := range traceFilePath {
				if strings.Contains(traceFilePath[index], ".etl") {
					log.Printf("[Info]-> ETL文件路径:%v\n", value)
					return value
				}
			}
		}
	}
	return ""
}

func StopGetFile() {
	log.Println("[Info]-> 停止抓包,这一步等待系统打包时间较长，请耐心等候")
	str, err := execCommand("netsh trace stop")
	if err != nil {
		log.Fatal("[Error]=> 停止抓包失败")
	}
	log.Printf("[Info]: %s", str)
}

func ETL2CSV(Etl string) string {
	log.Println("[Info]-> 文件转储")
	Csv := GetCsvFilePath(Etl)
	commandStr := "tracerpt " + Etl + " -o " + Csv + " -of CSV"
	str, err := execCommand(commandStr)
	fmt.Println(err)
	if err != nil {
		log.Fatal("[Info]->] 转储文件失败")
		return ""
	}
	log.Printf("[Info]-> %s", str)
	return Csv

}

func GetFileByPid(ptr *map[int]map[string]string) {
	for index, value := range *ptr {
		for i, v := range value {
			if i == "PID" {
				if (*ptr)[index][i] != "4" && (*ptr)[index][i] != "0" {
					commandStr := "wmic process get name,executablepath,processid|findstr " + v
					str, err := execCommand(commandStr)
					if err != nil {
						(*ptr)[index]["FilePath"] = "该进程已经结束，无法获取文件路径"
					} else {
						(*ptr)[index]["FilePath"] = str
					}
				} else {
					(*ptr)[index]["FilePath"] = "系统进程，无法获取到具体文件"
				}
			}

		}
	}
}
