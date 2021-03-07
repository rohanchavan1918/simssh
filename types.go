package main

import "golang.org/x/crypto/ssh"

//Connection is a struct with ssh connection and password
type Connection struct {
	*ssh.Client
	password string
}

//SSHHost is type for ssh host
type SSHHost struct {
	IPPort      string `json:"host"`
	Username    string `json:"username"`
	UsePassword bool   `json:"use_password"`
	Password    string `json:"password"`
	Publickey   string `json:"publickey"`
}

//HostList is list of hosts accepted by the json input file
type HostList struct {
	Hosts []SSHHost `json:"hosts"`
}

//CmdList slice of all commands returned from
type CmdList struct {
	Cmd []string
}
