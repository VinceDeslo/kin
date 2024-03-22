package main

import (
	"log/slog"
	"os"

	"github.com/charmbracelet/huh"
)

func main() {
	// mock data
	kubernetesClusters := []string{
		"polaris-pre-prod",
		"polaris-prod",
	}
	var selectedCluster string

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

	err := form.Run()
	if err != nil {
		slog.Error("Failed to run the CLI form")
		os.Exit(1)
	}

	if selectedCluster == "" {
		slog.Warn("Invalid cluster selection")
	}

	slog.Info("Selection received", "cluster:", selectedCluster)
	
	return
}
