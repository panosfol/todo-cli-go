package main

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
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
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"entry": &memdb.TableSchema{
				Name: "entry",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Title"},
					},
					"description": &memdb.IndexSchema{
						Name:    "description",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Description"},
					},
					"status": &memdb.IndexSchema{
						Name:    "status",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Status"},
					},
					"category": &memdb.IndexSchema{
						Name:    "category",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Category"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
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
						Usage:   "enter a category of the entry",
					},
				},
				Action: func(cCtx *cli.Context) error {
					//Create a write transaction with DB
					txn := db.Txn(true)
					var new_entry Entry
					fmt.Println("Enter the title of your entry:")
					util.Scanner(&new_entry.Title)
					fmt.Println("Enter a description:")
					util.Scanner(&new_entry.Description)
					//Setting the status active as a default for new entries
					new_entry.Status = "Active"
					new_entry.Category = cCtx.String("c")
					if err := txn.Insert("entry", new_entry); err != nil {
						panic(err)
					}
					txn.Commit()
					return nil
				},
			},

			// // 			{
			// // 				Name:    "delete",
			// // 				Aliases: []string{"d"},
			// // 				Usage:   "delete an entry from the list",
			// // 				Action: func(cCtx *cli.Context) error {
			// // 					fmt.Println("delete")
			// // 					return nil
			// // 				},
			// // 			},

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
					// fmt.Println("list")
					// switch cCtx.String("status") {
					// case "Abandoned":
					// 	fmt.Println("abandoned")
					// case "Active":
					// 	fmt.Println("active")
					// case "Done":
					// 	fmt.Println("done")
					// default:
					// 	fmt.Println("fetching all")
					// }
					// switch cCtx.String("category") {
					// case "Fun":
					// 	fmt.Println("fun")
					// case "Personal":
					// 	fmt.Println("personal")
					// case "Work":
					// 	fmt.Println("work")
					// default:
					// 	fmt.Println("fetching all")
					// }

					// Create read-only transaction
					txn := db.Txn(false)
					defer txn.Abort()

					it, err := txn.Get("entry", "id")
					if err != nil {
						panic(err)
					}
					fmt.Println("All the entries:")
					for obj := it.Next(); obj != nil; obj = it.Next() {
						p := obj.(Entry)
						fmt.Printf("%s\n%s\n%s\n%s\n", p.Title, p.Description, p.Status, p.Category)
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
