# Health & Fitness Club Management System

**Course:** COMP 3005 - Database Management Systems  
**Term:** Fall 2025  
**Student:** [Your Name]  
**Student ID:** [Your ID]  
**Group Size:** Solo Project

## Demo Video
**[Watch Full Demo Video Here](demo.mp4)**

## Project Overview
A comprehensive database-driven management system for fitness clubs supporting three user roles: Members, Trainers, and Administrative Staff. Built using PostgreSQL and Go with GORM ORM framework.

### Key Features
- **Member Portal:** Registration, health tracking, goal management, class enrollment, personal training booking
- **Trainer Portal:** Schedule management, client lookup, progress tracking
- **Admin Portal:** Class creation/management, capacity control, facility operations


### Relationships (6 total)
1. Member → HealthMetric (1:N, total)
2. Member → FitnessGoal (1:N, partial)
3. Member → TrainingSession (1:N)
4. Trainer → TrainingSession (1:N, partial)
5. Trainer → Class (1:N, partial)
6. Member ↔ Class (N:M via ClassEnrollment)

### Advanced Database Features
- **View:** `member_dashboard` - Aggregates member statistics (classes, sessions, goals)
- **Trigger:** `enforce_class_capacity` - Automatically prevents class overbooking
- **Index:** `idx_trainer_date` - Optimizes trainer schedule queries

## Technology Stack
- **Language:** Go 
- **Database:** PostgreSQL 
- **ORM:** GORM 

## Installation & Setup

### Prerequisites
```bash