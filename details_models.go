// Code generated from the observed 2026-10-06 registry response schemas.
// Values seen only as null or empty arrays remain json.RawMessage.
package fgis

import "encoding/json"

// CertificateDetails contains fields observed in public registry detail JSON.
type CertificateDetails struct {
	Annexes                   []CertificateDetailsAnnexesItem             `json:"annexes"`
	Applicant                 CertificateDetailsApplicant                 `json:"applicant"`
	ApplicantFilials          []json.RawMessage                           `json:"applicantFilials"`
	AssignRegNumber           bool                                        `json:"assignRegNumber"`
	AwaitForApprove           bool                                        `json:"awaitForApprove"`
	AwaitOperatorCheck        json.RawMessage                             `json:"awaitOperatorCheck"`
	BatchInspection           *bool                                       `json:"batchInspection"`
	BlankNumber               string                                      `json:"blankNumber"`
	CertEndDate               string                                      `json:"certEndDate"`
	CertRegDate               string                                      `json:"certRegDate"`
	CertificationAuthority    CertificateDetailsCertificationAuthority    `json:"certificationAuthority"`
	Changes                   json.RawMessage                             `json:"changes"`
	Documents                 CertificateDetailsDocuments                 `json:"documents"`
	EditApp                   bool                                        `json:"editApp"`
	Experts                   []CertificateDetailsExpertsItem             `json:"experts"`
	ExpiredInspectionControl  bool                                        `json:"expiredInspectionControl"`
	FirstName                 string                                      `json:"firstName"`
	IDBlank                   int64                                       `json:"idBlank"`
	IDCertBasis               json.RawMessage                             `json:"idCertBasis"`
	IDCertScheme              int64                                       `json:"idCertScheme"`
	IDCertType                int64                                       `json:"idCertType"`
	IDCertificate             int64                                       `json:"idCertificate"`
	IDDeclarationChecker      json.RawMessage                             `json:"idDeclarationChecker"`
	IDEmployee                int64                                       `json:"idEmployee"`
	IDFGISRA1                 json.RawMessage                             `json:"idFgisRaV1"`
	IDMRPA                    json.RawMessage                             `json:"idMrpa"`
	IDObjectCertType          int64                                       `json:"idObjectCertType"`
	IDProductSingleLists      []json.RawMessage                           `json:"idProductSingleLists"`
	IDSigner                  int64                                       `json:"idSigner"`
	IDStatus                  int64                                       `json:"idStatus"`
	IDTechnicalReglaments     []int64                                     `json:"idTechnicalReglaments"`
	InspectionControlPlanDate *string                                     `json:"inspectionControlPlanDate"`
	IsBlock                   bool                                        `json:"isBlock"`
	IsCanPublish              bool                                        `json:"isCanPublish"`
	IsTS                      bool                                        `json:"isTs"`
	ManufIsApplicant          bool                                        `json:"manufIsApplicant"`
	Manufacturer              CertificateDetailsManufacturer              `json:"manufacturer"`
	ManufacturerFilials       []CertificateDetailsManufacturerFilialsItem `json:"manufacturerFilials"`
	NoSanction                bool                                        `json:"noSanction"`
	Number                    string                                      `json:"number"`
	Patronymic                string                                      `json:"patronymic"`
	Product                   CertificateDetailsProduct                   `json:"product"`
	ProductGroups             []CertificateDetailsProductGroupsItem       `json:"productGroups"`
	ProductInfoOON106         json.RawMessage                             `json:"productInfoOon106"`
	ProductMatchOO106         json.RawMessage                             `json:"productMatchOON106"`
	ReplacedBy                *CertificateDetailsReplacedBy               `json:"replacedBy"`
	Replacement               *CertificateDetailsReplacement              `json:"replacement"`
	ShowInOpenPart            bool                                        `json:"showInOpenPart"`
	StatusChanges             []CertificateDetailsStatusChangesItem       `json:"statusChanges"`
	Surname                   string                                      `json:"surname"`
	TestingLabs               []CertificateDetailsTestingLabsItem         `json:"testingLabs"`
	Raw                       json.RawMessage                             `json:"-"`
}

// CertificateDetailsAnnexesItem contains fields observed in public registry detail JSON.
type CertificateDetailsAnnexesItem struct {
	AnnexBlanks []CertificateDetailsAnnexesItemAnnexBlanksItem `json:"annexBlanks"`
	IDAnnex     int64                                          `json:"idAnnex"`
	IDType      int64                                          `json:"idType"`
	Ord         int64                                          `json:"ord"`
	PageCount   *int64                                         `json:"pageCount"`
}

// CertificateDetailsAnnexesItemAnnexBlanksItem contains fields observed in public registry detail JSON.
type CertificateDetailsAnnexesItemAnnexBlanksItem struct {
	BlankNumber string `json:"blankNumber"`
	IDBlank     int64  `json:"idBlank"`
}

// CertificateDetailsApplicant contains fields observed in public registry detail JSON.
type CertificateDetailsApplicant struct {
	AddlRegInfo        *string                                    `json:"addlRegInfo"`
	Addresses          []CertificateDetailsApplicantAddressesItem `json:"addresses"`
	Contacts           []CertificateDetailsApplicantContactsItem  `json:"contacts"`
	FirstName          string                                     `json:"firstName"`
	FullName           string                                     `json:"fullName"`
	HeadPosition       string                                     `json:"headPosition"`
	IDApplicantType    int64                                      `json:"idApplicantType"`
	IDEGRUL            json.RawMessage                            `json:"idEgrul"`
	IDLegalForm        int64                                      `json:"idLegalForm"`
	IDLegalSubject     int64                                      `json:"idLegalSubject"`
	IDLegalSubjectType int64                                      `json:"idLegalSubjectType"`
	IDPerson           int64                                      `json:"idPerson"`
	IDPersonDoc        json.RawMessage                            `json:"idPersonDoc"`
	INN                string                                     `json:"inn"`
	IsEecRegister      bool                                       `json:"isEecRegister"`
	KPP                string                                     `json:"kpp"`
	OGRN               string                                     `json:"ogrn"`
	OGRNAssignDate     string                                     `json:"ogrnAssignDate"`
	PassportIssueDate  json.RawMessage                            `json:"passportIssueDate"`
	PassportIssuedBy   json.RawMessage                            `json:"passportIssuedBy"`
	PassportNum        json.RawMessage                            `json:"passportNum"`
	Patronymic         string                                     `json:"patronymic"`
	RegDate            string                                     `json:"regDate"`
	RegOrganName       string                                     `json:"regOrganName"`
	ShortName          string                                     `json:"shortName"`
	Surname            string                                     `json:"surname"`
	Transnational      []json.RawMessage                          `json:"transnational"`
}

// CertificateDetailsApplicantAddressesItem contains fields observed in public registry detail JSON.
type CertificateDetailsApplicantAddressesItem struct {
	Flat            json.RawMessage `json:"flat"`
	ForeignCity     json.RawMessage `json:"foreignCity"`
	ForeignDistrict json.RawMessage `json:"foreignDistrict"`
	ForeignHouse    json.RawMessage `json:"foreignHouse"`
	ForeignLocality json.RawMessage `json:"foreignLocality"`
	ForeignStreet   json.RawMessage `json:"foreignStreet"`
	FullAddress     string          `json:"fullAddress"`
	Gln             json.RawMessage `json:"gln"`
	Glonass         json.RawMessage `json:"glonass"`
	IDAddrType      int64           `json:"idAddrType"`
	IDAddress       int64           `json:"idAddress"`
	IDCity          json.RawMessage `json:"idCity"`
	IDCodeOksm      string          `json:"idCodeOksm"`
	IDDistrict      json.RawMessage `json:"idDistrict"`
	IDHouse         json.RawMessage `json:"idHouse"`
	IDLocality      json.RawMessage `json:"idLocality"`
	IDStreet        json.RawMessage `json:"idStreet"`
	IDSubject       json.RawMessage `json:"idSubject"`
	OksmShort       bool            `json:"oksmShort"`
	OtherGln        json.RawMessage `json:"otherGln"`
	PostCode        *string         `json:"postCode"`
	UniqueAddress   json.RawMessage `json:"uniqueAddress"`
}

