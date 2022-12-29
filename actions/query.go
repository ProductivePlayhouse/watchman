// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package actions

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Movie encapsulates data about a movie. Title and Year are the composite primary key
// of the movie in Amazon DynamoDB. Title is the sort key, Year is the partition key,
// and Info is additional data.
type Movie struct {
	Title string                 `dynamodbav:"title"`
	Year  int                    `dynamodbav:"year"`
	Info  map[string]interface{} `dynamodbav:"info"`
}

// Create a type struct with the following data types
// 	"SDNs": {
// 		L: []*dynamodb.AttributeValue{},
// 	},
// 	"altNames": {
// 		L: []*dynamodb.AttributeValue{},
// 	},
// 	"addresses": {
// 		NULL: true,
// 	},
// 	"deniedPersons": {
// 		L: []*dynamodb.AttributeValue{},
// 	},
// 	"bisEntities": {
// 		L: []*dynamodb.AttributeValue{},
// 	},
// 	"militaryEndUsers": {
// 		NULL: true,
// 	},
// 	"sectoralSanctions": {
// 		L: []*dynamodb.AttributeValue{},
// 	},
// 	"unverifiedCSL": {
// 		NULL: true,
// 	},
// 	"nonproliferationSanctions": {
// 		NULL: true,
// 	},
// 	"foreignSanctionsEvaders": {
// 		NULL: true,
// 	},
// 	"palestinianLegislativeCouncil": {
// 		NULL: true,
// 	},
// 	"captaList": {
// 		NULL: true,
// 	},
// 	"itarDebarred": {
// 		NULL: true,
// 	},
// 	"nonSDNChineseMilitaryIndustrialComplex": {
// 		NULL: true,
// 	},
// 	"nonSDNMenuBasedSanctionsList": {
// 		NULL: true,
// 	},
// 	"euConsolidatedSanctionsList": {
// 		L: []*dynamodb.AttributeValue{},
// 	},
// 	"ukConsolidatedSanctionsList": {
// 		L: []*dynamodb.AttributeValue{},
// 	},
// 	"refreshedAt": {
// 		S: "2022-12-27T21:34:37.812466469Z",
// 	},
type QueryResponse struct {
	Query						     string   `dynamodbav:"query"`
	Datetime 					   	 string   `dynamodbav:"datetime"`	
	SDNs                             []string `dynamodbav:"SDNs"`
	AltNames                         []string `dynamodbav:"altNames"`
	Addresses                        string   `dynamodbav:"addresses"`
	DeniedPersons                    []string `dynamodbav:"deniedPersons"`
	BisEntities                      []string `dynamodbav:"bisEntities"`
	MilitaryEndUsers                 string   `dynamodbav:"militaryEndUsers"`
	SectoralSanctions                []string `dynamodbav:"sectoralSanctions"`
	UnverifiedCSL                    string   `dynamodbav:"unverifiedCSL"`
	NonproliferationSanctions        string   `dynamodbav:"nonproliferationSanctions"`
	ForeignSanctionsEvaders          string   `dynamodbav:"foreignSanctionsEvaders"`
	PalestinianLegislativeCouncil    string   `dynamodbav:"palestinianLegislativeCouncil"`
	CaptaList                        string   `dynamodbav:"captaList"`
	ItarDebarred                     string   `dynamodbav:"itarDebarred"`
	NonSDNChineseMilitaryIndustrial  string   `dynamodbav:"nonSDNChineseMilitaryIndustrialComplex"`
	NonSDNMenuBasedSanctionsList     string   `dynamodbav:"nonSDNMenuBasedSanctionsList"`
	EuConsolidatedSanctionsList      []string `dynamodbav:"euConsolidatedSanctionsList"`
	UkConsolidatedSanctionsList      []string `dynamodbav:"ukConsolidatedSanctionsList"`
}



// GetKey returns the composite primary key of the movie in a format that can be
// sent to DynamoDB.
func (movie Movie) GetKey() map[string]types.AttributeValue {
	title, err := attributevalue.Marshal(movie.Title)
	if err != nil {
		panic(err)
	}
	year, err := attributevalue.Marshal(movie.Year)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"title": title, "year": year}
}

// String returns the title, year, rating, and plot of a movie, formatted for the example.
func (movie Movie) String() string {
	return fmt.Sprintf("%v\n\tReleased: %v\n\tRating: %v\n\tPlot: %v\n",
		movie.Title, movie.Year, movie.Info["rating"], movie.Info["plot"])
}

// IMovieSampler defines an interface that can be used to download sample movie data
// from a URL.
type IMovieSampler interface {
	GetURL() string
	GetSampleMovies() []Movie
}

// MovieSampler implements IMovieSampler to download movie data from URL.
type MovieSampler struct {
	URL string
}

// GetURL returns the URL of the sampler.
func (sampler MovieSampler) GetURL() string {
	return sampler.URL
}

// GetSampleMovies downloads a .zip file of movie data, unzips it in memory, and
// unmarshals it into a Movie slice.
func (sampler MovieSampler) GetSampleMovies() []Movie {
	resp, err := http.Get(sampler.URL)
	if err != nil {
		log.Panicf("Couldn't get movie data from %v. Here's why: %v\n", sampler.URL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Couldn't read body of response. Here's why: %v\n", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Panicf("Couldn't create .zip reader. Here's why: %v\n", err)
	}

	zipFile := zipReader.File[0]
	zf, err := zipFile.Open()
	if err != nil {
		log.Panicf("Couldn't open first archive in .zip file. Here's why: %v\n", err)
	}
	defer zf.Close()
	movieBytes, err := ioutil.ReadAll(zf)
	if err != nil {
		log.Panicf("Couldn't read bytes from .zip archive. Here's why: %v\n", err)
	}

	var movies []Movie
	err = json.Unmarshal(movieBytes, &movies)
	if err != nil {
		log.Panicf("Couldn't unmarshal JSON data from .zip archive. Here's why: %v\n", err)
	}

	return movies
}
