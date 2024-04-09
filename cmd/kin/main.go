package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

var selectedCluster string
const (
	proxyFlag = "proxy"
	proxyFlagUsage = "please provide a teleport proxy value"
	clusterFlag = "cluster"
	clusterFlagUsage = "please provide a teleport cluster value"

	tshCommand = "tsh"
	loginCommand = "login"
	kubeCommand = "kube"
	listCommand = "list"
	abbreviatedListCommand = "ls"

	proxyArg = "--proxy"
	argSeparator = "="

	clusterSelectPrompt = "Please choose your k8s cluster"

	invalidFlagsError = "provided flags are invalid, command format should be: kin <proxy> <cluster>"
	loginError = "failed to run the teleport login command"
	kubeListError = "failed to run the teleport kube list command"
	kubeLoginError = "failed to run the teleport kube login command"
	cliFormError = "failed to run the command line form"
	invalidClusterError = "selected cluster is invalid"

	clusterEntryDelimiter = " "
	emptyString = ""
)

func main() {
	var proxy string
	var cluster string
	flag.StringVar(&proxy, proxyFlag, emptyString, proxyFlagUsage)
	flag.StringVar(&cluster, clusterFlag, emptyString, clusterFlagUsage)
	flag.Parse()

	if proxy == emptyString || cluster == emptyString {
		throwFatal(invalidFlagsError)
	}

	tshLoginCmd := exec.Command(
		tshCommand,
		loginCommand,
		strings.Join([]string{proxyArg, proxy}, argSeparator),
		cluster,
	)

	loginOutput, err := tshLoginCmd.Output()
	if err != nil {
		throwFatal(loginError)
	}
	printCommandOutput(loginOutput)

	kubeListCmd := exec.Command(
		tshCommand,
		kubeCommand,
		abbreviatedListCommand,
	)
	listOutput, err := kubeListCmd.Output()
	if err != nil {
		throwFatal(kubeListError)
	}

	lines := strings.Split(string(listOutput), "\n")

	// Remove text headers
	entries := lines[2:]
	
	// Extract available clusters
	kubernetesClusters := []string{}
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
		throwFatal(cliFormError)
	}

	if len(selectedCluster) == 0 {
		throwFatal(invalidClusterError)
	}

	kubeLoginCmd := exec.Command(
		tshCommand,
		kubeCommand,
		loginCommand,
		selectedCluster,
	)
	kubeLoginOut, err := kubeLoginCmd.Output()
	if err != nil {
		throwFatal(kubeLoginError)
	}

	printCommandOutput(kubeLoginOut)
	os.Exit(0)
}

func throwFatal(errorMessage string){
	slog.Error(errorMessage)
	os.Exit(1)
}

func printCommandOutput(output []byte){
	fmt.Printf("%v\n", string(output))
}