// CertificateDetailsApplicantContactsItem contains fields observed in public registry detail JSON.
type CertificateDetailsApplicantContactsItem struct {
	IDContact     int64  `json:"idContact"`
	IDContactType int64  `json:"idContactType"`
	Value         string `json:"value"`
}

// CertificateDetailsCertificationAuthority contains fields observed in public registry detail JSON.
type CertificateDetailsCertificationAuthority struct {
	AccredOrgName            string                                                  `json:"accredOrgName"`
	Addresses                []CertificateDetailsCertificationAuthorityAddressesItem `json:"addresses"`
	AttestatEndDate          json.RawMessage                                         `json:"attestatEndDate"`
	AttestatRegDate          string                                                  `json:"attestatRegDate"`
	AttestatRegNumber        string                                                  `json:"attestatRegNumber"`
	Contacts                 []CertificateDetailsCertificationAuthorityContactsItem  `json:"contacts"`
	FirstName                string                                                  `json:"firstName"`
	FullName                 string                                                  `json:"fullName"`
	IDCertificationAuthority int64                                                   `json:"idCertificationAuthority"`
	IDPerson                 int64                                                   `json:"idPerson"`
	IDRAL                    int64                                                   `json:"idRal"`
	OGRN                     string                                                  `json:"ogrn"`
	Patronymic               string                                                  `json:"patronymic"`
	PrevAttestatRegNumber    json.RawMessage                                         `json:"prevAttestatRegNumber"`
	PrevIDRAL                json.RawMessage                                         `json:"prevIdRal"`
	Surname                  string                                                  `json:"surname"`
}

// CertificateDetailsCertificationAuthorityAddressesItem contains fields observed in public registry detail JSON.
type CertificateDetailsCertificationAuthorityAddressesItem struct {
	Flat            *string         `json:"flat"`
	ForeignCity     json.RawMessage `json:"foreignCity"`
	ForeignDistrict json.RawMessage `json:"foreignDistrict"`
	ForeignHouse    json.RawMessage `json:"foreignHouse"`
	ForeignLocality json.RawMessage `json:"foreignLocality"`
	ForeignStreet   json.RawMessage `json:"foreignStreet"`
	FullAddress     string          `json:"fullAddress"`
	Gln             json.RawMessage `json:"gln"`
	Glonass         json.RawMessage `json:"glonass"`
	IDAddrType      int64           `json:"idAddrType"`
	IDAddress       int64           `json:"idAddress"`
	IDCity          json.RawMessage `json:"idCity"`
	IDCodeOksm      *string         `json:"idCodeOksm"`
	IDDistrict      json.RawMessage `json:"idDistrict"`
	IDHouse         *string         `json:"idHouse"`
	IDLocality      json.RawMessage `json:"idLocality"`
	IDStreet        *string         `json:"idStreet"`
	IDSubject       *string         `json:"idSubject"`
	OksmShort       bool            `json:"oksmShort"`
	OtherGln        json.RawMessage `json:"otherGln"`
	PostCode        *string         `json:"postCode"`
	UniqueAddress   *string         `json:"uniqueAddress"`
}

// CertificateDetailsCertificationAuthorityContactsItem contains fields observed in public registry detail JSON.
type CertificateDetailsCertificationAuthorityContactsItem struct {
	IDContact     int64  `json:"idContact"`
	IDContactType int64  `json:"idContactType"`
	Value         string `json:"value"`
}

// CertificateDetailsDocuments contains fields observed in public registry detail JSON.
type CertificateDetailsDocuments struct {
	ApplicantOtherDocuments                []CertificateDetailsDocumentsApplicantOtherDocumentsItem                `json:"applicantOtherDocuments"`
	ApplicationCertificationDecision       CertificateDetailsDocumentsApplicationCertificationDecision             `json:"applicationCertificationDecision"`
	CommonDocuments                        map[string]json.RawMessage                                              `json:"commonDocuments"`
	ConformityAssessmentDocuments          []json.RawMessage                                                       `json:"conformityAssessmentDocuments"`
	ExpertConclusion                       CertificateDetailsDocumentsExpertConclusion                             `json:"expertConclusion"`
	ForeignManufacturerContract            CertificateDetailsDocumentsForeignManufacturerContract                  `json:"foreignManufacturerContract"`
	ManufacturingConditionAnalysisActs     []CertificateDetailsDocumentsManufacturingConditionAnalysisActsItem     `json:"manufacturingConditionAnalysisActs"`
	ProductCertificationContract           CertificateDetailsDocumentsProductCertificationContract                 `json:"productCertificationContract"`
	ProductDesignResearchConclusion        CertificateDetailsDocumentsProductDesignResearchConclusion              `json:"productDesignResearchConclusion"`
	ProductTypeCertificates                []CertificateDetailsDocumentsProductTypeCertificatesItem                `json:"productTypeCertificates"`
	ProductTypeResearchConclusion          CertificateDetailsDocumentsProductTypeResearchConclusion                `json:"productTypeResearchConclusion"`
	ProtocolTestingTechnicalService        CertificateDetailsDocumentsProtocolTestingTechnicalService              `json:"protocolTestingTechnicalService"`
	QMSCertificates                        []CertificateDetailsDocumentsQMSCertificatesItem                        `json:"qmsCertificates"`
	RawMaterialCertificates                []json.RawMessage                                                       `json:"rawMaterialCertificates"`
	ReportsOnTypeApprovalUnderUNRegulation []CertificateDetailsDocumentsReportsOnTypeApprovalUnderUNRegulationItem `json:"reportsOnTypeApprovalUnderUNRegulation"`
	SamplingAct                            CertificateDetailsDocumentsSamplingAct                                  `json:"samplingAct"`
	VehicleChassisTypeApprovals            []json.RawMessage                                                       `json:"vehicleChassisTypeApprovals"`
}

// CertificateDetailsDocumentsApplicantOtherDocumentsItem contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsApplicantOtherDocumentsItem struct {
	Annex                bool            `json:"annex"`
	Date                 string          `json:"date"`
	ID                   int64           `json:"id"`
	IDCategory           int64           `json:"idCategory"`
	IDTechnicalReglament json.RawMessage `json:"idTechnicalReglament"`
	Name                 string          `json:"name"`
	Number               string          `json:"number"`
}

// CertificateDetailsDocumentsApplicationCertificationDecision contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsApplicationCertificationDecision struct {
	Date                 string          `json:"date"`
	ID                   int64           `json:"id"`
	IDTechnicalReglament json.RawMessage `json:"idTechnicalReglament"`
	Number               string          `json:"number"`
}

// CertificateDetailsDocumentsExpertConclusion contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsExpertConclusion struct {
	AccreditedPersonName json.RawMessage `json:"accreditedPersonName"`
	Annex                bool            `json:"annex"`
	AttestateBeginDate   json.RawMessage `json:"attestateBeginDate"`
	AttestateEndDate     json.RawMessage `json:"attestateEndDate"`
	AttestateRegNumber   json.RawMessage `json:"attestateRegNumber"`
	Date                 json.RawMessage `json:"date"`
	ID                   int64           `json:"id"`
	IDCategory           json.RawMessage `json:"idCategory"`
	IDTechnicalReglament json.RawMessage `json:"idTechnicalReglament"`
	Name                 json.RawMessage `json:"name"`
	Number               json.RawMessage `json:"number"`
}

// CertificateDetailsDocumentsForeignManufacturerContract contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsForeignManufacturerContract struct {
	ApplicantResponsibility    json.RawMessage `json:"applicantResponsibility"`
	Date                       *string         `json:"date"`
	EndDate                    json.RawMessage `json:"endDate"`
	ID                         int64           `json:"id"`
	IDTechnicalReglament       json.RawMessage `json:"idTechnicalReglament"`
	ManufacturerResponsibility json.RawMessage `json:"manufacturerResponsibility"`
	Number                     *string         `json:"number"`
	Subject                    *string         `json:"subject"`
}

