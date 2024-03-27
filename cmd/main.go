package main

import (
	coredogcontroller "github.com/DomineCore/coredog/coredog-controller"
	coredogwatcher "github.com/DomineCore/coredog/coredog-watcher"
	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{}

	watcherBootstrap := cobra.Command{
		Use: "watcher",
		RunE: func(cmd *cobra.Command, args []string) error {
			coredogwatcher.WatchCorefile()
			return nil
		},
		Long: "start a watcher agent on host to watch corefile created.",
	}

	controllerBootstrap := cobra.Command{
		Use: "controller",
		RunE: func(cmd *cobra.Command, args []string) error {
			coredogcontroller.ListenAndServe()
			return nil
		},
		Long: "start a controller server to subscription the corefile and process.",
	}

	root.AddCommand(&watcherBootstrap, &controllerBootstrap)
	root.Execute()
}
