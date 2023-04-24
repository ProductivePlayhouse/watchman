// Copyright 2022 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	moovhttp "github.com/moov-io/base/http"
	"github.com/moov-io/base/log"
	"github.com/moov-io/watchman/pkg/csl"

	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

var (
	matchHist = prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Name:    "match_percentages",
		Help:    "Histogram representing the match percent of search routes",
		Buckets: []float64{0.0, 0.5, 0.8, 0.9, 0.99},
	}, []string{"type"})
)

// TODO: modify existing search endpoint with additional eu info and add an eu only endpoint
// PPH MODIFIED
func addSearchRoutes(logger log.Logger, r *mux.Router, searcher *searcher) {
	r.Methods("GET").Path("/search").Handler(withAuth(logger, http.HandlerFunc(search(logger, searcher))))
	r.Methods("GET").Path("/search/us-csl").Handler(withAuth(logger, http.HandlerFunc(searchUSCSL(logger, searcher))))
	r.Methods("GET").Path("/search/eu-csl").Handler(withAuth(logger, http.HandlerFunc(searchEUCSL(logger, searcher))))
	r.Methods("GET").Path("/search/uk-csl").Handler(withAuth(logger, http.HandlerFunc(searchUKCSL(logger, searcher))))
}

func extractSearchLimit(r *http.Request) int {
	limit := softResultsLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		n, _ := strconv.Atoi(v)
		if n > 0 {
			limit = n
		}
	}
	if limit > hardResultsLimit {
		limit = hardResultsLimit
	}
	return limit
}

func extractSearchMinMatch(r *http.Request) float64 {
	if v := r.URL.Query().Get("minMatch"); v != "" {
		n, _ := strconv.ParseFloat(v, 64)
		return n
	}
	return 0.00
}

type addressSearchRequest struct {
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Providence string `json:"providence"`
	Zip        string `json:"zip"`
	Country    string `json:"country"`
}

func (req addressSearchRequest) empty() bool {
	return req.Address == "" && req.City == "" && req.State == "" &&
		req.Providence == "" && req.Zip == "" && req.Country == ""
}

func readAddressSearchRequest(u *url.URL) addressSearchRequest {
	return addressSearchRequest{
		Address:    strings.ToLower(strings.TrimSpace(u.Query().Get("address"))),
		City:       strings.ToLower(strings.TrimSpace(u.Query().Get("city"))),
		State:      strings.ToLower(strings.TrimSpace(u.Query().Get("state"))),
		Providence: strings.ToLower(strings.TrimSpace(u.Query().Get("providence"))),
		Zip:        strings.ToLower(strings.TrimSpace(u.Query().Get("zip"))),
		Country:    strings.ToLower(strings.TrimSpace(u.Query().Get("country"))),
	}
}

func search(logger log.Logger, searcher *searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(logger, w, r)
		requestID := moovhttp.GetRequestID(r)

		// Search over all fields
		if q := strings.TrimSpace(r.URL.Query().Get("q")); q != "" {
			logger.Info().With(log.Fields{
				"requestID": log.String(requestID),
			}).Logf("searching all names and address for %s", q)
			searchViaQ(logger, searcher, q)(w, r)
			return
		}

		// Search by ID (found in an SDN's Remarks property)
		if id := strings.TrimSpace(r.URL.Query().Get("id")); id != "" {
			logger.Info().With(log.Fields{
				"requestID": log.String(requestID),
			}).Logf("searching SDNs by remarks ID for %s", id)
			searchByRemarksID(logger, searcher, id)(w, r)
			return
		}

		// Search by Name
		if name := strings.TrimSpace(r.URL.Query().Get("name")); name != "" {
			if req := readAddressSearchRequest(r.URL); !req.empty() {
				logger.Info().With(log.Fields{
					"requestID": log.String(requestID),
				}).Logf("searching SDN names='%s' and addresses", name)
				searchViaAddressAndName(logger, searcher, name, req)(w, r)
				return
			}

			logger.Info().With(log.Fields{
				"requestID": log.String(requestID),
			}).Logf("searching SDN names for %s", name)
			searchByName(logger, searcher, name)(w, r)
			return
		}

		// Search by Alt Name
		if alt := strings.TrimSpace(r.URL.Query().Get("altName")); alt != "" {
			logger.Info().With(log.Fields{
				"requestID": log.String(requestID),
			}).Logf("searching SDN alt names for %s", alt)
			searchByAltName(logger, searcher, alt)(w, r)
			return
		}

		// Search Addresses
		if req := readAddressSearchRequest(r.URL); !req.empty() {
			logger.Info().With(log.Fields{
				"requestID": log.String(requestID),
			}).Logf("searching address for %#v", req)
			searchByAddress(logger, searcher, req)(w, r)
			return
		}

		// Fallback if no search params were found
		moovhttp.Problem(w, errNoSearchParams)
	}
}

