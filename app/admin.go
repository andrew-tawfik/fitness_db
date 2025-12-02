package main

import (
	"bufio"
	"fitness_db/models"
	"fmt"
	"os"
	"strings"
	"time"
)

// CreateClass allows an admin to create a new fitness class.
func CreateClass() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Create New Class ===")

	// Get class name
	fmt.Print("Enter class name: ")
	className, _ := reader.ReadString('\n')
	className = strings.TrimSpace(className)
	if className == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	// Display available trainers
	fmt.Println("\nAvailable Trainers:")
	fmt.Println("─────────────────────────────────────────────────────────────")
	var trainers []models.Trainer
	DB.Find(&trainers)

	if len(trainers) == 0 {
		fmt.Println("No trainers in system. Class will be created without assigned trainer.")
	} else {
		for _, trainer := range trainers {
			fmt.Printf("ID: %d | %s %s | Specialization: %s\n",
				trainer.TrainerID, trainer.FirstName, trainer.LastName, trainer.Specialization)
		}
	}

	// Get trainer ID (optional)
	fmt.Print("\nEnter Trainer ID (or 0 for TBA): ")
	var trainerID uint
	_, err := fmt.Scan(&trainerID)
	if err != nil {
		return fmt.Errorf("invalid trainer ID")
	}

	// Validate trainer exists if provided
	var trainerPtr *uint
	if trainerID != 0 {
		var trainer models.Trainer
		if err := DB.First(&trainer, trainerID).Error; err != nil {
			return fmt.Errorf("trainer not found")
		}
		trainerPtr = &trainerID
	}

	// Get schedule time
	fmt.Print("Enter schedule date and time (YYYY-MM-DD HH:MM): ")
	scheduleStr, _ := reader.ReadString('\n')
	scheduleStr = strings.TrimSpace(scheduleStr)
	scheduleTime, err := time.Parse("2006-01-02 15:04", scheduleStr)
	if err != nil {
		return fmt.Errorf("invalid date/time format, use YYYY-MM-DD HH:MM")
	}

	scheduleTime = time.Date(
		scheduleTime.Year(),
		scheduleTime.Month(),
		scheduleTime.Day(),
		scheduleTime.Hour(),
		scheduleTime.Minute(),
		0, 0,
		time.Local)

	// Check if time is in the future
	if scheduleTime.Before(time.Now()) {
		return fmt.Errorf("cannot schedule classes in the past")
	}

	// Get duration
	fmt.Print("Enter duration (minutes): ")
	var duration int
	_, err = fmt.Scan(&duration)
	if err != nil || duration <= 0 || duration > 300 {
		return fmt.Errorf("invalid duration (must be between 1-300 minutes)")
	}

	// Get capacity
	fmt.Print("Enter class capacity: ")
	var capacity int
	_, err = fmt.Scan(&capacity)
	if err != nil || capacity <= 0 || capacity > 100 {
		return fmt.Errorf("invalid capacity (must be between 1-100)")
	}

	// Get room number
	fmt.Print("Enter room number: ")
	roomNumber, _ := reader.ReadString('\n')
	roomNumber = strings.TrimSpace(roomNumber)
	if roomNumber == "" {
		return fmt.Errorf("room number cannot be empty")
	}

	// Create the class
	class := models.Class{
		ClassName:         className,
		TrainerID:         trainerPtr,
		ScheduleTime:      scheduleTime,
		Duration:          duration,
		Capacity:          capacity,
		CurrentEnrollment: 0,
		RoomNumber:        roomNumber,
	}

	if err := DB.Create(&class).Error; err != nil {
		return fmt.Errorf("failed to create class: %v", err)
	}

	fmt.Printf("\nClass created successfully! Class ID: %d\n", class.ClassID)
	fmt.Printf("  Name: %s\n", class.ClassName)
	fmt.Printf("  Schedule: %s\n", scheduleTime.Format("Mon Jan 02, 2006 at 3:04 PM"))
	fmt.Printf("  Duration: %d minutes\n", duration)
	fmt.Printf("  Room: %s\n", roomNumber)
	fmt.Printf("  Capacity: %d\n", capacity)

	return nil
}