// CertificateDetailsDocumentsManufacturingConditionAnalysisActsItem contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsManufacturingConditionAnalysisActsItem struct {
	Addresses            []json.RawMessage                                                                  `json:"addresses"`
	AnalysisBeginDate    *string                                                                            `json:"analysisBeginDate"`
	AnalysisEndDate      *string                                                                            `json:"analysisEndDate"`
	AnalysisTakenDate    json.RawMessage                                                                    `json:"analysisTakenDate"`
	Annex                bool                                                                               `json:"annex"`
	CertExperts          []CertificateDetailsDocumentsManufacturingConditionAnalysisActsItemCertExpertsItem `json:"certExperts"`
	Date                 *string                                                                            `json:"date"`
	ID                   int64                                                                              `json:"id"`
	IDTechnicalReglament json.RawMessage                                                                    `json:"idTechnicalReglament"`
	Number               *string                                                                            `json:"number"`
}

// CertificateDetailsDocumentsManufacturingConditionAnalysisActsItemCertExpertsItem contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsManufacturingConditionAnalysisActsItemCertExpertsItem struct {
	FirstName  string `json:"firstName"`
	IDEmployee int64  `json:"idEmployee"`
	IDExpert   int64  `json:"idExpert"`
	IDPerson   int64  `json:"idPerson"`
	Patronimyc string `json:"patronimyc"`
	Surname    string `json:"surname"`
}

// CertificateDetailsDocumentsProductCertificationContract contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsProductCertificationContract struct {
	Date                 json.RawMessage `json:"date"`
	ID                   int64           `json:"id"`
	IDTechnicalReglament json.RawMessage `json:"idTechnicalReglament"`
	Number               json.RawMessage `json:"number"`
}

// CertificateDetailsDocumentsProductDesignResearchConclusion contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsProductDesignResearchConclusion struct {
	AccreditedPersonName json.RawMessage `json:"accreditedPersonName"`
	Annex                bool            `json:"annex"`
	AttestatEndDate      json.RawMessage `json:"attestatEndDate"`
	AttestatRegDate      json.RawMessage `json:"attestatRegDate"`
	AttestatRegNumber    json.RawMessage `json:"attestatRegNumber"`
	ID                   int64           `json:"id"`
	IDTechnicalReglament json.RawMessage `json:"idTechnicalReglament"`
	IssueDate            json.RawMessage `json:"issueDate"`
	Number               json.RawMessage `json:"number"`
}

// CertificateDetailsDocumentsProductTypeCertificatesItem contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsProductTypeCertificatesItem struct {
	AccreditedPersonName json.RawMessage `json:"accreditedPersonName"`
	Annex                bool            `json:"annex"`
	AttestatEndDate      json.RawMessage `json:"attestatEndDate"`
	AttestatRegDate      json.RawMessage `json:"attestatRegDate"`
	AttestatRegNumber    json.RawMessage `json:"attestatRegNumber"`
	ID                   int64           `json:"id"`
	IDTechnicalReglament int64           `json:"idTechnicalReglament"`
	IssueDate            json.RawMessage `json:"issueDate"`
	Number               json.RawMessage `json:"number"`
	SampleProduct        json.RawMessage `json:"sampleProduct"`
}

// CertificateDetailsDocumentsProductTypeResearchConclusion contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsProductTypeResearchConclusion struct {
	AccreditedPersonName json.RawMessage `json:"accreditedPersonName"`
	Annex                bool            `json:"annex"`
	AttestatEndDate      json.RawMessage `json:"attestatEndDate"`
	AttestatRegDate      json.RawMessage `json:"attestatRegDate"`
	AttestatRegNumber    json.RawMessage `json:"attestatRegNumber"`
	ID                   int64           `json:"id"`
	IDTechnicalReglament json.RawMessage `json:"idTechnicalReglament"`
	IssueDate            json.RawMessage `json:"issueDate"`
	Number               json.RawMessage `json:"number"`
	SampleProduct        json.RawMessage `json:"sampleProduct"`
}

// CertificateDetailsDocumentsProtocolTestingTechnicalService contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsProtocolTestingTechnicalService struct {
	Annex                bool            `json:"annex"`
	Date                 json.RawMessage `json:"date"`
	ID                   int64           `json:"id"`
	IDCategory           json.RawMessage `json:"idCategory"`
	IDTechnicalReglament json.RawMessage `json:"idTechnicalReglament"`
	Name                 json.RawMessage `json:"name"`
	Number               json.RawMessage `json:"number"`
	TechServiceName      json.RawMessage `json:"techServiceName"`
}

// CertificateDetailsDocumentsQMSCertificatesItem contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsQMSCertificatesItem struct {
	AccredEec                 bool              `json:"accredEec"`
	AccreditedPersonName      json.RawMessage   `json:"accreditedPersonName"`
	Annex                     bool              `json:"annex"`
	AttestatEndDate           json.RawMessage   `json:"attestatEndDate"`
	AttestatRegDate           json.RawMessage   `json:"attestatRegDate"`
	AttestatRegNumber         json.RawMessage   `json:"attestatRegNumber"`
	EndDate                   json.RawMessage   `json:"endDate"`
	ID                        int64             `json:"id"`
	IDAccredPlace             string            `json:"idAccredPlace"`
	IDTechnicalReglament      int64             `json:"idTechnicalReglament"`
	IDType                    json.RawMessage   `json:"idType"`
	IDTypeActivity            json.RawMessage   `json:"idTypeActivity"`
	IssueDate                 json.RawMessage   `json:"issueDate"`
	Number                    json.RawMessage   `json:"number"`
	QMSCertificationDocuments []json.RawMessage `json:"qmsCertificationDocuments"`
}

// CertificateDetailsDocumentsReportsOnTypeApprovalUnderUNRegulationItem contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsReportsOnTypeApprovalUnderUNRegulationItem struct {
	Annex                bool            `json:"annex"`
	Country              string          `json:"country"`
	Date                 json.RawMessage `json:"date"`
	ID                   int64           `json:"id"`
	IDCategory           json.RawMessage `json:"idCategory"`
	IDTechnicalReglament int64           `json:"idTechnicalReglament"`
	Name                 json.RawMessage `json:"name"`
	Number               json.RawMessage `json:"number"`
	OrgName              json.RawMessage `json:"orgName"`
}

// CertificateDetailsDocumentsSamplingAct contains fields observed in public registry detail JSON.
type CertificateDetailsDocumentsSamplingAct struct {
	Date                 *string         `json:"date"`
	ID                   int64           `json:"id"`
	IDTechnicalReglament json.RawMessage `json:"idTechnicalReglament"`
	Number               *string         `json:"number"`
}

// CertificateDetailsExpertsItem contains fields observed in public registry detail JSON.
type CertificateDetailsExpertsItem struct {
	FirstName  string `json:"firstName"`
	IDEmployee int64  `json:"idEmployee"`
	IDExpert   int64  `json:"idExpert"`
	IDPerson   int64  `json:"idPerson"`
	Patronimyc string `json:"patronimyc"`
	Surname    string `json:"surname"`
}

// CertificateDetailsManufacturer contains fields observed in public registry detail JSON.
type CertificateDetailsManufacturer struct {
	AddlRegInfo        *string                                       `json:"addlRegInfo"`
	Addresses          []CertificateDetailsManufacturerAddressesItem `json:"addresses"`
	Contacts           []json.RawMessage                             `json:"contacts"`
	FirstName          *string                                       `json:"firstName"`
	FullName           string                                        `json:"fullName"`
	HeadPosition       *string                                       `json:"headPosition"`
	IDEGRUL            json.RawMessage                               `json:"idEgrul"`
	IDLegalForm        json.RawMessage                               `json:"idLegalForm"`
	IDLegalSubject     int64                                         `json:"idLegalSubject"`
	IDLegalSubjectType int64                                         `json:"idLegalSubjectType"`
	IDPerson           int64                                         `json:"idPerson"`
	IDPersonDoc        json.RawMessage                               `json:"idPersonDoc"`
	INN                *string                                       `json:"inn"`
	IsEecRegister      bool                                          `json:"isEecRegister"`
	KPP                *string                                       `json:"kpp"`
	OGRN               *string                                       `json:"ogrn"`
	OGRNAssignDate     json.RawMessage                               `json:"ogrnAssignDate"`
	PassportIssueDate  json.RawMessage                               `json:"passportIssueDate"`
	PassportIssuedBy   json.RawMessage                               `json:"passportIssuedBy"`
	PassportNum        json.RawMessage                               `json:"passportNum"`
	Patronymic         *string                                       `json:"patronymic"`
	RegDate            json.RawMessage                               `json:"regDate"`
	RegOrganName       *string                                       `json:"regOrganName"`
	ShortName          *string                                       `json:"shortName"`
	Surname            *string                                       `json:"surname"`
	Transnational      []json.RawMessage                             `json:"transnational"`
}