type searchResponse struct {
	// DynamoDB - PPH Added
	Query   string       `json:"query" dynamodbav:"query"`
	Datetime int64       `json:"datetime" dynamodbav:"datetime"` // unix timestamp
	HighestMatch float64 `json:"highestMatch" dynamodbav:"highestMatch"`
	HighestMatchDetails string `json:"highestMatchDetails" dynamodbav:"highestMatchDetails"`

	// OFAC
	SDNs      []*SDN    `json:"SDNs" dynamodbav:"SDNs"`
	AltNames  []Alt     `json:"altNames" dynamodbav:"altNames"`
	Addresses []Address `json:"addresses" dynamodbav:"addresses"`

	// BIS
	DeniedPersons []DP `json:"deniedPersons" dynamodbav:"deniedPersons"`

	// Consolidated Screening List
	BISEntities                            []*Result[csl.EL]     `json:"bisEntities" dynamodbav:"bisEntities"`
	MilitaryEndUsers                       []*Result[csl.MEU]    `json:"militaryEndUsers" dynamodbav:"militaryEndUsers"`
	SectoralSanctions                      []*Result[csl.SSI]    `json:"sectoralSanctions" dynamodbav:"sectoralSanctions"`
	Unverified                             []*Result[csl.UVL]    `json:"unverifiedCSL" dynamodbav:"unverifiedCSL"`
	NonproliferationSanctions              []*Result[csl.ISN]    `json:"nonproliferationSanctions" dynamodbav:"nonproliferationSanctions"`
	ForeignSanctionsEvaders                []*Result[csl.FSE]    `json:"foreignSanctionsEvaders" dynamodbav:"foreignSanctionsEvaders"`
	PalestinianLegislativeCouncil          []*Result[csl.PLC]    `json:"palestinianLegislativeCouncil" dynamodbav:"palestinianLegislativeCouncil"`
	CaptaList                              []*Result[csl.CAP]    `json:"captaList" dynamodbav:"captaList"`
	ITARDebarred                           []*Result[csl.DTC]    `json:"itarDebarred" dynamodbav:"itarDebarred"`
	NonSDNChineseMilitaryIndustrialComplex []*Result[csl.CMIC]   `json:"nonSDNChineseMilitaryIndustrialComplex" dynamodbav:"nonSDNChineseMilitaryIndustrialComplex"`
	NonSDNMenuBasedSanctionsList           []*Result[csl.NS_MBS] `json:"nonSDNMenuBasedSanctionsList" dynamodbav:"nonSDNMenuBasedSanctionsList"`

	// EU - Consolidated Sanctions List
	EUCSL []*Result[csl.EUCSLRecord] `json:"euConsolidatedSanctionsList" dynamodbav:"euConsolidatedSanctionsList"`

	// UK - Consolidated Sanctions List
	UKCSL []*Result[csl.UKCSLRecord] `json:"ukConsolidatedSanctionsList" dynamodbav:"ukConsolidatedSanctionsList"`

	// Metadata
	RefreshedAt time.Time `json:"refreshedAt" dynamodbav:"refreshedAt"`
}

