# FindIcmpP
<a name="UjTPi"></a>
### 为什么有FindIcmpP
> 因为之前看到有人询问产品发现了ICMP隧道的告警如何定位到进程与文件,后来我花时间查了一下资料,这一方面的资料确实不多,所以尝试了一下.

<a name="im8so"></a>
### 工具的使用场景
当IDS设备出现ICMP隧道告警的时候,我们可以使用该工具在WINDOWS操作系统上进行进程与文件的定位,由于借助[netsh](https://docs.microsoft.com/zh-cn/windows-server/networking/technologies/netsh/netsh)工具,所以对WINDOWS操作系统有一定的要求,必须在win7以上,支持`netsh trace`功能.**该工具的本质是抓包,随后解析数据.所以可能会出现抓包的时候进程存活,但是解析完数据之后进程结束了从而无法定位到进程的现象.这一个问题我暂时么的解决方案.**
<a name="C6Ajl"></a>
### 为什么不支持Linux
因为ICMP本身是低层协议，在linux上的实现是使用的 `SOCK_RAW`<br />不论是`netstat` 还是高发行版本替代`netstat` 的 `ss` 都拥有查询原始套接字网络连接的功能，所以我们在Linux的主机上，可以通过特定的命令获取到Linux上发送`ICMP`数据包的进程.
<a name="sRl3O"></a>
#### 具体操作
可以根据对应的流量告警详情选择 `IP` 协议的版本,一下演示使用 `IPV4` 环境<br />`netstat -alpw4` 或`ss -alpw4`<br />![image.png](https://cdn.nlark.com/yuque/0/2022/png/2078172/1645181091633-64ebdf57-ea27-4031-8fe8-20d39d678be2.png#clientId=u3a505023-864d-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=70&id=u0d1f3465&margin=%5Bobject%20Object%5D&name=image.png&originHeight=94&originWidth=1030&originalType=binary&ratio=1&rotation=0&showTitle=false&size=13456&status=done&style=none&taskId=ud74da398-9626-441b-9b13-b50e33c0167&title=&width=761.983154296875)<br />将进程PID提取，并追踪查找父进程即可<br />![image.png](https://cdn.nlark.com/yuque/0/2022/png/2078172/1645181447166-53ab5b74-f719-4403-b373-97d77d488e7e.png#clientId=u3a505023-864d-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=169&id=u9abd9c0a&margin=%5Bobject%20Object%5D&name=image.png&originHeight=295&originWidth=1337&originalType=binary&ratio=1&rotation=0&showTitle=false&size=31984&status=done&style=none&taskId=u44fb1115-89ca-4491-92ab-3ed35e8bae4&title=&width=766.4957885742188)
<a name="TTYVj"></a>
#### 演示案例
```bash
客户端
pingtunnel.exe -type client -l 172.16.xx.xx:4455 -s 82.xx.xxx.xxx -t 82.xx.xxx.xxx:4455 -tcp 1

服务端
sudo ./pingtunnel -type server

服务端定位pingtunnel   -- 如果出现进程迁移，或者注入内存，可以通过pstree 跟踪父进程
root@VM-0-5-ubuntu:~# netstat -alpw4  //
Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name    
raw        0      0 0.0.0.0:icmp            0.0.0.0:*               7           4053491/./pingtunne 

root@VM-0-5-ubuntu:~# pstree -spna 4053491
systemd,1
  └─sshd,696
      └─sshd,4045273 
          └─sshd,4045382  
              └─bash,4045383
                  └─sudo,4045478 -i
                      └─bash,4045479
                          └─sudo,4053490 ./pingtunnel -type server
                              └─pingtunnel,4053491 -type server
                                  ├─{pingtunnel},4053492
                                  ├─{pingtunnel},4053493
                                  └─{pingtunnel},4053494

客户端定位与服务端一致
root@VM-24-8-ubuntu:~# clear
root@VM-24-8-ubuntu:~# netstat -alpw4
Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name    
raw        0      0 0.0.0.0:icmp            0.0.0.0:*               7           3104052/./pingtunne 
root@VM-24-8-ubuntu:~# pstree -snpa 3104052
systemd,1
  └─sshd,1106
      └─sshd,3102054 
          └─sshd,3102167  
              └─bash,3102168
                  └─sudo,3102306 -i
                      └─bash,3102307
                          └─pingtunnel,3104052 -type client -l 10.xx.xx.xx:4455 -s 82.xx.xx.xxx -t 82.1xx.xx.xx:4455 -tcp 1
                              ├─{pingtunnel},3104053
                              ├─{pingtunnel},3104054
                              ├─{pingtunnel},3104055
                              └─{pingtunnel},3104056
root@VM-24-8-ubuntu:~# 
```
<a name="IVZNv"></a>
### 工具的工作原理
通过`Netsh`抓取主机发送的`ICMP`数据包(目前只支持IPV4),随后对文件进行转储,解析.由于使用`Netsh`的`trace`功能需要管理员权限<br />**所以请使用管理员权限的**`**dos(cmd)**`**或者**`**powershell**`**运行该程序!**<br />**所以请使用管理员权限的**`**dos(cmd)**`**或者**`**powershell**`**运行该程序!**<br />**所以请使用管理员权限的**`**dos(cmd)**`**或者**`**powershell**`**运行该程序!**<br />详细流程:<br />![image.png](https://cdn.nlark.com/yuque/0/2022/png/2078172/1650381458421-0f67f04e-b85e-4ff6-8ebf-efbc9b625103.png#clientId=u51461c2d-5ad4-4&crop=0&crop=0&crop=1&crop=1&from=paste&height=569&id=u505145a0&margin=%5Bobject%20Object%5D&name=image.png&originHeight=569&originWidth=603&originalType=binary&ratio=1&rotation=0&showTitle=false&size=49495&status=done&style=none&taskId=u782af941-39a4-47c5-b6b7-f3369cff937&title=&width=603)
<a name="n04j6"></a>
### 使用方法
使用管理员权限启动控制台
```go
PS C:\Users\coder\GolandProjects\FindIcmpP> .\FindIcmpP.exe -h
Usage of C:\Users\coder\GolandProjects\FindIcmpP\FindIcmpP.exe:
  -c    默认模式下不会追踪启动进程的文件，如果不选择该参数，则不会有输出文件
  -f string
        选择-po的情况下，需要通过该参数来指定ETL文件路径
  -po
        如果选择该参数，则不会进行抓包，只会解析本地etl文件
  -t uint
        在主机抓包时长，默认10s,建议不超过30s (default 10)
```
<a name="M0kdL"></a>
#### 抓包解析模式
使用`FindIcmpP.exe -t 3 -c`即可在主机抓取`ICMPV4`通信的进程
<a name="XzicA"></a>
#### 本地文件解析模式
`FindIcmpP.exe -po -f etlFilePath -c` 即可解析本地的ETL文件，从而通过ETL文件获取到通信进程
<a name="Nklz5"></a>
#### 输出
会在`etl`的文件路径下生成一个`时间戳+result.csv`文件
### 致谢
感谢坤少🦸‍♂️与乐少🦸‍♂️对我的指导,没有两位师傅windows上排查会复杂很多