// CertificateDetailsManufacturerAddressesItem contains fields observed in public registry detail JSON.
type CertificateDetailsManufacturerAddressesItem struct {
	Flat            json.RawMessage `json:"flat"`
	ForeignCity     string          `json:"foreignCity"`
	ForeignDistrict *string         `json:"foreignDistrict"`
	ForeignHouse    json.RawMessage `json:"foreignHouse"`
	ForeignLocality json.RawMessage `json:"foreignLocality"`
	ForeignStreet   json.RawMessage `json:"foreignStreet"`
	FullAddress     string          `json:"fullAddress"`
	Gln             json.RawMessage `json:"gln"`
	Glonass         *string         `json:"glonass"`
	IDAddrType      int64           `json:"idAddrType"`
	IDAddress       int64           `json:"idAddress"`
	IDCity          json.RawMessage `json:"idCity"`
	IDCodeOksm      string          `json:"idCodeOksm"`
	IDDistrict      json.RawMessage `json:"idDistrict"`
	IDHouse         json.RawMessage `json:"idHouse"`
	IDLocality      json.RawMessage `json:"idLocality"`
	IDStreet        json.RawMessage `json:"idStreet"`
	IDSubject       json.RawMessage `json:"idSubject"`
	OksmShort       bool            `json:"oksmShort"`
	OtherGln        json.RawMessage `json:"otherGln"`
	PostCode        json.RawMessage `json:"postCode"`
	UniqueAddress   json.RawMessage `json:"uniqueAddress"`
}

// CertificateDetailsManufacturerFilialsItem contains fields observed in public registry detail JSON.
type CertificateDetailsManufacturerFilialsItem struct {
	Addresses   []CertificateDetailsManufacturerFilialsItemAddressesItem `json:"addresses"`
	Annex       bool                                                     `json:"annex"`
	Contacts    []json.RawMessage                                        `json:"contacts"`
	FullName    *string                                                  `json:"fullName"`
	IDFilial    int64                                                    `json:"idFilial"`
	IDLegalForm json.RawMessage                                          `json:"idLegalForm"`
	KPP         json.RawMessage                                          `json:"kpp"`
	ShortName   json.RawMessage                                          `json:"shortName"`
}

// CertificateDetailsManufacturerFilialsItemAddressesItem contains fields observed in public registry detail JSON.
type CertificateDetailsManufacturerFilialsItemAddressesItem struct {
	Flat            json.RawMessage `json:"flat"`
	ForeignCity     string          `json:"foreignCity"`
	ForeignDistrict *string         `json:"foreignDistrict"`
	ForeignHouse    json.RawMessage `json:"foreignHouse"`
	ForeignLocality json.RawMessage `json:"foreignLocality"`
	ForeignStreet   json.RawMessage `json:"foreignStreet"`
	FullAddress     string          `json:"fullAddress"`
	Gln             *string         `json:"gln"`
	Glonass         *string         `json:"glonass"`
	IDAddrType      int64           `json:"idAddrType"`
	IDAddress       int64           `json:"idAddress"`
	IDCity          json.RawMessage `json:"idCity"`
	IDCodeOksm      string          `json:"idCodeOksm"`
	IDDistrict      json.RawMessage `json:"idDistrict"`
	IDHouse         json.RawMessage `json:"idHouse"`
	IDLocality      json.RawMessage `json:"idLocality"`
	IDStreet        json.RawMessage `json:"idStreet"`
	IDSubject       json.RawMessage `json:"idSubject"`
	OksmShort       bool            `json:"oksmShort"`
	OtherGln        json.RawMessage `json:"otherGln"`
	PostCode        json.RawMessage `json:"postCode"`
	UniqueAddress   json.RawMessage `json:"uniqueAddress"`
}

// CertificateDetailsProduct contains fields observed in public registry detail JSON.
type CertificateDetailsProduct struct {
	BatchID          json.RawMessage                                `json:"batchId"`
	BatchSize        *string                                        `json:"batchSize"`
	FullName         string                                         `json:"fullName"`
	IDProduct        int64                                          `json:"idProduct"`
	IDProductOrigin  string                                         `json:"idProductOrigin"`
	Identification   json.RawMessage                                `json:"identification"`
	Identifications  []CertificateDetailsProductIdentificationsItem `json:"identifications"`
	Marking          json.RawMessage                                `json:"marking"`
	StorageCondition string                                         `json:"storageCondition"`
	UsageCondition   *string                                        `json:"usageCondition"`
	UsageScope       *string                                        `json:"usageScope"`
}

// CertificateDetailsProductGroupsItem contains fields observed in public registry detail JSON.
type CertificateDetailsProductGroupsItem struct {
	IDGroup        int64 `json:"idGroup"`
	IDProductGroup int64 `json:"idProductGroup"`
	IDTechReg      int64 `json:"idTechReg"`
}

// CertificateDetailsProductIdentificationsItem contains fields observed in public registry detail JSON.
type CertificateDetailsProductIdentificationsItem struct {
	Amount           json.RawMessage                                             `json:"amount"`
	Annex            bool                                                        `json:"annex"`
	Article          *string                                                     `json:"article"`
	Description      string                                                      `json:"description"`
	Documents        []CertificateDetailsProductIdentificationsItemDocumentsItem `json:"documents"`
	ExpiryDate       json.RawMessage                                             `json:"expiryDate"`
	FactoryNumber    json.RawMessage                                             `json:"factoryNumber"`
	GTIN             []string                                                    `json:"gtin"`
	IDIdentification int64                                                       `json:"idIdentification"`
	IDOkei           *int64                                                      `json:"idOkei"`
	IDOkpds          []json.RawMessage                                           `json:"idOkpds"`
	IDTnveds         []int64                                                     `json:"idTnveds"`
	LifeTime         json.RawMessage                                             `json:"lifeTime"`
	Model            *string                                                     `json:"model"`
	Name             string                                                      `json:"name"`
	ProductionDate   json.RawMessage                                             `json:"productionDate"`
	Sort             json.RawMessage                                             `json:"sort"`
	Standards        []CertificateDetailsProductIdentificationsItemStandardsItem `json:"standards"`
	StorageTime      *string                                                     `json:"storageTime"`
	TradeMark        json.RawMessage                                             `json:"tradeMark"`
	Type             json.RawMessage                                             `json:"type"`
}

// CertificateDetailsProductIdentificationsItemDocumentsItem contains fields observed in public registry detail JSON.
type CertificateDetailsProductIdentificationsItemDocumentsItem struct {
	Date   string `json:"date"`
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
}

// CertificateDetailsProductIdentificationsItemStandardsItem contains fields observed in public registry detail JSON.
type CertificateDetailsProductIdentificationsItemStandardsItem struct {
	AddlInfo       json.RawMessage `json:"addlInfo"`
	Annex          bool            `json:"annex"`
	Designation    string          `json:"designation"`
	IDDictStandard *int64          `json:"idDictStandard"`
	IDStandard     int64           `json:"idStandard"`
	IDStatus       *int64          `json:"idStatus"`
	Name           string          `json:"name"`
	Section        *string         `json:"section"`
}

// CertificateDetailsReplacedBy contains fields observed in public registry detail JSON.
type CertificateDetailsReplacedBy struct {
	CertSource              int64   `json:"certSource"`
	CertTarget              int64   `json:"certTarget"`
	ReplacementReasonIDList []int64 `json:"replacementReasonIdList"`
}

