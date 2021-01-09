## 开源端口扫描以及web(http,https)title 收集工具
## 使用说明:

### 扫描配置(默认)：
#### 程序直接运行会在本地生成配置文件，根据需要修改。
```
{
  "TargetFile":"list.txt",
  "ResultsFile":"results.txt",
  "Tasknum":1000,
  "ServerAndPorts":[
    {
      "Enable":true,
      "ServerType":"TCP",
      "ServerPort":"1-65535",
      "TimeOut":1
    },
    {
      "Enable":false,
      "ServerType":"TCP",
      "ServerPort":"21,22,23,25,53,80,110,135,137,138,139,443,1433,1434,1521,3306,3389,5000,5432,5632,6379,8000,8080,8081,8443,9090,10051,11211,27017",
      "TimeOut":1
    },
    {
      "Enable":false,
      "ServerType":"HTTP",
      "ServerPort":"80,8080",
      "TimeOut":10
    },
    {
      "Enable":true,
      "ServerType":"HTTPS",
      "ServerPort":"443,8443",
      "TimeOut":10
    }
  ]
}
```

### list.txt 示例:

```
[root@localhost portScan]# cat ip.txt 
192.168.0.1-8 
192.168.0.0/24
```



## 输出：
###端口开放信息

格式:
```
[ip]:[port]\t[opend]\n
```
结果示例:
```
192.xx.xx.92:111	opend
192.xx.xx.92:22	opend
192.xx.xx.92:3306	opend
192.xx.xx.92:7001	opend
192.xx.xx.97:80	opend
192.xx.xx.97:22	opend
192.xx.xx.97:8080	opend
```

### web titel扫描结果 :

格式:
```
[url]\t[statusCode]\t[title]\n
```
结果示例:
```cassandraql
https://10.xx.xx.10:443/	200	xxx系统
https://10.xx.xx.95:443/	200	xxx
https://10.xx.xx.9:443/	200	xxxx系统
https://10.xx.xx.6:443/	200	
https://10.xx.xx.12:443/	200	
https://10.xx.xx.91:443/	200	
https://10.xx.xx.7:443/	200	xxxx系统

```

## 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

禁止对本软件实施逆向工程、反编译、试图破译源代码等行为。

**如果发现上述禁止行为，我们将保留追究您法律责任的权利。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您**务必审慎阅读、充分理解各条款内容**，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。
除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。