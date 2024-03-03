package main

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

func SetUpDatabase() error {
	db = GetConnection()

	// Check the connection
	_, err := db.DB()
	if err != nil {
		return err
	}

	return nil
}

func GetUserWithCalendar() User {
	return User{
		Username:  "Younesious",
		FirstName: "Younes",
		LastName:  "Mahmoudi",
		Calendar: Calendar{
			Name: "QCalendar",
		},
	}
}

func GetUser(username string) User {
	return User{
		Username:  username,
		FirstName: "Younes",
		LastName:  "Mahmoudi",
	}
}

func GetAppointmentWithCalender(calender Calendar) Appointment {
	return Appointment{
		Subject:     "Meeting",
		Description: "Discuss project progress",
		StartTime:   time.Now(),
		Length:      60,
		OwnerID:     calender.ID,
		OwnerType:   "Calendar",
	}
}

func GetAppointmentWithTaskList(tasklist TaskList) Appointment {
	return Appointment{
		Subject:     "Meeting",
		Description: "Discuss project progress",
		StartTime:   time.Now(),
		Length:      60,
		OwnerID:     tasklist.ID,
		OwnerType:   "TaskList",
	}
}

func SeedMyUser(db *gorm.DB, username, firstName, lastName, calendarName, appointmentSubject string, startDate time.Time) (string, error) {
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

func TestDatabaseConnectionAndAutoMigrate(t *testing.T) {
	err := SetUpDatabase()
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	db = GetConnection()
	err = db.Migrator().AutoMigrate(&User{}, &Calendar{}, &Appointment{}, &TaskList{})
	if err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}
}

func TestUserTable(t *testing.T) {
	expectedFields := map[string]reflect.Kind{
		"Model":     reflect.Struct,
		"Username":  reflect.String,
		"FirstName": reflect.String,
		"LastName":  reflect.String,
		"Calendar":  reflect.Struct,
	}

	user := User{}
	userType := reflect.TypeOf(user)

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldName := field.Name
		fieldType := field.Type.Kind()

		expectedType, ok := expectedFields[fieldName]
		if !ok {
			t.Errorf("Unexpected field name: %s", fieldName)
		}
		if fieldType != expectedType {
			t.Errorf("Field %s has unexpected type: %s", fieldName, fieldType)
		}
	}
}

func TestCalendarTable(t *testing.T) {
	expectedFields := map[string]reflect.Kind{
		"Model":        reflect.Struct,
		"Name":         reflect.String,
		"UserID":       reflect.Uint,
		"Appointments": reflect.Slice,
	}

	calendar := Calendar{}
	userType := reflect.TypeOf(calendar)

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldName := field.Name
		fieldType := field.Type.Kind()
		expectedType, ok := expectedFields[fieldName]
		if !ok {
			t.Errorf("Unexpected field name: %s", fieldName)
		}
		if fieldType != expectedType {
			t.Errorf("Field %s has unexpected type: %s", fieldName, fieldType)
		}
	}
}

func TestAppointmentTable(t *testing.T) {
	expectedFields := map[string]reflect.Kind{
		"Model":       reflect.Struct,
		"Subject":     reflect.String,
		"Description": reflect.String,
		"StartTime":   reflect.Struct,
		"Length":      reflect.Uint,
		"OwnerID":     reflect.Uint,
		"OwnerType":   reflect.String,
		"Attendees":   reflect.Slice,
	}

	appointment := Appointment{}
	userType := reflect.TypeOf(appointment)

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldName := field.Name
		fieldType := field.Type.Kind()
		expectedType, ok := expectedFields[fieldName]
		if !ok {
			t.Errorf("Unexpected field name: %s", fieldName)
		}
		if fieldType != expectedType {
			t.Errorf("Field %s has unexpected type: %s", fieldName, fieldType)
		}
	}
}

func TestTaskListTable(t *testing.T) {
	expectedFields := map[string]reflect.Kind{
		"Model":        reflect.Struct,
		"Appointments": reflect.Slice,
	}

	taskList := TaskList{}
	userType := reflect.TypeOf(taskList)

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldName := field.Name
		fieldType := field.Type.Kind()
		expectedType, ok := expectedFields[fieldName]
		if !ok {
			t.Errorf("Unexpected field name: %s", fieldName)
		}
		if fieldType != expectedType {
			t.Errorf("Field %s has unexpected type: %s", fieldName, fieldType)
		}
	}
}

