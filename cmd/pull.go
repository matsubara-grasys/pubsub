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
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"sync/atomic"
	"time"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:     "pull",
	Short:   "pull subscription message",
	Long:    "pull subscription message",
	Example: fmt.Sprintln(" ", rootCmd.Use, "pull", "--project=project --sub=sub"),
	Run: func(cmd *cobra.Command, args []string) {
		subID, _ := cmd.Flags().GetString("sub")
		timeout, _ := cmd.Flags().GetInt("timeout")
		synchronous, _ := cmd.Flags().GetBool("sub")
		maxOutstandingMessages, _ := cmd.Flags().GetInt("MaxOutstandingMessages")
		number, _ := cmd.Flags().GetInt("number")

		if viper.GetString("ProjectID") == "" {
			fmt.Printf("required ProjectID\n")
			os.Exit(1)
			return
		}

		if subID == "" {
			fmt.Printf("required SubID\n")
			os.Exit(1)
			return
		}

		var ctx context.Context
		var cancel context.CancelFunc

		ctx = context.Background()
		client, err := pubsub.NewClient(ctx, viper.GetString("ProjectID"))
		if err != nil {
			fmt.Printf("pubsub.NewClient: %v\n", err)
			os.Exit(1)
			return
		}
		defer client.Close()

		sub := client.Subscription(subID)

		sub.ReceiveSettings.Synchronous = synchronous
		if maxOutstandingMessages > 0 {
			sub.ReceiveSettings.MaxOutstandingMessages = maxOutstandingMessages
		}

		if number > 0 {
			sub.ReceiveSettings.NumGoroutines = number
		}

		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
			defer cancel()
		}

		var received int32
		err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
			fmt.Printf("Got message: %q\n", string(msg.Data))
			atomic.AddInt32(&received, 1)
			msg.Ack()
		})
		if err != nil {
			fmt.Printf("sub.Receive: %v\n", err)
			os.Exit(1)
			return
		}
		fmt.Printf("Received %d messages\n", received)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringP("sub", "s", "", "SubID")
	pullCmd.Flags().Int("timeout", 0, "timeout")
	pullCmd.Flags().Bool("synchronous", false, "Use synchronous pull")
	pullCmd.Flags().Int("max-outstanding-messages", 0, "MaxOutstandingMessages")
	pullCmd.Flags().IntP("number", "n", 0, "Number of threads")
}