// CertificateDetailsReplacement contains fields observed in public registry detail JSON.
type CertificateDetailsReplacement struct {
	CertSource              int64   `json:"certSource"`
	CertTarget              int64   `json:"certTarget"`
	ReplacementReasonIDList []int64 `json:"replacementReasonIdList"`
}

// CertificateDetailsStatusChangesItem contains fields observed in public registry detail JSON.
type CertificateDetailsStatusChangesItem struct {
	BeginDate         string            `json:"beginDate"`
	Comment           *string           `json:"comment"`
	DocDate           *string           `json:"docDate"`
	DocName           *string           `json:"docName"`
	DocNumber         *string           `json:"docNumber"`
	EndDate           *string           `json:"endDate"`
	IDBasis           json.RawMessage   `json:"idBasis"`
	IDBasisSolutions  []json.RawMessage `json:"idBasisSolutions"`
	IDChangeStatus    int64             `json:"idChangeStatus"`
	IDFileScan        json.RawMessage   `json:"idFileScan"`
	IDMRPA            json.RawMessage   `json:"idMrpa"`
	IDOngoingEvent    json.RawMessage   `json:"idOngoingEvent"`
	IDStatus          int64             `json:"idStatus"`
	LevelStateControl json.RawMessage   `json:"levelStateControl"`
	LkType            string            `json:"lkType"`
	Publicated        bool              `json:"publicated"`
	PublishDate       string            `json:"publishDate"`
	TypeStateControl  json.RawMessage   `json:"typeStateControl"`
}

// CertificateDetailsTestingLabsItem contains fields observed in public registry detail JSON.
type CertificateDetailsTestingLabsItem struct {
	AccredEec                  bool                                                    `json:"accredEec"`
	ActDate                    json.RawMessage                                         `json:"actDate"`
	ActIdentificationDate      json.RawMessage                                         `json:"actIdentificationDate"`
	ActIdentificationNumber    json.RawMessage                                         `json:"actIdentificationNumber"`
	ActNumber                  json.RawMessage                                         `json:"actNumber"`
	Annex                      bool                                                    `json:"annex"`
	Basis                      json.RawMessage                                         `json:"basis"`
	BeginDate                  string                                                  `json:"beginDate"`
	DocConfirmCustom           []CertificateDetailsTestingLabsItemDocConfirmCustomItem `json:"docConfirmCustom"`
	EndDate                    json.RawMessage                                         `json:"endDate"`
	FullName                   string                                                  `json:"fullName"`
	IDAccredPlace              string                                                  `json:"idAccredPlace"`
	IDRAL                      json.RawMessage                                         `json:"idRal"`
	IDTestingLab               int64                                                   `json:"idTestingLab"`
	ImportedForResearchTesting bool                                                    `json:"importedForResearchTesting"`
	Protocols                  []CertificateDetailsTestingLabsItemProtocolsItem        `json:"protocols"`
	RegNumber                  string                                                  `json:"regNumber"`
}

// CertificateDetailsTestingLabsItemDocConfirmCustomItem contains fields observed in public registry detail JSON.
type CertificateDetailsTestingLabsItemDocConfirmCustomItem struct {
	CustomInfo                         []CertificateDetailsTestingLabsItemDocConfirmCustomItemCustomInfoItem `json:"customInfo"`
	IDDocConfirmCustom                 int64                                                                 `json:"idDocConfirmCustom"`
	IDDocConfirmCustomType             int64                                                                 `json:"idDocConfirmCustomType"`
	OtherDocs                          []json.RawMessage                                                     `json:"otherDocs"`
	ReasonNonRegistrCustomsDeclaration *string                                                               `json:"reasonNonRegistrCustomsDeclaration"`
}

// CertificateDetailsTestingLabsItemDocConfirmCustomItemCustomInfoItem contains fields observed in public registry detail JSON.
type CertificateDetailsTestingLabsItemDocConfirmCustomItemCustomInfoItem struct {
	CustomDeclNumber string `json:"customDeclNumber"`
	IDCustomInfo     int64  `json:"idCustomInfo"`
}

// CertificateDetailsTestingLabsItemProtocolsItem contains fields observed in public registry detail JSON.
type CertificateDetailsTestingLabsItemProtocolsItem struct {
	Date              string            `json:"date"`
	IDProtocol        int64             `json:"idProtocol"`
	IDProtocolRPI     *int64            `json:"idProtocolRpi"`
	IsProtocolInvalid bool              `json:"isProtocolInvalid"`
	Number            string            `json:"number"`
	Standards         []json.RawMessage `json:"standards"`
}

// DeclarationDetails contains fields observed in public registry detail JSON.
type DeclarationDetails struct {
	AccreditationBody         json.RawMessage                             `json:"accreditationBody"`
	Annexes                   []DeclarationDetailsAnnexesItem             `json:"annexes"`
	Applicant                 DeclarationDetailsApplicant                 `json:"applicant"`
	ApplicantFilials          []json.RawMessage                           `json:"applicantFilials"`
	ApplicationDate           json.RawMessage                             `json:"applicationDate"`
	ApplicationNumber         json.RawMessage                             `json:"applicationNumber"`
	ApplicationSubmissionDate json.RawMessage                             `json:"applicationSubmissionDate"`
	AssignRegNumber           bool                                        `json:"assignRegNumber"`
	AwaitForApprove           bool                                        `json:"awaitForApprove"`
	CertificationAuthority    json.RawMessage                             `json:"certificationAuthority"`
	ChangePublishDate         string                                      `json:"changePublishDate"`
	Changes                   json.RawMessage                             `json:"changes"`
	DeclEndDate               string                                      `json:"declEndDate"`
	DeclRegDate               string                                      `json:"declRegDate"`
	DeclStage                 DeclarationDetailsDeclStage                 `json:"declStage"`
	DeclarationRegInsteadNum  json.RawMessage                             `json:"declarationRegInsteadNum"`
	DeclarationReplacedNum    json.RawMessage                             `json:"declarationReplacedNum"`
	Documents                 DeclarationDetailsDocuments                 `json:"documents"`
	DraftCreationDate         string                                      `json:"draftCreationDate"`
	EditApp                   bool                                        `json:"editApp"`
	Experts                   []json.RawMessage                           `json:"experts"`
	FirstName                 json.RawMessage                             `json:"firstName"`
	Functions                 json.RawMessage                             `json:"functions"`
	IDApplication             json.RawMessage                             `json:"idApplication"`
	IDApplicationStatus       json.RawMessage                             `json:"idApplicationStatus"`
	IDDeclScheme              int64                                       `json:"idDeclScheme"`
	IDDeclType                int64                                       `json:"idDeclType"`
	IDDeclaration             int64                                       `json:"idDeclaration"`
	IDDeclarationChecker      json.RawMessage                             `json:"idDeclarationChecker"`
	IDDeclarationRegInstead   json.RawMessage                             `json:"idDeclarationRegInstead"`
	IDDeclarationReplaced     json.RawMessage                             `json:"idDeclarationReplaced"`
	IDGroups                  []int64                                     `json:"idGroups"`
	IDMRPA                    *string                                     `json:"idMrpa"`
	IDObjectDeclType          int64                                       `json:"idObjectDeclType"`
	IDProductSingleLists      []json.RawMessage                           `json:"idProductSingleLists"`
	IDReplacementReason       []json.RawMessage                           `json:"idReplacementReason"`
	IDSigner                  json.RawMessage                             `json:"idSigner"`
	IDSignerEmployee          json.RawMessage                             `json:"idSignerEmployee"`
	IDStatus                  int64                                       `json:"idStatus"`
	IDStatusTestingLabs       int64                                       `json:"idStatusTestingLabs"`
	IDTechnicalReglaments     []int64                                     `json:"idTechnicalReglaments"`
	IsSrd                     bool                                        `json:"isSrd"`
	LastUpdate                string                                      `json:"lastUpdate"`
	ManufIsApplicant          bool                                        `json:"manufIsApplicant"`
	Manufacturer              DeclarationDetailsManufacturer              `json:"manufacturer"`
	ManufacturerFilials       []DeclarationDetailsManufacturerFilialsItem `json:"manufacturerFilials"`
	NoSanction                bool                                        `json:"noSanction"`
	Number                    string                                      `json:"number"`
	Patronymic                json.RawMessage                             `json:"patronymic"`
	Product                   DeclarationDetailsProduct                   `json:"product"`
	PublishDate               string                                      `json:"publishDate"`
	ScanCopy                  []DeclarationDetailsScanCopyItem            `json:"scanCopy"`
	ShowInOpenPart            bool                                        `json:"showInOpenPart"`
	SNILS                     json.RawMessage                             `json:"snils"`
	StatusChanges             []DeclarationDetailsStatusChangesItem       `json:"statusChanges"`
	SubmissionDate            json.RawMessage                             `json:"submissionDate"`
	Surname                   json.RawMessage                             `json:"surname"`
	TempNumber                string                                      `json:"tempNumber"`
	TestingLabs               []DeclarationDetailsTestingLabsItem         `json:"testingLabs"`
	UnpublishedChanges        bool                                        `json:"unpublishedChanges"`
	ViolationPublishDate      json.RawMessage                             `json:"violationPublishDate"`
	Raw                       json.RawMessage                             `json:"-"`
}

