/*================================================================
*
*  文件名称：create_securitygroup_task.go
*  创 建 者: mongia
*  创建日期：2021年01月27日
*
================================================================*/

package securitygroupsvc

import (
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	sg "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"golang.org/x/net/context"

	"iaas-api-server/common"
	"iaas-api-server/proto/securitygroup"
)

type CreateSecurityGroupRPCTask struct {
	req *securitygroup.CreateSecurityGroupReq
	res *securitygroup.SecurityGroupRes
	err error
}

// TODO 后续直接调用common接口，把此函数删掉
func getProvider() *gophercloud.ProviderClient {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.66.128/identity",
		Username:         "admin",
		Password:         "secret",
		DomainName:       "default",
		TenantID:         "2b5149b85977407ea11667f580a62d80",
	}

	// Option 2: Use a utility function to retrieve all your environment variables
	//      opts, err := openstack.AuthOptionsFromEnv()
	//      if nil != err {
	//              fmt.Printf("openstack auth from env failed, %s", err)
	//              return
	//      }

	providers, err := openstack.AuthenticatedClient(opts)
	if nil != err {
		return nil
	}

	return providers
}

func (this *CreateSecurityGroupRPCTask) Run(context.Context) {
	if err := this.checkParam(); nil != err {
		this.err = err
		return
	}

	if !common.APIAuth(this.req.Apikey, this.req.TenantId, this.req.PlatformUserid) {
		this.err = errors.New("api auth failed.")
		return
	}

	providers := getProvider()
	if nil == providers {
		this.err = errors.New("get provider failed.")
		return
	}

	this.execute(providers)
}

func (this *CreateSecurityGroupRPCTask) execute(providers *gophercloud.ProviderClient) {
	client, err := openstack.NewNetworkV2(providers, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if nil != err {
		this.err = err
		return
	}

	createOpts := sg.CreateOpts{
		Name:        this.req.SecurityGroupName,
		Description: this.req.SecurityGroupDesc,
	}

	group, err := sg.Create(client, createOpts).Extract()
	if nil != err {
		this.err = err
		return
	}

	//TODO 返回参数
	this.res.Code = 0
	this.res.Msg = "success"
	this.res.SecurityGroup.SecurityGroupId = group.ID

}

func (this *CreateSecurityGroupRPCTask) checkParam() error {
	return nil
}
