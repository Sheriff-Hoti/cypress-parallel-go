/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/Sheriff-Hoti/cypress-parallel-go/core"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cypress-parallel-go",
	Short: "Run cypress in parallel",
	Long: `Run cypress in parallel:

cypress-parallel-go is a binary written in go that runs cypress specs in "parallel".`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		tool, tool_err := cmd.Flags().GetString("tool")
		dir, dir_err := cmd.Flags().GetString("dir")
		cyArgs, args_err := cmd.Flags().GetString("args")
		script, script_err := cmd.Flags().GetString("script")
		cores, cores_err := cmd.Flags().GetInt16("cores")

		if tool_err != nil {
			return tool_err
		}
		if dir_err != nil {
			return dir_err
		}

		if args_err != nil {
			return args_err
		}

		if script_err != nil {
			return script_err
		}

		if cores_err != nil {
			return cores_err
		}

		start := time.Now()

		if err := core.Run(tool, dir, script, cyArgs, cores); err != nil {
			return err
		}

		elapsed := time.Since(start)
		fmt.Printf("The time it took %s", elapsed)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cypress-parallel-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "tg", false, "Help message for toggle")
	rootCmd.Flags().StringP("tool", "t", "", "specify the tool to execute your Cypress command(required)")
	rootCmd.Flags().StringP("dir", "d", "", "Cypress specs directory (required)")
	rootCmd.Flags().StringP("script", "s", "cypress run", "Your npm Cypress command")
	rootCmd.Flags().Int16P("cores", "c", 2, "Number of cores")
	rootCmd.Flags().StringP("args", "a", "cypress run", "Your npm Cypress command arguments")

	// rootCmd.RegisterFlagCompletionFunc("tool", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// 	return []string{"docker", "yarn"}, cobra.ShellCompDirectiveNoFileComp
	// })

	rootCmd.MarkFlagRequired("tool")
	rootCmd.MarkFlagRequired("dir")
}
