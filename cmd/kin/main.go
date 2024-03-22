package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

var selectedCluster string

func main() {
	kubernetesClusters := []string{}

	// Login to the teleport proxy
	tshLogin := exec.Command(
		"tsh",
		"login",
		"--proxy=snyk.teleport.sh:443",
		"snyk.teleport.sh",
	)
	loginOutput, err := tshLogin.Output()
	if err != nil {
		slog.Error("Failed to run the Teleport login command")
		os.Exit(1)
	}
	fmt.Printf("%v\n", string(loginOutput))

	// List kubernetes clusters
	tshCmd := exec.Command("tsh", "kube", "ls")
	listOutput, err := tshCmd.Output()
	if err != nil {
		slog.Error("Failed to run the Teleport list command")
		os.Exit(1)
	}

	lines := strings.Split(string(listOutput), "\n")

	// Remove headers
	entries := lines[2:]
	
	// Extract available clusters
	for _, entry := range entries {
		clusterName := strings.Split(entry, " ")[0]
		if clusterName != "" {
			kubernetesClusters = append(kubernetesClusters, clusterName)
		}
	}

	// Build user selection options
	selectOptions := []huh.Option[string]{}
	for _, cluster := range kubernetesClusters {
		option := huh.NewOption(cluster, cluster)
		selectOptions = append(selectOptions, option)
	}

	selectComponent := huh.NewSelect[string]().
		Title("Choose your k8s cluster").
		Options(selectOptions...).
		Value(&selectedCluster)

	form := huh.NewForm(
		huh.NewGroup(
			selectComponent,
		),
	)

	err = form.Run()
	if err != nil {
		slog.Error("Failed to run the CLI form")
		os.Exit(1)
	}

	if selectedCluster == "" {
		slog.Warn("Invalid cluster selection")
	}

	clusterLoginCmd := exec.Command("tsh", "kube", "login", selectedCluster)
	clusterLoginOut, err := clusterLoginCmd.Output()
	if err != nil {
		slog.Error("Failed to run the Teleport cluster login command")
		os.Exit(1)
	}

	fmt.Printf("%v\n", string(clusterLoginOut))
	os.Exit(0)
}
