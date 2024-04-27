package cmd

import (
	"Edot/cmd/seeds"
	"Edot/packages/logger"
	"Edot/packages/postgres"

	"github.com/spf13/cobra"
)

var seed = &cobra.Command{
	Use:   "seed",
	Short: "Seeder",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger()
		postgres := postgres.NewPostgres(log)
		option := ""

		if len(args) == 2 && (args[1] != "fresh") {
			log.Error("please insert argument seed name")
			return
		} else if len(args) == 2 && (args[1] != "") {
			option = "fresh"
		}

		switch args[0] {
		case "user":
			if err := seeds.User(postgres, option); err != nil {
				log.Error(err)
				return
			}
			return
		case "product":
			if err := seeds.Product(postgres, option); err != nil {
				log.Error(err)
				return
			}
		case "warehouse":
			if err := seeds.Warehouse(postgres, option); err != nil {
				log.Error(err)
				return
			}
		case "stock":
			if err := seeds.Stock(postgres, option); err != nil {
				log.Error(err)
				return
			}
		case "all":
			if err := seeds.User(postgres, option); err != nil {
				log.Error(err)
				return
			}
			if err := seeds.Warehouse(postgres, option); err != nil {
				log.Error(err)
				return
			}
			if err := seeds.Product(postgres, option); err != nil {
				log.Error(err)
				return
			}
			if err := seeds.Stock(postgres, option); err != nil {
				log.Error(err)
				return
			}
		default:
			log.Error("argument not found")
		}
	},
}

func init() {
	rootCmd.AddCommand(seed)
}