// DeclarationDetailsAnnexesItem contains fields observed in public registry detail JSON.
type DeclarationDetailsAnnexesItem struct {
	IDAnnex   int64  `json:"idAnnex"`
	IDType    int64  `json:"idType"`
	Ord       string `json:"ord"`
	PageCount int64  `json:"pageCount"`
}

// DeclarationDetailsApplicant contains fields observed in public registry detail JSON.
type DeclarationDetailsApplicant struct {
	AddlRegInfo           string                                     `json:"addlRegInfo"`
	Addresses             []DeclarationDetailsApplicantAddressesItem `json:"addresses"`
	ApplicantHead         bool                                       `json:"applicantHead"`
	Contacts              []DeclarationDetailsApplicantContactsItem  `json:"contacts"`
	FirstName             string                                     `json:"firstName"`
	FullName              string                                     `json:"fullName"`
	HeadPosition          string                                     `json:"headPosition"`
	IDApplicantType       int64                                      `json:"idApplicantType"`
	IDEGRUL               json.RawMessage                            `json:"idEgrul"`
	IDKopf                json.RawMessage                            `json:"idKopf"`
	IDLegalForm           int64                                      `json:"idLegalForm"`
	IDLegalSubject        int64                                      `json:"idLegalSubject"`
	IDLegalSubjectType    int64                                      `json:"idLegalSubjectType"`
	IDPerson              int64                                      `json:"idPerson"`
	IDResponsiblePerson   int64                                      `json:"idResponsiblePerson"`
	INN                   string                                     `json:"inn"`
	IsEecRegister         bool                                       `json:"isEecRegister"`
	KPP                   string                                     `json:"kpp"`
	NameLegalForm         string                                     `json:"nameLegalForm"`
	OGRN                  string                                     `json:"ogrn"`
	OGRNAssignDate        string                                     `json:"ogrnAssignDate"`
	Patronymic            string                                     `json:"patronymic"`
	RegDate               string                                     `json:"regDate"`
	RegOrganName          string                                     `json:"regOrganName"`
	ResponsibleContacts   []json.RawMessage                          `json:"responsibleContacts"`
	ResponsibleDocDate    json.RawMessage                            `json:"responsibleDocDate"`
	ResponsibleDocName    json.RawMessage                            `json:"responsibleDocName"`
	ResponsibleDocNumber  json.RawMessage                            `json:"responsibleDocNumber"`
	ResponsibleFirstName  string                                     `json:"responsibleFirstName"`
	ResponsiblePatronymic string                                     `json:"responsiblePatronymic"`
	ResponsiblePosition   string                                     `json:"responsiblePosition"`
	ResponsibleSurname    string                                     `json:"responsibleSurname"`
	ShortName             string                                     `json:"shortName"`
	SNILS                 string                                     `json:"snils"`
	Surname               string                                     `json:"surname"`
	Transnational         []json.RawMessage                          `json:"transnational"`
}

// DeclarationDetailsApplicantAddressesItem contains fields observed in public registry detail JSON.
type DeclarationDetailsApplicantAddressesItem struct {
	Flat            json.RawMessage `json:"flat"`
	ForeignCity     json.RawMessage `json:"foreignCity"`
	ForeignDistrict json.RawMessage `json:"foreignDistrict"`
	ForeignHouse    json.RawMessage `json:"foreignHouse"`
	ForeignLocality json.RawMessage `json:"foreignLocality"`
	ForeignStreet   json.RawMessage `json:"foreignStreet"`
	FullAddress     string          `json:"fullAddress"`
	Gln             json.RawMessage `json:"gln"`
	Glonass         json.RawMessage `json:"glonass"`
	IDAddrType      int64           `json:"idAddrType"`
	IDAddress       int64           `json:"idAddress"`
	IDCity          json.RawMessage `json:"idCity"`
	IDCodeOksm      string          `json:"idCodeOksm"`
	IDDistrict      json.RawMessage `json:"idDistrict"`
	IDHouse         json.RawMessage `json:"idHouse"`
	IDLocality      json.RawMessage `json:"idLocality"`
	IDStreet        json.RawMessage `json:"idStreet"`
	IDSubject       *string         `json:"idSubject"`
	OksmShort       bool            `json:"oksmShort"`
	OtherGln        json.RawMessage `json:"otherGln"`
	PostCode        *string         `json:"postCode"`
	UniqueAddress   json.RawMessage `json:"uniqueAddress"`
}

// DeclarationDetailsApplicantContactsItem contains fields observed in public registry detail JSON.
type DeclarationDetailsApplicantContactsItem struct {
	IDContact     int64  `json:"idContact"`
	IDContactType int64  `json:"idContactType"`
	Value         string `json:"value"`
}

// DeclarationDetailsDeclStage contains fields observed in public registry detail JSON.
type DeclarationDetailsDeclStage struct {
	IDStage int64 `json:"idStage"`
	IsValid bool  `json:"isValid"`
}

// DeclarationDetailsDocuments contains fields observed in public registry detail JSON.
type DeclarationDetailsDocuments struct {
	ApplicantOtherDocuments                []json.RawMessage                                        `json:"applicantOtherDocuments"`
	CommonDocuments                        map[string]json.RawMessage                               `json:"commonDocuments"`
	ComplianceStandards                    []json.RawMessage                                        `json:"complianceStandards"`
	ComponentCertificates                  []json.RawMessage                                        `json:"componentCertificates"`
	ForeignManufacturerContract            json.RawMessage                                          `json:"foreignManufacturerContract"`
	ProductTypeApplicationDecision         json.RawMessage                                          `json:"productTypeApplicationDecision"`
	ProductTypeCertificates                []DeclarationDetailsDocumentsProductTypeCertificatesItem `json:"productTypeCertificates"`
	ProductTypeContract                    json.RawMessage                                          `json:"productTypeContract"`
	ProductTypeResearchConclusion          json.RawMessage                                          `json:"productTypeResearchConclusion"`
	ProductionControlResults               []json.RawMessage                                        `json:"productionControlResults"`
	ProtocolTestingTechnicalService        json.RawMessage                                          `json:"protocolTestingTechnicalService"`
	QMSCertificates                        []json.RawMessage                                        `json:"qmsCertificates"`
	ReportsOnTypeApprovalUnderUNRegulation []json.RawMessage                                        `json:"reportsOnTypeApprovalUnderUNRegulation"`
	VehicleChassisTypeApproval             json.RawMessage                                          `json:"vehicleChassisTypeApproval"`
}

// DeclarationDetailsDocumentsProductTypeCertificatesItem contains fields observed in public registry detail JSON.
type DeclarationDetailsDocumentsProductTypeCertificatesItem struct {
	AccreditedPersonName  json.RawMessage `json:"accreditedPersonName"`
	Annex                 json.RawMessage `json:"annex"`
	AttestatEndDate       json.RawMessage `json:"attestatEndDate"`
	AttestatRegDate       json.RawMessage `json:"attestatRegDate"`
	AttestatRegNumber     json.RawMessage `json:"attestatRegNumber"`
	IDAccredPlace         json.RawMessage `json:"idAccredPlace"`
	IDFile                json.RawMessage `json:"idFile"`
	IDPtCertificate       int64           `json:"idPtCertificate"`
	IDTechnicalRegulation int64           `json:"idTechnicalRegulation"`
	IsAccredEec           bool            `json:"isAccredEec"`
	IssueDate             json.RawMessage `json:"issueDate"`
	Number                json.RawMessage `json:"number"`
	SampleProduct         json.RawMessage `json:"sampleProduct"`
}

