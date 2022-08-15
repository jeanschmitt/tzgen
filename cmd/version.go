package cmd

import (
	"github.com/spf13/cobra"
	"text/template"
)

var Version = "v0.3.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print tzgen version",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := template.New("versionTmpl")
		template.Must(t.Parse(cmd.VersionTemplate()))
		return t.Execute(cmd.OutOrStdout(), cmd.Root())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
