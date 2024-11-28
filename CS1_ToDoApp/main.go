package main

import (
	"CS1_ToDoApp/database"
	"CS1_ToDoApp/routes"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	database.InitDB()
	defer database.CloseDB()

	r := routes.SetupRouter()

	err := r.Run(":8888")
	if err != nil {
		log.Fatal("Cann't run server", err)
	}
}
