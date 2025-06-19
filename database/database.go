package database

import (
	"database/sql"
	"time"
	_ "github.com/lib/pq"
)

// ReadDataForDate
// You need to implement reading the required data from the database
// TODO change function return values if needed

type AppointmentData struct {
	CenterID          int
	CenterName        string
	DoctorID          int
	DoctorName        string
	DoctorMobile      string
	PatientName       string
	PatientSalutation string
	TreatmentCategory string
	StartTime         time.Time
	EndTime           time.Time
}

func ReadDataForDate(date string) ([]AppointmentData, error) {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=appointment_db sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT a.CenterID, c.CenterName AS CenterName,
	d.DoctorID, d.Name AS DoctorName, d.Mobile AS DoctorMobile,
	p.Salutation, p.Name AS PatientName, a.TreatmentCategory,
	a.AppointmentStartdttm, a.AppointmentEndddttm
	FROM Appointment a
	JOIN Center c ON a.CenterID = c.CenterID
	JOIN DoctorStaff d ON a.DoctorStaffID = d.DoctorID
	JOIN Patient p ON a.PatientID = p.PatientID
	WHERE DATE(a.AppointmentStartdttm) = $1 AND a.AppointmentStatus = 'S'
	ORDER BY d.DoctorID, a.CenterID, a.AppointmentStartdttm; 
	`

	rows, err := db.Query(query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []AppointmentData
	for rows.Next() {
		var data AppointmentData
		err := rows.Scan(
			&data.CenterID, &data.CenterName,
			&data.DoctorID, &data.DoctorName, &data.DoctorMobile,
			&data.PatientSalutation, &data.PatientName, &data.TreatmentCategory,
			&data.StartTime, &data.EndTime,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}
	return results, nil
}
