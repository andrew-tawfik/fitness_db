package main

import (
	"fitness_db/models"
	"fmt"
	"time"
)

func SeedDatabase() {
	var count int64
	DB.Model(&models.Member{}).Count(&count)
	if count > 0 {
		fmt.Println("Database already seeded, skipping...")
		return
	}

	fmt.Println("Seeding database with realistic test data...")
	members := []models.Member{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Gender:      "Male", Phone: "613-555-0001"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com",
			DateOfBirth: time.Date(1992, 5, 15, 0, 0, 0, 0, time.UTC),
			Gender:      "Female", Phone: "613-555-0002"},
		{FirstName: "Bob", LastName: "Wilson", Email: "bob@example.com",
			DateOfBirth: time.Date(1988, 8, 20, 0, 0, 0, 0, time.UTC),
			Gender:      "Male", Phone: "613-555-0003"},
		{FirstName: "Alice", LastName: "Brown", Email: "alice@example.com",
			DateOfBirth: time.Date(1995, 3, 10, 0, 0, 0, 0, time.UTC),
			Gender:      "Female", Phone: "613-555-0004"},
		{FirstName: "Charlie", LastName: "Davis", Email: "charlie@example.com",
			DateOfBirth: time.Date(1985, 12, 5, 0, 0, 0, 0, time.UTC),
			Gender:      "Male", Phone: "613-555-0005"},
	}
	DB.Create(&members)

	trainers := []models.Trainer{
		{FirstName: "Mike", LastName: "Johnson", Email: "mike@fitness.com",
			Specialization: "Strength Training", Phone: "613-555-1001"},
		{FirstName: "Sarah", LastName: "Williams", Email: "sarah@fitness.com",
			Specialization: "Yoga", Phone: "613-555-1002"},
		{FirstName: "Tom", LastName: "Anderson", Email: "tom@fitness.com",
			Specialization: "Cardio & HIIT", Phone: "613-555-1003"},
	}
	DB.Create(&trainers)

	admins := []models.Admin{
		{FirstName: "Alice", LastName: "Manager", Email: "alice@admin.com",
			Phone: "613-555-2001", Role: "manager"},
	}
	DB.Create(&admins)

	now := time.Now()

	classes := []models.Class{
		// SCENARIO 1: Empty class (no enrollments yet)
		{ClassName: "Morning Yoga", TrainerID: &trainers[1].TrainerID,
			ScheduleTime: now.Add(24 * time.Hour), Duration: 60,
			Capacity: 20, CurrentEnrollment: 0, RoomNumber: "Studio A"},

		// SCENARIO 2: Partially filled class (5/15 spots)
		{ClassName: "Evening Strength", TrainerID: &trainers[0].TrainerID,
			ScheduleTime: now.Add(48 * time.Hour), Duration: 90,
			Capacity: 15, CurrentEnrollment: 0, RoomNumber: "Gym 1"},

		// SCENARIO 3: Nearly full class (2/3 spots taken)
		{ClassName: "HIIT Cardio", TrainerID: &trainers[2].TrainerID,
			ScheduleTime: now.Add(72 * time.Hour), Duration: 45,
			Capacity: 3, CurrentEnrollment: 0, RoomNumber: "Studio B"},

		// SCENARIO 4: Class in 1 week
		{ClassName: "Advanced Yoga", TrainerID: &trainers[1].TrainerID,
			ScheduleTime: now.Add(7 * 24 * time.Hour), Duration: 75,
			Capacity: 12, CurrentEnrollment: 0, RoomNumber: "Studio A"},

		// SCENARIO 5: Class without trainer (TBA)
		{ClassName: "Beginner Pilates", TrainerID: nil,
			ScheduleTime: now.Add(96 * time.Hour), Duration: 60,
			Capacity: 15, CurrentEnrollment: 0, RoomNumber: "Studio C"},
	}
	DB.Create(&classes)

	enrollments := []models.ClassEnrollment{
		{MemberID: members[0].MemberID, ClassID: classes[1].ClassID,
			EnrollmentDate: now.AddDate(0, 0, -2), Status: "active"},
		{MemberID: members[1].MemberID, ClassID: classes[1].ClassID,
			EnrollmentDate: now.AddDate(0, 0, -2), Status: "active"},
		{MemberID: members[2].MemberID, ClassID: classes[1].ClassID,
			EnrollmentDate: now.AddDate(0, 0, -1), Status: "active"},
		{MemberID: members[3].MemberID, ClassID: classes[1].ClassID,
			EnrollmentDate: now.AddDate(0, 0, -1), Status: "active"},
		{MemberID: members[4].MemberID, ClassID: classes[1].ClassID,
			EnrollmentDate: now, Status: "active"},

		{MemberID: members[0].MemberID, ClassID: classes[2].ClassID,
			EnrollmentDate: now.AddDate(0, 0, -1), Status: "active"},
		{MemberID: members[2].MemberID, ClassID: classes[2].ClassID,
			EnrollmentDate: now, Status: "active"},

		{MemberID: members[1].MemberID, ClassID: classes[3].ClassID,
			EnrollmentDate: now.AddDate(0, 0, -3), Status: "active"},
	}

	for _, e := range enrollments {
		DB.Exec(`INSERT INTO class_enrollments (member_id, class_id, enrollment_date, status) 
			VALUES (?, ?, ?, ?)`, e.MemberID, e.ClassID, e.EnrollmentDate, e.Status)

		DB.Exec(`UPDATE classes SET current_enrollment = current_enrollment + 1 WHERE class_id = ?`, e.ClassID)
	}

	sessions := []models.TrainingSession{
		{MemberID: members[0].MemberID, TrainerID: &trainers[0].TrainerID,
			Date:      now.AddDate(0, 0, 2),
			StartTime: time.Date(now.Year(), now.Month(), now.Day()+2, 10, 0, 0, 0, now.Location()),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day()+2, 11, 0, 0, 0, now.Location()),
			Status:    "scheduled"},
		{MemberID: members[0].MemberID, TrainerID: &trainers[0].TrainerID,
			Date:      now.AddDate(0, 0, 5),
			StartTime: time.Date(now.Year(), now.Month(), now.Day()+5, 10, 0, 0, 0, now.Location()),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day()+5, 11, 0, 0, 0, now.Location()),
			Status:    "scheduled"},

		{MemberID: members[1].MemberID, TrainerID: &trainers[1].TrainerID,
			Date:      now.AddDate(0, 0, 3),
			StartTime: time.Date(now.Year(), now.Month(), now.Day()+3, 14, 0, 0, 0, now.Location()),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day()+3, 15, 0, 0, 0, now.Location()),
			Status:    "scheduled"},

		{MemberID: members[1].MemberID, TrainerID: &trainers[1].TrainerID,
			Date:      now.AddDate(0, 0, 6),
			StartTime: time.Date(now.Year(), now.Month(), now.Day()+6, 14, 0, 0, 0, now.Location()),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day()+6, 15, 0, 0, 0, now.Location()),
			Status:    "scheduled"},

		{MemberID: members[2].MemberID, TrainerID: &trainers[2].TrainerID,
			Date:      now.AddDate(0, 0, 4),
			StartTime: time.Date(now.Year(), now.Month(), now.Day()+4, 16, 0, 0, 0, now.Location()),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day()+4, 17, 0, 0, 0, now.Location()),
			Status:    "scheduled"},

		{MemberID: members[3].MemberID, TrainerID: &trainers[0].TrainerID,
			Date:      now.AddDate(0, 0, 2),
			StartTime: time.Date(now.Year(), now.Month(), now.Day()+2, 15, 0, 0, 0, now.Location()),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day()+2, 16, 0, 0, 0, now.Location()),
			Status:    "scheduled"},
	}
	DB.Create(&sessions)

	metrics := []models.HealthMetric{
		{MemberID: members[0].MemberID, MetricID: 1, Weight: 185.0,
			Height: 70.0, HeartRate: 75, BodyFatPct: 22.0,
			RecordedDate: now.AddDate(0, 0, -28)},
		{MemberID: members[0].MemberID, MetricID: 2, Weight: 182.5,
			Height: 70.0, HeartRate: 73, BodyFatPct: 21.2,
			RecordedDate: now.AddDate(0, 0, -21)},
		{MemberID: members[0].MemberID, MetricID: 3, Weight: 180.5,
			Height: 70.0, HeartRate: 72, BodyFatPct: 20.5,
			RecordedDate: now.AddDate(0, 0, -14)},
		{MemberID: members[0].MemberID, MetricID: 4, Weight: 178.0,
			Height: 70.0, HeartRate: 70, BodyFatPct: 19.8,
			RecordedDate: now.AddDate(0, 0, -7)},

		{MemberID: members[1].MemberID, MetricID: 1, Weight: 145.0,
			Height: 65.0, HeartRate: 68, BodyFatPct: 22.0,
			RecordedDate: now.AddDate(0, 0, -14)},
		{MemberID: members[1].MemberID, MetricID: 2, Weight: 144.5,
			Height: 65.0, HeartRate: 67, BodyFatPct: 21.8,
			RecordedDate: now.AddDate(0, 0, -7)},

		{MemberID: members[2].MemberID, MetricID: 1, Weight: 175.0,
			Height: 72.0, HeartRate: 70, BodyFatPct: 18.0,
			RecordedDate: now.AddDate(0, 0, -21)},
		{MemberID: members[2].MemberID, MetricID: 2, Weight: 178.5,
			Height: 72.0, HeartRate: 68, BodyFatPct: 17.5,
			RecordedDate: now.AddDate(0, 0, -7)},

		{MemberID: members[3].MemberID, MetricID: 1, Weight: 155.0,
			Height: 66.0, HeartRate: 72, BodyFatPct: 25.0,
			RecordedDate: now.AddDate(0, 0, -3)},

		{MemberID: members[4].MemberID, MetricID: 1, Weight: 200.0,
			Height: 68.0, HeartRate: 80, BodyFatPct: 28.0,
			RecordedDate: now.AddDate(0, 0, -60)},
	}
	DB.Create(&metrics)

	goals := []models.FitnessGoal{
		{MemberID: members[0].MemberID, GoalID: 1, GoalType: "Weight Loss",
			TargetWeight: 170.0, TargetDate: now.AddDate(0, 3, 0), Status: "active"},

		{MemberID: members[1].MemberID, GoalID: 1, GoalType: "Maintain Weight",
			TargetWeight: 145.0, TargetDate: now.AddDate(0, 6, 0), Status: "active"},

		{MemberID: members[2].MemberID, GoalID: 1, GoalType: "Muscle Gain",
			TargetWeight: 185.0, TargetDate: now.AddDate(0, 6, 0), Status: "active"},

		{MemberID: members[3].MemberID, GoalID: 1, GoalType: "Weight Loss",
			TargetWeight: 140.0, TargetDate: now.AddDate(0, 4, 0), Status: "active"},

		{MemberID: members[4].MemberID, GoalID: 1, GoalType: "Weight Loss",
			TargetWeight: 180.0, TargetDate: now.AddDate(0, -1, 0), Status: "active"},
	}
	DB.Create(&goals)
}
