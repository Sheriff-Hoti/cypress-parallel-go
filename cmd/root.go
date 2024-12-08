/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
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
	Run: func(cmd *cobra.Command, args []string) {
		tool, tool_err := cmd.Flags().GetString("tool")
		dir, dir_err := cmd.Flags().GetString("dir")
		if tool_err != nil {
			log.Panic(tool_err)
		}
		if dir_err != nil {
			log.Panic(dir_err)
		}

		if tool != "docker" && tool != "yarn" {
			log.Fatal("invalid value for --tool or -t; must be 'docker' or 'yarn'")
		}
		start := time.Now()

		core.Run(tool, dir)

		elapsed := time.Since(start)
		fmt.Printf("The time it took %s", elapsed)
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
	rootCmd.Flags().StringP("tool", "t", "", "specify the tool to execute the cypress run command(required)")
	rootCmd.Flags().StringP("dir", "d", "", "Cypress specs directory (required)")

	// rootCmd.RegisterFlagCompletionFunc("tool", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// 	return []string{"docker", "yarn"}, cobra.ShellCompDirectiveNoFileComp
	// })

	rootCmd.MarkFlagRequired("tool")
	rootCmd.MarkFlagRequired("dir")
}
