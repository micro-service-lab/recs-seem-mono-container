package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/micro-service-lab/recs-seem-mono-container/app/batch"
)

var (
	force     bool
	diff      bool
	noDelete  bool
	deepEqual bool
)

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
var seedAttendStatusesCmd = &cobra.Command{
	Use:   "attend_status",
	Short: "Inserts attend statuses",
	Long: `Inserting attend statuses will insert the data necessary for the application to operate into the database.

The seed command is executed when the application is started for the first time.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		color.HiCyan("seed attend statuses called...")
		s := spinner.New(spinner.CharSets[11], spinnerFrequency*time.Millisecond)
		s.Start()

		ctx := cmd.Context()
		if !force && !diff {
			count, err := AppContainer.ServiceManager.GetAttendStatusesCount(ctx, "")
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to get attend statuses count: %w", err).Error())
				return
			}
			if count > 0 {
				s.Stop()
				color.Yellow("Attend statuses already exist. Use --force to seed again")
				return
			}
		}
		b := batch.InitAttendStatuses{
			Manager: &AppContainer.ServiceManager,
		}
		if diff {
			err := b.RunDiff(ctx, noDelete, deepEqual)
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to insert attend statuses: %w", err).Error())
				return
			}
			s.Stop()
			color.Green("Completed filling in the difference on attend statuses")
			return
		}
		err := b.Run(ctx)
		if err != nil {
			s.Stop()
			color.Red(fmt.Errorf("failed to insert attend statuses: %w", err).Error())
			return
		}
		s.Stop()
		color.Green("Attend statuses inserted successfully")
	},
}

// seedAttendanceTypesCmd inserts attendance types.
var seedAttendanceTypesCmd = &cobra.Command{
	Use:   "attendance_type",
	Short: "Inserts attendance types",
	Long: `Inserting attendance types will insert the data necessary for the application to operate into the database.

The seed command is executed when the application is started for the first time.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		color.HiCyan("seed attendance types called...")
		s := spinner.New(spinner.CharSets[11], spinnerFrequency*time.Millisecond)
		s.Start()

		ctx := cmd.Context()
		if !force && !diff {
			count, err := AppContainer.ServiceManager.GetAttendanceTypesCount(ctx, "")
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to get attendance types count: %w", err).Error())
				return
			}
			if count > 0 {
				s.Stop()
				color.Yellow("Attendance types already exist. Use --force to seed again")
				return
			}
		}
		b := batch.InitAttendanceTypes{
			Manager: &AppContainer.ServiceManager,
		}
		if diff {
			err := b.RunDiff(ctx, noDelete, deepEqual)
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to insert attendance types: %w", err).Error())
				return
			}
			s.Stop()
			color.Green("Completed filling in the difference on attendance types")
			return
		}
		err := b.Run(ctx)
		if err != nil {
			s.Stop()
			color.Red(fmt.Errorf("failed to insert attendance types: %w", err).Error())
			return
		}
		s.Stop()
		color.Green("Attendance types inserted successfully")
	},
}

