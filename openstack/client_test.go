package openstack

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

func TestAuthenticatedClientV3(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	const ID = "0123456789"

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": {
					"values": [
						{
							"status": "stable",
							"id": "v3.0",
							"links": [
								{ "href": "%s", "rel": "self" }
							]
						},
						{
							"status": "stable",
							"id": "v2.0",
							"links": [
								{ "href": "%s", "rel": "self" }
							]
						}
					]
				}
			}
		`, testhelper.Endpoint()+"v3/", testhelper.Endpoint()+"v2.0/")
	})

	testhelper.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Subject-Token", ID)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{ "token": { "expires_at": "2013-02-02T18:30:59.000000Z" } }`)
	})

	options := gophercloud.AuthOptions{
		UserID:           "me",
		Password:         "secret",
		IdentityEndpoint: testhelper.Endpoint(),
	}
	client, err := AuthenticatedClient(options)

	if err != nil {
		t.Fatalf("Unexpected error from AuthenticatedClient: %s", err)
	}

	if client.TokenID != ID {
		t.Errorf("Expected token ID to be [%s], but was [%s]", ID, client.TokenID)
	}
}

func TestAuthenticatedClientV2(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"versions": {
					"values": [
						{
							"status": "experimental",
							"id": "v3.0",
							"links": [
								{ "href": "%s", "rel": "self" }
							]
						},
						{
							"status": "stable",
							"id": "v2.0",
							"links": [
								{ "href": "%s", "rel": "self" }
							]
						}
					]
				}
			}
		`, testhelper.Endpoint()+"v3/", testhelper.Endpoint()+"v2.0/")
	})

	testhelper.Mux.HandleFunc("/v2/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `
			{
				"access": {
					"token": {
						"id": "01234567890"
					},
					"serviceCatalog": [
						{
							"name": "Cloud Servers",
							"type": "compute",
							"endpoints": [
								{
									"tenantId": "t1000",
									"publicURL": "https://compute.north.host.com/v1/t1000",
									"internalURL": "https://compute.north.internal/v1/t1000",
									"region": "North",
									"versionId": "1",
									"versionInfo": "https://compute.north.host.com/v1/",
									"versionList": "https://compute.north.host.com/"
								},
								{
									"tenantId": "t1000",
									"publicURL": "https://compute.north.host.com/v1.1/t1000",
									"internalURL": "https://compute.north.internal/v1.1/t1000",
									"region": "North",
									"versionId": "1.1",
									"versionInfo": "https://compute.north.host.com/v1.1/",
									"versionList": "https://compute.north.host.com/"
								}
							],
							"endpoints_links": []
						},
						{
							"name": "Cloud Files",
							"type": "object-store",
							"endpoints": [
								{
									"tenantId": "t1000",
									"publicURL": "https://storage.north.host.com/v1/t1000",
									"internalURL": "https://storage.north.internal/v1/t1000",
									"region": "North",
									"versionId": "1",
									"versionInfo": "https://storage.north.host.com/v1/",
									"versionList": "https://storage.north.host.com/"
								},
								{
									"tenantId": "t1000",
									"publicURL": "https://storage.south.host.com/v1/t1000",
									"internalURL": "https://storage.south.internal/v1/t1000",
									"region": "South",
									"versionId": "1",
									"versionInfo": "https://storage.south.host.com/v1/",
									"versionList": "https://storage.south.host.com/"
								}
							]
						}
					]
				}
			}
		`)
	})

	options := gophercloud.AuthOptions{
		UserID:           "me",
		Password:         "secret",
		IdentityEndpoint: testhelper.Endpoint(),
	}
	client, err := AuthenticatedClient(options)

	if err != nil {
		t.Fatalf("Unexpected error from AuthenticatedClient: %s", err)
	}

	if client.TokenID != "01234567890" {
		t.Errorf("Expected token ID to be [01234567890], but was [%s]", client.TokenID)
	}
}
