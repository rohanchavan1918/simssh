package main

import (
	"flag"
	"os"

	"github.com/fatih/color"
)

func main() {

	hosts := flag.String("hosts", "", "Input json file, for more info check :- https://github.com/rohanchavan1918/simssh ")
	mode := flag.String("mode", "interactive", "Mode to run. Available modes { interactive | batch }")
	cmdFile := flag.String("cmd_file", "", "The batch of commands to run on all servers")

	flag.Parse()

	// Check if the hosts file exists
	switch *hosts {
	case "":
		color.Red("[Error] Hosts file not passed.")
		flag.PrintDefaults()
		os.Exit(0)
	default:
		if !DoFileExists(*hosts) {
			color.Red("[ERROR] Commands/batch file passsed in argument doesnot exist.")
			os.Exit(1)
		}
	}

	// Check mode exists
	switch *mode {
	case "":
		color.Cyan("[INFO] Starting with interactive mode (default)")
	case "interactive":
		color.Cyan("[INFO] Starting with interactive mode")
	case "batch":
		color.Cyan("[INFO] Starting with batch mode")
		// Check if the file path is passed in the args
		if *cmdFile == "" {
			color.Red("[ERROR] Commands/batch file not given in arguments.")
			os.Exit(0)
		}
		if !DoFileExists(*cmdFile) {
			color.Red("[ERROR] Commands/batch file passsed in argument doesnot exist.")
			os.Exit(0)
		}
	}

	if *mode == "interactive" {
		RunInteractiveMode(*hosts)
	} else if *mode == "batch" {
		RunBatchMode(*hosts, *cmdFile)
	}

	// reader := bufio.NewReader(os.Stdin)
	// targets := getTargets(*hosts)
	// cmdIp := color.New(color.FgCyan)
	// for {
	// 	cmdIp.Print("[cmd]> ")
	// 	cmd, _ := reader.ReadString('\n')
	// 	cmd = strings.TrimRight(cmd, "\r\n")
	// 	if cmd == "quit" {
	// 		break
	// 	}
	// 	executeBatchCommands(cmd, targets)
	// }

}
