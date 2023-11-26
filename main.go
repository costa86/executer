package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	port     int
	host     string
	username string
	password string
	command  string
	version  string = "1.0.0"
)

func handleFailure(e error) {
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "executer",
	Short: fmt.Sprintf("Executes SSH commands on a remote host. Version %s", version),
	Run: func(cmd *cobra.Command, args []string) {

		config := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)

		handleFailure(err)

		defer client.Close()

		session, err := client.NewSession()

		handleFailure(err)

		defer session.Close()

		_, err = session.CombinedOutput(command)

		handleFailure(err)

		fmt.Printf("Command '%s' executed on '%s'", command, host)

		os.Exit(0)
	}}

func init() {
	rootCmd.Flags().StringVarP(&host, "host", "t", "", "SSH server host")
	rootCmd.Flags().IntVarP(&port, "port", "p", 22, "SSH server port")
	rootCmd.Flags().StringVarP(&username, "username", "u", "", "SSH username")
	rootCmd.Flags().StringVarP(&password, "password", "w", "", "SSH password")
	rootCmd.Flags().StringVarP(&command, "cmd", "c", "", "Command to execute")

	rootCmd.MarkFlagRequired("host")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")
	rootCmd.MarkFlagRequired("cmd")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		handleFailure(err)
	}
}
