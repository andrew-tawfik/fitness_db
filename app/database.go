package main

import (
	"fmt"
	"log"

	"fitness_db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := "host=localhost user=postgres password=archSQL dbname=fitness_club port=5433 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully")

	// Auto-migrate (creates tables)
	err = DB.AutoMigrate(
		&models.Member{},
		&models.Trainer{},
		&models.Class{},
		&models.Admin{},
		&models.HealthMetric{},
		&models.FitnessGoal{},
		&models.TrainingSession{},
		&models.ClassEnrollment{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully")
}

// CreateViews creates database views for common query patterns.
// This view provides a summary dashboard for each member, getting their
// class enrollments, training sessions, and active fitness goals.
func CreateViews() {
	err := DB.Exec(`
		create or replace view member_dashboard as
		select 
			m.member_id,
			m.first_name,
			m.last_name,
			m.email,
			count(distinct ce.enrollment_id) as total_classes,
			count(distinct ts.session_id) as total_sessions,
			(select count(*) 
			 from fitness_goals fg 
			 where fg.member_id = m.member_id and fg.status = 'active') as active_goals
		from members m
		left join class_enrollments ce on m.member_id = ce.member_id
		left join training_sessions ts on m.member_id = ts.member_id
		group by m.member_id, m.first_name, m.last_name, m.email
	`).Error

	if err != nil {
		log.Printf("Warning: Failed to create member_dashboard view: %v", err)
	} else {
		fmt.Println("Views created successfully")
	}
}

// CreateTriggers creates database triggers for automatic constraint enforcement.
// This trigger prevents class overbooking by checking capacity before enrollment
// and automatically updating the current_enrollment counter when a member enrolls.
// Triggers ensure data integrity at the database level.
func CreateTriggers() {
	// Function to check class capacity and update enrollment count
	// Raises an exception if class is at full capacity
	err := DB.Exec(`
		create or replace function check_class_capacity()
		returns trigger as $$
		begin
			if (select current_enrollment from classes where class_id = new.class_id) 
			   >= (select capacity from classes where class_id = new.class_id) then
				raise exception 'Class is full - cannot enroll more members';
			end if;
			
			update classes 
			set current_enrollment = current_enrollment + 1 
			where class_id = new.class_id;
			
			return new;
		end;
		$$ language plpgsql;
	`).Error

	if err != nil {
		log.Printf("Warning: Failed to create check_class_capacity function: %v", err)
		return
	}

	// Trigger that fires before each insert into class_enrollments
	// Automatically enforces capacity limits 
	err = DB.Exec(`
		drop trigger if exists enforce_class_capacity on class_enrollments;
		create trigger enforce_class_capacity
		before insert on class_enrollments
		for each row execute function check_class_capacity();
	`).Error

	if err != nil {
		log.Printf("Warning: Failed to create enforce_class_capacity trigger: %v", err)
	} else {
		fmt.Println("Triggers created successfully")
	}
}