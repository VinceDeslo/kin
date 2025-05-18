package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

const (
	proxyFlag = "proxy"
	proxyFlagUsage = "please provide a teleport proxy value"
	clusterFlag = "cluster"
	clusterFlagUsage = "please provide a teleport cluster value"

	clusterSelectPrompt = "Please choose your k8s cluster"

	invalidFlagsError = "provided flags are invalid, command format should be: kin <proxy> <cluster>"
	loginError = "failed to run the teleport login command"
	kubeListError = "failed to run the teleport kube list command"
	unmarshalError = "failed to unmarshal teleport cluster list"
	kubeLoginError = "failed to run the teleport kube login command"
	cliFormError = "failed to run the command line form"
	invalidClusterError = "selected cluster is invalid"
)

var selectedCluster string

type Cluster struct {
	Name string `json:"kube_cluster_name"`
}

func main() {
	var proxy string
	var cluster string
	flag.StringVar(&proxy, proxyFlag, "", proxyFlagUsage)
	flag.StringVar(&cluster, clusterFlag, "", clusterFlagUsage)
	flag.Parse()

	if proxy == "" || cluster == "" {
		throwFatal(invalidFlagsError)
	}

	proxyArg := strings.Join([]string{"--proxy", proxy}, "=")
	tshLoginCmd := exec.Command("tsh", "login", proxyArg, cluster)
	loginOutput, err := tshLoginCmd.Output()
	if err != nil {
		throwFatal(loginError)
	}
	printCommandOutput(loginOutput)

	kubeListCmd := exec.Command("tsh", "kube", "ls", "--format=json")
	listOutput, err := kubeListCmd.Output()
	if err != nil {
		throwFatal(kubeListError)
	}

	var clusters []Cluster
	err = json.Unmarshal(listOutput, &clusters)
	if err != nil {
		throwFatal(unmarshalError)
	}

	selectOptions := []huh.Option[string]{}
	for _, cluster := range clusters {
		option := huh.NewOption(cluster.Name, cluster.Name)
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

	kubeLoginCmd := exec.Command("tsh", "kube", "login", selectedCluster)
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
