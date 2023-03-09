# pop3crack

## 使用方法  

```bash
  -h help  
        语法信息  
  -d string  
        请输入邮箱域名或者ip地址,eg:mail.test.com.cn (default "192.168.0.1")
  -m string
        指定模式，eg:pop3/pop3s (default "pop3")
  -o string
        output result (default "result.txt")
  -p string
        请输入端口,eg:110、995 (default "110")
  -pd string
        请输入密码字典
  -pwd string
        请输入密码
  -ud string
        请输入用户名字典
  -user string
        请输入用户名
       
```  

## 交叉编译（mac） 

```bash
# mac
CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o pop3crack-darwin-amd64 main.go
# windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o op3crack-darwin-amd64.exe main.go
linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o pop3crack-linux-amd64 main.go 

```
