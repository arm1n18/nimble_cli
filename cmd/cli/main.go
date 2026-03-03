package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/arm1n18/nimble_cli/commands"
	"github.com/arm1n18/nimble_cli/parser"
	"github.com/arm1n18/nimble_cli/protocol"
)

type Config struct {
	host     string
	port     string
	user     string
	password string
}

func safeFatal(conn net.Conn, msg string) {
	protocol.Error(msg)
	fmt.Fprintln(conn, "EXIT")
	conn.Close()
	fmt.Println("\nConnection closed")
	os.Exit(1)
}

func connectToServer(config *Config) {
	conn, err := net.Dial("tcp", config.host+":"+config.port)
	if err != nil {
		protocol.FatalError("Failed to connect %v\n", err)
	}

	defer func() {
		fmt.Fprintln(conn, "EXIT")
		conn.Close()
		fmt.Println("Connection closed")
		os.Exit(1)
	}()

	serverReader := bufio.NewReader(conn)
	reader := bufio.NewReader(os.Stdin)

	if config.user != "" {
		fmt.Fprintf(conn, "AUTH %s %s\n", config.user, config.password)
		resp, _ := serverReader.ReadString('\n')

		value, ok := parser.GetParsedValue(strings.TrimSpace(resp))
		if !ok {
			safeFatal(conn, value)
		} else {
			fmt.Println(value)
		}
	}

	for {
		fmt.Printf(protocol.Bold + protocol.Cyan + "nimble =# " + protocol.Reset)
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.ToLower(line) == "exit" {
			break
		}

		if strings.ToLower(line) == "\\help" {
			state := struct {
				cat int
				cmd int
			}{0, 0}

			maxLen := struct {
				lenCmd  int
				lenDesc int
			}{12, 55}

			reader := bufio.NewReader(os.Stdin)
			fmt.Println("AVAILABLE COMMANDS:")

			done := false
			for !done {
				if state.cat >= len(commands.Commands) {
					break
				}

				cat := commands.Commands[state.cat]

				if state.cmd == 0 {
					fmt.Fprint(os.Stdout, "\033[1A\033[2K\r")
					fmt.Printf("\n%s\n", strings.Repeat("-", 75))
					fmt.Println()
					protocol.BoldText(fmt.Sprintf("%*s", (75+len(cat.Cat))/2, cat.Cat))
					fmt.Println()
				}

				if state.cmd >= len(cat.Cmds) {
					state.cat++
					state.cmd = 0
					continue
				}

				c := cat.Cmds[state.cmd]
				if state.cmd > 0 {
					fmt.Fprint(os.Stdout, "\033[1A\033[2K\r")
				}
				protocol.Commad(c.Cmd, c.Desc, c.Syntax, maxLen.lenCmd+2, maxLen.lenDesc)

				fmt.Print("Press Enter to continue, q to quit =# ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				if strings.ToLower(input) == "q" {
					done = true
					break
				}

				state.cmd++
			}
			continue
		}

		if strings.ToLower(line) == "cls" {
			protocol.ClearScreen()
			continue
		}

		fmt.Fprintf(conn, "%s\n", line)
		resp, err := serverReader.ReadString('\n')
		if err != nil {
			break
		}

		value, ok := parser.GetParsedValue(resp)
		if !ok {
			protocol.Error(value)
		} else {
			fmt.Print(value)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf(protocol.Bold + protocol.Cyan + "Host [localhost]: " + protocol.Reset)
	hostInput, _ := reader.ReadString('\n')
	hostInput = strings.TrimSpace(hostInput)
	if hostInput == "" {
		hostInput = "localhost"
	}

	fmt.Printf(protocol.Bold + protocol.Cyan + "Port [8085]: " + protocol.Reset)
	portInput, _ := reader.ReadString('\n')
	portInput = strings.TrimSpace(portInput)
	if portInput == "" {
		portInput = "8085"
	}

	fmt.Printf(protocol.Bold + protocol.Cyan + "User [nimble]: " + protocol.Reset)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	if userInput == "" {
		userInput = "nimble"
	}

	fmt.Printf(protocol.Bold + protocol.Cyan + "Password [required]: " + protocol.Reset)
	passwordInput, _ := reader.ReadString('\n')
	passwordInput = strings.TrimSpace(passwordInput)

	if strings.Trim(passwordInput, "") == "" {
		protocol.FatalError("no password provided")
	}

	connectToServer(&Config{
		host:     hostInput,
		port:     portInput,
		user:     userInput,
		password: passwordInput,
	})
}
