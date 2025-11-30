package main

import "fmt"

func MainMenu() {
	for {
		fmt.Println("\n=== Health & Fitness Club Management ===")
		fmt.Println("1. Member Portal")
		fmt.Println("2. Trainer Portal")
		fmt.Println("3. Admin Portal")
		fmt.Println("4. Register New Member")
		fmt.Println("5. Exit")
		fmt.Print("\nChoice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			MemberMenu()
		case 2:
			TrainerMenu()
		case 3:
			AdminMenu()
		case 4:
			if err := RegisterNewMember(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 5:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func MemberMenu() {
	for {
		fmt.Println("\n--- Member Menu ---")
		fmt.Println("1. View Dashboard")
		fmt.Println("2. Log Health Metrics")
		fmt.Println("3. Book Personal Training")
		fmt.Println("4. Enroll in Class")
		fmt.Println("5. Back")
		fmt.Print("\nChoice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if err := ViewDashboard(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 2:
			if err := AddHealthMetric(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 3:
			if err := BookPersonalTraining(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 4:
			if err := EnrollClass(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 5:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func TrainerMenu() {
	for {
		fmt.Println("\n--- Trainer Menu ---")
		fmt.Println("1. View My Schedule")
		fmt.Println("2. Member Lookup")
		fmt.Println("3. Back")
		fmt.Print("\nChoice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if err := ViewTrainerSchedule(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 2:
			if err := MemberLookup(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 3:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func AdminMenu() {
	for {
		fmt.Println("\n--- Admin Menu ---")
		fmt.Println("1. Create New Class")
		fmt.Println("2. Update/Cancel Class")
		fmt.Println("3. Back")
		fmt.Print("\nChoice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if err := CreateClass(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 2:
			if err := UpdateClass(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 3:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}