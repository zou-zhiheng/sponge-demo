package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * @Author: zhm
 * @Description: ssh登录操作包
 * @Version: 1.0.0
 * @Date: 2023/10/16 17:53
 */
type SSHClient struct {
	IP         string      //IP地址
	Username   string      //用户名
	Password   string      //ssh密码 or 秘钥登录密码
	Port       uint32      //端口号
	client     *ssh.Client //ssh客户端
	LastResult string      //最近一次Run的结果
	PrivKey    []byte      //私钥 若是私钥不会空，则说明使用秘钥登录
}

// 创建命令行对象
// @param ip IP地址
// @param username 用户名
// @param password 密码
// @param port 端口号,默认22
func NewCli(ip string, username string, password string, port ...uint32) *SSHClient {
	cli := new(SSHClient)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}
	return cli
}
func NewCliByPrivate(ip string, username string, password string, pemPrivate []byte, port ...uint32) *SSHClient {
	cli := new(SSHClient)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}
	cli.PrivKey = pemPrivate
	return cli
}

func (c *SSHClient) SetPrivKey(privKey []byte) {
	c.PrivKey = privKey
}

func (c *SSHClient) NewSSHClientWithPrivkey(ip string, port int, privkey []byte) (cli *SSHClient, err error) {
	signer, err := ssh.ParsePrivateKey(privkey)
	if err != nil {
		return
	}
	config := ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		//Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	cli = &SSHClient{
		IP:     ip,
		Port:   uint32(port),
		client: conn,
	}
	return
}

// @Title Connect
// @Description   Connect ssh连接
// @Author  zhm
// @Date 2023-10-16 18:26:18
// @Return error 错误信息
func (c *SSHClient) Connect() error {
	var config ssh.ClientConfig
	if len(c.PrivKey) > 0 {
		var signer ssh.Signer
		var err error
		if len(c.Password) > 0 {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(c.PrivKey, []byte(c.Password))
			if err != nil {
				return err
			}
		} else {
			signer, err = ssh.ParsePrivateKey(c.PrivKey)
			if err != nil {
				return err
			}
		}

		config = ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
		}
	} else {
		config = ssh.ClientConfig{
			User: c.Username,
			Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
		}
	}

	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}

// 执行shell
// @param shell shell脚本命令
func (c SSHClient) Run(shell string) (string, error) {
	if c.client == nil {
		if err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer c.client.Close()
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err
}

// 执行带交互的命令
func (c *SSHClient) RunTerminal(shell string, stdout, stderr io.Writer) error {
	if c.client == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = stdout
	session.Stderr = stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}

	session.Run(shell)
	return nil
}

func GroupUsedDiskStringSpit(groupUsedDiskString string) map[string]int {
	resultString := make(map[string]int)
	for _, v := range strings.Split(groupUsedDiskString, "\n") {
		info := strings.Fields(v)
		if len(info) != 2 {
			continue
		}

		disk, err := strconv.Atoi(info[0])
		if err != nil {
			resultString[info[1]] = 0
		} else {
			resultString[info[1]] = disk
		}
	}

	return resultString
}

func FileUsedStringSpit(fileDiskString string, startDiskString string, endDiskString string) map[string]int64 {

	resultString := make(map[string]int64)
	var totalDiskCapacity int64
	var usedDiskCapacity int64
	//对结果按行进行分割
	for _, v := range strings.Split(fileDiskString, "\n") {
		if strings.HasPrefix(v, startDiskString) && strings.HasSuffix(v, endDiskString) {
			//按空格进行分割
			info := strings.Fields(v)

			if len(info) > 3 {
				totalDiskCapacity = DisposeUnit(info[1])
				usedDiskCapacity = DisposeUnit(info[2])
			}
		}
	}
	resultString["totalDiskCapacity"] = totalDiskCapacity
	resultString["usedDiskCapacity"] = usedDiskCapacity

	return resultString
}

func DisposeUnit(content string) int64 {
	s := strings.ToUpper(content)

	lastStr := s[len(s)-1:]
	num := s[0:(len(s) - 1)]
	parseFloat, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0
	}

	switch lastStr {
	case "M":
		return int64(parseFloat)
	case "G":
		return int64(math.Ceil(parseFloat * float64(1024)))
	case "T":
		return int64(math.Ceil(parseFloat * float64(1024) * float64(1024)))
	case "K":
	case "KB":
		return int64(math.Ceil(parseFloat / float64(1024)))
	default:
		return 0
	}

	return 0

}
