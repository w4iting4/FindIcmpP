package flag

import (
	"FindIcmpP/utils"
	"flag"
	"log"
	"time"
)

func Getflag() (uint, bool, string, bool) {
	var (
		Times        uint
		ParseEtlOnly bool
		EtlFilePath  string
		Catch        bool
	)
	flag.UintVar(&Times, "t", 10, "在主机抓包时长，默认10s,建议不超过30s")
	flag.BoolVar(&ParseEtlOnly, "po", false, "如果选择该参数，则不会进行抓包，只会解析本地etl文件")
	flag.StringVar(&EtlFilePath, "f", "", "选择-p的情况下，需要通过该参数来指定ETL文件路径")
	//是否追踪进程文件 -catch
	flag.BoolVar(&Catch, "c", false, "默认模式下不会追踪启动进程的文件，如果选择该参数，则会追踪含有ICMP协议通信的文件")
	flag.Parse() // 必须得parse 不然无法取值
	return Times, ParseEtlOnly, EtlFilePath, Catch
}
func ParseFlag() {
	times, ParseEtlOnly, etlFilePath, catch := Getflag()
	//init console
	utils.Chcp65001()
	//fmt.Println(times, ParseEtlOnly, etlFilePath, catch)
	if ParseEtlOnly {
		log.Println("[Info]-> 本地文件解析模式")
		//判断文件是否为ETL格式且是否存在
		jug := utils.CheckFile(etlFilePath)
		if jug {
			//1. 获取csv文件路径 + 文件转储
			csvFilePath := utils.ETL2CSV(etlFilePath)
			//3. 转储文件解析 得到存储map地址
			mapPtr := utils.ParseFile(csvFilePath)
			if len(*mapPtr) > 0 {
				log.Println("[Info] 存在ICMP通信的进程")
				if catch {
					//判断是否存在ICMP进程数据
					//log.Println("[Info] 存在ICMP通信的进程，正在解析数据")
					utils.GetFileByPid(mapPtr)
					resFilePath := utils.GetResultFilePath(csvFilePath)
					utils.OutPutFile(resFilePath, mapPtr)
				}
			} else {
				log.Fatal("[Info]-> 没有ICMP通信的进程，程序即将退出，请稍后再尝试~")
				log.Fatal("[Info]-> bye~")
			}
		} else {
			log.Fatal("[Error]=> ETL文件不存在，请检查文件路径！")
		}
	} else {
		log.Println("[Info]-> 抓包解析模式")
		etlFilePath := utils.StartGetFile()
		if etlFilePath != "" {
			time.Sleep(time.Second * time.Duration(times))
			utils.StopGetFile()
			csvFilePath := utils.ETL2CSV(etlFilePath)
			//3. 转储文件解析 得到存储map地址
			mapPtr := utils.ParseFile(csvFilePath)
			if len(*mapPtr) > 0 {
				log.Println("[Info]-> 存在ICMP通信的进程")
				if catch {
					//判断是否存在ICMP进程数据
					//log.Println("[Info]-> 存在ICMP通信的进程，正在解析数据")
					utils.GetFileByPid(mapPtr)
					resFilePath := utils.GetResultFilePath(csvFilePath)
					utils.OutPutFile(resFilePath, mapPtr)
				}
			} else {
				log.Fatal("[Info]-> 没有ICMP通信的进程，请再次尝试~")
			}
		} else {
			log.Fatal("[Erorr]=> 路径获取失败")
		}

	}
}