func TestUserCalendarRelationship(t *testing.T) {
	user := GetUserWithCalendar()
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	var retrievedUser User
	if err := db.Preload("Calendar").First(&retrievedUser, user.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if retrievedUser.Calendar.ID != user.Calendar.ID {
		t.Errorf("User's calendar ID doesn't match expected value")
	}
	if retrievedUser.Calendar.Name != user.Calendar.Name {
		t.Errorf("User's calendar name doesn't match expected value")
	}
}

func TestCalendarAppointmentRelationship(t *testing.T) {
	user := GetUserWithCalendar()
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	calendar := Calendar{
		Name:   "QCalendar",
		UserID: user.ID,
	}
	if err := db.Create(&calendar).Error; err != nil {
		t.Fatalf("Failed to create calendar: %v", err)
	}
	appointment := GetAppointmentWithCalender(calendar)
	if err := db.Create(&appointment).Error; err != nil {
		t.Fatalf("Failed to create appointment: %v", err)
	}

	var retrievedAppointment Appointment
	if err := db.First(&retrievedAppointment, appointment.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve appointment: %v", err)
	}
	if retrievedAppointment.OwnerID != calendar.ID || retrievedAppointment.OwnerType != "Calendar" {
		t.Errorf("Appointment doesn't belong to the correct calendar")
	}

	var retrievedCalendar Calendar
	if err := db.Preload("Appointments").First(&retrievedCalendar, user.Calendar.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve calendar: %v", err)
	}

}

func TestTaskListAppointmentRelationship(t *testing.T) {
	taskList := TaskList{}
	if err := db.Create(&taskList).Error; err != nil {
		t.Fatalf("Failed to create task list: %v", err)
	}

	appointment := GetAppointmentWithTaskList(taskList)
	if err := db.Create(&appointment).Error; err != nil {
		t.Fatalf("Failed to create appointment: %v", err)
	}

	var retrievedAppointment Appointment
	if err := db.First(&retrievedAppointment, appointment.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve appointment: %v", err)
	}
	if retrievedAppointment.OwnerID != taskList.ID || retrievedAppointment.OwnerType != "TaskList" {
		t.Errorf("Appointment doesn't belong to the correct calendar")
	}

}

func TestUserAppointmentRelationship(t *testing.T) {
	user1 := GetUser("younesious")
	user2 := GetUser("roozbehious")
	db.Create(&user1)
	db.Create(&user2)

	appointment := Appointment{
		Subject:     "Team Meeting",
		Description: "Discuss progress on project X",
		StartTime:   time.Now(),
		Length:      60,
		OwnerID:     user1.ID,
		OwnerType:   "User",
	}
	db.Create(&appointment)

	appointment.Attendees = []User{user2}
	db.Save(&appointment)

	var retrievedAppointment Appointment
	if err := db.Preload("Attendees").First(&retrievedAppointment, appointment.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve appointment: %v", err)
	}
	if len(retrievedAppointment.Attendees) != 1 {
		t.Errorf("Expected 1 attendee, but got %d", len(retrievedAppointment.Attendees))
	}
	if retrievedAppointment.Attendees[0].Username != "roozbehious" {
		t.Errorf("Expected attendee username to be 'roozbehious', but got '%s'", retrievedAppointment.Attendees[0].Username)
	}

	var retrievedUser1 User
	if err := db.Preload("Calendar.Appointments").First(&retrievedUser1, user1.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve user1: %v", err)
	}

	var retrievedUser2 User
	if err := db.Preload("Calendar.Appointments").First(&retrievedUser2, user2.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve user2: %v", err)
	}
}

func TestRefreshDatabase(t *testing.T) {
	tables := []interface{}{
		&User{},
		&Calendar{},
		&Appointment{},
		&TaskList{},
	}

	message, err := RefreshDatabase(db, tables)
	if err != nil {
		t.Fatalf("Failed to refresh database: %v", err)
	}

	if message != "Refresh database successfully done" {
		t.Errorf("Expected message %q, but got %q", "Refresh database successfully done", message)
	}

	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("Table %T does not exist in database", table)
		}
	}

	var userCount int64
	if err := db.Model(&User{}).Count(&userCount).Error; err != nil {
		t.Fatalf("Failed to count users: %v", err)
	}
	if userCount != 0 {
		t.Errorf("Expected users table to be empty, but it has %d rows", userCount)
	}

	var calendarCount int64
	if err := db.Model(&Calendar{}).Count(&calendarCount).Error; err != nil {
		t.Fatalf("Failed to count calendars: %v", err)
	}
	if calendarCount != 0 {
		t.Errorf("Expected calendars table to be empty, but it has %d rows", calendarCount)
	}

	var appointmentCount int64
	if err := db.Model(&Appointment{}).Count(&appointmentCount).Error; err != nil {
		t.Fatalf("Failed to count appointments: %v", err)
	}
	if appointmentCount != 0 {
		t.Errorf("Expected appointments table to be empty, but it has %d rows", appointmentCount)
	}

	var taskListCount int64
	if err := db.Model(&TaskList{}).Count(&taskListCount).Error; err != nil {
		t.Fatalf("Failed to count task lists: %v", err)
	}
	if taskListCount != 0 {
		t.Errorf("Expected task lists table to be empty, but it has %d rows", taskListCount)
	}
}

