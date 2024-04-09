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
const (
	tshCommand = "tsh"
	loginCommand = "login"
	kubeCommand = "kube"
	listCommand = "list"
	abbreviatedListCommand = "ls"

	proxyFlag = "--proxy"
	argSeparator = "="

	clusterSelectPrompt = "Please choose your k8s cluster"

	loginError = "failed to run the teleport login command"
	kubeListError = "failed to run the teleport kube list command"
	kubeLoginError = "failed to run the teleport kube login command"
	cliFormError = "failed to run the command line form"
	invalidClusterError = "invalid cluster selection"

	clusterEntryDelimiter = " "
	emptyString = ""
)

func main() {
	kubernetesClusters := []string{}

	// Input by CLI tool
	cluster := "snyk.teleport.sh"
	proxy :=  "snyk.teleport.sh:443"

	// Login to the teleport proxy
	tshLoginCmd := exec.Command(
		tshCommand,
		loginCommand,
		strings.Join([]string{proxyFlag, proxy}, argSeparator),
		cluster,
	)

	loginOutput, err := tshLoginCmd.Output()
	if err != nil {
		slog.Error(loginError)
		os.Exit(1)
	}
	fmt.Printf("%v\n", string(loginOutput))

	kubeListCmd := exec.Command(
		tshCommand,
		kubeCommand,
		abbreviatedListCommand,
	)
	listOutput, err := kubeListCmd.Output()
	if err != nil {
		slog.Error(kubeListError)
		os.Exit(1)
	}

	lines := strings.Split(string(listOutput), "\n")

	// Remove headers
	entries := lines[2:]
	
	// Extract available clusters
	for _, entry := range entries {
		clusterName := strings.Split(entry, clusterEntryDelimiter)[0]
		if len(clusterName) > 0 {
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
		Title(clusterSelectPrompt).
		Options(selectOptions...).
		Value(&selectedCluster)

	group := huh.NewGroup(selectComponent)
	form := huh.NewForm(group)

	err = form.Run()
	if err != nil {
		slog.Error(cliFormError)
		os.Exit(1)
	}

	if len(selectedCluster) == 0 {
		slog.Warn(invalidClusterError)
	}

	kubeLoginCmd := exec.Command(
		tshCommand,
		kubeCommand,
		loginCommand,
		selectedCluster,
	)
	kubeLoginOut, err := kubeLoginCmd.Output()
	if err != nil {
		slog.Error(kubeLoginError)
		os.Exit(1)
	}

	fmt.Printf("%v\n", string(kubeLoginOut))
	os.Exit(0)
}
