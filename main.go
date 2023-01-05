package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"github.com/urfave/cli/v2"
	"gitlab.com/go-classroom/todo/util"
	"log"
	"os"
)

type Entry struct {
	Description string
	Category    string
	Status      string
}

func main() {
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}
	db.CreateIndex("categories", "*", buntdb.IndexJSON("category"))
	db.CreateIndex("status", "*", buntdb.IndexJSON("status"))

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
					var new_title, new_description string
					fmt.Println("Enter the title of your entry:")
					util.Scanner(&new_title)
					fmt.Println("Enter a description:")
					util.Scanner(&new_description)
					//Setting the status active as a default for new entries
					db.Update(func(tx *buntdb.Tx) error {
						tx.Set(new_title, fmt.Sprintf(`{"description" : "%s", "status": "Active", "category": "%s"}`, new_description, cCtx.String("c")), nil)
						return nil
					})
					defer db.Close()
					return nil
				},
			},

			// {
			// 	Name:    "delete",
			// 	Aliases: []string{"d"},
			// 	Usage:   "delete an entry from the list",
			// 	Action: func(cCtx *cli.Context) error {
			// 		var title_del string
			// 		fmt.Println("Enter the title of the entry to be delete:")
			// 		util.Scanner(&title_del)
			// 		return nil
			// 	},
			// },

			// // 			{
			// // 				Name:    "edit",
			// // 				Aliases: []string{"e"},
			// // 				Usage:   "edit the title, description and/or category",
			// // 				Flags: []cli.Flag{
			// // 					&cli.StringFlag{
			// // 						Name:    "field",
			// // 						Aliases: []string{"f"},
			// // 						Value:   "all",
			// // 						Usage:   "specify which field to edit",
			// // 					},
			// // 				},
			// // 				Action: func(cCtx *cli.Context) error {
			// // 					fmt.Println("edit")
			// // 					switch cCtx.String("field") {
			// // 					case "title":
			// // 						fmt.Println("Title will be edited")
			// // 					case "description":
			// // 						fmt.Println("Description will be edited")
			// // 					case "category":
			// // 						fmt.Println("Category of entry will be edited")
			// // 					default:
			// // 						fmt.Println("error")
			// // 					}
			// // 					return nil
			// // 				},
			// // 			},

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
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Value:   "all",
						Usage:   "Filter the returning list through the category given",
					},
				},
				Action: func(cCtx *cli.Context) error {
					switch cCtx.String("c") {
					case "all":
						db.View(func(tx *buntdb.Tx) error {
							tx.Ascend("", func(key, value string) bool {
								entry := Entry{}
								if err := json.Unmarshal([]byte(value), &entry); err != nil {
									panic(err)
								}
								fmt.Println(entry)
								return true
							})
							return nil
						})
					case "Fun":
						db.View(func(tx *buntdb.Tx) error {
							tx.Ascend("", func(key, value string) bool {
								entry := Entry{}
								if err := json.Unmarshal([]byte(value), &entry); err != nil {
									panic(err)
								}
								if entry.Category == "Fun" {
									fmt.Println(entry)
								}
								return true
							})
							return nil
						})
					case "Personal":
						db.View(func(tx *buntdb.Tx) error {
							tx.Ascend("", func(key, value string) bool {
								entry := Entry{}
								if err := json.Unmarshal([]byte(value), &entry); err != nil {
									panic(err)
								}
								if entry.Category == "Personal" {
									fmt.Println(entry)
								}
								return true
							})
							return nil
						})
					case "Work":
						db.View(func(tx *buntdb.Tx) error {
							tx.Ascend("", func(key, value string) bool {
								entry := Entry{}
								if err := json.Unmarshal([]byte(value), &entry); err != nil {
									panic(err)
								}
								if entry.Category == "Work" {
									fmt.Println(entry)
								}
								return true
							})
							return nil
						})
					default:
						panic("Please enter a correct category")
					}
					return nil
				},
			},

			// // 			{
			// // 				Name:    "status",
			// // 				Aliases: []string{"s"},
			// // 				Usage:   "change the status of an entry",
			// // 				Action: func(cCtx *cli.Context) error {
			// // 					fmt.Println("status")
			// // 					return nil
			// // 				},
			// // 			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