func buildAddressCompares(req addressSearchRequest) []func(*Address) *item {
	var compares []func(*Address) *item
	if req.Address != "" {
		compares = append(compares, topAddressesAddress(req.Address))
	}
	if req.City != "" {
		compares = append(compares, topAddressesCityState(req.City))
	}
	if req.State != "" {
		compares = append(compares, topAddressesCityState(req.State))
	}
	if req.Providence != "" {
		compares = append(compares, topAddressesCityState(req.Providence))
	}
	if req.Zip != "" {
		compares = append(compares, topAddressesCityState(req.Zip))
	}
	if req.Country != "" {
		compares = append(compares, topAddressesCountry(req.Country))
	}
	return compares
}

func searchByAddress(logger log.Logger, searcher *searcher, req addressSearchRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if req.empty() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := searchResponse{
			RefreshedAt: searcher.lastRefreshedAt,
		}
		limit := extractSearchLimit(r)
		minMatch := extractSearchMinMatch(r)

		// Perform our ranking across all accumulated compare functions
		//
		// TODO(adam): Is there something in the (SDN?) files which signal to block an entire country? (i.e. Needing to block Iran all together)
		// https://www.treasury.gov/resource-center/sanctions/CivPen/Documents/20190327_decker_settlement.pdf
		compares := buildAddressCompares(req)

		filtered := searcher.FilterCountries(req.Country)
		resp.Addresses = TopAddressesFn(limit, minMatch, filtered, multiAddressCompare(compares...))

		// record Prometheus metrics
		if len(resp.Addresses) > 0 {
			matchHist.With("type", "address").Observe(resp.Addresses[0].match)
		} else {
			matchHist.With("type", "address").Observe(0.0)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func searchViaQ(logger log.Logger, searcher *searcher, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name = strings.TrimSpace(name)
		if name == "" {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}
		limit := extractSearchLimit(r)
		minMatch := extractSearchMinMatch(r)

		// Perform multiple searches over the set of SDNs
		resp := buildFullSearchResponse(searcher, buildFilterRequest(r.URL), limit, minMatch, name)

		// record Prometheus metrics
		if len(resp.SDNs) > 0 {
			matchHist.With("type", "q").Observe(resp.SDNs[0].match)
		} else {
			matchHist.With("type", "q").Observe(0.0)
		}

		// Build our big response object
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// searchGather performs an inmem search with *searcher and mutates *searchResponse by setting a specific field
type searchGather func(searcher *searcher, filters filterRequest, limit int, minMatch float64, name string, resp *searchResponse)

var (
	baseGatherings = append([]searchGather{
		// OFAC SDN Search
		func(s *searcher, filters filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			sdns := s.FindSDNsByRemarksID(limit, name)
			if len(sdns) == 0 {
				sdns = s.TopSDNs(limit, minMatch, name, keepSDN(filters))
			}
			resp.SDNs = filterSDNs(sdns, filters)
		},
		// OFAC SDN Alt Names
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.AltNames = s.TopAltNames(limit, minMatch, name)
		},
		// OFAC Addresses
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.Addresses = s.TopAddresses(limit, minMatch, name)
		},

		// BIS Denied Persons
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.DeniedPersons = s.TopDPs(limit, minMatch, name)
		},
	}, cslGatherings...)

	// Consolidated Screening List Results
	cslGatherings = []searchGather{
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.BISEntities = s.TopBISEntities(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.MilitaryEndUsers = s.TopMEUs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.SectoralSanctions = s.TopSSIs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.Unverified = s.TopUVLs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.NonproliferationSanctions = s.TopISNs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.ForeignSanctionsEvaders = s.TopFSEs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.PalestinianLegislativeCouncil = s.TopPLCs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.CaptaList = s.TopCAPs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.ITARDebarred = s.TopDTCs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.NonSDNChineseMilitaryIndustrialComplex = s.TopCMICs(limit, minMatch, name)
		},
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.NonSDNMenuBasedSanctionsList = s.TopNS_MBS(limit, minMatch, name)
		},
	}

	// eu - consolidated sanctions list
	euGatherings = []searchGather{
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.EUCSL = s.TopEUCSL(limit, minMatch, name)
		},
	}

	// eu - consolidated sanctions list
	ukGatherings = []searchGather{
		func(s *searcher, _ filterRequest, limit int, minMatch float64, name string, resp *searchResponse) {
			resp.UKCSL = s.TopUKCSL(limit, minMatch, name)
		},
	}

	allGatherings = append(baseGatherings, euGatherings...)
)

