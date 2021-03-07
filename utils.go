package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

//DoFileExists checks if the file exists in the provided path.
func DoFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//Connect connects to the host and returns the connection.
func Connect(addr, user, password, publickey string, usePassword bool) (*Connection, error) {
	if usePassword == true {
		sshConfig := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		}
		conn, err := ssh.Dial("tcp", addr, sshConfig)
		if err != nil {
			fmt.Println("err while transport.")
			return nil, err
		}

		return &Connection{conn, password}, nil
	} else {
		sshConfig := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				activatepublicKey(publickey),
			},
			HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		}
		conn, err := ssh.Dial("tcp", addr, sshConfig)
		if err != nil {
			fmt.Println("err while transport.")
			return nil, err
		}

		return &Connection{conn, password}, nil
	}

}

func activatepublicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

//getTargets reads the json file and returns a slice of struct.
func getTargets(file string) []SSHHost {
	var hostlist HostList
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		log.Fatal("[!] ERROR > ", err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &hostlist)
	return hostlist.Hosts
}

//SendCommands sends the command
func (conn *Connection) SendCommands(cmds ...string) ([]byte, error) {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	in, err := session.StdinPipe()
	if err != nil {
		fmt.Println("err in stdinPipe", err.Error())

		// log.Fatal(err)
	}

	out, err := session.StdoutPipe()
	if err != nil {
		fmt.Println("err in stdoutPipe", err.Error())
		// log.Fatal(err)
	}

	var output []byte
	_, _ = in, out

	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)

			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)

			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(conn.password + "\n"))
				if err != nil {
					// break
					fmt.Println("ERR while writing sudo pwd - ", err.Error())
				}
			}
		}
	}(in, out, &output)

	cmd := strings.Join(cmds, "; ")
	_, err = session.Output(cmd)
	if err != nil {
		return []byte{}, err
	}

	return output, nil

}

//executeCommand runs the individual command
func executeCommand(host string, username string, password string, publickey string, usePassword bool, cmd string) {
	conn, err := Connect(host, username, password, publickey, usePassword)
	// conn, err := Connect(ip, username, password)
	if err != nil {
		log.Fatal("[!] Error connecting to host")
		log.Fatal(err)
	}
	// hostclr := color.New(color.FgGreen)
	red := color.New(color.FgHiWhite)
	whiteBackground := red.Add(color.BgHiBlue)
	whiteBackground.Printf("[ %s ]>", host)
	output, err := conn.SendCommands(cmd)
	if err != nil {
		fmt.Println("err in stdout", err)
	}
	fmt.Println(string(output))

}

//executeBatchCommands itereate through the host and returs commands
func executeBatchCommands(cmd string, targets []SSHHost) {
	for _, target := range targets {
		executeCommand(target.IPPort, target.Username, target.Password, target.Publickey, target.UsePassword, cmd)
	}
	color.Yellow("[!] Batch Completed")
	color.Yellow("--------------------------------")
}

// RunInteractiveMode runs the client in interactive mode.
func RunInteractiveMode(hostpath string) {
	reader := bufio.NewReader(os.Stdin)
	targets := getTargets(hostpath)
	cmdIp := color.New(color.FgCyan)
	for {
		cmdIp.Print("[cmd]> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimRight(cmd, "\r\n")
		if cmd == "quit" {
			break
		}
		executeBatchCommands(cmd, targets)
	}
}

// RunBatchMode runs the client in batchmode
func RunBatchMode(hostpath string, cmdFile string) {
	file, err := os.Open(cmdFile)
	targets := getTargets(hostpath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cmdIp := color.New(color.FgCyan)
	for scanner.Scan() {
		cmd := scanner.Text()
		cmdIp.Println("[cmd]> ", cmd)
		executeBatchCommands(cmd, targets)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
