package cmd

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/micro-service-lab/recs-seem-mono-container/app/batch"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
)

var (
	force     bool
	diff      bool
	noDelete  bool
	deepEqual bool
)

type getCount func(ctx context.Context) (int64, error)

func seedCmdGenerator(
	name, logName, capitalName string,
	getCount getCount,
	b batch.Batch,
) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Inserts %s", logName),
		Long: fmt.Sprintf(`Inserting %s will insert the data necessary for the application to operate into the database.

The seed command is executed when the application is started for the first time.`, logName),
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			color.HiCyan(fmt.Sprintf("seed %s called...", logName))
			s := spinner.New(spinner.CharSets[11], spinnerFrequency*time.Millisecond)
			s.Start()

			ctx := cmd.Context()
			if !force && !diff {
				count, err := getCount(ctx)
				if err != nil {
					s.Stop()
					color.Red(fmt.Errorf("failed to get %s count: %w", logName, err).Error())
					return
				}
				if count > 0 {
					s.Stop()
					color.Yellow(fmt.Sprintf("%s already exist. Use --force to seed again", capitalName))
					return
				}
			}
			if diff {
				err := b.RunDiff(ctx, noDelete, deepEqual)
				if err != nil {
					s.Stop()
					color.Red(fmt.Errorf("failed to insert %s: %w", logName, err).Error())
					return
				}
				s.Stop()
				color.Green(fmt.Sprintf("Completed filling in the difference on %s", logName))
				return
			}
			err := b.Run(ctx)
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to insert %s: %w", logName, err).Error())
				return
			}
			s.Stop()
			color.Green(fmt.Sprintf("%s inserted successfully", capitalName))
		},
	}
}

// ValidTarget is a list of valid seed targets.
var ValidTarget = []string{"attend_status"}

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Inserts initial data into the database",
	//nolint: lll
	Long: `Inserting initial data into the database will insert the data necessary for the application to operate into the database.

The seed command is executed when the application is started for the first time.`,
}

// seedAttendStatusesCmd inserts attend statuses.
var seedAttendStatusesCmd *cobra.Command

// seedAttendanceTypesCmd inserts attendance types.
var seedAttendanceTypesCmd *cobra.Command

// seedEventTypesCmd inserts event types.
var seedEventTypesCmd *cobra.Command

// seedPermissionCategoriesCmd inserts permission categories.
var seedPermissionCategoriesCmd *cobra.Command

// seedPolicyCategoriesCmd inserts policy categories.
var seedPolicyCategoriesCmd *cobra.Command

// seedMimeTypesCmd inserts mime types.
var seedMimeTypesCmd *cobra.Command

// seedRecordTypesCmd inserts record types.
var seedRecordTypesCmd *cobra.Command

// seedPermissionsCmd inserts permissions.
var seedPermissionsCmd *cobra.Command

// seedPoliciesCmd inserts policies.
var seedPoliciesCmd *cobra.Command

// seedDefaultImagesCmd inserts default images.
var seedDefaultImagesCmd *cobra.Command

// seedGradesCmd inserts grades.
var seedGradesCmd *cobra.Command

// seedGroupsCmd inserts groups.
var seedGroupsCmd *cobra.Command

// seedWholeOrganizationCmd inserts whole organization.
var seedWholeOrganizationCmd *cobra.Command

// seedChatRoomActionTypesCmd inserts chat room action types.
var seedChatRoomActionTypesCmd *cobra.Command

