package fgis

import (
	"encoding/json"
	"fmt"
)

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

// DeclarationFilter follows the observed declaration list requests. RawMessage
// preserves compatibility with callers that supply custom search conditions.
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

// DeclarationColumnSearch is the text-column filter used by the declaration
// registry. The certificate registry uses a different JSON shape.
type DeclarationColumnSearch struct {
	Name   string `json:"name"`
	Search string `json:"search"`
	Type   int    `json:"type"`
}

// AddColumnSearch appends a typed condition to the existing raw filter API.
func (f *DeclarationFilter) AddColumnSearch(condition DeclarationColumnSearch) {
	encoded, _ := json.Marshal(condition)
	f.ColumnsSearch = append(f.ColumnsSearch, encoded)
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
	DeclarationReplacedNumber          string `json:"declReplacedNumber"`
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

// CertificateSearchRequest is the original extensible search form. For the
// observed certificate filter and typed rows, use CertificateQuery instead.
type CertificateSearchRequest struct {
	Size        int            `json:"size"`
	Page        int            `json:"page"`
	Count       int            `json:"count"`
	Filter      map[string]any `json:"filter"`
	ColumnsSort []ColumnSort   `json:"columnsSort"`
}

// Certificate keeps the original raw list-row API for compatibility.
type Certificate map[string]json.RawMessage

// CertificateSummary is the observed certificate list-row schema. Nullable
// replacement numbers remain pointers so null is not mistaken for an empty
// number. The upstream API spells manufacturer keys "manufacter".
type CertificateSummary struct {
	ID                                      int     `json:"id"`
	IDRALCertificationAuthority             int     `json:"idRalCertificationAuthority"`
	IDStatus                                int     `json:"idStatus"`
	Number                                  string  `json:"number"`
	BlankNumber                             string  `json:"blankNumber"`
	Date                                    string  `json:"date"`
	EndDate                                 string  `json:"endDate"`
	CertificateType                         string  `json:"certType"`
	CertificateObjectType                   string  `json:"certObjectType"`
	CertificateRegInsteadNumber             *string `json:"certRegInsteadNumber"`
	CertificateReplacedNumber               *string `json:"certReplacedNumber"`
	CertificationAuthorityAttestatRegNumber string  `json:"certificationAuthorityAttestatRegNumber"`
	ApplicantName                           string  `json:"applicantName"`
	ApplicantOPF                            string  `json:"applicantOpf"`
	ApplicantType                           string  `json:"applicantType"`
	ApplicantLegalSubjectType               string  `json:"applicantLegalSubjectType"`
	ApplicantFilialFullNames                string  `json:"applicantFilialFullNames"`
	ManufacturerName                        string  `json:"manufacterName"`
	ManufacturerOPF                         string  `json:"manufacterOpf"`
	ManufacturerType                        string  `json:"manufacterType"`
	ManufacturerLegalSubjectType            string  `json:"manufacterLegalSubjectType"`
	ManufacturerFilialFullNames             string  `json:"manufacterFilialFullNames"`
	ExpertName                              string  `json:"expertFio"`
	ExpertSNILS                             string  `json:"expertSnils"`
	Group                                   string  `json:"group"`
	ProductBatchSize                        string  `json:"productBatchSize"`
	ProductFullName                         string  `json:"productFullName"`
	ProductIdentificationArticle            string  `json:"productIdentificationArticle"`
	ProductIdentificationGTIN               string  `json:"productIdentificationGtin"`
	ProductIdentificationModel              string  `json:"productIdentificationModel"`
	ProductIdentificationName               string  `json:"productIdentificationName"`
	ProductIdentificationSort               string  `json:"productIdentificationSort"`
	ProductIdentificationTrademark          string  `json:"productIdentificationTrademark"`
	ProductIdentificationType               string  `json:"productIdentificationType"`
	ProductOrigin                           string  `json:"productOrig"`
	RecordSign                              string  `json:"recordSign"`
	TechnicalReglaments                     string  `json:"technicalReglaments"`
}

// CertificateColumnSearch uses the certificate registry's column/search
// shape; it has no declaration-style name/type fields.
type CertificateColumnSearch struct {
	Column string `json:"column"`
	Search string `json:"search"`
}

// CertificateDateRange sends explicit nulls for unset bounds, as observed in
// the certificate list request. The server format for nonempty bounds is not
// yet confirmed.
type CertificateDateRange struct {
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
}

// CertificateFilter models only confirmed fields. Extra allows callers to
// supply other upstream filters without letting them replace known fields.
type CertificateFilter struct {
	Status           []int                      `json:"status,omitempty"`
	IDCertType       []int                      `json:"idCertType,omitempty"`
	IDCertObjectType []int                      `json:"idCertObjectType,omitempty"`
	IDCertScheme     []int                      `json:"idCertScheme"`
	RegDate          CertificateDateRange       `json:"regDate"`
	EndDate          CertificateDateRange       `json:"endDate"`
	ColumnsSearch    []CertificateColumnSearch  `json:"columnsSearch"`
	Extra            map[string]json.RawMessage `json:"-"`
}

func (f CertificateFilter) MarshalJSON() ([]byte, error) {
	type known CertificateFilter
	base, err := json.Marshal(known(f))
	if err != nil {
		return nil, err
	}
	if len(f.Extra) == 0 {
		return base, nil
	}
	var merged map[string]json.RawMessage
	if err := json.Unmarshal(base, &merged); err != nil {
		return nil, err
	}
	for key, value := range f.Extra {
		switch key {
		case "status", "idCertType", "idCertObjectType", "idCertScheme", "regDate", "endDate", "columnsSearch":
			return nil, fmt.Errorf("reserved certificate filter field %q", key)
		}
		merged[key] = value
	}
	return json.Marshal(merged)
}

// CertificateQuery is the typed certificate search request. Unlike the
// declaration request, it has no count field.
type CertificateQuery struct {
	Size        int               `json:"size"`
	Page        int               `json:"page"`
	Filter      CertificateFilter `json:"filter"`
	ColumnsSort []ColumnSort      `json:"columnsSort"`
}

func DefaultCertificateQuery() CertificateQuery {
	return CertificateQuery{
		Size: 10,
		Filter: CertificateFilter{
			IDCertScheme:  []int{},
			ColumnsSearch: []CertificateColumnSearch{},
		},
		ColumnsSort: []ColumnSort{{Column: "date", Sort: "DESC"}},
	}
}