// UpdateClass allows an admin to modify or cancel an existing class.
func UpdateClass() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Update/Cancel Class ===")

	// Display existing classes
	fmt.Println("\nExisting Classes:")
	fmt.Println("─────────────────────────────────────────────────────────────")

	var classes []models.Class
	DB.Order("schedule_time").Find(&classes)

	if len(classes) == 0 {
		return fmt.Errorf("no classes in the system")
	}

	for _, class := range classes {
		trainerName := "TBA"
		if class.TrainerID != nil {
			var trainer models.Trainer
			if err := DB.First(&trainer, *class.TrainerID).Error; err == nil {
				trainerName = trainer.FirstName + " " + trainer.LastName
			}
		}
		fmt.Printf("ID: %d | %s | Trainer: %s\n", class.ClassID, class.ClassName, trainerName)
		fmt.Printf("  Time: %s | Duration: %d min | Room: %s\n",
			class.ScheduleTime.Format("Mon Jan 02, 3:04 PM"), class.Duration, class.RoomNumber)
		fmt.Printf("  Enrollment: %d/%d\n\n",
			class.CurrentEnrollment, class.Capacity)
	}

	// Get class ID to update
	fmt.Print("Enter Class ID to update: ")
	var classID uint
	_, err := fmt.Scan(&classID)
	if err != nil {
		return fmt.Errorf("invalid class ID")
	}

	// Fetch the class
	var class models.Class
	if err := DB.First(&class, classID).Error; err != nil {
		return fmt.Errorf("class not found")
	}

	fmt.Printf("\nCurrent class details:\n")
	fmt.Printf("  Name: %s\n", class.ClassName)
	fmt.Printf("  Schedule: %s\n", class.ScheduleTime.Format("Mon Jan 02, 2006 at 3:04 PM"))
	fmt.Printf("  Capacity: %d (Current enrollment: %d)\n", class.Capacity, class.CurrentEnrollment)
	fmt.Printf("  Room: %s\n", class.RoomNumber)

	// Update menu
	fmt.Println("\nWhat would you like to update?")
	fmt.Println("1. Change schedule time")
	fmt.Println("2. Change trainer")
	fmt.Println("3. Change capacity")
	fmt.Println("4. Change room number")
	fmt.Println("5. Cancel class (delete)")
	fmt.Println("7. Back to menu")
	fmt.Print("\nEnter choice: ")

	var choice int
	_, err = fmt.Scan(&choice)
	if err != nil {
		return fmt.Errorf("invalid choice")
	}

	switch choice {
	case 1: // Change schedule time
		fmt.Print("Enter new schedule time (YYYY-MM-DD HH:MM): ")
		timeStr, _ := reader.ReadString('\n')
		timeStr = strings.TrimSpace(timeStr)
		newTime, err := time.Parse("2006-01-02 15:04", timeStr)
		if err != nil {
			return fmt.Errorf("invalid date/time format")
		}

		newTime = time.Date(
			newTime.Year(),
			newTime.Month(),
			newTime.Day(),
			newTime.Hour(),
			newTime.Minute(),
			0, 0,
			time.Local)

		if newTime.Before(time.Now()) {
			return fmt.Errorf("cannot schedule classes in the past")
		}

		if err := DB.Model(&class).Update("schedule_time", newTime).Error; err != nil {
			return fmt.Errorf("failed to update schedule: %v", err)
		}
		fmt.Printf("\nSchedule updated to %s\n", newTime.Format("Mon Jan 02, 2006 at 3:04 PM"))

	case 2: // Change trainer
		fmt.Println("\nAvailable Trainers:")
		var trainers []models.Trainer
		DB.Find(&trainers)
		for _, t := range trainers {
			fmt.Printf("ID: %d | %s %s | %s\n", t.TrainerID, t.FirstName, t.LastName, t.Specialization)
		}

		fmt.Print("\nEnter new Trainer ID (or 0 for TBA): ")
		var newTrainerID uint
		fmt.Scan(&newTrainerID)

		var trainerPtr *uint
		if newTrainerID != 0 {
			var trainer models.Trainer
			if err := DB.First(&trainer, newTrainerID).Error; err != nil {
				return fmt.Errorf("trainer not found")
			}
			trainerPtr = &newTrainerID
		}

		if err := DB.Model(&class).Update("trainer_id", trainerPtr).Error; err != nil {
			return fmt.Errorf("failed to update trainer: %v", err)
		}
		fmt.Println("\nTrainer updated successfully")

	case 3: // Change capacity
		fmt.Printf("Current enrollment: %d\n", class.CurrentEnrollment)
		fmt.Print("Enter new capacity (must be >= current enrollment): ")
		var newCapacity int
		fmt.Scan(&newCapacity)

		if newCapacity < class.CurrentEnrollment {
			return fmt.Errorf("cannot reduce capacity below current enrollment (%d)", class.CurrentEnrollment)
		}
		if newCapacity <= 0 || newCapacity > 100 {
			return fmt.Errorf("invalid capacity (must be between 1-100)")
		}

		if err := DB.Model(&class).Update("capacity", newCapacity).Error; err != nil {
			return fmt.Errorf("failed to update capacity: %v", err)
		}
		fmt.Printf("\nCapacity updated to %d\n", newCapacity)

	case 4: // Change room number
		fmt.Print("Enter new room number: ")
		newRoom, _ := reader.ReadString('\n')
		newRoom = strings.TrimSpace(newRoom)
		if newRoom == "" {
			return fmt.Errorf("room number cannot be empty")
		}

		if err := DB.Model(&class).Update("room_number", newRoom).Error; err != nil {
			return fmt.Errorf("failed to update room: %v", err)
		}
		fmt.Printf("\nRoom updated to %s\n", newRoom)

	case 5: // Cancel/Delete class
		if class.CurrentEnrollment > 0 {
			fmt.Printf("\nWARNING: This class has %d enrolled members.\n", class.CurrentEnrollment)
			fmt.Print("Are you sure you want to cancel? (yes/no): ")
			confirm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(strings.ToLower(confirm))
			if confirm != "yes" {
				fmt.Println("Cancellation aborted.")
				return nil
			}
		}

		if err := DB.Delete(&class).Error; err != nil {
			return fmt.Errorf("failed to cancel class: %v", err)
		}
		fmt.Println("\nClass cancelled and removed from system")

	case 7:
		return nil

	default:
		return fmt.Errorf("invalid choice")
	}

	return nil
}
