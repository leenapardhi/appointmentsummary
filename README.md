# appointmentsummary
# 🩺 Appointment Summary Generator – Golang + PostgreSQL

This project reads appointment data from a PostgreSQL database and generates daily summary messages for doctors and centers based on appointment information.

---

## 📦 Project Structure

AppointmentSummary_Assignment/
│
├── main.go # Entry point
├── database/
│ └── database.go # Reads appointment data from PostgreSQL
├── sender/
│ └── sender.go # Creates summary messages and stores them
├── data_files.zip/
│ ├── appointment.csv
│ ├── center.csv
│ ├── doctorstaff.csv
│ └── patient.csv

---

## 🛠️ Setup Instructions

### 1. ✅ Prerequisites

- Go 1.18+
- PostgreSQL 13 or later
- pgAdmin or `psql` CLI
- `github.com/lib/pq` Go PostgreSQL driver

---

### 2. 🗃️ Create the Database

Create a PostgreSQL database:

```sql
CREATE DATABASE appointment_db;

3. 🧱 Create Tables
Run the following SQL to create necessary tables:

CREATE TABLE Center (
    CenterID INT PRIMARY KEY,
    Name TEXT
);

CREATE TABLE DoctorStaff (
    DoctorStaffID INT PRIMARY KEY,
    Name TEXT,
    Mobile TEXT
);

CREATE TABLE Patient (
    PatientID INT PRIMARY KEY,
    Salutation TEXT,
    Name TEXT
);

CREATE TABLE Appointment (
    AppointmentID INT PRIMARY KEY,
    CenterID INT,
    DoctorStaffID INT,
    PatientID INT,
    AppointmentStartDttm TIMESTAMP,
    AppointmentEndDttm TIMESTAMP,
    AppointmentStatus TEXT,
    TreatmentCategory TEXT
);

CREATE TABLE DoctorMessages (
    ID SERIAL PRIMARY KEY,
    DoctorID INT,
    Phone TEXT,
    Message TEXT
);

CREATE TABLE CenterMessages (
    ID SERIAL PRIMARY KEY,
    CenterID INT,
    Message TEXT
);

4. 📥 Import CSV Data Using pgAdmin or psql
Use pgAdmin's Import/Export option OR run this inside psql:
\copy Center FROM 'C:/path/center.csv' DELIMITER ',' CSV HEADER;
\copy DoctorStaff FROM 'C:/path/doctorstaff.csv' DELIMITER ',' CSV HEADER;
\copy Patient FROM 'C:/path/patient.csv' DELIMITER ',' CSV HEADER;
\copy Appointment FROM 'C:/path/appointment.csv' DELIMITER ',' CSV HEADER;
Ensure your file paths are correct and use double backslashes (\\) on Windows if needed.

5. ⚙️ Update DB Connection Strings
In database/database.go and sender/sender.go, update:
sql.Open("postgres", "user=postgres password=yourpassword dbname=appointment_db sslmode=disable")
Replace yourpassword with your actual PostgreSQL password.

6. 🚀 Run the Project
Use: go run main.go 2025-05-12
This will:
Read all appointments for that date with status 'S'
Generate summary messages for doctors and centers
Insert them into DoctorMessages and CenterMessages tables

📊 Sample Output (Database)
SELECT * FROM DoctorMessages;
SELECT * FROM CenterMessages;

✅ Features
Fetches and groups appointment data
Generates doctor- and center-level summaries
Inserts generated messages into summary tables
Structured and modular Go code

## 📦 CSV Files
The actual CSV files are compressed and included as `data_files.zip` in the root folder.
To extract:
- Unzip `data_files.zip`
- Place the extracted files into the `CSV Files/` folder

Make sure your PostgreSQL `COPY` or pgAdmin import path points to the unzipped CSV files.


🧪 Test Case Tip
To test quickly, insert sample data manually:

INSERT INTO Appointment VALUES (1, 1, 101, 201, '2025-05-12 10:00:00', '2025-05-12 10:30:00', 'S', 'Dental');
-- Add matching entries in DoctorStaff, Center, and Patient

👩‍💻 Developer Notes
Go modules required: go mod init, go get github.com/lib/pq
Uses time.Time for start and end timestamps
Format messages with proper grouping logic

