package fgis

import "encoding/json"

// Page is the envelope returned by both registry list endpoints.
type Page[T any] struct {
	Items []T `json:"items"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// DateRange uses the string format accepted by the registry API. An empty
// value means that the corresponding bound is not set.
type DateRange struct {
	MinDate string `json:"minDate"`
	MaxDate string `json:"maxDate"`
}

type ColumnSort struct {
	Column string `json:"column"`
	Sort   string `json:"sort"`
}

// DeclarationFilter follows the request captured for the declaration list.
// The precise shape of non-empty columnsSearch entries has not been observed,
// so RawMessage keeps those entries available without inventing a schema.
type DeclarationFilter struct {
	Status           []int             `json:"status"`
	IDDeclType       []int             `json:"idDeclType"`
	IDCertObjectType []int             `json:"idCertObjectType"`
	IDProductType    []int             `json:"idProductType"`
	IDGroupRU        []int             `json:"idGroupRU"`
	IDGroupEEU       []int             `json:"idGroupEEU"`
	IDTechReg        []int             `json:"idTechReg"`
	IDApplicantType  []int             `json:"idApplicantType"`
	RegDate          DateRange         `json:"regDate"`
	EndDate          DateRange         `json:"endDate"`
	ColumnsSearch    []json.RawMessage `json:"columnsSearch"`
}

type DeclarationSearchRequest struct {
	Size        int               `json:"size"`
	Page        int               `json:"page"`
	Count       int               `json:"count"`
	Filter      DeclarationFilter `json:"filter"`
	ColumnsSort []ColumnSort      `json:"columnsSort"`
}

// DefaultDeclarationSearchRequest reproduces the empty first-page query from
// the captured browser request.
func DefaultDeclarationSearchRequest() DeclarationSearchRequest {
	return DeclarationSearchRequest{
		Size: 10,
		Filter: DeclarationFilter{
			Status:           []int{},
			IDDeclType:       []int{},
			IDCertObjectType: []int{},
			IDProductType:    []int{},
			IDGroupRU:        []int{},
			IDGroupEEU:       []int{},
			IDTechReg:        []int{},
			IDApplicantType:  []int{},
			ColumnsSearch:    []json.RawMessage{},
		},
		ColumnsSort: []ColumnSort{{Column: "declDate", Sort: "DESC"}},
	}
}

// Declaration is the observed list-row schema. The upstream API uses the
// spelling "manufacter" in its JSON keys; field names are normalized here.
// Dates and identifiers remain strings because that is how the API returned
// them in the captured response.
type Declaration struct {
	ID                                 int    `json:"id"`
	IDStatus                           int    `json:"idStatus"`
	IDApplicantLegalSubjectType        int    `json:"idApplicantLegalSubjectType"`
	IDManufacturerLegalSubjectType     int    `json:"idManufacterLegalSubjectType"`
	IsProtocolInvalid                  int    `json:"isProtocolInvalid"`
	SRD                                bool   `json:"srd"`
	Number                             string `json:"number"`
	DeclarationNumber                  string `json:"declNumber"`
	DeclarationTempNumber              string `json:"declTempNumber"`
	CustomDeclarationNumber            string `json:"customDeclNumber"`
	DocumentDeclarationNumber          string `json:"documentDeclarationNumber"`
	DeclarationDate                    string `json:"declDate"`
	DeclarationDraftDate               string `json:"declDraftDate"`
	DeclarationEndDate                 string `json:"declEndDate"`
	DeclarationType                    string `json:"declType"`
	DeclarationObjectType              string `json:"declObjectType"`
	DeclarationScheme                  string `json:"declSchema"`
	ApplicantName                      string `json:"applicantName"`
	ApplicantAddress                   string `json:"applicantAddress"`
	ApplicantINN                       string `json:"applicantInn"`
	ApplicantOGRN                      string `json:"applicantOgrn"`
	ApplicantOPF                       string `json:"applicantOpf"`
	ApplicantType                      string `json:"applicantType"`
	ApplicantLegalSubjectType          string `json:"applicantLegalSubjectType"`
	ApplicantFilialFullNames           string `json:"applicantFilialFullNames"`
	ManufacturerName                   string `json:"manufacterName"`
	ManufacturerAddress                string `json:"manufacterAddress"`
	ManufacturerINN                    string `json:"manufacterInn"`
	ManufacturerOGRN                   string `json:"manufacterOgrn"`
	ManufacturerOPF                    string `json:"manufacterOpf"`
	ManufacturerType                   string `json:"manufacterType"`
	ManufacturerLegalSubjectType       string `json:"manufacterLegalSubjectType"`
	ManufacturerFilialFullNames        string `json:"manufacterFilialFullNames"`
	CreatorINN                         string `json:"creatorInn"`
	CreatorOGRN                        string `json:"creatorOgrn"`
	ExpertName                         string `json:"expertFio"`
	ExpertSNILS                        string `json:"expertSnils"`
	Group                              string `json:"group"`
	AwaitForApprove                    string `json:"awaitForApprove"`
	EditApp                            string `json:"editApp"`
	ProductBatchSize                   string `json:"productBatchSize"`
	ProductFullName                    string `json:"productFullName"`
	ProductIdentificationArticle       string `json:"productIdentificationArticle"`
	ProductIdentificationFactoryNumber string `json:"productIdentificationFactoryNumber"`
	ProductIdentificationGTIN          string `json:"productIdentificationGtin"`
	ProductIdentificationModel         string `json:"productIdentificationModel"`
	ProductIdentificationName          string `json:"productIdentificationName"`
	ProductIdentificationOKPDTNVED     string `json:"productIdentificationOkpdTnved"`
	ProductIdentificationSort          string `json:"productIdentificationSort"`
	ProductIdentificationTrademark     string `json:"productIdentificationTrademark"`
	ProductIdentificationType          string `json:"productIdentificationType"`
	ProductOrigin                      string `json:"productOrig"`
	ProductSingleList                  string `json:"productSingleList"`
	ProductStandard                    string `json:"productStandard"`
	StatusTestingLabs                  string `json:"statusTestingLabs"`
	TechnicalReglaments                string `json:"technicalReglaments"`
	TestLabAccreditationNumber         string `json:"testLabAccredNumber"`
	TestLabProtocolDate                string `json:"testLabProtocolDate"`
	TestLabProtocolNumber              string `json:"testLabProtocolNumber"`
}

// CertificateSearchRequest models the confirmed list envelope. The attached
// HAR does not contain a certificate search, so Filter remains extensible.
type CertificateSearchRequest struct {
	Size        int            `json:"size"`
	Page        int            `json:"page"`
	Count       int            `json:"count"`
	Filter      map[string]any `json:"filter"`
	ColumnsSort []ColumnSort   `json:"columnsSort"`
}

// Certificate holds the raw fields returned by a certificate list row until
// a real certificate response can establish their names and types.
type Certificate map[string]json.RawMessage
