/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/juparave/genmodels/cmd/package/generate"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "genmodels",
	Short: "A utility to generate gorm models from database schema",
	Long:  `A utility to generate gorm models from database schema. `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		generate.CreateModels(
			cmd.Flag("user").Value.String(),
			cmd.Flag("password").Value.String(),
			cmd.Flag("host").Value.String(),
			cmd.Flag("port").Value.String(),
			cmd.Flag("database").Value.String())
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.genmodels.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// need to get the database name, user, password, host, port, sslmode
	rootCmd.Flags().StringP("database", "d", "", "Database name")
	rootCmd.Flags().StringP("user", "u", "", "Database user")
	rootCmd.Flags().StringP("password", "p", "", "Database password")
	rootCmd.Flags().StringP("host", "H", "localhost", "Database host")
	rootCmd.Flags().StringP("port", "P", "3306", "Database port")
}
