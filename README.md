```
[root@localhost portScan]# ./portScan -h


                               $$\      $$$$$$\                               
                               $$ |    $$  __$$\                              
 $$$$$$\   $$$$$$\   $$$$$$\ $$$$$$\   $$ /  \__| $$$$$$$\ $$$$$$\  $$$$$$$\  
$$  __$$\ $$  __$$\ $$  __$$\\_$$  _|  \$$$$$$\  $$  _____|\____$$\ $$  __$$\ 
$$ /  $$ |$$ /  $$ |$$ |  \__| $$ |     \____$$\ $$ /      $$$$$$$ |$$ |  $$ |
$$ |  $$ |$$ |  $$ |$$ |       $$ |$$\ $$\   $$ |$$ |     $$  __$$ |$$ |  $$ |
$$$$$$$  |\$$$$$$  |$$ |       \$$$$  |\$$$$$$  |\$$$$$$$\\$$$$$$$ |$$ |  $$ |
$$  ____/  \______/ \__|        \____/  \______/  \_______|\_______|\__|  \__|
$$ |                                                                          
$$ |                                                                          
\__|                                                                          


Usage of ./portScan:
  -file
    	Use file mode to specify ip address .
  -full
    	Scan all TCP and UDP ports in full scan mode. The default is off. By default, only common TCP ports are scanned.
  -ip string
    	IP to be scanned, supports three formats:
    	192.168.0.1 
    	192.168.0.1-8 
    	192.168.0.0/24
  -p string
    	Port to be scanned, supports three formats:
    	80
    	22,80 
    	22:65535,21
  -t int
    	Maximum number of threads (default 10000)
```



## 使用说明:

### 文件指定扫描ip列表：

```
./portScan --file -ip ip.txt -p 1:65535 -t 1000
```

### ip.txt 示例:

```
[root@localhost portScan]# cat ip.txt 
192.168.0.1-8 
192.168.0.0/24
```

### 不使用文件方式：

```
./portScan -ip 192.168.0.21-100 -p 1:65535 -t 1000
```



## 输出：

```
192.168.0.92:111	opend
192.168.0.92:22	opend
192.168.0.92:3306	opend
192.168.0.92:7001	opend
192.168.0.97:80	opend
192.168.0.97:22	opend
192.168.0.97:8080	opend
```