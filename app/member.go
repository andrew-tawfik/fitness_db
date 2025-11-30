package main

import (
	"bufio"
	"errors"
	"fitness_db/models"
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

func RegisterNewMember() error {
	var member models.Member
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Member Registration ===")

	// Get first name
	fmt.Print("Enter first name: ")
	firstName, _ := reader.ReadString('\n')
	member.FirstName = strings.TrimSpace(firstName)
	if member.FirstName == "" {
		return fmt.Errorf("first name cannot be empty")
	}

	// Get last name
	fmt.Print("Enter last name: ")
	lastName, _ := reader.ReadString('\n')
	member.LastName = strings.TrimSpace(lastName)
	if member.LastName == "" {
		return fmt.Errorf("last name cannot be empty")
	}

	// Get email
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	member.Email = strings.TrimSpace(email)
	if member.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// Check if email already exists
	var existingMember models.Member
	if err := DB.Where("email = ?", member.Email).First(&existingMember).Error; err == nil {
		return fmt.Errorf("email already registered")
	}

	// Get date of birth
	fmt.Print("Enter date of birth (YYYY-MM-DD): ")
	dobStr, _ := reader.ReadString('\n')
	dobStr = strings.TrimSpace(dobStr)
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	member.DateOfBirth = dob

	// Get gender
	fmt.Print("Enter gender (Male/Female/Other): ")
	gender, _ := reader.ReadString('\n')
	member.Gender = strings.TrimSpace(gender)

	// Get phone
	fmt.Print("Enter phone number: ")
	phone, _ := reader.ReadString('\n')
	member.Phone = strings.TrimSpace(phone)

	// Set join date to current time
	member.JoinDate = time.Now()

	// Create member in database
	if err := DB.Create(&member).Error; err != nil {
		return fmt.Errorf("failed to register member: %v", err)
	}

	fmt.Printf("\n✓ Member registered successfully! Member ID: %d\n", member.MemberID)
	return nil
}

// ViewDashboard displays a member's dashboard summary using the database view.
// Demonstrates use of the member_dashboard view created in the database.
func ViewDashboard() error {
    
    fmt.Print("\nEnter Member ID: ")
	var memberID uint
	fmt.Scan(&memberID)

	// Query the member_dashboard view
	var result struct {
		FirstName     string
		LastName      string
		TotalClasses  int64
		TotalSessions int64
		ActiveGoals   int64
	}

	err := DB.Raw(`
		SELECT first_name, last_name, total_classes, total_sessions, active_goals 
		FROM member_dashboard 
		WHERE member_id = ?`, memberID).Scan(&result).Error

	if err != nil {
		return fmt.Errorf("member not found")
	}

	fmt.Printf("\n=== Dashboard for %s %s ===\n", result.FirstName, result.LastName)
	fmt.Printf("Classes Enrolled: %d\n", result.TotalClasses)
	fmt.Printf("Training Sessions: %d\n", result.TotalSessions)
	fmt.Printf("Active Goals: %d\n", result.ActiveGoals)

	return nil
}

// AddHealthMetric allows a member to log their health metrics (weight, height, heart rate, body fat %).
func AddHealthMetric() error {
	var metric models.HealthMetric

	fmt.Println("\n=== Log Health Metrics ===")

	// Get member ID
	fmt.Print("Enter your Member ID: ")
	var memberID uint
	_, err := fmt.Scan(&memberID)
	if err != nil {
		return fmt.Errorf("invalid member ID")
	}

	// Verify member exists
	var member models.Member
	if err := DB.First(&member, memberID).Error; err != nil {
		return fmt.Errorf("member not found")
	}

	metric.MemberID = memberID

	// Get the next metric ID for this member
	var maxMetricID uint
	DB.Model(&models.HealthMetric{}).
		Where("member_id = ?", memberID).
		Select("COALESCE(MAX(metric_id), 0)").
		Scan(&maxMetricID)
	metric.MetricID = maxMetricID + 1

	// Get weight
	fmt.Print("Enter weight (lbs): ")
	_, err = fmt.Scan(&metric.Weight)
	if err != nil || metric.Weight <= 0 {
		return fmt.Errorf("invalid weight")
	}

	// Get height
	fmt.Print("Enter height (inches): ")
	_, err = fmt.Scan(&metric.Height)
	if err != nil || metric.Height <= 0 {
		return fmt.Errorf("invalid height")
	}

	// Get heart rate
	fmt.Print("Enter heart rate (bpm): ")
	_, err = fmt.Scan(&metric.HeartRate)
	if err != nil || metric.HeartRate <= 0 {
		return fmt.Errorf("invalid heart rate")
	}

	// Get body fat percentage
	fmt.Print("Enter body fat percentage: ")
	_, err = fmt.Scan(&metric.BodyFatPct)
	if err != nil || metric.BodyFatPct < 0 || metric.BodyFatPct > 100 {
		return fmt.Errorf("invalid body fat percentage")
	}

	// Set recorded date to now
	metric.RecordedDate = time.Now()

	// Save to database
	if err := DB.Create(&metric).Error; err != nil {
		return fmt.Errorf("failed to log health metric: %v", err)
	}

	fmt.Printf("\n✓ Health metrics logged successfully! (Metric ID: %d)\n", metric.MetricID)
	return nil
}

// EnrollClass allows a member to register for a group fitness class.
func EnrollClass() error {

	fmt.Println("\n=== Enroll in Group Class ===")

	// Get member ID
	fmt.Print("Enter your Member ID: ")
	var memberID uint
	_, err := fmt.Scan(&memberID)
	if err != nil {
		return fmt.Errorf("invalid member ID")
	}
	// Verify member exists
	var member models.Member
	if err := DB.First(&member, memberID).Error; err != nil {
		return fmt.Errorf("member not found")
	}

	// Display available classes
	fmt.Println("\nAvailable Classes:")
	fmt.Println("─────────────────────────────────────────────────────────────")

	var classes []models.Class
	if err := DB.Where("schedule_time > ?", time.Now()).
		Order("schedule_time").
		Find(&classes).Error; err != nil {
		return fmt.Errorf("failed to load classes: %v", err)
	}

	if len(classes) == 0 {
		return fmt.Errorf("no upcoming classes available")
	}

	for _, class := range classes {
		trainerName := "TBA"

		// Manually fetch trainer if TrainerID is set
		if class.TrainerID != nil {
			var trainer models.Trainer
			if err := DB.First(&trainer, *class.TrainerID).Error; err == nil {
				trainerName = trainer.FirstName + " " + trainer.LastName
			}
		}

		available := class.Capacity - class.CurrentEnrollment
		fmt.Printf("ID: %d | %s | Trainer: %s\n", class.ClassID, class.ClassName, trainerName)
		fmt.Printf("  Time: %s | Duration: %d min | Room: %s\n",
			class.ScheduleTime.Format("Mon Jan 02, 3:04 PM"), class.Duration, class.RoomNumber)
		fmt.Printf("  Capacity: %d/%d (Available: %d)\n\n",
			class.CurrentEnrollment, class.Capacity, available)
	}

	// Get class selection
	fmt.Print("Enter Class ID to enroll: ")
	var classID uint
	_, err = fmt.Scan(&classID)
	if err != nil {
		return fmt.Errorf("invalid class ID")
	}

	// Verify class exists and has capacity
	var selectedClass models.Class
	if err := DB.First(&selectedClass, classID).Error; err != nil {
		return fmt.Errorf("class not found")
	}

	// Check if already enrolled - use errors.Is with gorm.ErrRecordNotFound
	var existingEnrollment models.ClassEnrollment
	err = DB.Where("member_id = ? AND class_id = ?", memberID, classID).
		First(&existingEnrollment).Error

	if err == nil {
		// Record found - already enrolled
		return fmt.Errorf("you are already enrolled in this class")
	}

	// Check if it's not a "record not found" error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check enrollment status: %v", err)
	}

	// Record not found means not enrolled - proceed with enrollment

	// Create enrollment (trigger will check capacity)
	enrollment := models.ClassEnrollment{
		MemberID:       memberID,
		ClassID:        classID,
		EnrollmentDate: time.Now(),
		Status:         "active",
	}

	if err := DB.Create(&enrollment).Error; err != nil {
		if strings.Contains(err.Error(), "Class is full") {
			return fmt.Errorf("class is full - cannot enroll")
		}
		return fmt.Errorf("failed to enroll in class: %v", err)
	}

	fmt.Printf("\n✓ Successfully enrolled in %s!\n", selectedClass.ClassName)
	return nil
}

