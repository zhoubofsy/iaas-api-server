package common

import (
	"errors"

	"github.com/gophercloud/gophercloud"
)

func GetOpenstackClient(apikey string, tenant_id string, platform_userid string, resource_id ...string) (*gophercloud.ProviderClient, error) {
	return nil, errors.New("provider is nil")
}
