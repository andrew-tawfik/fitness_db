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

	fmt.Print("\nEnter your Trainer ID: ")
	var trainerID uint
	fmt.Scanln(&trainerID)

	// Verify trainer exists
	var trainer models.Trainer
	if err := DB.First(&trainer, trainerID).Error; err != nil {
		return fmt.Errorf("trainer not found")
	}

	fmt.Print("Enter member name to search (from your clients): ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Get only members this trainer has sessions with
	var members []models.Member
	searchTerm := "%" + strings.ToLower(name) + "%"
	
	err := DB.Table("members m").
		Select("DISTINCT m.*").
		Joins("JOIN training_sessions ts ON m.member_id = ts.member_id").
		Where("ts.trainer_id = ?", trainerID).
		Where("LOWER(m.first_name) LIKE ? OR LOWER(m.last_name) LIKE ?", searchTerm, searchTerm).
		Find(&members).Error

	if err != nil || len(members) == 0 {
		return fmt.Errorf("no clients found matching '%s'", name)
	}

	fmt.Printf("\nFound %d client(s):\n", len(members))

	for _, m := range members {
		fmt.Printf("\n%s %s (ID: %d)\n", m.FirstName, m.LastName, m.MemberID)

		// Show goal
		var goal models.FitnessGoal
		err := DB.Where("member_id = ? AND status = 'active'", m.MemberID).First(&goal).Error
		if err == nil {
			fmt.Printf("  Goal: %s - Target: %.1f lbs by %s\n", 
				goal.GoalType, goal.TargetWeight, goal.TargetDate.Format("Jan 02"))
		} else {
			fmt.Printf("  Goal: None set\n")
		}

		// Show metric history (last 3 entries)
		var metrics []models.HealthMetric
		DB.Where("member_id = ?", m.MemberID).
			Order("recorded_date DESC").
			Limit(3).
			Find(&metrics)

		if len(metrics) > 0 {
			fmt.Printf("  Recent Metrics:\n")
			for _, mt := range metrics {
				fmt.Printf("    %s: %.1f lbs, %d bpm, %.1f%% BF\n",
					mt.RecordedDate.Format("Jan 02"),
					mt.Weight, mt.HeartRate, mt.BodyFatPct)
			}
		} else {
			fmt.Printf("  Metrics: None recorded\n")
		}
	}

	return nil
}