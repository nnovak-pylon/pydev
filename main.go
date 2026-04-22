package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	workspaceFolder = "/home/nicholas/git/braid/"
	configName      = "/home/nicholas/git/braid/.devcontainer/local/devcontainer.json"
	containerName   = "pylon-dev"
)

var ContainerNotRunningError = fmt.Errorf("No output from command, container not running")

type DockerStatus struct {
	State  string `json:"State"`
	Status string `json:"Status"`
	Image  string `json:"Image"`
}

func containerStatus() (DockerStatus, error) {
	statusCmd := exec.Command("docker", "ps", "-a", "-f", fmt.Sprintf("name=%s", containerName), "--format", "json")
	statusCmd.Stderr = os.Stderr

	var status DockerStatus

	if out, err := statusCmd.Output(); err != nil {
		return status, err
	} else {

		if len(out) == 0 {
			return status, ContainerNotRunningError
		}

		if err := json.Unmarshal(out, &status); err != nil {
			return status, err
		}
	}

	return status, nil
}

var BuildDevcontainer = &cobra.Command{
	Use:   "build",
	Short: "Builds the devcontainer container",
	RunE: func(cmd *cobra.Command, args []string) error {
		startCmd := exec.Command("devcontainer", "--workspace-folder", workspaceFolder, "--config", configName, "build")
		startCmd.Stderr = os.Stderr
		startCmd.Stdout = os.Stdout
		if err := startCmd.Run(); err != nil {
			return err
		}

		return nil
	},
}

var StartDevcontainer = &cobra.Command{
	Use:     "up",
	Short:   "Starts the development environment",
	Aliases: []string{"start"},
	RunE: func(cmd *cobra.Command, args []string) error {

		startCmd := exec.Command("devcontainer", "up", "--workspace-folder", workspaceFolder, "--config", configName)
		startCmd.Stderr = os.Stderr
		startCmd.Stdout = os.Stdout
		if err := startCmd.Run(); err != nil {
			return err
		}

		return nil
	},
}

var StopDevcontainer = &cobra.Command{
	Use:     "down",
	Short:   "Stops the development environment",
	Aliases: []string{"stop"},
	RunE: func(cmd *cobra.Command, args []string) error {

		stopCmd := exec.Command("docker", "stop", containerName)

		stopCmd.Stdout = os.Stdout
		stopCmd.Stderr = os.Stderr

		if err := stopCmd.Run(); err != nil {
			return err
		}

		return nil
	},
}

var ExecCommand = &cobra.Command{
	Use:     "exec",
	Short:   "Starts a shell in the local development environment",
	Aliases: []string{"shell"},
	RunE: func(cmd *cobra.Command, args []string) error {

		status, err := containerStatus()
		if err != nil {
			slog.Error("error getting status of devcontainer", "error", err)
		}

		switch status.State {
		case "exited":
		case "stopped":
			slog.Info("Devcontainer was found in a halted status, do you want to start it? [yY]", "status", status.State)
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune() // Reads one Unicode character

			if err != nil {
				slog.Error("error reading input from stdin", "error", err)
			}

			switch char {
			case 'y':
			case 'Y':
				err := StartDevcontainer.RunE(cmd, args)
				if err != nil {
					slog.Error("error starting devcontainer", "error", err)
				}
			}
		}

		execCommand := exec.Command("devcontainer", "exec", "--workspace-folder", workspaceFolder, "--config", configName, "bash")

		execCommand.Stderr = os.Stderr
		execCommand.Stdout = os.Stdout

		execCommand.Stdin = os.Stdin

		if err := execCommand.Run(); err != nil {
			return err
		}

		return nil
	},
}

var StatusDevcontainer = &cobra.Command{
	Use:     "status",
	Short:   "Displays the current status of the development environment",
	Aliases: []string{"ps"},
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("Getting status of container...")

		status, err := containerStatus()
		if err != nil {
			return err
		}

		fmt.Printf("State: %s, Status: %s\n", status.State, status.Status)
		return nil
	},
}

var CleanContainerCommand = &cobra.Command{
	Use: "clean",
	RunE: func(cmd *cobra.Command, args []string) error {

		status, err := containerStatus()
		if err != nil {
			return err
		}

		slog.Info("Stopping development container", "containerName", containerName)

		cleanCmd := exec.Command("docker", "rm", containerName)

		if out, err := cleanCmd.CombinedOutput(); err != nil {
			if strings.Contains(string(out), "No such container") {
				slog.Warn("Could not find container to remove, skipping")
			} else {
				return err
			}
		}

		slog.Info("Removing image %s for container", "imageName", status.Image)

		if err := exec.Command("docker", "rmi", status.Image).Run(); err != nil {
			return err
		}

		return nil
	},
}

func main() {
	rootCmd := &cobra.Command{
		Use: "pydev",
	}

	rootCmd.AddCommand(BuildDevcontainer)
	rootCmd.AddCommand(StartDevcontainer)
	rootCmd.AddCommand(StopDevcontainer)
	rootCmd.AddCommand(StatusDevcontainer)
	rootCmd.AddCommand(ExecCommand)
	rootCmd.AddCommand(CleanContainerCommand)

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Error executing command", "error", err)
	}
}
