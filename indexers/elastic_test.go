package indexers

import (
	"net/http"
	"net/http/httptest"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tests for elastic.go", func() {
	Context("Tests for new()", func() {
		payload := []byte(`{
			"name" : "0bcd132328f2f0c8ee451d471960750e",
			"cluster_name" : "415909267177:perfscale-dev",
			"cluster_uuid" : "Xz2IU4etSieAeaO2j-QCUw",
			"version" : {
			  "number" : "7.10.2",
			  "build_type" : "tar",
			  "build_hash" : "unknown",
			  "build_date" : "2023-03-22T14:16:51.874273Z",
			  "build_snapshot" : false,
			  "lucene_version" : "9.3.0",
			  "minimum_wire_compatibility_version" : "7.10.0",
			  "minimum_index_compatibility_version" : "7.0.0"
			},
			"tagline" : "The OpenSearch Project: https://opensearch.org/"
		  }`)
		testcase := struct {
			indexerConfig IndexerConfig
			mockServer    *httptest.Server
		}{
			indexerConfig: IndexerConfig{Type: "elastic",
				Servers:            []string{},
				Index:              "go-commons-test",
				InsecureSkipVerify: true,
			},
			mockServer: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(payload)
			})),
		}

		It("Returns error status bad request", func() {
			var indexer Elastic
			indexer.index = "go-commons-test"
			indexer.getnewclient = func(c elasticsearch.Config) (*elasticsearch.Client, error) {
				client, err := elasticsearch.NewDefaultClient()
				return client, err
			}
			defer testcase.mockServer.Close()
			testcase.indexerConfig.Servers = []string{testcase.mockServer.URL}
			err := indexer.new(testcase.indexerConfig)
			Expect(err).NotTo(BeNil())
		})

	})

	Context("Tests for Index()", func() {
		testcase := struct {
			documents []interface{}
			opts      IndexingOpts
		}{
			documents: []interface{}{
				"example document",
				42,
				3.14,
				false,
				struct {
					Name string
					Age  int
				}{
					Name: "John Doe",
					Age:  25,
				},
				map[string]interface{}{
					"key1": "value1",
					"key2": 123,
					"key3": true,
				}},
			opts: IndexingOpts{
				MetricName: "placeholder",
				JobName:    "placeholder",
			},
		}

		It("No err returned", func() {
			var indexer Elastic
			_, err := indexer.Index(testcase.documents, testcase.opts)
			Expect(err).To(BeNil())
		})

		It("err returned docs not processed", func() {
			testcase.documents = append(testcase.documents, make(chan string))
			var indexer Elastic
			_, err := indexer.Index(testcase.documents, testcase.opts)
			Expect(err).NotTo(BeNil())
		})

	})
})
