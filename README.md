# todo-cli

<p>A cli for managing todos. This project uses an inmem database and its not safe from security breaches. Use at your own risk.

You can:
- Create a new entry
- Edit your entry(title, description and/or category of entry)
- Delete your entry
- Change the status of the entry
- Fetch and read all the entries

## Prerequisites

[Go](https://go.dev/doc/install)

## Setup 

```sh
go mod tidy
go build
```

## Example usage

*Use todo help to see specific details of the commands*

#### New
```sh
todo new || n
todo -c fun || personal || work
```
#### Delete
```sh
todo delete || delete
todo delete -id 1
```
#### Edit
```sh
todo edit || e
todo edit -id 1 -f title || category || description
```
#### Read entries
```sh
todo list || ls
todo list -s abandoned || active || done
todo list -c work || fun || personal
todo list -id 1
```
#### Update status of entry
```sh
todo update || u
todo update -id 1 -s done || active || abandoned
```


