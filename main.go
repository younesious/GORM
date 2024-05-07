package main

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func main() {

}

func RefreshDatabase(db *gorm.DB, tables []interface{}) (string, error) {
	for _, table := range tables {
		err := db.Migrator().DropTable(table)
		if err != nil {
			return "", err
		}
		err = db.AutoMigrate(table)
		if err != nil {
			return "", err
		}
	}
	return "Refresh database successfully done", nil
}

func SeedUser(db *gorm.DB, username, firstName, lastName, calendarName, appointmentSubject string, startDate time.Time) (string, error) {
	user := User{
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Calendar: Calendar{
			Name: calendarName,
			Appointments: []Appointment{
				{Subject: appointmentSubject, StartTime: startDate},
			},
		},
	}

	err := db.Create(&user).Error
	if err != nil {
		return "", err
	}

	return "Seeding database successfully done", nil
}

func UserWithRangeAppointment(startDate, endDate time.Time, subject string, calendarTableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Calendar.Appointments.Attendees").
			Joins("JOIN calendars ON users.id = calendars.user_id").
			Joins("JOIN appointments ON calendars.id = appointments.owner_id AND appointments.owner_type = 'calendars'").
			Joins("JOIN appointment_user ON appointments.id = appointment_user.appointment_id").
			Where("appointments.start_time BETWEEN ? AND ?", startDate, endDate).
			Where("appointments.subject = ?", subject).
			Where("calendars.name LIKE ?", fmt.Sprintf("%%%s%%", calendarTableName)).
			Distinct("username")
	}
}

func updateAppointment(db *gorm.DB, calendarName string, startTime time.Time, endTime time.Time, keyword string) error {
	var appointments []Appointment
	if err := db.Joins("JOIN calendars ON appointments.owner_id = calendars.id  AND appointments.owner_type = 'calendars'").
		Where("calendars.name = ? AND start_time BETWEEN ? AND ? AND subject LIKE ?", calendarName, startTime, endTime, fmt.Sprintf("%%%s%%", keyword)).
		Find(&appointments).Error; err != nil {
		return err
	}

	for i := range appointments {
		appointments[i].StartTime = appointments[i].StartTime.Add(time.Hour)
		appointments[i].Description = appointments[i].Subject + " event"
		if err := db.Save(&appointments[i]).Error; err != nil {
			return err
		}
	}

	return nil
}


type User struct {
	gorm.Model
	Username  string
	FirstName string
	LastName  string
	Calendar  Calendar
}

type Calendar struct {
	gorm.Model
	Name         string
	UserID       uint
	Appointments []Appointment `gorm:"polymorphic:Owner;"`
}

type Appointment struct {
	gorm.Model
	Subject     string
	Description string
	StartTime   time.Time
	Length      uint
	OwnerID     uint
	OwnerType   string
	Attendees   []User `gorm:"many2many:appointment_user;"`
}

type TaskList struct {
	gorm.Model
	Appointments []Appointment `gorm:"polymorphic:Owner;"`
}
