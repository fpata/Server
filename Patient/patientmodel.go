package Patient

import "clinic_server/types"

type Patient struct {
	Id                       types.NullInt64  `gorm:"PrimaryKey"`
	FirstName                types.NullString `json:"FirstName"`
	LastName                 types.NullString `json:"LastName"`
	MiddleName               types.NullString `json:"MiddleName"`
	Age                      types.NullInt64  `json:"Age"`
	Gender                   types.NullInt64  `json:"Gender"`
	Role                     types.NullString `json:"Role"`
	PermAddress1             types.NullString `json:"PermAddress1"`
	PermAddress2             types.NullString `json:"PermAddress2"`
	PermCity                 types.NullString `json:"PermCity"`
	PermState                types.NullString `json:"PermState"`
	PermCountry              types.NullString `json:"PermCountry"`
	PermPostalCode           types.NullString `json:"PermPostalCode"`
	CorrAddress1             types.NullString `json:"CorrAddress1"`
	CorrAddress2             types.NullString `json:"CorrAddress2"`
	CorrCity                 types.NullString `json:"CorrCity"`
	CorrState                types.NullString `json:"CorrState"`
	CorrCountry              types.NullString `json:"CorrCountry"`
	CorrPostalCode           types.NullString `json:"CorrPostalCode"`
	PrimaryPhone             types.NullString `json:"PrimaryPhone"`
	PrimaryEmail             types.NullString `json:"PrimaryEmail"`
	SecondaryPhone           types.NullString `json:"SecondaryPhone"`
	SecondaryEmail           types.NullString `json:"SecondaryEmail"`
	EmergencyContactName     types.NullString `json:"EmergencyContactName"`
	EmergencyContactEmail    types.NullString `json:"EmergencyContactEmail"`
	EmergencyContactPhone    types.NullString `json:"EmergencyContactPhone"`
	EmergencyContactRelation types.NullString `json:"EmergencyContactRelation"`
	ExistingDieases          types.NullString `json:"ExistingDieases"`
	Medications              types.NullString `json:"Medications"`
	Allergies                types.NullString `json:"Allergies"`
	FatherMedicalHistory     types.NullString `json:"FatherMedicalHistory"`
	MotherMedicalHistory     types.NullString `json:"MotherMedicalHistory"`
	PatientReports           []*PatientReport
	PatientTreatments        []*PatientTreatment
	PatientAppointments      []*PatientAppointment
}

type PatientReport struct {
	Id            types.NullInt64  `gorm:"PrimaryKey"`
	PatientID     types.NullInt64  `json:"PatientId"`
	ReportDate    types.NullString `json:"ReportDate"`
	ReportName    types.NullString `json:"ReportName"`
	ReportFinding types.NullString `json:"RepoprtFinding"`
	DoctorName    types.NullString `json:"DoctorName"`
}

type PatientTreatment struct {
	Id                      types.NullInt64  `gorm:"PrimaryKey"`
	PatientID               types.NullInt64  `json:"PatientId"`
	ChiefComplaint          types.NullString `json:"ChiefComplaint"`
	Observation             types.NullString `json:"Observation"`
	TreatmentPlan           types.NullString `json:"TreatmentPlan"`
	PatientTreatmentDetails []*PatientTreatmentDetail
}

type PatientTreatmentDetail struct {
	Id                 types.NullInt64  `gorm:"PrimaryKey"`
	PatientID          types.NullInt64  `json:"PatientId"`
	PatientTreatmentID types.NullInt64  `json:"PatientTreatmentID"`
	Tooth              types.NullString `json:"Tooth"`
	Procedure          types.NullString `json:"Procedure"`
	Advice             types.NullString `json:"Advice"`
	TreatmentDate      types.NullString `json:"TreatmentDate"`
}

type PatientAppointment struct {
	Id            types.NullInt64  `gorm:"PrimaryKey"`
	PatientID     types.NullInt64  `json:"PatientId"`
	PatientName   types.NullString `json:"PatientName"`
	ApptDate      types.NullString `json:"ApptDate"`
	ApptTime      types.NullString `json:"ApptTime"`
	TreatmentName types.NullString `json:"TreatmentName"`
	DoctorName    types.NullString `json:"DoctorName"`
	DoctorID      types.NullInt64  `json:"DoctorId"`
}

type SearchResult struct {
	Id           types.NullInt64  `gorm:"PrimaryKey"`
	FirstName    types.NullString `json:"FirstName"`
	LastName     types.NullString `json:"LastName"`
	PrimaryPhone types.NullString `json:"PrimaryPhone"`
	PrimaryEmail types.NullString `json:"PrimaryEmail"`
	PermCity     types.NullString `json:"PermCity"`
}
