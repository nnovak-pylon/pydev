package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	workspaceFolder = "/home/nicholas/git/braid/"
	configName      = "/home/nicholas/git/braid/.devcontainer/local/devcontainer.json"
	containerName   = "pylon-dev"
)

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

type DockerStatus struct {
	State  string `json:"State"`
	Status string `json:"Status"`
}

var StatusDevcontainer = &cobra.Command{
	Use:     "status",
	Short:   "Displays the current status of the development environment",
	Aliases: []string{"ps"},
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("Getting status of container...")

		statusCmd := exec.Command("docker", "ps", "-f", fmt.Sprintf("name=%s", containerName), "--format", "json")
		statusCmd.Stderr = os.Stderr

		var status DockerStatus

		if out, err := statusCmd.Output(); err != nil {
			return err
		} else {

			if len(out) == 0 {
				fmt.Println("No output from command, container not running")
				return nil
			}

			if err := json.Unmarshal(out, &status); err != nil {
				return err
			}

			fmt.Printf("State: %s, Status: %s\n", status.State, status.Status)
		}

		return nil
	},
}

var CleanContainerCommand = &cobra.Command{
	Use: "clean",
	RunE: func(cmd *cobra.Command, args []string) error {

		cleanCmd := exec.Command("docker", "rm", containerName)
		cleanCmd.Stderr = os.Stderr
		cleanCmd.Stdout = os.Stdout

		if err := cleanCmd.Run(); err != nil {
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
