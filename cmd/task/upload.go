package task

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	pkgstorage "github.com/szpp-dev-team/szpp-judge-tool/internal/storage"
	"github.com/szpp-dev-team/szpp-judge-tool/internal/task"
	backendv1 "github.com/szpp-dev-team/szpp-judge/proto-gen/go/backend/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var uploadCmd = &cobra.Command{
	Use: "upload",
	RunE: func(cmd *cobra.Command, args []string) error {
		storagePath := os.Getenv("STORAGE_PATH")
		if storagePath == "" {
			return errors.New("the environment value \"STORAGE_PATH\" must be set")
		}
		storage, err := pkgstorage.LoadOrInit(storagePath)
		if err != nil {
			return err
		}

		taskPath, err := os.Getwd()
		if err != nil {
			return err
		}
		if len(args) > 0 {
			taskPath = args[0]
		}

		controller, err := task.Load(taskPath)
		if err != nil {
			return err
		}
		if err := controller.Validate(); err != nil {
			return err
		}
		checker, err := controller.ReadChecker()
		if err != nil {
			return err
		}

		mutationTask := &backendv1.MutationTask{
			Title:           controller.Config.Title,
			Statement:       controller.Statement,
			ExecTimeLimit:   int32(controller.Config.TimeLimitMs),
			ExecMemoryLimit: int32(controller.Config.MemoryLimitMb),
			Difficulty:      backendv1.Difficulty(backendv1.Difficulty_value[controller.Config.Difficulty]),
			IsPublic:        false, // this field must be changed from only web
			Checker:         checker,
		}

		mutationTestcaseSets := make([]*backendv1.MutationTestcaseSet, 0, len(controller.Config.TestcaseSets))
		for tsName, ts := range controller.Config.TestcaseSets {
			mutationTestcaseSets = append(mutationTestcaseSets, &backendv1.MutationTestcaseSet{
				Slug:          tsName,
				ScoreRatio:    int32(ts.ScoreRatio),
				IsSample:      ts.IsSample,
				TestcaseSlugs: ts.TestcaseSlugs,
			})
		}

		mutationTestcases := make([]*backendv1.MutationTestcase, 0, len(controller.Config.Testcases))
		for _, meta := range controller.Config.Testcases {
			testcase, err := controller.ReadTestcase(meta.Slug)
			if err != nil {
				return err
			}
			mutationTestcases = append(mutationTestcases, &backendv1.MutationTestcase{
				Slug:        meta.Slug,
				Description: &meta.Description,
				Input:       testcase.In,
				Output:      testcase.Out,
			})
		}

		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return err
		}
		conn, err := grpc.Dial(os.Getenv("BACKEND_GRPC_ADDR"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})))
		if err != nil {
			return err
		}
		defer conn.Close()
		c := backendv1.NewTaskServiceClient(conn)

		taskID, ok := storage.GetTaskID(taskPath)
		if ok {
			mutationTask.Id = lo.ToPtr(int32(taskID))
			resp, err := c.UpdateTask(cmd.Context(), &backendv1.UpdateTaskRequest{
				TaskId: int32(taskID),
				Task:   mutationTask,
			})
			if err != nil {
				return err
			}
			cmd.Printf("task was updated(taskID: %d)\n", taskID)
			cmd.Println(resp)
		} else {
			resp, err := c.CreateTask(cmd.Context(), &backendv1.CreateTaskRequest{
				Task: mutationTask,
			})
			if err != nil {
				return err
			}
			cmd.Printf("task was created\n")
			cmd.Println(resp)
			storage.SetTaskID(taskPath, int(resp.Task.Id))
		}

		resp, err := c.SyncTestcaseSets(cmd.Context(), &backendv1.SyncTestcaseSetsRequest{
			TaskId:       int32(taskID),
			TestcaseSets: mutationTestcaseSets,
			Testcases:    mutationTestcases,
		})
		if err != nil {
			return err
		}
		cmd.Println("testcases and testcase_sets were upserted")
		cmd.Println(resp)

		if err := storage.Save(); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		return nil
	},
}