// BookPersonalTraining allows a member to schedule a one-on-one training session with a trainer.
func BookPersonalTraining() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Book Personal Training Session ===")

	// Get member ID
	fmt.Print("Enter your Member ID: ")
	var memberID uint
	_, err := fmt.Scan(&memberID)
	if err != nil {
		return fmt.Errorf("invalid member ID")
	}

	// Verify member exists
	var member models.Member
	if err := DB.First(&member, memberID).Error; err != nil {
		return fmt.Errorf("member not found")
	}

	// Display available trainers
	fmt.Println("\nAvailable Trainers:")
	fmt.Println("─────────────────────────────────────────────────────────────")

	var trainers []models.Trainer
	DB.Find(&trainers)

	if len(trainers) == 0 {
		return fmt.Errorf("no trainers available")
	}

	for _, trainer := range trainers {
		fmt.Printf("ID: %d | %s %s | Specialization: %s\n",
			trainer.TrainerID, trainer.FirstName, trainer.LastName, trainer.Specialization)
	}

	// Get trainer selection
	fmt.Print("\nEnter Trainer ID: ")
	var trainerID uint
	_, err = fmt.Scan(&trainerID)
	if err != nil {
		return fmt.Errorf("invalid trainer ID")
	}

	// Verify trainer exists
	var trainer models.Trainer
	if err := DB.First(&trainer, trainerID).Error; err != nil {
		return fmt.Errorf("trainer not found")
	}

	// Get session date
	fmt.Print("Enter session date (YYYY-MM-DD): ")
	dateStr, _ := reader.ReadString('\n')
	dateStr = strings.TrimSpace(dateStr)
	sessionDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	// Check if date is in the future
	if sessionDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return fmt.Errorf("cannot book sessions in the past")
	}

	// Get start time
	fmt.Print("Enter start time (HH:MM, 24-hour format): ")
	startTimeStr, _ := reader.ReadString('\n')
	startTimeStr = strings.TrimSpace(startTimeStr)
	startTimeParsed, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		return fmt.Errorf("invalid time format, use HH:MM")
	}

	// Combine date and time
	startTime := time.Date(sessionDate.Year(), sessionDate.Month(), sessionDate.Day(),
		startTimeParsed.Hour(), startTimeParsed.Minute(), 0, 0, sessionDate.Location())

	// Get duration
	fmt.Print("Enter duration (minutes): ")
	var duration int
	_, err = fmt.Scan(&duration)
	if err != nil || duration <= 0 || duration > 240 {
		return fmt.Errorf("invalid duration (must be between 1-240 minutes)")
	}

	endTime := startTime.Add(time.Duration(duration) * time.Minute)

	// Check for trainer conflicts
	var trainerConflicts int64
	DB.Model(&models.TrainingSession{}).
		Where("trainer_id = ? AND date = ? AND status != ?", trainerID, sessionDate, "cancelled").
		Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?)",
			endTime, startTime, endTime, startTime).
		Count(&trainerConflicts)

	if trainerConflicts > 0 {
		return fmt.Errorf("trainer is not available at this time")
	}

	// Check for member conflicts
	var memberConflicts int64
	DB.Model(&models.TrainingSession{}).
		Where("member_id = ? AND date = ? AND status != ?", memberID, sessionDate, "cancelled").
		Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?)",
			endTime, startTime, endTime, startTime).
		Count(&memberConflicts)

	if memberConflicts > 0 {
		return fmt.Errorf("you already have a session booked at this time")
	}

	// Create training session
	session := models.TrainingSession{
		MemberID:  memberID,
		TrainerID: &trainerID,
		Date:      sessionDate,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    "scheduled",
	}

	if err := DB.Create(&session).Error; err != nil {
		return fmt.Errorf("failed to book training session: %v", err)
	}

	fmt.Printf("\n✓ Training session booked successfully!\n")
	fmt.Printf("  Trainer: %s %s\n", trainer.FirstName, trainer.LastName)
	fmt.Printf("  Date: %s\n", sessionDate.Format("Mon Jan 02, 2006"))
	fmt.Printf("  Time: %s - %s\n", startTime.Format("3:04 PM"), endTime.Format("3:04 PM"))

	return nil
}