func TestSeedUser(t *testing.T) {

	username := "Younesious"
	firstName := "Younes"
	lastName := "Mahmoudi"
	calendarName := "QCalendar"
	appointmentSubject := "HamkaranSystem"
	startDate := time.Now()

	_, err := SeedUser(db, username, firstName, lastName, calendarName, appointmentSubject, startDate)
	if err != nil {
		t.Fatalf("Failed to seed database: %v", err)

	}

	var retrievedUser User
	if err := db.Preload("Calendar").Preload("Calendar.Appointments").
		First(&retrievedUser, "username = ?", username).Error; err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if retrievedUser.Username != username || retrievedUser.FirstName != firstName || retrievedUser.LastName != lastName {
		t.Errorf("User is not correctly created")
	}

	if retrievedUser.Calendar.Name != calendarName || len(retrievedUser.Calendar.Appointments) != 1 ||
		retrievedUser.Calendar.Appointments[0].Subject != appointmentSubject ||
		retrievedUser.Calendar.Appointments[0].StartTime.Unix() != startDate.Unix() {
		t.Errorf("Calendar and appointment are not correctly created")
	}
}

func TestUserWithRangeAppointmentScope(t *testing.T) {
	tables := []interface{}{&User{}, &Calendar{}, &Appointment{}, &TaskList{}}
	for _, table := range tables {
		err := db.Migrator().DropTable(table)
		if err != nil {
			t.Fatalf("Failed to refresh database: %v", err)
		}
		err = db.AutoMigrate(table)
		if err != nil {
			t.Fatalf("Failed to refresh database: %v", err)
		}
	}

	if _, err := SeedMyUser(db, "Younesious", "Younes", "Mahmoudi",
		"Contest events", "HamkaranSystem",
		time.Date(2023, time.May, 2, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}

	if _, err := SeedMyUser(db, "Roozbehiano", "Roozbeh", "SharifN",
		"Contest events", "HamkaranSystem",
		time.Date(2023, time.May, 25, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}
	if _, err := SeedMyUser(db, "Matiniano", "Matin", "Moeenie",
		"Contest events", "HamkaranSystem",
		time.Date(2023, time.April, 2, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}
	if _, err := SeedMyUser(db, "Moieenious", "Alice", "Jones",
		"My Calendar", "HamkaranSystem",
		time.Date(2023, time.April, 2, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}
	if _, err := SeedMyUser(db, "Ali", "Ali", "Jones",
		"Contest events", "quera",
		time.Date(2023, time.April, 2, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}

	startTime := time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, time.May, 31, 0, 0, 0, 0, time.UTC)
	subject := "HamkaranSystem"
	calendarTableName := "Contest events"

	var users []User
	if err := db.Scopes(UserWithRangeAppointment(startTime, endTime, subject, calendarTableName)).Find(&users).Error; err != nil {
		t.Fatalf("Failed to retrieve users with range appointments: %v", err)
	}

	if len(users) != 2 || users[1].Username != "Younesious" {
		t.Errorf("Expected 1 user but got %d", len(users))
	}
	if len(users) != 2 || users[0].Username != "Roozbehiano" {
		t.Errorf("Expected 1 user but got %d", len(users))
	}
}

func TestUpdateAppointment(t *testing.T) {
	startTime := time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, time.May, 31, 0, 0, 0, 0, time.UTC)

	_, err := SeedMyUser(db, "Younesious", "Younes", "Mahmoudi",
		"test_calendar", "test_appointment", startTime)
	if err != nil {
		t.Errorf("Failed to seed the database: %v", err)
		return
	}

	calendarName := "test_calendar"
	keyword := "test_appointment"
	err = updateAppointment(db, calendarName, startTime, endTime, keyword)
	if err != nil {
		t.Errorf("Failed to update the appointment: %v", err)
		return
	}

	var appointment Appointment
	if err := db.Where("subject = ?", "test_appointment").First(&appointment).Error; err != nil {
		t.Errorf("Failed to find the updated appointment: %v", err)
		return
	}
	if !appointment.StartTime.Equal(startTime.Add(time.Hour)) {
		t.Errorf("Appointment start time not updated correctly")
		return
	}
	if appointment.Description != "test_appointment event" {
		t.Errorf("Appointment description not updated correctly")
		return
	}

	TearDown()
}

func TearDown() {
	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}(sqlDB)
}