// DeclarationDetailsManufacturer contains fields observed in public registry detail JSON.
type DeclarationDetailsManufacturer struct {
	AddlRegInfo           string                                        `json:"addlRegInfo"`
	Addresses             []DeclarationDetailsManufacturerAddressesItem `json:"addresses"`
	ApplicantHead         bool                                          `json:"applicantHead"`
	Contacts              []DeclarationDetailsManufacturerContactsItem  `json:"contacts"`
	FirstName             string                                        `json:"firstName"`
	FullName              string                                        `json:"fullName"`
	HeadPosition          string                                        `json:"headPosition"`
	IDEGRUL               json.RawMessage                               `json:"idEgrul"`
	IDKopf                json.RawMessage                               `json:"idKopf"`
	IDLegalForm           *int64                                        `json:"idLegalForm"`
	IDLegalSubject        int64                                         `json:"idLegalSubject"`
	IDLegalSubjectType    int64                                         `json:"idLegalSubjectType"`
	IDPerson              int64                                         `json:"idPerson"`
	IDResponsiblePerson   json.RawMessage                               `json:"idResponsiblePerson"`
	INN                   string                                        `json:"inn"`
	IsEecRegister         bool                                          `json:"isEecRegister"`
	KPP                   string                                        `json:"kpp"`
	NameLegalForm         *string                                       `json:"nameLegalForm"`
	OGRN                  string                                        `json:"ogrn"`
	OGRNAssignDate        *string                                       `json:"ogrnAssignDate"`
	Patronymic            string                                        `json:"patronymic"`
	RegDate               *string                                       `json:"regDate"`
	RegOrganName          string                                        `json:"regOrganName"`
	ResponsibleContacts   json.RawMessage                               `json:"responsibleContacts"`
	ResponsibleDocDate    json.RawMessage                               `json:"responsibleDocDate"`
	ResponsibleDocName    json.RawMessage                               `json:"responsibleDocName"`
	ResponsibleDocNumber  json.RawMessage                               `json:"responsibleDocNumber"`
	ResponsibleFirstName  json.RawMessage                               `json:"responsibleFirstName"`
	ResponsiblePatronymic json.RawMessage                               `json:"responsiblePatronymic"`
	ResponsiblePosition   json.RawMessage                               `json:"responsiblePosition"`
	ResponsibleSurname    json.RawMessage                               `json:"responsibleSurname"`
	ShortName             string                                        `json:"shortName"`
	SNILS                 string                                        `json:"snils"`
	Surname               string                                        `json:"surname"`
	Transnational         []json.RawMessage                             `json:"transnational"`
}

// DeclarationDetailsManufacturerAddressesItem contains fields observed in public registry detail JSON.
type DeclarationDetailsManufacturerAddressesItem struct {
	Flat            *string         `json:"flat"`
	ForeignCity     json.RawMessage `json:"foreignCity"`
	ForeignDistrict *string         `json:"foreignDistrict"`
	ForeignHouse    json.RawMessage `json:"foreignHouse"`
	ForeignLocality json.RawMessage `json:"foreignLocality"`
	ForeignStreet   json.RawMessage `json:"foreignStreet"`
	FullAddress     string          `json:"fullAddress"`
	Gln             json.RawMessage `json:"gln"`
	Glonass         json.RawMessage `json:"glonass"`
	IDAddrType      int64           `json:"idAddrType"`
	IDAddress       int64           `json:"idAddress"`
	IDCity          json.RawMessage `json:"idCity"`
	IDCodeOksm      string          `json:"idCodeOksm"`
	IDDistrict      json.RawMessage `json:"idDistrict"`
	IDHouse         json.RawMessage `json:"idHouse"`
	IDLocality      json.RawMessage `json:"idLocality"`
	IDStreet        json.RawMessage `json:"idStreet"`
	IDSubject       json.RawMessage `json:"idSubject"`
	OksmShort       bool            `json:"oksmShort"`
	OtherGln        json.RawMessage `json:"otherGln"`
	PostCode        json.RawMessage `json:"postCode"`
	UniqueAddress   json.RawMessage `json:"uniqueAddress"`
}

// DeclarationDetailsManufacturerContactsItem contains fields observed in public registry detail JSON.
type DeclarationDetailsManufacturerContactsItem struct {
	IDContact     int64  `json:"idContact"`
	IDContactType int64  `json:"idContactType"`
	Value         string `json:"value"`
}

// DeclarationDetailsManufacturerFilialsItem contains fields observed in public registry detail JSON.
type DeclarationDetailsManufacturerFilialsItem struct {
	Addresses     []DeclarationDetailsManufacturerFilialsItemAddressesItem `json:"addresses"`
	Annex         bool                                                     `json:"annex"`
	Contacts      []json.RawMessage                                        `json:"contacts"`
	EGRULSelected bool                                                     `json:"egrulSelected"`
	FullName      *string                                                  `json:"fullName"`
	IDFilial      int64                                                    `json:"idFilial"`
	IDKopf        json.RawMessage                                          `json:"idKopf"`
	IDLegalForm   json.RawMessage                                          `json:"idLegalForm"`
	KPP           json.RawMessage                                          `json:"kpp"`
	NameLegalForm json.RawMessage                                          `json:"nameLegalForm"`
	ShortName     json.RawMessage                                          `json:"shortName"`
}

// DeclarationDetailsManufacturerFilialsItemAddressesItem contains fields observed in public registry detail JSON.
type DeclarationDetailsManufacturerFilialsItemAddressesItem struct {
	Flat            *string         `json:"flat"`
	ForeignCity     json.RawMessage `json:"foreignCity"`
	ForeignDistrict *string         `json:"foreignDistrict"`
	ForeignHouse    json.RawMessage `json:"foreignHouse"`
	ForeignLocality json.RawMessage `json:"foreignLocality"`
	ForeignStreet   json.RawMessage `json:"foreignStreet"`
	FullAddress     string          `json:"fullAddress"`
	Gln             json.RawMessage `json:"gln"`
	Glonass         json.RawMessage `json:"glonass"`
	IDAddrType      int64           `json:"idAddrType"`
	IDAddress       int64           `json:"idAddress"`
	IDCity          json.RawMessage `json:"idCity"`
	IDCodeOksm      string          `json:"idCodeOksm"`
	IDDistrict      json.RawMessage `json:"idDistrict"`
	IDHouse         json.RawMessage `json:"idHouse"`
	IDLocality      json.RawMessage `json:"idLocality"`
	IDStreet        json.RawMessage `json:"idStreet"`
	IDSubject       *string         `json:"idSubject"`
	OksmShort       bool            `json:"oksmShort"`
	OtherGln        json.RawMessage `json:"otherGln"`
	PostCode        *string         `json:"postCode"`
	UniqueAddress   json.RawMessage `json:"uniqueAddress"`
}

// DeclarationDetailsProduct contains fields observed in public registry detail JSON.
type DeclarationDetailsProduct struct {
	BatchID          json.RawMessage                                `json:"batchId"`
	BatchSize        json.RawMessage                                `json:"batchSize"`
	FullName         string                                         `json:"fullName"`
	IDProduct        int64                                          `json:"idProduct"`
	IDProductOrigin  string                                         `json:"idProductOrigin"`
	Identification   json.RawMessage                                `json:"identification"`
	Identifications  []DeclarationDetailsProductIdentificationsItem `json:"identifications"`
	Marking          json.RawMessage                                `json:"marking"`
	StorageCondition string                                         `json:"storageCondition"`
	UsageCondition   json.RawMessage                                `json:"usageCondition"`
	UsageScope       json.RawMessage                                `json:"usageScope"`
}

