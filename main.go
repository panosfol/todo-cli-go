package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"n"},
				Usage:   "add a new entry to the list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Value:   "Other",
						Usage:   "enter a category of the entry",
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("new")
					switch cCtx.String("category") {
					case "Fun" :
						fmt.Println("gg fun")
					case "fun" :
						fmt.Println("gg fun")
					case "Personal":
						fmt.Println("gg per")
					case "personal" :
						fmt.Println("gg per")
					case "Work" :
						fmt.Println("gg work")
					case "work" :
						fmt.Println("gg work")
					default :
						fmt.Println("error")
					}
					return nil
				},
			},

			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "delete an entry from the list",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("delete")
					return nil
				},
			},

			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "edit the title, description and/or category",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "field",
						Aliases: []string{"f"},
						Value:   "all",
						Usage:   "specify which field to edit",
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("edit")
					switch cCtx.String("field") {
					case "title" :
						fmt.Println("Title will be edited")
					case "description" :
						fmt.Println("Description will be edited")
					case "category" :
						fmt.Println("Category of entry will be edited")
					default :
						fmt.Println("error")
					}
					return nil
				},
			},

			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list all the entries",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "status",
						Aliases: []string{"s"},
 						Value:   "all",
						Usage:   "Filter the returning list through the status given",
					},
					&cli.StringFlag {
						Name:    "category",
						Aliases: []string{"c"},
						Value:   "all",
						Usage:   "Filter the returning list through the category given",
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("list")
					switch cCtx.String("status") {
					case "Abandoned" :
						fmt.Println("abandoned")
					case "Active" :
						fmt.Println("active")
					case "Done" :
						fmt.Println("done")
					default :
						fmt.Println("fetching all")
					}
					switch cCtx.String("category") {
					case "Fun" :
						fmt.Println("fun")
					case "Personal" :
						fmt.Println("personal")
					case "Work" :
						fmt.Println("work")
						default :
						fmt.Println("fetching all")
					}
					return nil
				},
			},

			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "change the status of an entry",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("status")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
