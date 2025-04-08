CREATE TABLE users(
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    first_name NVARCHAR(100),
    last_name NVARCHAR(100),
    middle_name NVARCHAR(100),
    age INT,
    gender INT,
	user_type NVARCHAR(100),
	role NVARCHAR(50) NOT NULL,
	status NVARCHAR(20) DEFAULT 'active',
)

-- Login/Users table
CREATE TABLE user_auth (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    username NVARCHAR(100) UNIQUE NOT NULL,
    password NVARCHAR(255) NOT NULL,
    last_login DATETIME,
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
	[user_id] BIGINT,
	FOREIGN KEY ([user_id]) REFERENCES users(id) ON DELETE SET NULL
);
GO

-- New table for addresses
CREATE TABLE addresses (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,	
    address_type NVARCHAR(100),
	address1 NVARCHAR(255),
    address2 NVARCHAR(255),
    city NVARCHAR(100),
    state NVARCHAR(100),
    country NVARCHAR(100),
    postal_code NVARCHAR(20),
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
	[user_id] BIGINT,
	FOREIGN KEY ([user_id]) REFERENCES users(id) ON DELETE SET NULL
);
GO

-- New table for emergency contacts
CREATE TABLE emergency_contacts (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    name NVARCHAR(100),
    email NVARCHAR(100),
    phone NVARCHAR(20),
    relation NVARCHAR(50),
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
	[user_id] BIGINT,
	FOREIGN KEY ([user_id]) REFERENCES users(id) ON DELETE SET NULL
);
GO
-- Update patients table to reference new tables
CREATE TABLE patients (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    patient_type NVARCHAR(50),
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
	[user_id] BIGINT,
	FOREIGN KEY ([user_id]) REFERENCES users(id) ON DELETE SET NULL
);
GO
-- New table for medical histories
CREATE TABLE medical_histories (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    patient_id BIGINT,
    existing_diseases NVARCHAR(MAX),
    medications NVARCHAR(MAX),
    allergies NVARCHAR(MAX),
    father_medical_history NVARCHAR(MAX),
    mother_medical_history NVARCHAR(MAX),
    created_at DATETIME DEFAULT GETDATE(),
    updated_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
GO

CREATE TABLE user_contacts (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    perm_address_id BIGINT,
    corr_address_id BIGINT,
    emergency_contact_id BIGINT,
    primary_phone NVARCHAR(20),
    primary_email NVARCHAR(100),
    secondary_phone NVARCHAR(20),
    secondary_email NVARCHAR(100),
	[user_id] BIGINT,
	FOREIGN KEY ([user_id]) REFERENCES users(id) ON DELETE SET NULL
)



CREATE TABLE patient_reports (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    patient_id BIGINT,
    report_date DATE,
    report_name NVARCHAR(255),
    report_finding NVARCHAR(MAX),
    doctor_name NVARCHAR(100),
    created_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
GO

CREATE TABLE patient_treatments (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    patient_id BIGINT,
    chief_complaint NVARCHAR(MAX),
    observation NVARCHAR(MAX),
    treatment_plan NVARCHAR(MAX),
    created_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
GO

-- Remove redundant columns from patient_treatment_details
CREATE TABLE patient_treatment_details (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    patient_id BIGINT,
    patient_treatment_id BIGINT,
    tooth NVARCHAR(50),
    procedure_name NVARCHAR(MAX),
    advice NVARCHAR(MAX),
    treatment_date DATE,
    created_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE,
    FOREIGN KEY (patient_treatment_id) REFERENCES patient_treatments(id) 
);
GO

CREATE TABLE patient_appointments (
    id BIGINT IDENTITY(1,1) PRIMARY KEY,
    patient_id BIGINT,
    patient_name NVARCHAR(200),
    appt_date DATE,
    appt_time TIME,
    treatment_name NVARCHAR(255),
    doctor_name NVARCHAR(100),
    doctor_id BIGINT,
    created_at DATETIME DEFAULT GETDATE(),
    modified_by NVARCHAR(100),
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE
);
GO

-- Create indexes for better query performance
CREATE INDEX idx_user_contacts_email ON user_contacts(primary_email);
CREATE INDEX idx_user_contacts_phone ON user_contacts(primary_phone);
CREATE INDEX idx_appointment_date ON patient_appointments(appt_date);
CREATE INDEX idx_treatment_date ON patient_treatment_details(treatment_date);
GO