// DeclarationDetailsProductIdentificationsItem contains fields observed in public registry detail JSON.
type DeclarationDetailsProductIdentificationsItem struct {
	Amount           json.RawMessage                                             `json:"amount"`
	Annex            bool                                                        `json:"annex"`
	Article          json.RawMessage                                             `json:"article"`
	Description      json.RawMessage                                             `json:"description"`
	Documents        []DeclarationDetailsProductIdentificationsItemDocumentsItem `json:"documents"`
	ExpiryDate       json.RawMessage                                             `json:"expiryDate"`
	FactoryNumber    json.RawMessage                                             `json:"factoryNumber"`
	GTIN             []string                                                    `json:"gtin"`
	IDIdentification int64                                                       `json:"idIdentification"`
	IDOkei           json.RawMessage                                             `json:"idOkei"`
	IDOkpds          []int64                                                     `json:"idOkpds"`
	IDTnveds         []int64                                                     `json:"idTnveds"`
	LifeTime         json.RawMessage                                             `json:"lifeTime"`
	Model            json.RawMessage                                             `json:"model"`
	Name             string                                                      `json:"name"`
	ProductionDate   json.RawMessage                                             `json:"productionDate"`
	Sort             json.RawMessage                                             `json:"sort"`
	Standards        []json.RawMessage                                           `json:"standards"`
	StorageTime      json.RawMessage                                             `json:"storageTime"`
	TradeMark        json.RawMessage                                             `json:"tradeMark"`
	Type             json.RawMessage                                             `json:"type"`
}

// DeclarationDetailsProductIdentificationsItemDocumentsItem contains fields observed in public registry detail JSON.
type DeclarationDetailsProductIdentificationsItemDocumentsItem struct {
	Date                 json.RawMessage `json:"date"`
	IDProductionDocument int64           `json:"idProductionDocument"`
	Name                 string          `json:"name"`
	Number               json.RawMessage `json:"number"`
}

// DeclarationDetailsScanCopyItem contains fields observed in public registry detail JSON.
type DeclarationDetailsScanCopyItem struct {
	IDFile string `json:"idFile"`
	IDType int64  `json:"idType"`
	Name   string `json:"name"`
}

// DeclarationDetailsStatusChangesItem contains fields observed in public registry detail JSON.
type DeclarationDetailsStatusChangesItem struct {
	BeginDate         string            `json:"beginDate"`
	Comment           json.RawMessage   `json:"comment"`
	DocDate           json.RawMessage   `json:"docDate"`
	DocName           json.RawMessage   `json:"docName"`
	DocNumber         json.RawMessage   `json:"docNumber"`
	EndDate           json.RawMessage   `json:"endDate"`
	IDBasis           json.RawMessage   `json:"idBasis"`
	IDBasisSolutions  []json.RawMessage `json:"idBasisSolutions"`
	IDChangeStatus    int64             `json:"idChangeStatus"`
	IDFileScan        json.RawMessage   `json:"idFileScan"`
	IDMRPA            json.RawMessage   `json:"idMrpa"`
	IDOngoingEvent    json.RawMessage   `json:"idOngoingEvent"`
	IDStatus          int64             `json:"idStatus"`
	LevelStateControl json.RawMessage   `json:"levelStateControl"`
	LkType            json.RawMessage   `json:"lkType"`
	Prescription      json.RawMessage   `json:"prescription"`
	Publicated        bool              `json:"publicated"`
	PublishDate       string            `json:"publishDate"`
	TypeStateControl  json.RawMessage   `json:"typeStateControl"`
}

// DeclarationDetailsTestingLabsItem contains fields observed in public registry detail JSON.
type DeclarationDetailsTestingLabsItem struct {
	AccredEec                  bool                                                    `json:"accredEec"`
	Address                    *DeclarationDetailsTestingLabsItemAddress               `json:"address"`
	Annex                      bool                                                    `json:"annex"`
	Basis                      json.RawMessage                                         `json:"basis"`
	BeginDate                  *string                                                 `json:"beginDate"`
	DocConfirmCustom           []DeclarationDetailsTestingLabsItemDocConfirmCustomItem `json:"docConfirmCustom"`
	EndDate                    json.RawMessage                                         `json:"endDate"`
	FullName                   string                                                  `json:"fullName"`
	IDAccredPlace              string                                                  `json:"idAccredPlace"`
	IDRAL                      *int64                                                  `json:"idRal"`
	IDTestingLab               int64                                                   `json:"idTestingLab"`
	ImportedForResearchTesting bool                                                    `json:"importedForResearchTesting"`
	INN                        json.RawMessage                                         `json:"inn"`
	IsAccredEec                bool                                                    `json:"isAccredEec"`
	OGRN                       json.RawMessage                                         `json:"ogrn"`
	Protocols                  []DeclarationDetailsTestingLabsItemProtocolsItem        `json:"protocols"`
	RegNumber                  *string                                                 `json:"regNumber"`
}

// DeclarationDetailsTestingLabsItemAddress contains fields observed in public registry detail JSON.
type DeclarationDetailsTestingLabsItemAddress struct {
	Flat            string          `json:"flat"`
	ForeignCity     json.RawMessage `json:"foreignCity"`
	ForeignDistrict json.RawMessage `json:"foreignDistrict"`
	ForeignHouse    json.RawMessage `json:"foreignHouse"`
	ForeignLocality json.RawMessage `json:"foreignLocality"`
	ForeignStreet   json.RawMessage `json:"foreignStreet"`
	FullAddress     string          `json:"fullAddress"`
	Gln             json.RawMessage `json:"gln"`
	Glonass         json.RawMessage `json:"glonass"`
	IDAddrType      int64           `json:"idAddrType"`
	IDAddress       int64           `json:"idAddress"`
	IDCity          string          `json:"idCity"`
	IDCodeOksm      string          `json:"idCodeOksm"`
	IDDistrict      json.RawMessage `json:"idDistrict"`
	IDHouse         json.RawMessage `json:"idHouse"`
	IDLocality      json.RawMessage `json:"idLocality"`
	IDStreet        string          `json:"idStreet"`
	IDSubject       string          `json:"idSubject"`
	OksmShort       bool            `json:"oksmShort"`
	OtherGln        json.RawMessage `json:"otherGln"`
	PostCode        json.RawMessage `json:"postCode"`
	UniqueAddress   json.RawMessage `json:"uniqueAddress"`
}

// DeclarationDetailsTestingLabsItemDocConfirmCustomItem contains fields observed in public registry detail JSON.
type DeclarationDetailsTestingLabsItemDocConfirmCustomItem struct {
	CustomInfo                         []json.RawMessage `json:"customInfo"`
	IDDocConfirmCustom                 int64             `json:"idDocConfirmCustom"`
	IDDocConfirmCustomType             int64             `json:"idDocConfirmCustomType"`
	OtherDocs                          []json.RawMessage `json:"otherDocs"`
	ReasonNonRegistrCustomsDeclaration string            `json:"reasonNonRegistrCustomsDeclaration"`
}

// DeclarationDetailsTestingLabsItemProtocolsItem contains fields observed in public registry detail JSON.
type DeclarationDetailsTestingLabsItemProtocolsItem struct {
	Date              string            `json:"date"`
	IDFile            string            `json:"idFile"`
	IDProtocol        int64             `json:"idProtocol"`
	IDProtocolRPI     *int64            `json:"idProtocolRpi"`
	IsProtocolInvalid bool              `json:"isProtocolInvalid"`
	Number            string            `json:"number"`
	Standards         []json.RawMessage `json:"standards"`
}

// UnmarshalJSON preserves the complete response for fields whose type is unconfirmed.
func (d *DeclarationDetails) UnmarshalJSON(data []byte) error {
	type plain DeclarationDetails
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*d = DeclarationDetails(value)
	d.Raw = append([]byte(nil), data...)
	return nil
}

// UnmarshalJSON preserves the complete response for fields whose type is unconfirmed.
func (d *CertificateDetails) UnmarshalJSON(data []byte) error {
	type plain CertificateDetails
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*d = CertificateDetails(value)
	d.Raw = append([]byte(nil), data...)
	return nil
}
