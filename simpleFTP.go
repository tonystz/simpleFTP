// test project main.go
package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	BUFF_SIZE     = 512
	READ_DEADLINE = 10
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func pass227(str string) string {
	re := regexp.MustCompile(".*[(](.*)[)].*")

	sp := strings.Split(re.FindStringSubmatch(str)[1], ",")
	p1, _ := strconv.Atoi(sp[4])
	p2, _ := strconv.Atoi(sp[5])
	addr := fmt.Sprintf("%s.%s.%s.%s:%d", sp[0], sp[1], sp[2], sp[3], p1*256+p2)

	fmt.Println("Server is listen on ", addr)

	return addr
}
func recv(r io.Reader) string {
	var err error
	var cnt, n int

	data := bytes.NewBuffer(nil)
	buff := make([]byte, BUFF_SIZE)

	for true {
		n, err = r.Read(buff[0:])
		cnt += n
		data.Write(buff[:n])
		if n < BUFF_SIZE {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	trim := strings.TrimSpace(data.String())
	fmt.Printf("Response:size=%v data=%v\n", cnt, trim)
	return trim

}

func (ftp *FTP) Open(addr string) net.Conn {
	c, err := net.Dial("tcp", addr)
	check(err)
	return c
}

type FTP struct {
	Conn net.Conn
	err  error
}

func (ftp *FTP) SendCmd(cmd string) (string, int) {
	fmt.Printf("Request:%s\n", cmd)
	fmt.Fprintf(ftp.Conn, "%s\r\n", cmd)
	ftp.Conn.SetReadDeadline(time.Now().Add(READ_DEADLINE * time.Second))

	res := recv(ftp.Conn)
	code, _ := strconv.Atoi(res[0:3])
	return res, code
}

func (ftp *FTP) SetPasv() net.Conn {

	res, code := ftp.SendCmd("PASV")

	if 227 == code {
		c := ftp.Open(pass227(res))
		fmt.Println("Parse PASV 227 code done")
		return c
	} else {
		panic("Cannot get 227 request")
		return nil
	}
}

func (ftp *FTP) List(dir string) {

	c := ftp.SetPasv()
	ftp.SendCmd("LIST " + dir)
	recv(c)

	defer func() {
		c.Close()
	}()
}
func (ftp *FTP) GetFile(file string) {
	ftp.SendCmd("SIZE " + file)
	ftp.SendCmd("TYPE I")
	c := ftp.SetPasv()
	ftp.SendCmd("RETR " + file)

	out, _ := os.Create(path.Base(file))
	n, err := io.Copy(out, c)
	check(err)
	fmt.Printf("Store %s[%d]\n", path.Base(file), n)
	defer func() {
		out.Close()
		c.Close()
	}()
}

func (ftp *FTP) New(args ...string) {
	var addr, user, pass string

	switch len(args) {
	case 1:
		addr = args[0] + ":21"
		user = "anonymous"
		pass = ""
	case 3:
		addr = args[0] + ":21"
		user = args[1]
		pass = args[2]
	default:
		panic("error,please specify the FTP server host at least")
	}
	ftp.Conn = ftp.Open(addr)
	recv(ftp.Conn)
	ftp.SendCmd("USER " + user)
	ftp.SendCmd("PASS " + pass)

}

func (ftp *FTP) Close() {

	ftp.SendCmd("QUIT")
	ftp.Conn.Close()
}

func main() {

	ftp := FTP{}
	ftp.New("10.64.70.73")
	ftp.List("/")
	ftp.Close()

	ftp.New("10.64.70.73")
	ftp.GetFile("/pub/atop-1.27-3.tar.gz")
	ftp.Close()
}
