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
	"io/ioutil"
	"os"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:     "publish",
	Short:   "Publish to Cloud Pub/Sub",
	Long:    "Publish to Cloud Pub/Sub",
	Example: "pubsub publish --project=project --topic=topic --message=message",
	Run: func(cmd *cobra.Command, args []string) {
		topicID, _ := cmd.Flags().GetString("topic")
		message, _ := cmd.Flags().GetString("message")
		filePath, _ := cmd.Flags().GetString("file")
		number, _ := cmd.Flags().GetInt("number")

		if config.ProjectID == "" {
			fmt.Printf("required ProjectID\n")
			os.Exit(1)
			return
		}

		if topicID == "" {
			fmt.Printf("required TopicID\n")
			os.Exit(1)
			return
		}

		if message == "" && filePath == "" {
			fmt.Printf("required message or filePath\n")
			os.Exit(1)
			return
		}

		if filePath != "" {
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			message = string(data)
		}

		ctx := context.Background()
		client, err := pubsub.NewClient(ctx, config.ProjectID)
		if err != nil {
			fmt.Printf("pubsub.NewClient: %v\n", err)
			os.Exit(1)
			return
		}
		defer client.Close()

		t := client.Topic(topicID)
		if number > 0 {
			t.PublishSettings.NumGoroutines = number
		}

		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte(message),
		})
		id, err := result.Get(ctx)
		if err != nil {
			fmt.Printf("Get: %v\n", err)
			os.Exit(1)
			return
		}
		fmt.Printf("Published a message; msg ID: %v\n", id)
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringP("topic", "t", "", "TopicID")
	publishCmd.Flags().StringP("message", "m", "test", "message")
	publishCmd.Flags().StringP("file", "f", "", "Set contents of file to message")
	publishCmd.Flags().IntP("number", "n", 0, "Number of threads")
}
