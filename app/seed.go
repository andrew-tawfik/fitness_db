package main

import (
    "fitness_db/models"
    "time"
)

func SeedDatabase() {
    // Members
    members := []models.Member{
        {FirstName: "John", LastName: "Doe", Email: "john@example.com", 
         DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), 
         Gender: "Male", Phone: "613-555-0001"},
        {FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", 
         DateOfBirth: time.Date(1992, 5, 15, 0, 0, 0, 0, time.UTC), 
         Gender: "Female", Phone: "613-555-0002"},
        {FirstName: "Bob", LastName: "Wilson", Email: "bob@example.com", 
         DateOfBirth: time.Date(1988, 8, 20, 0, 0, 0, 0, time.UTC), 
         Gender: "Male", Phone: "613-555-0003"},
    }
    DB.Create(&members)

    // Trainers
    trainers := []models.Trainer{
        {FirstName: "Mike", LastName: "Johnson", Email: "mike@fitness.com", 
         Specialization: "Strength Training", Phone: "613-555-1001"},
        {FirstName: "Sarah", LastName: "Williams", Email: "sarah@fitness.com", 
         Specialization: "Yoga", Phone: "613-555-1002"},
    }
    DB.Create(&trainers)

    // Classes
    classes := []models.Class{
        {ClassName: "Morning Yoga", TrainerID: &trainers[1].TrainerID, 
         ScheduleTime: time.Now().Add(24 * time.Hour), Duration: 60, 
         Capacity: 20, RoomNumber: "Studio A"},
        {ClassName: "Evening Strength", TrainerID: &trainers[0].TrainerID, 
         ScheduleTime: time.Now().Add(48 * time.Hour), Duration: 90, 
         Capacity: 15, RoomNumber: "Gym 1"},
    }
    DB.Create(&classes)

    // Health Metrics
    metrics := []models.HealthMetric{
        {MemberID: members[0].MemberID, MetricID: 1, Weight: 180.5, 
         Height: 70.0, HeartRate: 72, BodyFatPct: 18.5},
    }
    DB.Create(&metrics)

    // Fitness Goals
    goals := []models.FitnessGoal{
        {MemberID: members[0].MemberID, GoalID: 1, GoalType: "Weight Loss", 
         TargetWeight: 170.0, TargetDate: time.Now().AddDate(0, 3, 0), 
         Status: "active"},
    }
    DB.Create(&goals)
}