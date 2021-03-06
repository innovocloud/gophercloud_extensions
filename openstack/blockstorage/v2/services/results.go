package services

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	osTime "github.com/innovocloud/gophercloud_extensions/openstack/time"
)

// ServicePage abstracts the raw results of making a List() request against
// the API. As OpenStack extensions may freely alter the response bodies of
// structures returned to the client, you may only safely access the data
// provided through the ExtractServices call.
type ServicePage struct {
	pagination.LinkedPageBase
}

// Service represents a blockstorage service in the OpenStack cloud.
type Service struct {
	Status         string               `json:"status"`
	Binary         string               `json:"binary"`
	Host           string               `json:"host"`
	Zone           string               `json:"zone"`
	State          string               `json:"state"`
	DisabledReason string               `json:"disabled_reason"`
	Updated        osTime.OpenStackTime `json:"updated_at"`
}

// IsEmpty returns true if a page contains no Services results.
func (r ServicePage) IsEmpty() (bool, error) {
	services, err := ExtractServices(r)
	return len(services) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r ServicePage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"services_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractServices interprets the results of a single page from a List() call,
// producing a slice of Services entities.
func ExtractServices(r pagination.Page) ([]Service, error) {
	var s []Service
	err := r.(ServicePage).Result.ExtractIntoSlicePtr(&s, "services")
	return s, err
}