// seedEventTypesCmd inserts event types.
var seedEventTypesCmd = &cobra.Command{
	Use:   "event_type",
	Short: "Inserts event types",
	Long: `Inserting event types will insert the data necessary for the application to operate into the database.

The seed command is executed when the application is started for the first time.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		color.HiCyan("seed event types called...")
		s := spinner.New(spinner.CharSets[11], spinnerFrequency*time.Millisecond)
		s.Start()

		ctx := cmd.Context()
		if !force && !diff {
			count, err := AppContainer.ServiceManager.GetEventTypesCount(ctx, "")
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to get event types count: %w", err).Error())
				return
			}
			if count > 0 {
				s.Stop()
				color.Yellow("Event types already exist. Use --force to seed again")
				return
			}
		}
		b := batch.InitEventTypes{
			Manager: &AppContainer.ServiceManager,
		}
		if diff {
			err := b.RunDiff(ctx, noDelete, deepEqual)
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to insert event types: %w", err).Error())
				return
			}
			s.Stop()
			color.Green("Completed filling in the difference on event types")
			return
		}
		err := b.Run(ctx)
		if err != nil {
			s.Stop()
			color.Red(fmt.Errorf("failed to insert event types: %w", err).Error())
			return
		}
		s.Stop()
		color.Green("Event types inserted successfully")
	},
}

// seedPermissionCategoriesCmd inserts permission categories.
var seedPermissionCategoriesCmd = &cobra.Command{
	Use:   "permission_category",
	Short: "Inserts permission categories",
	Long: `Inserting permission categories will insert the data necessary for the application to operate into the database.

The seed command is executed when the application is started for the first time.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		color.HiCyan("seed permission categories called...")
		s := spinner.New(spinner.CharSets[11], spinnerFrequency*time.Millisecond)
		s.Start()

		ctx := cmd.Context()
		if !force && !diff {
			count, err := AppContainer.ServiceManager.GetPermissionCategoriesCount(ctx, "")
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to get permission categories count: %w", err).Error())
				return
			}
			if count > 0 {
				s.Stop()
				color.Yellow("Permission Categories already exist. Use --force to seed again")
				return
			}
		}
		b := batch.InitPermissionCategories{
			Manager: &AppContainer.ServiceManager,
		}
		if diff {
			err := b.RunDiff(ctx, noDelete, deepEqual)
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to insert permission categories: %w", err).Error())
				return
			}
			s.Stop()
			color.Green("Completed filling in the difference on permission categories")
			return
		}
		err := b.Run(ctx)
		if err != nil {
			s.Stop()
			color.Red(fmt.Errorf("failed to insert permission categories: %w", err).Error())
			return
		}
		s.Stop()
		color.Green("Permission Categories inserted successfully")
	},
}

// seedPolicyCategoriesCmd inserts policy categories.
var seedPolicyCategoriesCmd = &cobra.Command{
	Use:   "policy_category",
	Short: "Inserts policy categories",
	Long: `Inserting policy categories will insert the data necessary for the application to operate into the database.

The seed command is executed when the application is started for the first time.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		color.HiCyan("seed policy categories called...")
		s := spinner.New(spinner.CharSets[11], spinnerFrequency*time.Millisecond)
		s.Start()

		ctx := cmd.Context()
		if !force && !diff {
			count, err := AppContainer.ServiceManager.GetPolicyCategoriesCount(ctx, "")
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to get policy categories count: %w", err).Error())
				return
			}
			if count > 0 {
				s.Stop()
				color.Yellow("Policy Categories already exist. Use --force to seed again")
				return
			}
		}
		b := batch.InitPolicyCategories{
			Manager: &AppContainer.ServiceManager,
		}
		if diff {
			err := b.RunDiff(ctx, noDelete, deepEqual)
			if err != nil {
				s.Stop()
				color.Red(fmt.Errorf("failed to insert policy categories: %w", err).Error())
				return
			}
			s.Stop()
			color.Green("Completed filling in the difference on policy categories")
			return
		}
		err := b.Run(ctx)
		if err != nil {
			s.Stop()
			color.Red(fmt.Errorf("failed to insert policy categories: %w", err).Error())
			return
		}
		s.Stop()
		color.Green("Policy Categories inserted successfully")
	},
}

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
			seedPermissionCategoriesCmd.Run,
			seedPolicyCategoriesCmd.Run,
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

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.AddCommand(seedAllCmd)
	seedCmd.AddCommand(seedAttendStatusesCmd)
	seedCmd.AddCommand(seedAttendanceTypesCmd)
	seedCmd.AddCommand(seedEventTypesCmd)
	seedCmd.AddCommand(seedPermissionCategoriesCmd)
	seedCmd.AddCommand(seedPolicyCategoriesCmd)

	seedCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force seed")
	seedCmd.PersistentFlags().BoolVarP(&diff, "diff", "d", false, "Seed only if there is a difference")
	seedCmd.PersistentFlags().BoolVarP(&noDelete, "no-delete", "n", false, "Do not delete, only insert.This option is valid only when --diff is specified")  //nolint: lll
	seedCmd.PersistentFlags().BoolVarP(&deepEqual, "deep-equal", "e", false, "Use deep equal comparison.This option is valid only when --diff is specified") //nolint: lll
}