func buildFullSearchResponse(searcher *searcher, filters filterRequest, limit int, minMatch float64, name string) *searchResponse {
	return buildFullSearchResponseWith(searcher, allGatherings, filters, limit, minMatch, name)
}

func buildFullSearchResponseWith(searcher *searcher, searchGatherings []searchGather, filters filterRequest, limit int, minMatch float64, name string) *searchResponse {
	resp := searchResponse{
		RefreshedAt: searcher.lastRefreshedAt,
	}
	var wg sync.WaitGroup
	wg.Add(len(searchGatherings))
	for i := range searchGatherings {
		go func(i int) {
			searchGatherings[i](searcher, filters, limit, minMatch, name, &resp)
			wg.Done()
		}(i)
	}
	wg.Wait()



	return &resp
}

func searchViaAddressAndName(logger log.Logger, searcher *searcher, name string, req addressSearchRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name = strings.TrimSpace(name)
		if name == "" || req.empty() {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		limit, minMatch := extractSearchLimit(r), extractSearchMinMatch(r)

		resp := &searchResponse{
			RefreshedAt: searcher.lastRefreshedAt,
		}

		resp.SDNs = searcher.TopSDNs(limit, minMatch, name, keepSDN(buildFilterRequest(r.URL)))

		compares := buildAddressCompares(req)
		filtered := searcher.FilterCountries(req.Country)
		resp.Addresses = TopAddressesFn(limit, minMatch, filtered, multiAddressCompare(compares...))

		// record Prometheus metrics
		if len(resp.SDNs) > 0 && len(resp.Addresses) > 0 {
			matchHist.With("type", "addressname").Observe(math.Max(resp.SDNs[0].match, resp.Addresses[0].match))
		} else {
			matchHist.With("type", "addressname").Observe(0.0)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func searchByRemarksID(logger log.Logger, searcher *searcher, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if id == "" {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		limit := extractSearchLimit(r)

		sdns := searcher.FindSDNsByRemarksID(limit, id)
		sdns = filterSDNs(sdns, buildFilterRequest(r.URL))

		// record Prometheus metrics
		if len(sdns) > 0 {
			matchHist.With("type", "remarksID").Observe(sdns[0].match)
		} else {
			matchHist.With("type", "remarksID").Observe(0.0)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&searchResponse{
			SDNs:        sdns,
			RefreshedAt: searcher.lastRefreshedAt,
		})
	}
}

func searchByName(logger log.Logger, searcher *searcher, nameSlug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nameSlug = strings.TrimSpace(nameSlug)
		if nameSlug == "" {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		limit := extractSearchLimit(r)
		minMatch := extractSearchMinMatch(r)

		// Grab the SDN's and then filter any out based on query params
		sdns := searcher.TopSDNs(limit, minMatch, nameSlug, keepSDN(buildFilterRequest(r.URL)))

		// record Prometheus metrics
		if len(sdns) > 0 {
			matchHist.With("type", "name").Observe(sdns[0].match)
		} else {
			matchHist.With("type", "name").Observe(0.0)
		}

		altNames := searcher.TopAltNames(limit, minMatch, nameSlug)
		sectoralSanctions := searcher.TopSSIs(limit, minMatch, nameSlug)
		deniedPersons := searcher.TopDPs(limit, minMatch, nameSlug)
		bisEntities := searcher.TopBISEntities(limit, minMatch, nameSlug)
		eucsl := searcher.TopEUCSL(limit, minMatch, nameSlug)
		ukcsl := searcher.TopUKCSL(limit, minMatch, nameSlug)

		// PPH Changed
		highestMatch := 0.0
		highestMatchDetails := ""

		if len(sdns) > 0 && sdns[0].match > highestMatch {
			highestMatch = sdns[0].match
			matchJSON, err := json.Marshal(sdns[0])
			if err == nil {
				highestMatchDetails = string(matchJSON)
			}
			else {
				highestMatchDetails = ""
				logger.Error().LogErrorf("Error marshalling JSON: %v", err)
			}
		}
		if len(altNames) > 0 && altNames[0].match > highestMatch {
			highestMatch = altNames[0].match
			matchJSON, err := json.Marshal(altNames[0])
			if err == nil {
				highestMatchDetails = string(matchJSON)
			}
			else {
				highestMatchDetails = ""
				logger.Error().LogErrorf("Error marshalling JSON: %v", err)
			}
		}
		if len(sectoralSanctions) > 0 && sectoralSanctions[0].match > highestMatch {
			highestMatch = sectoralSanctions[0].match
			matchJSON, err := json.Marshal(sectoralSanctions[0])
			if err == nil {
				highestMatchDetails = string(matchJSON)
			}
			else {
				highestMatchDetails = ""
				logger.Error().LogErrorf("Error marshalling JSON: %v", err)
			}
		}
		if len(deniedPersons) > 0 && deniedPersons[0].match > highestMatch {
			highestMatch = deniedPersons[0].match
			matchJSON, err := json.Marshal(deniedPersons[0])
			if err == nil {
				highestMatchDetails = string(matchJSON)
			}
			else {
				highestMatchDetails = ""
				logger.Error().LogErrorf("Error marshalling JSON: %v", err)
			}
		}
		if len(bisEntities) > 0 && bisEntities[0].match > highestMatch {
			highestMatch = bisEntities[0].match
			matchJSON, err := json.Marshal(bisEntities[0])
			if err == nil {
				highestMatchDetails = string(matchJSON)
			}
			else {
				highestMatchDetails = ""
				logger.Error().LogErrorf("Error marshalling JSON: %v", err)
			}
		}
		if len(eucsl) > 0 && eucsl[0].match > highestMatch {
			highestMatch = eucsl[0].match
			matchJSON, err := json.Marshal(eucsl[0])
			if err == nil {
				highestMatchDetails = string(matchJSON)
			}
			else {
				highestMatchDetails = ""
				logger.Error().LogErrorf("Error marshalling JSON: %v", err)
			}
		}
		if len(ukcsl) > 0 && ukcsl[0].match > highestMatch {
			highestMatch = ukcsl[0].match
			matchJSON, err := json.Marshal(ukcsl[0])
			if err == nil {
				highestMatchDetails = string(matchJSON)
			}
			else {
				highestMatchDetails = ""
				logger.Error().LogErrorf("Error marshalling JSON: %v", err)
			}
		}
	
		// PPH CHANGED
		resp := &searchResponse{
			// DynamoDB
			Query: nameSlug,
			Datetime: searcher.lastRefreshedAt.Unix(),
			HighestMatch: highestMatch,
			HighestMatchDetails: highestMatchDetails,
			// OFAC
			SDNs:              sdns,
			AltNames:          altNames,
			SectoralSanctions: sectoralSanctions,
			// BIS
			DeniedPersons: deniedPersons,
			BISEntities:   bisEntities,
			// EUCSL
			EUCSL: eucsl,
			// UKCSL
			UKCSL: ukcsl,
			// Metadata
			RefreshedAt: searcher.lastRefreshedAt,			
		}
		WriteQueryToDB(logger, *resp)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(&resp)
	}
}

func searchByAltName(logger log.Logger, searcher *searcher, altSlug string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		altSlug = strings.TrimSpace(altSlug)
		if altSlug == "" {
			moovhttp.Problem(w, errNoSearchParams)
			return
		}

		limit := extractSearchLimit(r)
		minMatch := extractSearchMinMatch(r)
		alts := searcher.TopAltNames(limit, minMatch, altSlug)

		// record Prometheus metrics
		if len(alts) > 0 {
			matchHist.With("type", "altName").Observe(alts[0].match)
		} else {
			matchHist.With("type", "altName").Observe(0.0)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&searchResponse{
			AltNames:    alts,
			RefreshedAt: searcher.lastRefreshedAt,
		})
	}
}
