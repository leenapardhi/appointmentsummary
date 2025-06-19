package main

import (
	"AppointmentSummmary_Assignment/database"
	"AppointmentSummmary_Assignment/sender"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Invalid number of arguments. Expected usage: ./main <date>")
		return
	}
	date := os.Args[1]
	appointments, err := database.ReadDataForDate(date)
	if err != nil {
		panic(err)
	}
	if err = sender.CreateAndScheduleSummaryAppointmentMessages(appointments); err != nil {
		panic(err)
	}
	fmt.Println("Appointment summaries generated and saved successfully.") 
}
