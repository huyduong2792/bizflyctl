/*
Copyright © 2020 BizFly Cloud

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
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/bizflycloud/bizflycli/common"
	"github.com/bizflycloud/gobizfly"
)

var (
	snapshotHeaderList = []string{"ID", "Name", "Status", "Size", "Type", "Created At", "Volume ID"}
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "BizFly Cloud Snapshot Interaction",
	Long: `BizFly Cloud Server Action: Create, List, Delete, Snapshot`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snapshot called")
	},
}

// createSnapshotCmd represents the create snapshot command
var createSnapshotCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snapshot",
	Long: `Create a new snapshot
Exmaple: bizfly snapshot create <volume_id>.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

// deleteCmd represents the delete command
var deleteSnapshotCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete snapshots",
	Long: `Delete a snapshot or list of snapshots.
Example: bizfly snapshot delete <snapshot_id> <snapshot_id>`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := apiClientForContext(cmd)
		for _, snapshotID := range args {
			fmt.Printf("Deleting snapshot %s \n", snapshotID)
			err := client.Snapshot.Delete(ctx, snapshotID)
			if err != nil {
				if errors.Is(err, gobizfly.ErrNotFound) {
					fmt.Printf("Snapshot %s is not found", snapshotID)
					return
				}
			}
		}
	},
}

// getSnapshotCmd represents the get command
var getSnapshotCmd = &cobra.Command{
	Use:   "get",
	Short: "Get detail a snapshot",
	Long: `Get detail a snapshot
Example: bizfly snapshot get <snapshot_id>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Printf("Unknow variable %s", strings.Join(args[1:], ""))
		}
		client, ctx := apiClientForContext(cmd)

		snap, err := client.Snapshot.Get(ctx, args[0])
		if err != nil {
			if errors.Is(err, gobizfly.ErrNotFound) {
				fmt.Printf("Snapshot %s not found.", args[0])
				return
			}
			log.Fatal(err)
		}
		var data [][]string
		data = append(data, []string{snap.Id, snap.Name, snap.Status, strconv.Itoa(snap.Size), snap.VolumeTypeId, snap.CreateAt, snap.VolumeId})
		common.Output(snapshotHeaderList, data)
	},
}

var listSnapshotCmd = &cobra.Command{
	Use: "list",
	Short: "List all snapshots in your account",
	Long: `List all snapshots in your account
Example: bizfly snapshot list
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx := apiClientForContext(cmd)
		snapshots, err := client.Snapshot.List(ctx, &gobizfly.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var data [][]string
		for _, snap:= range snapshots {
			data = append(data, []string{
				snap.Id, snap.Name, snap.Status, strconv.Itoa(snap.Size), snap.VolumeTypeId, snap.CreateAt, snap.VolumeId})
		}
		common.Output(snapshotHeaderList, data)
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.AddCommand(createSnapshotCmd)
	snapshotCmd.AddCommand(deleteSnapshotCmd)
	snapshotCmd.AddCommand(getSnapshotCmd)
	snapshotCmd.AddCommand(listSnapshotCmd)
}
