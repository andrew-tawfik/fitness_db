package main

import (
	"bufio"
	"fitness_db/models"
	"fmt"
	"os"
	"strings"
	"time"
)

// ViewTrainerSchedule shows upcoming sessions for a trainer.
// Demonstrates querying training sessions and classes by trainer ID.
func ViewTrainerSchedule() error {
	fmt.Print("\nEnter Trainer ID: ")
	var trainerID uint
	fmt.Scan(&trainerID)

	// Verify trainer exists
	var trainer models.Trainer
	if err := DB.First(&trainer, trainerID).Error; err != nil {
		return fmt.Errorf("trainer not found")
	}

	fmt.Printf("\n=== Schedule for %s %s ===\n", trainer.FirstName, trainer.LastName)

	// Get upcoming training sessions
	var sessions []models.TrainingSession
	DB.Where("trainer_id = ? AND date >= ?", trainerID, time.Now()).
		Order("date, start_time").
		Find(&sessions)

	fmt.Printf("\nPersonal Training Sessions: %d\n", len(sessions))
	for _, s := range sessions {
		fmt.Printf("  - %s at %s (Status: %s)\n",
			s.Date.Format("Jan 02"),
			s.StartTime.Format("3:04 PM"),
			s.Status)
	}

	// Get upcoming classes
	var classes []models.Class
	DB.Where("trainer_id = ? AND schedule_time >= ?", trainerID, time.Now()).
		Order("schedule_time").
		Find(&classes)

	fmt.Printf("\nGroup Classes: %d\n", len(classes))
	for _, c := range classes {
		fmt.Printf("  - %s on %s (%d enrolled)\n",
			c.ClassName,
			c.ScheduleTime.Format("Jan 02 3:04 PM"),
			c.CurrentEnrollment)
	}

	return nil
}

// MemberLookup allows trainers to search for members by name.
// Displays member's current goal and latest health metric (read-only).
func MemberLookup() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nEnter member name to search: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Case-insensitive search
	var members []models.Member
	searchTerm := "%" + strings.ToLower(name) + "%"
	DB.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?", searchTerm, searchTerm).
		Find(&members)

	if len(members) == 0 {
		return fmt.Errorf("no members found")
	}

	fmt.Printf("\nFound %d member(s):\n", len(members))

	for _, m := range members {
		fmt.Printf("\n%s %s (ID: %d)\n", m.FirstName, m.LastName, m.MemberID)

		// Get current active goal
		var goal models.FitnessGoal
		err := DB.Where("member_id = ? AND status = 'active'", m.MemberID).First(&goal).Error
		if err == nil {
			fmt.Printf("  Goal: %s - Target: %.1f lbs\n", goal.GoalType, goal.TargetWeight)
		} else {
			fmt.Printf("  Goal: None\n")
		}

		// Get latest metric
		var metric models.HealthMetric
		err = DB.Where("member_id = ?", m.MemberID).
			Order("recorded_date DESC").
			First(&metric).Error
		if err == nil {
			fmt.Printf("  Last Metric: %.1f lbs, %d bpm\n", metric.Weight, metric.HeartRate)
		} else {
			fmt.Printf("  Last Metric: None recorded\n")
		}
	}

	return nil
}