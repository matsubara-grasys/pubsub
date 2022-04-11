/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// createConfigCmd represents the createConfig command
var createConfigCmd = &cobra.Command{
	Use:     "createConfig",
	Short:   "Write config file",
	Long:    "Write config file",
	Example: fmt.Sprintln(" ", rootCmd.Use, "createConfig", "--project=project"),
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.SafeWriteConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createConfigCmd)
}
