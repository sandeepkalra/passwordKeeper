package main

import (
	"bufio"
	"fmt"
	"os"

	ui "github.com/manifoldco/promptui"
)

/**********/
type dbReader interface {
	Read(sqlQuery string, sqlParams []interface{}) error
}
type dbWriter interface {
	Write(sqlQuery string, sqlParams []interface{}) error
}
type DBInterface interface {
	Init(sqlName, sqlDBSchema, sqlUser, sqlPassword string, sqlPort int) error
	dbReader
	dbWriter
	Close() error
}

/**********/

type mysqlDb struct {
}

func (sql *mysqlDb) Init(sqlName, sqlDBSchema, sqlUser, sqlPassword string, sqlPort int) error {
	return nil
}
func (sql *mysqlDb) Read(sqlQuery string, sqlParams []interface{}) error {
	return nil
}
func (sql *mysqlDb) Write(sqlQuery string, sqlParams []interface{}) error {
	return nil
}
func (sql *mysqlDb) Close() error {
	return nil
}

func askString(s string, isHidden bool, v interface{}) (result string, e error) {
	label := "Enter " + s + " :"
	mask := rune(0)

	validate := func(s string) error {
		if len(s) < 2 {
			return fmt.Errorf("too small")
		}
		return nil
	}

	if v != nil {
		validate = v.(func(string) error)
	}

	if isHidden {
		mask = rune(0x220E)
	}
	prompt := ui.Prompt{
		Label:    label,
		Mask:     mask,
		Validate: validate,
	}
	result, e = prompt.Run()
	return
}

func askSelect(s []string, isHidden bool) (index int, result string, e error) {
	prompt := ui.SelectWithAdd{
		Label:    "Please select one",
		Items:    s,
		AddLabel: "Other",
	}
	index, result, e = prompt.Run()
	return
}

func performDBInit(dbIf DBInterface) error {
	return dbIf.Init("", "", "", "", 0)
}

func performDBOperation(dbIf DBInterface, name, password, operation, key, value string) {
	switch operation {

	case "store":
		values := []string{key, value, name, password}
		args := make([]interface{}, 4)
		for k, v := range values {
			args[k] = v
		}
		dbIf.Write("insert into db values(?,?) where name=?, password=?", args)
		fmt.Println(" ... Stored")

	case "show":
		values := []string{value, name, password, key}
		args := make([]interface{}, 4)
		for k, v := range values {
			args[k] = v
		}
		dbIf.Read("select value from db where name=?, password = ?, key = ?", args)
		fmt.Println(" ... Done")
	case "update":
		values := []string{value, name, password, key}
		args := make([]interface{}, 4)
		for k, v := range values {
			args[k] = v
		}
		fmt.Println(" ... Updated")
	case "delete":
		values := []string{name, password, key}
		args := make([]interface{}, 3)
		for k, v := range values {
			args[k] = v
		}
		fmt.Println(" ... Deleted")
	default:
		return
	}
}

func main() {
	name := ""
	password := ""
	option := ""
	index := 0
	var e error
	var dbIf mysqlDb
	keys := []string{"gmail", "twitter", "facebook", "workmail", "personalmail"}
	operations := []string{"show", "update", "delete"}

	if e = performDBInit(&dbIf); e != nil {
		return
	}

	if name, e = askString("Name", true, nil); e != nil {
		fmt.Println(e.Error())
		return
	}
	if password, e = askString("Password", true,
		func(s string) error {
			if len(s) < 4 {
				return fmt.Errorf("too small")
			}
			return nil
		}); e != nil {
		fmt.Println(e.Error())
		return
	}

	if index, option, e = askSelect(keys, true); e == nil {
		if index == -1 {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Give Value to Store:")
			text, _ := reader.ReadString('\n')
			performDBOperation(&dbIf, name, password, "store", option, text)
		} else {
			prmpt := ui.Select{
				Label: "\tSelect Operation",
				Items: operations,
				Size:  3,
			}
			_, result, err := prmpt.Run()
			if err != nil {
				fmt.Println("failed :", err.Error())
				return
			}
			performDBOperation(&dbIf, name, password, result, option, "")

		}
	}
}
