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
	Title       string
	Description string
	Category    string
	Status      string
}

func main() {
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}

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
						Usage:   "Enter a category of the entry: Fun, Work or Personal",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.String("c") != "Fun" && cCtx.String("c") != "Personal" && cCtx.String("c") != "Work" && cCtx.String("c") != "Other" {
						panic("Please enter a correct category for the entry")
					}
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
					&cli.StringFlag{
						Name:    "title",
						Aliases: []string{"t"},
						Value:   "all",
						Usage:   "Return information about a specific title",
					},
				},
				Action: func(cCtx *cli.Context) error {
					entries := []Entry{}
					db.View(func(tx *buntdb.Tx) error {
						fetched_entry := Entry{}
						tx.Ascend("", func(key, value string) bool {
							if err := json.Unmarshal([]byte(value), &fetched_entry); err != nil {
								panic(err)
							}
							fetched_entry.Title = key
							entries = append(entries, fetched_entry)
							return true
						})
						return nil
					})
					//Checking if the user used a title flag
					if cCtx.String("t") == "all" {
						//Using 2 new variables to potentially filter twice through the list of entries, once for category flag and once for status flag
						entries1 := []Entry{}
						entries2 := []Entry{}
						//Checking if the user put the correct input
						if cCtx.String("c") != "Fun" && cCtx.String("c") != "Work" && cCtx.String("c") != "all" && cCtx.String("c") != "Personal" {
							panic("Wrong category parameter given")
						}
						if cCtx.String("s") != "Active" && cCtx.String("s") != "Done" && cCtx.String("s") != "all" && cCtx.String("s") != "Abandoned" {
							panic("Wrong status parameter given")
						}
						//Filtering through all the entries using the category flag
						if cCtx.String("c") == "Fun" || cCtx.String("c") == "Personal" || cCtx.String("c") == "Work" {
							for _, v := range entries {
								if v.Category == cCtx.String("c") {
									entries1 = append(entries1, v)
								}
							}
						} else {
							//If no category flag was provided then the list remains the same
							entries1 = entries

						}
						//Filtering through all the entries using the status flag
						if cCtx.String("s") == "Abandoned" || cCtx.String("s") == "Done" || cCtx.String("s") == "Active" {
							for _, v := range entries1 {
								if v.Status == cCtx.String("s") {
									entries2 = append(entries2, v)
								}
							}
						} else {
							entries2 = entries1
						}
						fmt.Printf("Here is a list of your entries: \n")
						for _, v := range entries2 {
							fmt.Println(v.Title)
						}

					} else {
						title_found := false
						for _, v := range entries {
							if cCtx.String("t") == v.Title {
								title_found = true
								fmt.Printf("Title: %s\nDescription: %s\nCategory: %s\nStatus: %s\n", v.Title, v.Description, v.Category, v.Status)
								break
							}
						}
						if title_found == false {
							fmt.Println("Title was not found")
						}
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
