package fgis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// NSISelection requests named fields for a set of numeric or string NSI IDs.
// IDs must use the representation of the corresponding upstream dictionary.
type NSISelection struct {
	ID     []any    `json:"id"`
	Fields []string `json:"fields"`
}

// NSIQuery is the observed body of POST /nsi/api/multi.
type NSIQuery struct {
	Items map[string][]NSISelection `json:"items"`
}

// NSITNVED and NSIOKPD2 keep the internal numeric ID separate from the
// classification code. An internal ID is never itself a TNVED or OKPD2 code.
type NSITNVED struct {
	ID       int64  `json:"id"`
	MasterID string `json:"masterId"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Hidden   bool   `json:"hidden"`
}

type NSIOKPD2 struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Hidden bool   `json:"hidden"`
}

type NSIValidationFormNormDoc struct {
	ID       int64  `json:"id"`
	MasterID string `json:"masterId"`
	Name     string `json:"name"`
	DocNum   string `json:"docNum"`
}

type NSIValidationScheme struct {
	ID                       int64  `json:"id"`
	MasterID                 string `json:"masterId"`
	Name                     string `json:"name"`
	ValidityTerm             int64  `json:"validityTerm"`
	IsSeriesProduction       bool   `json:"isSeriesProduction"`
	IsBatchProduction        bool   `json:"isBatchProduction"`
	IsOneOffProduction       bool   `json:"isOneOffProduction"`
	IsProductSampleTesting   bool   `json:"isProductSampleTesting"`
	IsBatchProductTesting    bool   `json:"isBatchProductTesting"`
	IsOneOffProductTesting   bool   `json:"isOneOffProductTesting"`
	IsAccreditationLab       bool   `json:"isAccreditationLab"`
	IsApplicantManufacturer  bool   `json:"isApplicantManufacturer"`
	IsApplicantProvider      bool   `json:"isApplicantProvider"`
	IsPresenceOfProxy        bool   `json:"isPresenceOfProxy"`
	IsApplicantForeign       bool   `json:"isApplicantForeign"`
	IsApplicantEEUMember     bool   `json:"isApplicantEeuMember"`
	IsICProductionAnalysis   *bool  `json:"isIcProductionAnalysis"`
	IsICProductSampleTesting *bool  `json:"isIcProductSampleTesting"`
	IsICQMSValidation        *bool  `json:"isIcQmsValidation"`
}

type NSIStatus struct {
	ID       int64  `json:"id"`
	MasterID string `json:"masterId"`
	Name     string `json:"name"`
}

type NSICountry struct {
	ID        string `json:"id"`
	MasterID  int64  `json:"masterId"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Alpha2    string `json:"alpha2"`
	EEUMember *bool  `json:"eeuMember"`
}

// NSIResponse decodes known groups and retains every group, including future
// and unrecognised dictionaries, in Groups.
type NSIResponse struct {
	TNVED                 []NSITNVED                 `json:"-"`
	OKPD2                 []NSIOKPD2                 `json:"-"`
	ValidationFormNormDoc []NSIValidationFormNormDoc `json:"-"`
	ValidationScheme2     []NSIValidationScheme      `json:"-"`
	Status                []NSIStatus                `json:"-"`
	OKSM                  []NSICountry               `json:"-"`
	Groups                map[string]json.RawMessage `json:"-"`
}

func (r *NSIResponse) UnmarshalJSON(data []byte) error {
	var groups map[string]json.RawMessage
	if err := json.Unmarshal(data, &groups); err != nil {
		return err
	}
	var value NSIResponse
	value.Groups = groups
	for name, target := range map[string]any{
		"tnved":                 &value.TNVED,
		"okpd2":                 &value.OKPD2,
		"validationFormNormDoc": &value.ValidationFormNormDoc,
		"validationScheme2":     &value.ValidationScheme2,
		"status":                &value.Status,
		"oksm":                  &value.OKSM,
	} {
		if raw, ok := groups[name]; ok {
			if err := json.Unmarshal(raw, target); err != nil {
				return fmt.Errorf("decode NSI group %s: %w", name, err)
			}
		}
	}
	*r = value
	return nil
}

// GetNSIMulti looks up several dictionary groups in one request.
func (c *Client) GetNSIMulti(ctx context.Context, query NSIQuery) (NSIResponse, error) {
	if query.Items == nil {
		return NSIResponse{}, fmt.Errorf("NSI items must be set")
	}
	return requestJSON[NSIResponse](ctx, c, http.MethodPost, "/nsi/api/multi", c.startPath, query)
}

// ResolvedCode reports whether a dictionary ID was actually found. Code and
// Name remain nil when NSI has no matching record.
type ResolvedCode struct {
	ID   int64
	Code *string
	Name *string
}

type ProductCodes struct {
	TNVED []ResolvedCode
	OKPD2 []ResolvedCode
}

// ResolveProductCodes fetches both classifications in one NSI request and
// restores the input order and duplicates in the result.
func (c *Client) ResolveProductCodes(ctx context.Context, tnvedIDs, okpdIDs []int64) (ProductCodes, error) {
	result := ProductCodes{TNVED: make([]ResolvedCode, len(tnvedIDs)), OKPD2: make([]ResolvedCode, len(okpdIDs))}
	if len(tnvedIDs) == 0 && len(okpdIDs) == 0 {
		return result, nil
	}
	for _, ids := range [][]int64{tnvedIDs, okpdIDs} {
		for _, id := range ids {
			if id <= 0 {
				return ProductCodes{}, fmt.Errorf("classification ID must be positive: %d", id)
			}
		}
	}
	query := NSIQuery{Items: make(map[string][]NSISelection)}
	if ids := uniqueIDs(tnvedIDs); len(ids) > 0 {
		query.Items["tnved"] = []NSISelection{{ID: ids, Fields: []string{"id", "masterId", "name", "code", "hidden"}}}
	}
	if ids := uniqueIDs(okpdIDs); len(ids) > 0 {
		query.Items["okpd2"] = []NSISelection{{ID: ids, Fields: []string{"id", "name", "code", "hidden"}}}
	}
	response, err := c.GetNSIMulti(ctx, query)
	if err != nil {
		return ProductCodes{}, err
	}
	tnved := make(map[int64]NSITNVED, len(response.TNVED))
	for _, item := range response.TNVED {
		tnved[item.ID] = item
	}
	okpd := make(map[int64]NSIOKPD2, len(response.OKPD2))
	for _, item := range response.OKPD2 {
		okpd[item.ID] = item
	}
	for i, id := range tnvedIDs {
		result.TNVED[i].ID = id
		if item, ok := tnved[id]; ok {
			result.TNVED[i].Code, result.TNVED[i].Name = &item.Code, &item.Name
		}
	}
	for i, id := range okpdIDs {
		result.OKPD2[i].ID = id
		if item, ok := okpd[id]; ok {
			result.OKPD2[i].Code, result.OKPD2[i].Name = &item.Code, &item.Name
		}
	}
	return result, nil
}

func uniqueIDs(ids []int64) []any {
	seen := make(map[int64]struct{}, len(ids))
	result := make([]any, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
