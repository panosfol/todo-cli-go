package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"github.com/urfave/cli/v2"
	"gitlab.com/go-classroom/todo/util"
	"log"
	"os"
	"strconv"
)

type Entry struct {
	Id          string
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
				Usage:   "Add a new entry to the list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Value:   "other",
						Usage:   "Enter a category of the entry: fun, work or personal",
					},
				},
				Action: func(cCtx *cli.Context) error {
					//Checking if the user input is correct
					if cCtx.String("c") != "fun" && cCtx.String("c") != "personal" && cCtx.String("c") != "work" && cCtx.String("c") != "other" {
						fmt.Println("Please enter a correct category for the entry")
						os.Exit(0)
					}
					//Using custom made key that auto increments according to the key of the last entry
					new_key := 0
					var titles []string
					db.View(func(tx *buntdb.Tx) error {
						fetched_entry := Entry{}
						tx.Ascend("", func(key, value string) bool {
							if err := json.Unmarshal([]byte(value), &fetched_entry); err != nil {
								panic(err)
							}
							titles = append(titles, fetched_entry.Title)
							new_key, _ = strconv.Atoi(key)
							return true
						})
						return nil
					})
					//If not incremented it will be the same on as the one of the last entry
					new_key++
					var new_title, new_description string
					fmt.Println("Enter the title of your entry:")
					util.Scanner(&new_title)
					//Check if the title provided already exists
					for _, v := range titles {
						if new_title == v {
							fmt.Println("Title already exists")
							os.Exit(0)
						}
					}
					fmt.Println("Enter a description:")
					util.Scanner(&new_description)
					//Setting the status "active" as a default for new entries
					db.Update(func(tx *buntdb.Tx) error {
						tx.Set(fmt.Sprintf("%d", new_key), fmt.Sprintf(
							`{"id": "%s", "title": "%s", "description" : "%s", "status": "active", "category": "%s"}`,
							fmt.Sprintf("%d", new_key), new_title, new_description, cCtx.String("c")), nil)
						return nil
					})
					defer db.Close()
					return nil
				},
			},

			{
				Name:    "delete",
				Aliases: []string{"del"},
				Usage:   "Delete an entry from the list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id",
						Value: "",
						Usage: "Enter the id of the entry to be delete",
					},
				},
				Action: func(cCtx *cli.Context) error {
					//If the id flag is empty the process should terminate
					if cCtx.String("id") == "" {
						fmt.Println("Please enter the title of the entry to be delete")
						os.Exit(0)
					}
					//Fetching all the entries to match the id given with the correct entry
					entries := []Entry{}
					db.View(func(tx *buntdb.Tx) error {
						fetched_entry := Entry{}
						tx.Ascend("", func(key, value string) bool {
							if err := json.Unmarshal([]byte(value), &fetched_entry); err != nil {
								panic(err)
							}
							entries = append(entries, fetched_entry)
							return true
						})
						return nil
					})
					var del_key string
					for _, v := range entries {
						if cCtx.String("id") == v.Id {
							del_key = v.Id
						}
					}
					db.Update(func(tx *buntdb.Tx) error {
						_, err := tx.Delete(del_key)
						if err != nil {
							fmt.Printf("Entry with id \"%s\" was not found\n", cCtx.String("id"))
						} else {
							fmt.Printf("Entry with id \"%s\" has been succesfully deleted\n", del_key)
						}
						return nil
					})
					defer db.Close()
					return nil
				},
			},

			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "Edit the title, description and/or category",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id",
						Value: "",
						Usage: "Id of the entry to be deleted",
					},
					&cli.StringFlag{
						Name:    "field",
						Aliases: []string{"f"},
						Value:   "all",
						Usage:   "Specify which field to edit: title, description or category",
					},
				},
				Action: func(cCtx *cli.Context) error {
					//If the id flag is empty the process should terminate
					if cCtx.String("id") == "" {
						fmt.Println("Please enter the id of the entry to be delete")
						os.Exit(0)
					}
					//If the field flag was not correctly given the process should terminate
					if cCtx.String("f") != "title" && cCtx.String("f") != "description" && cCtx.String("f") != "category" && cCtx.String("f") != "all" {
						fmt.Println("Please enter a correct field to edit")
						os.Exit(0)
					}
					//Fetching all the entries to match the id given with the entry to be edited
					entries := []Entry{}
					db.View(func(tx *buntdb.Tx) error {
						fetched_entry := Entry{}
						tx.Ascend("", func(key, value string) bool {
							if err := json.Unmarshal([]byte(value), &fetched_entry); err != nil {
								panic(err)
							}
							entries = append(entries, fetched_entry)
							return true
						})
						return nil
					})
					edit_entry := Entry{}
					id_found := false
					for _, v := range entries {
						if cCtx.String("id") == v.Id {
							edit_entry.Title = v.Title
							edit_entry.Description = v.Description
							edit_entry.Category = v.Category
							edit_entry.Status = v.Status
							id_found = true
						}
					}
					if id_found == false {
						fmt.Println("The id given doesn't exist")
						os.Exit(0)
					}
					switch cCtx.String("f") {
					case "title":
						var new_title string
						fmt.Println("Enter the new title: ")
						util.Scanner(&new_title)
						edit_entry.Title = new_title
					case "description":
						var new_desc string
						fmt.Println("Enter the new description: ")
						util.Scanner(&new_desc)
						edit_entry.Description = new_desc
					case "category":
						var new_cat string
						fmt.Println("Enter the new category: ")
						util.Scanner(&new_cat)
						//Using an infinite loop to ensure the user input is correct
						for {
							if new_cat != "work" && new_cat != "fun" && new_cat != "personal" {
								fmt.Println("Please enter the correct category(use only lowercase letters): ")
								util.Scanner(&new_cat)
							} else {
								edit_entry.Category = new_cat
								break
							}
						}
					case "all":
						var new_title, new_desc, new_cat string
						fmt.Println("Enter the new title: ")
						util.Scanner(&new_title)
						edit_entry.Title = new_title
						fmt.Println("Enter the new description: ")
						util.Scanner(&new_desc)
						edit_entry.Description = new_desc
						fmt.Println("Enter the new category: ")
						util.Scanner(&new_cat)
						for {
							if new_cat != "work" && new_cat != "fun" && new_cat != "personal" {
								fmt.Println("Please enter the correct category: ")
								util.Scanner(&new_cat)
							} else {
								edit_entry.Category = new_cat
								break
							}
						}
					default:
						fmt.Println("Please enter a correct field to edit")
					}
					db.Update(func(tx *buntdb.Tx) error {
						tx.Set(cCtx.String("id"), fmt.Sprintf(
							`{"id": "%s", "title": "%s", "description" : "%s", "status": "%s", "category": "%s"}`,
							cCtx.String("id"), edit_entry.Title, edit_entry.Description, edit_entry.Status, edit_entry.Category), nil)
						return nil
					})
					defer db.Close()
					return nil
				},
			},

			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List all the entries",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "status",
						Aliases: []string{"s"},
						Value:   "all",
						Usage:   "Filter the returning list through the status given: abandoned, active, done",
					},
					&cli.StringFlag{
						Name:    "category",
						Aliases: []string{"c"},
						Value:   "all",
						Usage:   "Filter the returning list through the category given: fun, work, personal",
					},
					&cli.StringFlag{
						Name:  "id",
						Value: "all",
						Usage: "Return information about a specific entry",
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
							entries = append(entries, fetched_entry)
							return true
						})
						return nil
					})
					//Checking if the user used a title flag
					if cCtx.String("id") == "all" {
						//Using 2 new variables to potentially filter twice through the list of entries, once for category flag and once for status flag
						entries1 := []Entry{}
						entries2 := []Entry{}
						//Checking if the user put the correct input
						if cCtx.String("c") != "fun" && cCtx.String("c") != "work" && cCtx.String("c") != "all" && cCtx.String("c") != "personal" {
							panic("Wrong category parameter given")
						}
						if cCtx.String("s") != "active" && cCtx.String("s") != "done" && cCtx.String("s") != "all" && cCtx.String("s") != "abandoned" {
							panic("Wrong status parameter given")
						}
						//Filtering through all the entries using the category flag
						if cCtx.String("c") == "fun" || cCtx.String("c") == "personal" || cCtx.String("c") == "work" {
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
						if cCtx.String("s") == "abandoned" || cCtx.String("s") == "done" || cCtx.String("s") == "active" {
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
							fmt.Printf("%s: %s\n", v.Id, v.Title)
						}

					} else {
						id_found := false
						for _, v := range entries {
							if cCtx.String("id") == v.Id {
								id_found = true
								fmt.Printf("Title: %s\nDescription: %s\nCategory: %s\nStatus: %s\n", v.Title, v.Description, v.Category, v.Status)
								break
							}
						}
						if id_found == false {
							fmt.Printf("Entry with id \"%s\" was not found\n", cCtx.String("id"))
						}
					}
					return nil
				},
			},

			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "Update the status of an entry",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id",
						Value: "",
						Usage: "Enter the id of the entry",
					},
					&cli.StringFlag{
						Name:    "status",
						Aliases: []string{"s"},
						Value:   "",
						Usage:   "Enter the new status of the entry: done, active, abandoned",
					},
				},
				Action: func(cCtx *cli.Context) error {
					//If the id is empty the process should terminate
					if cCtx.String("id") == "" {
						fmt.Println("Please enter the id of the entry")
						os.Exit(0)
					}
					//If the status given is empty or not correct the process should terminate
					if cCtx.String("s") != "done" && cCtx.String("s") != "abandoned" && cCtx.String("s") != "active" {
						fmt.Println("Please enter a correct status")
						os.Exit(0)
					}
					//Fetching all the entries to match the id given with the entry to be updated
					entries := []Entry{}
					db.View(func(tx *buntdb.Tx) error {
						fetched_entry := Entry{}
						tx.Ascend("", func(key, value string) bool {
							if err := json.Unmarshal([]byte(value), &fetched_entry); err != nil {
								panic(err)
							}
							entries = append(entries, fetched_entry)
							return true
						})
						return nil
					})
					id_found := false
					var entry_title, entry_desc, entry_cat string
					for _, v := range entries {
						if cCtx.String("id") == v.Id {
							entry_title = v.Title
							entry_desc = v.Description
							entry_cat = v.Category
							id_found = true
							break
						}
					}
					//If the id given was not found the process should terminate
					if id_found == false {
						fmt.Printf("Entry with id \"%s\" was not found\n", cCtx.String("id"))
						os.Exit(0)
					}
					db.Update(func(tx *buntdb.Tx) error {
						tx.Set(cCtx.String("id"), fmt.Sprintf(
							`{"id": "%s", "title": "%s", "description" : "%s", "status": "%s", "category": "%s"}`,
							cCtx.String("id"), entry_title, entry_desc, cCtx.String("s"), entry_cat), nil)
						return nil
					})
					defer db.Close()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
