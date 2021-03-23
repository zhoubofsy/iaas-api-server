/*================================================================
*
*  文件名称：APIAuth.go
*  创 建 者: tiantingting
*  创建日期：2021年02月18日
*
================================================================*/
package common

func APIAuth(apikey string, tenant_id string, platform_userid string, resource_id ...string) bool {
	var result bool
	//TODO 调用向明接口校验是否有操作权限,后续可能要加option用于表示是哪个模块，有什么操作权限，到时需要确认
	//查询本地是否存在相关租户信息（根据apiKey和tenant_id查询租户信息）
	tenantInfo,_:=QueryTenantInfoByTenantIdAndApikey(tenant_id,apikey)
	if !tenantInfo.IsEmpty(){
		result= true
	}
	return result
}