// seedAllCmd inserts all seed data.
var seedAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Inserts all seed data",
	Long: `Inserting all seed data will insert the data necessary for the application to operate into the database.

The seed command is executed when the application is started for the first time.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		color.HiCyan("seed all called...")
		cmds := []func(cmd *cobra.Command, args []string){
			seedAttendStatusesCmd.Run,
			seedAttendanceTypesCmd.Run,
			seedEventTypesCmd.Run,
			func(cmd *cobra.Command, args []string) {
				seedPermissionCategoriesCmd.Run(cmd, args)
				seedPermissionsCmd.Run(cmd, args)
			},
			func(cmd *cobra.Command, args []string) {
				seedPolicyCategoriesCmd.Run(cmd, args)
				seedPoliciesCmd.Run(cmd, args)
			},
			func(cmd *cobra.Command, args []string) {
				seedMimeTypesCmd.Run(cmd, args)
				seedDefaultImagesCmd.Run(cmd, args)
				seedGradesCmd.Run(cmd, args)
				seedGroupsCmd.Run(cmd, args)
				seedWholeOrganizationCmd.Run(cmd, args)
			},
			seedRecordTypesCmd.Run,
			seedChatRoomActionTypesCmd.Run,
		}
		var wg sync.WaitGroup
		wg.Add(len(cmds))
		for _, c := range cmds {
			go func(c func(cmd *cobra.Command, args []string)) {
				defer wg.Done()
				c(cmd, args)
			}(c)
		}
		wg.Wait()
		color.Green("All seed data inserted successfully")
	},
}

// SeedInit initializes seed commands.
//
//nolint:wrapcheck
func seedInit() {
	seedAttendStatusesCmd = seedCmdGenerator(
		"attend_status",
		"attend statuses",
		"Attend Statuses",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetAttendStatusesCount(ctx, "")
		},
		&batch.InitAttendStatuses{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedAttendanceTypesCmd = seedCmdGenerator(
		"attendance_type",
		"attendance types",
		"Attendance Types",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetAttendanceTypesCount(ctx, "")
		},
		&batch.InitAttendanceTypes{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedEventTypesCmd = seedCmdGenerator(
		"event_type",
		"event types",
		"Event Types",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetEventTypesCount(ctx, "")
		},
		&batch.InitEventTypes{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedPermissionCategoriesCmd = seedCmdGenerator(
		"permission_category",
		"permission categories",
		"Permission Categories",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetPermissionCategoriesCount(ctx, "")
		},
		&batch.InitPermissionCategories{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedPolicyCategoriesCmd = seedCmdGenerator(
		"policy_category",
		"policy categories",
		"Policy Categories",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetPolicyCategoriesCount(ctx, "")
		},
		&batch.InitPolicyCategories{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedMimeTypesCmd = seedCmdGenerator(
		"mime_type",
		"mime types",
		"Mime Types",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetMimeTypesCount(ctx, "")
		},
		&batch.InitMimeTypes{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedRecordTypesCmd = seedCmdGenerator(
		"record_type",
		"record types",
		"Record Types",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetRecordTypesCount(ctx, "")
		},
		&batch.InitRecordTypes{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedPermissionsCmd = seedCmdGenerator(
		"permission",
		"permissions",
		"Permissions",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetPermissionsCount(ctx, "", []uuid.UUID{})
		},
		&batch.InitPermissions{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedPoliciesCmd = seedCmdGenerator(
		"policy",
		"policies",
		"Policies",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetPoliciesCount(ctx, "", []uuid.UUID{})
		},
		&batch.InitPolicies{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedDefaultImagesCmd = seedCmdGenerator(
		"default_image",
		"default images",
		"Default Images",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetImagesCount(ctx)
		},
		&batch.InitDefaultImages{
			Manager: &AppContainer.ServiceManager,
		},
	)
	seedGradesCmd = seedCmdGenerator(
		"grade",
		"grades",
		"Grades",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetGradesCount(ctx)
		},
		&batch.InitGrades{
			Manager: &AppContainer.ServiceManager,
			Storage: AppContainer.Storage,
		},
	)
	seedGroupsCmd = seedCmdGenerator(
		"group",
		"groups",
		"Groups",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetGroupsCount(ctx)
		},
		&batch.InitGroups{
			Manager: &AppContainer.ServiceManager,
			Storage: AppContainer.Storage,
		},
	)
	seedWholeOrganizationCmd = seedCmdGenerator(
		"whole_organization",
		"whole organization",
		"Whole Organization",
		func(ctx context.Context) (int64, error) {
			_, err := AppContainer.ServiceManager.FindWholeOrganization(ctx)
			if err != nil {
				var e errhandle.ModelNotFoundError
				if !errors.As(err, &e) {
					return 0, err
				}
				return 0, nil
			}
			return 1, nil
		},
		&batch.InitWholeOrganization{
			Manager: &AppContainer.ServiceManager,
			Storage: AppContainer.Storage,
		},
	)
	seedChatRoomActionTypesCmd = seedCmdGenerator(
		"chat_room_action_type",
		"chat room action types",
		"Chat Room Action Types",
		func(ctx context.Context) (int64, error) {
			return AppContainer.ServiceManager.GetChatRoomActionTypesCount(ctx, "")
		},
		&batch.InitChatRoomActionTypes{
			Manager: &AppContainer.ServiceManager,
		},
	)

	rootCmd.AddCommand(seedCmd)
	seedCmd.AddCommand(seedAllCmd)
	seedCmd.AddCommand(seedAttendStatusesCmd)
	seedCmd.AddCommand(seedAttendanceTypesCmd)
	seedCmd.AddCommand(seedEventTypesCmd)
	seedCmd.AddCommand(seedPermissionCategoriesCmd)
	seedCmd.AddCommand(seedPolicyCategoriesCmd)
	seedCmd.AddCommand(seedMimeTypesCmd)
	seedCmd.AddCommand(seedRecordTypesCmd)
	seedCmd.AddCommand(seedPermissionsCmd)
	seedCmd.AddCommand(seedPoliciesCmd)
	seedCmd.AddCommand(seedDefaultImagesCmd)
	seedCmd.AddCommand(seedGradesCmd)
	seedCmd.AddCommand(seedGroupsCmd)
	seedCmd.AddCommand(seedWholeOrganizationCmd)
	seedCmd.AddCommand(seedChatRoomActionTypesCmd)

	seedCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force seed")
	seedCmd.PersistentFlags().BoolVarP(&diff, "diff", "d", false, "Seed only if there is a difference")
	seedCmd.PersistentFlags().BoolVarP(&noDelete, "no-delete", "n", false, "Do not delete, only insert.This option is valid only when --diff is specified")  //nolint: lll
	seedCmd.PersistentFlags().BoolVarP(&deepEqual, "deep-equal", "e", false, "Use deep equal comparison.This option is valid only when --diff is specified") //nolint: lll
}
