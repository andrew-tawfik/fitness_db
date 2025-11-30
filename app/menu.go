package main

import "fmt"

func MainMenu() {
    for {
        fmt.Println("\n=== Health & Fitness Club Management ===")
        fmt.Println("1. Member Login")
        fmt.Println("2. Trainer Login")
        fmt.Println("3. Admin Login")
        fmt.Println("4. Register New Member")
        fmt.Println("5. Exit")
        
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
            RegisterNewMember()
        case 5:
            return
        }
    }
}

func MemberMenu() {

}

func TrainerMenu() {

}

func AdminMenu() {
	
}