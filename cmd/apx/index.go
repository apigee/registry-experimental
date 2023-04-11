// Code generated. DO NOT EDIT.

package main

import (
	"github.com/spf13/cobra"

	"fmt"

	"github.com/golang/protobuf/jsonpb"

	"os"

	rpcpb "github.com/apigee/registry-experimental/rpc"
)

var IndexInput rpcpb.IndexRequest

var IndexFromFile string

var IndexFollow bool

var IndexPollOperation string

func init() {
	SearchServiceCmd.AddCommand(IndexCmd)

	IndexCmd.Flags().StringVar(&IndexInput.ResourceName, "resource_name", "", "")

	IndexCmd.Flags().StringVar(&IndexFromFile, "from_file", "", "Absolute path to JSON file containing request payload")

	IndexCmd.Flags().BoolVar(&IndexFollow, "follow", false, "Block until the long running operation completes")

	SearchServiceCmd.AddCommand(IndexPollCmd)

	IndexPollCmd.Flags().BoolVar(&IndexFollow, "follow", false, "Block until the long running operation completes")

	IndexPollCmd.Flags().StringVar(&IndexPollOperation, "operation", "", "Required. Operation name to poll for")

	IndexPollCmd.MarkFlagRequired("operation")

}

var IndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Add a resource to the search index.",
	Long:  "Add a resource to the search index.",
	PreRun: func(cmd *cobra.Command, args []string) {

		if IndexFromFile == "" {

		}

	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		in := os.Stdin
		if IndexFromFile != "" {
			in, err = os.Open(IndexFromFile)
			if err != nil {
				return err
			}
			defer in.Close()

			err = jsonpb.Unmarshal(in, &IndexInput)
			if err != nil {
				return err
			}

		}

		if Verbose {
			printVerboseInput("Search", "Index", &IndexInput)
		}
		resp, err := SearchClient.Index(ctx, &IndexInput)
		if err != nil {
			return err
		}

		if !IndexFollow {
			var s interface{}
			s = resp.Name()

			if OutputJSON {
				d := make(map[string]string)
				d["operation"] = resp.Name()
				s = d
			}

			printMessage(s)
			return err
		}

		result, err := resp.Wait(ctx)
		if err != nil {
			return err
		}

		if Verbose {
			fmt.Print("Output: ")
		}
		printMessage(result)

		return err
	},
}

var IndexPollCmd = &cobra.Command{
	Use:   "poll-index",
	Short: "Poll the status of a IndexOperation by name",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		op := SearchClient.IndexOperation(IndexPollOperation)

		if IndexFollow {
			resp, err := op.Wait(ctx)
			if err != nil {
				return err
			}

			if Verbose {
				fmt.Print("Output: ")
			}
			printMessage(resp)
			return err
		}

		resp, err := op.Poll(ctx)
		if err != nil {
			return err
		} else if resp != nil {
			if Verbose {
				fmt.Print("Output: ")
			}

			printMessage(resp)
			return
		}

		if op.Done() {
			fmt.Println(fmt.Sprintf("Operation %s is done", op.Name()))
		} else {
			fmt.Println(fmt.Sprintf("Operation %s not done", op.Name()))
		}

		return err
	},
}
