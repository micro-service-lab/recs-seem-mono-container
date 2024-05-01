package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/micro-service-lab/recs-seem-mono-container/app/batch"
)

var force bool

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
		if !force {
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

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.AddCommand(seedAttendStatusesCmd)

	seedCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force seed")
}
