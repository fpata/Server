
-- Insert dummy data into users table
INSERT INTO users (first_name, last_name, middle_name, age, gender, user_type, role, status)
VALUES 
('John', 'Doe', 'A', 30, 1, 'admin', 'Admin', 'active'),
('Jane', 'Smith', 'B', 25, 2, 'staff', 'Nurse', 'active'),
('Alice', 'Brown', 'C', 40, 2, 'patient', 'Patient', 'inactive');

-- Insert dummy data into user_auth table
INSERT INTO user_auth (username, password, last_login, modified_by, user_id)
VALUES 
('johndoe', 'hashed_password1', NULL, 'System', 1),
('janesmith', 'hashed_password2', NULL, 'System', 2),
('alicebrown', 'hashed_password3', NULL, 'System', 3);

-- Insert dummy data into addresses table
INSERT INTO addresses (address_type, address1, address2, city, state, country, postal_code, modified_by, user_id)
VALUES 
('Permanent', '123 Main St', 'Apt 1', 'New York', 'NY', 'USA', '10001', 'System', 1),
('Correspondence', '456 Elm St', NULL, 'Los Angeles', 'CA', 'USA', '90001', 'System', 2),
('Permanent', '789 Oak St', 'Suite 5', 'Chicago', 'IL', 'USA', '60601', 'System', 3);

-- Insert dummy data into emergency_contacts table
INSERT INTO emergency_contacts (name, email, phone, relation, modified_by, user_id)
VALUES 
('Mary Doe', 'mary.doe@example.com', '1234567890', 'Spouse', 'System', 1),
('Robert Smith', 'robert.smith@example.com', '0987654321', 'Parent', 'System', 2),
('Emma Brown', 'emma.brown@example.com', '1122334455', 'Sibling', 'System', 3);

-- Insert dummy data into patients table
INSERT INTO patients (patient_type, modified_by, user_id)
VALUES 
('Outpatient', 'System', 3),
('Inpatient', 'System', NULL);

-- Insert dummy data into medical_histories table
INSERT INTO medical_histories (patient_id, existing_diseases, medications, allergies, father_medical_history, mother_medical_history, modified_by)
VALUES 
(1, 'Diabetes', 'Metformin', 'Peanuts', 'Heart Disease', 'Hypertension', 'System'),
(2, 'Hypertension', 'Amlodipine', 'None', 'None', 'Diabetes', 'System');

-- Insert dummy data into user_contacts table
INSERT INTO user_contacts (perm_address_id, corr_address_id, emergency_contact_id, primary_phone, primary_email, secondary_phone, secondary_email, user_id)
VALUES 
(1, 2, 1, '1234567890', 'johndoe@example.com', '0987654321', 'john.alt@example.com', 1),
(2, NULL, 2, '2233445566', 'janesmith@example.com', NULL, NULL, 2);

-- Insert dummy data into patient_reports table
INSERT INTO patient_reports (patient_id, report_date, report_name, report_finding, doctor_name, modified_by)
VALUES 
(1, '2023-01-01', 'Blood Test', 'Normal', 'Dr. Adams', 'System'),
(2, '2023-02-01', 'X-Ray', 'Fracture', 'Dr. Baker', 'System');

-- Insert dummy data into patient_treatments table
INSERT INTO patient_treatments (patient_id, chief_complaint, observation, treatment_plan, modified_by)
VALUES 
(1, 'Toothache', 'Cavity in molar', 'Filling', 'System'),
(2, 'Back Pain', 'Muscle strain', 'Physical Therapy', 'System');

-- Insert dummy data into patient_treatment_details table
INSERT INTO patient_treatment_details (patient_id, patient_treatment_id, tooth, procedure_name, advice, treatment_date, modified_by)
VALUES 
(1, 1, 'Molar', 'Filling', 'Avoid hard foods', '2023-01-15', 'System'),
(2, 2, NULL, 'Physical Therapy', 'Stretch daily', '2023-02-10', 'System');

-- Insert dummy data into patient_appointments table
INSERT INTO patient_appointments (patient_id, patient_name, appt_date, appt_time, treatment_name, doctor_name, doctor_id, modified_by)
VALUES 
(1, 'Alice Brown', '2023-03-01', '10:00:00', 'Dental Checkup', 'Dr. Adams', 1, 'System'),
(2, 'John Doe', '2023-03-02', '14:00:00', 'Physiotherapy', 'Dr. Baker', 2, 'System');
