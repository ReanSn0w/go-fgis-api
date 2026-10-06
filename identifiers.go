package fgis

import (
	"context"
	"encoding/json"
	"net/http"
)

// IdentifierEntry is a named dictionary value with a numeric ID. Raw retains
// upstream attributes that differ between individual dictionaries.
type IdentifierEntry struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	IDGroup   *int64          `json:"idGroup"`
	MasterID  json.RawMessage `json:"masterId"`
	ShortName *string         `json:"shortName"`
	Raw       json.RawMessage `json:"-"`
}

func (v *IdentifierEntry) UnmarshalJSON(data []byte) error {
	type plain IdentifierEntry
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*v = IdentifierEntry(value)
	v.Raw = append([]byte(nil), data...)
	return nil
}

// CountryIdentifier keeps the string OKSM ID distinct from numeric dictionary IDs.
type CountryIdentifier struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Alpha2    string          `json:"alpha2"`
	ShortName *string         `json:"shortName"`
	Raw       json.RawMessage `json:"-"`
}

func (v *CountryIdentifier) UnmarshalJSON(data []byte) error {
	type plain CountryIdentifier
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*v = CountryIdentifier(value)
	v.Raw = append([]byte(nil), data...)
	return nil
}

// DeclarationSchemeIdentifier retains all scheme-specific flags in Raw.
// Only ID and Name are common enough to type from the observed response.
type DeclarationSchemeIdentifier struct {
	ID   int64           `json:"id"`
	Name string          `json:"name"`
	Raw  json.RawMessage `json:"-"`
}

func (v *DeclarationSchemeIdentifier) UnmarshalJSON(data []byte) error {
	type plain DeclarationSchemeIdentifier
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*v = DeclarationSchemeIdentifier(value)
	v.Raw = append([]byte(nil), data...)
	return nil
}

// DeclarationIdentifiers preserves the registry's dictionary names, including
// notActive in its status map. Entries not covered by a common value schema
// remain raw JSON.
type DeclarationIdentifiers struct {
	ActiveID           int64                                  `json:"activeId"`
	AddrTypes          map[string]IdentifierEntry             `json:"addrTypes"`
	AnnexTypes         map[string]IdentifierEntry             `json:"annexTypes"`
	ApplicantTypes     map[string]IdentifierEntry             `json:"applicantTypes"`
	CanceledID         int64                                  `json:"canceledId"`
	ContactType        map[string]IdentifierEntry             `json:"contactType"`
	CountryCode        map[string]IdentifierEntry             `json:"countryCode"`
	DeclObjectTypes    map[string]IdentifierEntry             `json:"declObjectTypes"`
	DeclScheme         map[string]DeclarationSchemeIdentifier `json:"declScheme"`
	DeclType           map[string]IdentifierEntry             `json:"declType"`
	DocTypeGroups      map[string]IdentifierEntry             `json:"docTypeGroups"`
	DocTypes           map[string]IdentifierEntry             `json:"docTypes"`
	LegalSubjectTypes  map[string]IdentifierEntry             `json:"legalSubjectTypes"`
	NormDocStatuses    map[string]IdentifierEntry             `json:"normDocStatuses"`
	NormDocTypes       map[string]IdentifierEntry             `json:"normDocTypes"`
	NormDocs           map[string]json.RawMessage             `json:"normDocs"`
	NotActiveID        int64                                  `json:"notActiveId"`
	OKSM               map[string]CountryIdentifier           `json:"oksm"`
	RALType            map[string]IdentifierEntry             `json:"ralType"`
	ReplacementReason  map[string]IdentifierEntry             `json:"replacementReason"`
	ScanCopyTypes      map[string]json.RawMessage             `json:"scanCopyTypes"`
	Status             map[string]IdentifierEntry             `json:"status"`
	StatusChangeReason map[string]IdentifierEntry             `json:"statusChangeReason"`
	StatusTestingLabs  map[string]IdentifierEntry             `json:"statusTestingLabs"`
	SuspendedID        int64                                  `json:"suspendedId"`
	ValidationForms    map[string]IdentifierEntry             `json:"validationForms"`
	Raw                json.RawMessage                        `json:"-"`
}

func (v *DeclarationIdentifiers) UnmarshalJSON(data []byte) error {
	type plain DeclarationIdentifiers
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*v = DeclarationIdentifiers(value)
	v.Raw = append([]byte(nil), data...)
	return nil
}

// CertificateIdentifiers is separate because its dictionary set and spelling
// of some keys differ from DeclarationIdentifiers.
type CertificateIdentifiers struct {
	AddrTypes              map[string]IdentifierEntry   `json:"addrTypes"`
	AnnexTypes             map[string]IdentifierEntry   `json:"annexTypes"`
	ApplicantTypes         map[string]IdentifierEntry   `json:"applicantTypes"`
	Basis                  map[string]IdentifierEntry   `json:"basis"`
	BlankTypes             map[string]IdentifierEntry   `json:"blankTypes"`
	CertObjectTypes        map[string]IdentifierEntry   `json:"certObjectTypes"`
	CertStatusChangeReason map[string]IdentifierEntry   `json:"certStatusChangeReason"`
	CertType               map[string]IdentifierEntry   `json:"certType"`
	ContactType            map[string]IdentifierEntry   `json:"contactType"`
	CountryCode            map[string]IdentifierEntry   `json:"countryCode"`
	DocTypeGroups          map[string]IdentifierEntry   `json:"docTypeGroups"`
	DocTypes               map[string]IdentifierEntry   `json:"docTypes"`
	InspectionControlTypes map[string]IdentifierEntry   `json:"inspectionControlTypes"`
	LegalSubjectTypes      map[string]IdentifierEntry   `json:"legalSubjectTypes"`
	NormDocStatuses        map[string]IdentifierEntry   `json:"normDocStatuses"`
	NormDocTypes           map[string]IdentifierEntry   `json:"normDocTypes"`
	NormDocs               map[string]json.RawMessage   `json:"normDocs"`
	OKSM                   map[string]CountryIdentifier `json:"oksm"`
	RALType                map[string]IdentifierEntry   `json:"ralType"`
	Status                 map[string]IdentifierEntry   `json:"status"`
	ValidationForms        map[string]IdentifierEntry   `json:"validationForms"`
	Raw                    json.RawMessage              `json:"-"`
}

func (v *CertificateIdentifiers) UnmarshalJSON(data []byte) error {
	type plain CertificateIdentifiers
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*v = CertificateIdentifiers(value)
	v.Raw = append([]byte(nil), data...)
	return nil
}

func (c *Client) GetDeclarationIdentifiers(ctx context.Context) (DeclarationIdentifiers, error) {
	return requestJSON[DeclarationIdentifiers](ctx, c, http.MethodGet, "/api/v1/rds/common/identifiers", "/rds/declaration", nil)
}

func (c *Client) GetCertificateIdentifiers(ctx context.Context) (CertificateIdentifiers, error) {
	return requestJSON[CertificateIdentifiers](ctx, c, http.MethodGet, "/api/v1/rss/common/identifiers", "/rss/certificate", nil)
}
