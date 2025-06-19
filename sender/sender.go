package sender

import (
	"AppointmentSummmary_Assignment/database"
	"database/sql"
	"fmt"
	"strings"
	"time"
	_ "github.com/lib/pq"
)

type DoctorMessage struct {
	DoctorID    int
	DoctorPhone string
	Message     string
}

type CenterMessage struct {
	CenterID int
	Message  string
}

func CreateAndScheduleSummaryAppointmentMessages(appointments []database.AppointmentData) error {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=appointment_db sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	doctorCenterMap := make(map[string][]database.AppointmentData)
	centerSummary := make(map[int]map[int]int)

	// Group appointments
	for _, appt := range appointments {
		// Doctor-wise
		key := fmt.Sprintf("%d|%d", appt.DoctorID, appt.CenterID)
		doctorCenterMap[key] = append(doctorCenterMap[key], appt)

		// Center-wise
		if _, ok := centerSummary[appt.CenterID]; !ok {
			centerSummary[appt.CenterID] = make(map[int]int)
		}
		centerSummary[appt.CenterID][appt.DoctorID]++
	}

	var doctorMessages []DoctorMessage
	var centerMessages []CenterMessage

	// Doctor-wise summaries
	for _, appts := range doctorCenterMap {
		first := appts[0]
		dateStr := first.StartTime.Format("2nd January, 2006")
		header := fmt.Sprintf("Dr. %s's appointments on <%s> at %s: %d",
			first.DoctorName, dateStr, first.CenterName, len(appts))

		lines := []string{header}
		for _, appt := range appts {
			timeStr := appt.StartTime.Format("3:04 PM")
			durStr := formatDuration(appt.EndTime.Sub(appt.StartTime))
			patientFullName := fmt.Sprintf("%s %s", appt.PatientSalutation, appt.PatientName)
			category := ""
			if appt.TreatmentCategory != "Not Specified" {
				category = fmt.Sprintf(" (%s)", appt.TreatmentCategory)
			}
			lines = append(lines, fmt.Sprintf("%s, %s: %s%s", timeStr, durStr, patientFullName, category))
		}

		message := strings.Join(lines, "\n")
		doctorMessages = append(doctorMessages, DoctorMessage{
			DoctorID:    first.DoctorID,
			DoctorPhone: first.DoctorMobile,
			Message:     message,
		})
	}

	// Center-wise summaries
	for centerID, doctorCounts := range centerSummary {
		var centerName string
		var dateStr string
		total := 0
		lines := []string{}

		// Use any appointment from the center to get name/date
		for _, appt := range appointments {
			if appt.CenterID == centerID {
				centerName = appt.CenterName
				dateStr = appt.StartTime.Format("2nd January, 2006")
				break
			}
		}

		for doctorID, count := range doctorCounts {
			total += count
			var doctorName string
			for _, appt := range appointments {
				if appt.DoctorID == doctorID {
					doctorName = appt.DoctorName
					break
				}
			}
			lines = append(lines, fmt.Sprintf("Dr. %s: %d", doctorName, count))
		}

		header := fmt.Sprintf("Summary of appointments at %s on %s: %d", centerName, dateStr, total)
		message := header + "\n" + strings.Join(lines, "\n")

		centerMessages = append(centerMessages, CenterMessage{
			CenterID: centerID,
			Message:  message,
		})
	}

	// Store to DB
	if err := storeDoctorMessages(db, doctorMessages); err != nil {
		return err
	}
	if err := storeCenterMessages(db, centerMessages); err != nil {
		return err
	}
	return nil
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 && m > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	} else if h > 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dm", m)
}

func storeDoctorMessages(db *sql.DB, messages []DoctorMessage) error {
	for _, msg := range messages {
		_, err := db.Exec(`INSERT INTO DoctorMessages (DoctorID, Phone, Message) VALUES ($1, $2, $3)`,
			msg.DoctorID, msg.DoctorPhone, msg.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func storeCenterMessages(db *sql.DB, messages []CenterMessage) error {
	for _, msg := range messages {
		_, err := db.Exec(`INSERT INTO CenterMessages (CenterID, Message) VALUES ($1, $2)`,
			msg.CenterID, msg.Message)
		if err != nil {
			return err
		}
	}
	return nil
}
