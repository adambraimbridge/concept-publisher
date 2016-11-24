package main

import (
	"testing"
	"net/url"
	"reflect"
	"github.com/pkg/errors"
	"strings"
)

func TestGetJobIds_Empty(t *testing.T) {
	url, err := url.Parse("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	var mockQueueSer queueServiceI = mockedQueueService{}
	var mockHttpSer httpServiceI = mockedHttpService{}
	pubService := newPublishService(url, &mockQueueSer, &mockHttpSer)
	actualIds := pubService.getJobIds()
	expectedIds := []string{}
	if !reflect.DeepEqual(actualIds, expectedIds) {
		t.Errorf("wrong ids list: got %v want %v", actualIds, expectedIds)
	}
}

func TestCreateJob(t *testing.T) {
	tests := []struct {
		clusterUrl     string
		baseUrl        string
		conceptType    string
		ids            []string
		throttle       int
		createErr      error
		idToTID        map[string]string
	}{
		{
			"http://localhost:8080",
			"http://localhost:8080/__special-reports-transformer/transformers/special-reports/",
			"special-reports",
			[]string{},
			1,
			nil,
			make(map[string]string),
		},
		{
			"http://localhost:8080",
			"http://localhost:8080/__special-reports-transformer/transformers/special-reports",
			"special-reports",
			[]string{},
			1,
			errors.New(`message="Can't find concept type in URL. Must be like the following __special-reports-transformer/transformers/special-reports/`),
			make(map[string]string),
		},
	}
	for _, test := range tests {
		clusterUrl, err := url.Parse(test.clusterUrl)
		if err != nil {
			t.Fatal(err)
		}
		var mockQueueSer queueServiceI = mockedQueueService{}
		var mockHttpSer httpServiceI = mockedHttpService{}
		pubService := newPublishService(clusterUrl, &mockQueueSer, &mockHttpSer)
		baseUrl, err := url.Parse(test.baseUrl)
		if err != nil {
			t.Fatalf("unexpected error. %v", err)
		}

		actualJob, err := pubService.createJob(test.ids, *baseUrl, test.throttle)

		if (err != nil) && (test.createErr != nil) {
			if !strings.HasPrefix(err.Error(), test.createErr.Error()) {
				t.Fatalf("unexpected error. diff got vs want:\n%v\n%v", err, test.createErr)
			}
			return
		}
		expectedJob := job{
			JobID: actualJob.JobID,
			ConceptType: "special-reports",
			IDToTID: test.idToTID,
			URL: *baseUrl,
			Throttle: test.throttle,
			Progress: 0,
			Status: defined,
			FailedIDs: []string{},
		}
		if !reflect.DeepEqual(*actualJob, expectedJob) {
			t.Fatalf("wrong job. diff got vs want:\n%v\n%v", *actualJob, expectedJob)
		}
	}
}

type mockedQueueService struct {
}

func (q mockedQueueService) sendMessage(id string, conceptType string, tid string, payload []byte) error {
	return nil
}

type mockedHttpService struct {
}

func (h mockedHttpService) reload(url string, authorization string) error {
	return nil
}

func (h mockedHttpService) getIds(url string, authorization string) ([]byte, *failure) {
	return []byte{}, nil
}

func (h mockedHttpService) getCount(url string, authorization string) (int, error) {
	return -1, nil
}

func (h mockedHttpService) fetchConcept(conceptID string, url string, authorization string) ([]byte, *failure) {
	return []byte{}, nil
}