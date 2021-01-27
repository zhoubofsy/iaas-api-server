# IaaS API设计文档

## 租户相关服务

```jsx
service TenantService{

//创建租户

rpc CreateTenant(CreateTenantReq) returns(CreateTenantRes);

}

message CreateTenantReq{

string tenant_name;

}

message CreateTenantRes{

string tenant_id;

string apikey;

}
```

### 租户注册

**说明**：为云管平台租户，生成调用本API的key，并为其在openstack中创建独立的域，项目、用户、角色等。本功能为grpc接口

**要点**：

1. 验证租户名，不允许重复
2. 生成tenant_id，id格式：t-10位流水号，例如：t-0000000001
3. 调用openstack api，创建openstack中的域、项目、用户、角色，并生成api key；域、项目、用户、角色的名称均按照tenant_id命名；openstack用户密码、api key随机生成，

相关信息保存到租户信息（tenant_info）表：

| 名称 | 数据类型 | 描述 | 主键 |
| ---  | --- | --- | --- |
| tenant_id | varchar2(20) | 租户id | 
| tenant_name | varchar2(100) | 租户名 | 
| openstack_domainname | varchar2(20) | openstack域名，取值=租户id | 
| openstack_domainid | varchar2(40) | openstack域ID | 
| openstack_projectname | varchar2(20) | openstack项目名，取值=租户id |
| openstack_projectid | varchar2(40) | openstack项目id |
| openstack_username | varchar2(20) | openstack用户名，取值=租户id | 
| openstack_userid | varchar2(40) | openstack用户id | 
| openstack_password | varchar2(20) | openstack用户密码 | 
| openstack_rolename | varchar2(20) | openstack角色，取值=租户id | 
| openstack_roleid | varchar2(40) | openstack角色id | 
| apikey | varchar2(20) | api key | 

### Api key验证

**说明**：非grpc接口，一个公用的函数，其他各grpc服务的方法中首先调用本函数

```jsx
func APIAuth(apikey string, tenant_id string, platform_userid string, resource_id ...string) (result boolen) {

//your logic

}
```

**要点**：

1. 根据传入的apikey、tenant_id到租户信息（tenant_info）表中验证api key合法性；
2. 如果传入了rrsource_id，根据传入的tenant_id、platform_userid、resource_id调用云管平台用户认证接口反向二次验证该平台用户是否有操作该资源的权限（接口待定）,resource_id是变长参数，可为0个或多个

以上二者均符合条件，方可返回true

## 规格相关服务

```jsx
service FlavorService{

//获取规格列表

rpc ListFlavors(ListFlavorsReq) returns(ListFlavorsRes);

//获取规格信息

rpc GetFlavor(GetFlavorReq) returns(Flavor);

}

message ListFlavorsReq {

string apikey;

string tenant_id;

string platform_userid;

int32 page_number;

int32 page_size;

}

message ListFlavorsRes {

repeated Flavor flavors;

}

message GetFlavorReq {

string apikey;

string tenant_id;

string platform_userid;

string flavor_id;

}

message Flavor {

string flavor_id;

string flavor_name;

string flavor_vcpus;

string flavor_ram;

string flavor_disk;

}
```

### 获取规格列表

**要点**：

1. 调用APIAuth函数鉴权
2. 根据tenant_id查tenant_info表，拿到该租户的openstack用户名、密码、project_id，domain_id去调openstack API或sdk
3. 为防止flavor数据太多，要根据传入的pagesize和page number
4. Flavor相关API在compute API v2中

### 获取规格信息

**要点**：

同上

## 镜像相关服务

```jsx
service ImageService{

//获取镜像列表

rpc ListImages(ListImagesReq) returns(ListImagesRes);

//获取镜像信息

rpc GetImage(GetImageReq) returns(Image);

}

rpc ListImages(ListImagesReq) returns(ListImagesRes);

message ListImagesReq {

string apikey;

string tenant_id;

string platform_userid;

int32 page_number;

int32 page_size;

}

message Image {

string image_id;

string image_name;

string image_diskformat;

string image_containerformat;

}

message ListImagesRes {

repeated Image images;

}

rpc GetImage(GetImageReq) returns(Image);

message GetImageReq {

string apikey;

string tenant_id;

string platform_userid;

string image_id;

}
```

### 获取镜像列表

**要点**

1. 鉴权 && 查表获取该租户的openstack连接参数
2. 分页处理
3. Image相关在compute api v2

### 获取镜像信息

**要点**：

同上

## 云主机相关服务

```jsx
service InstanceService{

//创建云主机

rpc CreateInstance(CreateInstanceReq) returns(Instance);

//获取云主机信息

rpc GetInstance(GetInstanceReq) returns(Instance);

//修改云主机规格

rpc UpdateInstanceFlavor(UpdateInstanceFlavorReq) returns(Instance);

//删除云主机

rpc DeleteInstance(DeleteInstanceReq) returns(DeleteInstanceRes);

//启动-停止-挂起-重启云主机

rpc OperateInstance(OperateInstanceReq) returns(OperateInstanceRes);

}

message CloudDiskConf {

string volume_type;

int32 size_in_g;

}

message Instance {

string instance_id;

string instance_status;

string instance_addr;

string region;

string availability_zone;

Flavor flavor;

string image_ref;

CloudDiskConf system_disk;

repeated CloudDiskConf data_disks;

string network_uuid;

repeated string security_group_name;

string instance_name;

string hypervisor_hostname;

string created_time;

string updated_time;

}

message CreateInstanceReq {

string apikey;

string tenant_id;

string platform_userid;

string region;

string availability_zone;

string flavor_ref;

string image_ref;

CloudDiskConf system_disk;

repeated CloudDiskConf data_disks;

string network_uuid;

repeated string security_group_name;

string instance_name;

string hypervisor_hostname;
}

message GetInstanceReq {

string apikey;

string tenant_id;

string platform_userid;

string instance_id;

}

message UpdateInstanceFlavorReq {

string apikey;

string tenant_id;

string platform_userid;

string instance_id;

string flavor_ref;

}

message DeleteInstanceReq {

string apikey;

string tenant_id;

string platform_userid;

string instance_id;

}

message DeleteInstanceRes {

string instance_id;

string deleted_time;

}

message OperateInstanceReq {

string apikey;

string tenant_id;

string platform_userid;

string instance_id;

string operate_type; //start/stop/softreboot/hardreboot/suspend

}

message OperateInstanceRes {

string instance_id;

string operated_time;

string operate_type;

}
```

### 创建云主机

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数
2. 对传入的volume_type暂不做处理，咱们目前还没有Qos差异化的云硬盘服务
3. 根据data_disks数组，创建多块数据盘volume;可调用本API种CreateCloudDisk接口
4. 系统盘（根卷），不先创建，创建server时通过field：block storage maping v2 那段设定，
5. 创建 server
6. 数据盘volume调用nova api的attachvolume 挂到server上

**讨论**：3、5两步各开一个[goroutine](https://www.google.com/search?client=safari&rls=en&q=goroutine&spell=1&sa=X&ved=2ahUKEwienPSlgLfuAhUUP30KHUC3CsoQkeECKAB6BAgKEDU)，主程序等signal，是否可行？

### 获取云主机信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数
2. Volume size的查找方法：volumes_attached.id，然后到block storage api中找

### 修改云主机规格

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数
2. Compute API 中 有resize server接口

### 删除云主机

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数
2. Compute api v2中直接有delete server的接口
3. 无需考虑挂载的数据volume连带删除问题

### 启动-停止-挂起-重启云主机

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数
2. Compute api v2中直接有这些接口

## 安全组相关服务

```jsx
service SecurityGroupService{

//创建安全组

rpc CreateSecurityGroup(CreateSecurityGroupReq) returns(SecurityGroup);

//获取安全组信息

rpc GetSecurityGroup(GetSecurityGroupReq) returns(SecurityGroup);

//修改安全组

rpc UpdateSecurityGroup(UpdateSecurityGroupReq) returns(SecurityGroup);

//删除安全组

rpc DeleteSecurityGroup(DeleteSecurityGroupReq) returns(DeleteSecurityGroupRes);

//安全组关联/取关云主机

rpc OperateSecurityGroup(OperateSecurityGroupReq) returns(OperateSecurityGroupRes);

}

message SecurityGroup {

string security_group_id;

string security_group_name;

string security_group_desc;

message SecurityGroupRule {

string rule_id;

string rule_desc;

string direction;

string protocol;

int32 port_range_min;

int32 port_range_max;

string remote_ip_prefix;

string security_group_id;

string created_time;

string updated_time;

}

repeated SecurityGroupRule security_group_rules;

string created_time;

string updated_time;

}

message SecurityGroupRuleSet {

string rule_desc;

string direction;

string protocol;

int32 port_range_min;

int32 port_range_max;

string remote_ip_prefix;

}

message CreateSecurityGroupReq {

string apikey;

string tenant_id;

string platform_userid;

string security_group_name;

string security_group_desc;

repeated SecurityGroupRuleSet security_group_rule_sets;

}

message GetSecurityGroupReq {

string apikey;

string tenant_id;

string platform_userid;

string security_group_id;

}

message UpdateSecurityGroupReq {

string apikey;

string tenant_id;

string platform_userid;

string security_group_id;

string security_group_name;

string security_group_desc;

repeated SecurityGroupRuleSet security_group_rule_sets;

}

message DeleteSecurityGroupReq {

string apikey;

string tenant_id;

string platform_userid;

string security_group_id;

}

message DeleteSecurityGroupRes {

string security_group_id;

string deleted_time;

}

message OperateSecurityGroupReq {

string apikey;

string tenant_id;

string platform_userid;

string security_group_id;

repeated string instance_ids;

string ops_type; //attach or detach

}

message OperateSecurityGroupRes {

string security_group_id;

string ops_type;

string operateed_time;

}
```

### 创建安全组

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数
2. 根据传入的security_group_rule_sets数组，多次调用security_group_rule API，创建security_group_rule数组
3. 调用create openstack security_group API，创建安全组，关联上一步做好的rules数组

### 获取安全组信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. Networking API v2中直接有get security_group detail的接口

### 修改安全组

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 调用 get security_group API，把安全组当前的rules数组拿到
3. 根据传入的security_group_rule_sets数组，多次调用security_group_rule API，创建security_group_rule数组
4. 调用修改 openstack security_group API，修改安全组，关联第3步做好的rules数组
5. 根据第2步拿到的旧rules数组，调用删除security_group_rule API；

### 删除安全组

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 删除安全组前应该验证是否有instance关联着这个安全组，这个验证交由云管平台做，底层API暂定不再做（因为要实现就需要遍历instaces）；
3. Openstack delete security_group api文档显示，已经实现了连带删除与之关联的security_group_rule，可以验证下，如无效，我们自行做连带删除

### 安全组关联/取关云主机

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. compute api v2里有add/remove security_group to server 接口
3. 根据传入的instance_ids多次调用 add/remove security_group to server接口

## 块存储相关服务

```jsx
service CloudDiskService{

//创建云硬盘

rpc CreateCloudDisk(CreateCloudDiskReq) returns(CloudDisk);

//获取云硬盘信息

rpc GetCloudDisk(GetCloudDiskReq) returns(CloudDisk);

//云硬盘扩容

rpc ReqizeCloudDisk(ReqizeCloudDiskReq) returns(CloudDisk);

//修改云硬盘信息

rpc ModifyCloudDiskInfo(ModifyCloudDiskInfoReq) returns(CloudDisk);

//云主机挂载/卸载云硬盘

rpc OperateCloudDisk(OperateCloudDiskReq) returns(CloudDisk);

//删除云硬盘

rpc DeleteCloudDisk(DeleteCloudDiskReq) returns(DeleteCloudDiskRes);

}

message CloudDisk {

string volume_id;

string volume_name;

string volume_desc;

string region;

string availability_zone;

CloudDiskConf cloud_disk_conf;

string volume_status;

string created_time;

string updated_time;

string attach_instance_id;

string attach_instance_device;

string attached_time;

}

message CreateCloudDiskReq {

string apikey;

string tenant_id;

string platform_userid;

string volume_name;

string volume_desc;

string region;

string availability_zone;

CloudDiskConf cloud_disk_conf;

}

message GetCloudDiskReq {

string apikey;

string tenant_id;

string platform_userid;

string volume_id;

}

message ReqizeCloudDiskReq {

string apikey;

string tenant_id;

string platform_userid;

string volume_id;

CloudDiskConf cloud_disk_conf;

}

message ModifyCloudDiskInfoReq {

string apikey;

string tenant_id;

string platform_userid;

string volume_id;

string volume_name;

string volume_desc;

}

message OperateCloudDiskReq {

string apikey;

string tenant_id;

string platform_userid;

string volume_id;

string instance_id;

string ops_type; //attach or detach

}

message DeleteCloudDiskReq {

string apikey;

string tenant_id;

string platform_userid;

string volume_id;

}

message DeleteCloudDiskRes {

string volume_id;

string deleted_time;

}
```

### 创建云硬盘

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. Block strorage api中有create volume方法；

### 获取云硬盘信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. Block strorage api中有show volume detail方法；

### 云硬盘扩容

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. Block strorage api 中有Extend a volume size方法

### 修改云硬盘信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. Block strorage api 中有update volume方法

### 云主机挂载/卸载云硬盘

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 使用compute API中的Attach/detach a volume to an instance接口似乎更靠谱

### 删除云硬盘

**要点**；

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. Block storage API中有delete volume的方法

## NAS存储相关服务

```jsx
service NasDiskService {

//创建NAS盘

rpc CreateNasDisk(CreateNasDiskReq) returns(NasDisk);

//删除NAS盘

rpc DeleteNasDisk(DeleteNasDiskReq) returns(DeleteNasDiskRes);

//查看挂载客户端

rpc GetMountClients(GetMountClientsReq) returns(MountClients);

}

message NasDisk {

string share_id;

string share_name;

string share_desci

string share_proto;

int32 share_size_in_g;

string region;

string network_id;

string share_status;

string share_progress;

string share_server_id;

string share_server_host;

string share_network_id;

string created_time;

string updated_time;

}

message CreateNasDiskReq {

string apikey;

string tenant_id;

string platform_userid;

string share_name;

string share_desc;

string share_proto;

int32 share_size_in_g;

string region;

string network_id;

}

message DeleteNasDiskReq {

string apikey;

string tenant_id;

string platform_userid;

string share_id;

}

message DeleteNasDiskRes {

string share_id;

string deleted_time;

}

message GetMountClientsReq {

string apikey;

string tenant_id;

string platform_userid;

string share_id;

}

message MountClients {

repeated string instance_id;

}
```

**注：以下描述是基于使用manila方案的实现，只看文档没使用过，不能保证准确性和完整性。**

### 创建NAS盘

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据传入的network_id，调用share filesystem api API，创建share_network，获得share_network_id;
3. 根据拿到的share_network_id，调用share filesystem api API，创建share，创建出的share中包含share_server_id;
4. 根据拿到的share_server_id，调用share filesystem api AP，获取share_server_host；

### 删除NAS盘

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据输入的share_id，share filesystem api中有delete share的方法，看文档说会连带删除share_server、share network资源，如无效，需我们自己做连带删除；

### 查看挂载客户端

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 还未确定实现方法，待补充；

## 对象存储相关服务

```jsx
service OSSService {

//创建Oss账号和bucket

rpc CreateUserAndBucket(CreateUserAndBucketReq) returns(CreateUserAndBucketRes);

//查看具体一个bucket详情

rpc GetBucketInfo(GetBucketInfoReq) returns(OssBucket);

//列出云管平台用户在一个ceph集群下的bucet列表

rpc ListBucketsInfo(ListBucketsInfoReq) returns(ListBucketsInfoRes);

//扩容oss_user配额

rpc SetOssUserQuota(SetOssUserQuotaReq) returns(OssUser);

//找回key

rpc RecoverKey(RecoverKeyReq) returns(RecoverKeyRes)

}

message OssUser {

string oss_uid;

string oss_user_created_time;

int32 user_max_size_in_g;

int32 user_max_objects;

int32 user_use_size_in_g;

int32 user_use_objects;

string user_created_time;

int32 total_buckets;

}

message CreateUserAndBucketReq {

string apikey;

string tenant_id;

string platform_userid;

string region;

string bucket_name;

string storege_type;

int32 user_max_size_in_g;

int32 user_max_objects;

string bucket_policy;

}

message CreateUserAndBucketRes {

string oss_endpoint;

string oss_access_key;

string oss_secret_key;

OssUser oss_user;

OssBucket oss_bucket;

}

message GetBucketInfoReq {

string apikey;

string tenant_id;

string platform_userid;

string region;

string oss_uid;

string bucket_name;

}

message ListBucketsInfoReq {

string apikey;

string tenant_id;

string platform_userid;

string region;

string oss_uid;

int32 page_number;

int32 page_size;

}

message ListBucketsInfoRes {

repeated OssBucket oss_buckets;

}

message SetOssUserQuotaReq {

string apikey;

string tenant_id;

string platform_userid;

string region;

string oss_uid;

int32 user_max_size_in_g;

int32 user_max_objects;

}

message RecoverKeyReq {

string apikey;

string tenant_id;

string platform_userid;

string region;

string oss_uid;

}

message RecoverKeyRes {

string oss_endpoint;

string oss_access_key;

string oss_secret_key;

}
```

### 首次创建账号并创建1个bucket

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 使用传入的platform_userid，作为要创建的S3用户的用户名，到ceph集群中查找，该userid是否存在
4. 如果不存在，创建该S3账户，并为其赋缺省配额限制，并根据传入的bucket名，bucket policy创建一个bucket
5. 如果存在，只创建bucket。

### 查看具体一个bucket详情

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 根据传入的oss_uid和bucket名查询

### 列出云管平台用户在一个ceph集群下的bucet列表

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 根据传入的uid，遍历出其下所有bucker信息
4. 根据传入的page_number、page_size做分页处理

### 扩容oss_user配额

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 根据传入的S3 uid 配额指标，修改该S3用户配额。

### 找回key

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 根据传入的uid，返回key信息

## 专有网络相关服务

```jsx
service VpcService {

//创建vpc

rpc CreateVpc(CreateVpcReq) returns(Vpc);

//查看vpc信息

rpc GetVpcInfo(GetVpcInfoReq) returns(Vpc)

//修改vpc信息

rpc SetVpcInfo(SetVpcInfoReq) returns(Vpc)

}

message Vpc {

string vpc_id;

string vpc_name;

string vpc_desc;

string region;

repeated string subnet;

string vpc_status;

string created_time;

}

message CreateVpcReq {

string apikey;

string tenant_id;

string platform_userid;

string vpc_name;

string vpc_desc;

string region;

string subnet;

}

message GetVpcInfoReq {

string apikey;

string tenant_id;

string platform_userid;

string vpc_id;

}

message SetVpcInfoReq {

string apikey;

string tenant_id;

string platform_userid;

string vpc_id;

string vpc_name;

string vpc_desc;

}
```

### 创建vpc

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据传入的vpc名称、描述，调用network API v2中networks create接口，拿到network_id
3. 根据查到的network_id，和输入的subnet，调用network API v2中subnet create接口

### 查看vpc信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据传入的vpc_id，调用network API v2中show network detail接口；

### 修改vpc信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 调用network API v2中network update接口

## 路由相关服务（待补充）

```jsx
service RouteService {

//创建专有网络路由器

//编辑路由表条目

}
```

## 统一注意事项

1. 调用openstack API，ceph API，有的响应比较慢，应有超时处理
2. 收到底层API抛出error，都触发向上（云管平台）抛出error，error信息可以隐藏细节信息
3. 尽量使用gophercloud sdk，减少原生openstack API调用
4. 返回参数值一定是从openstack查出来的真实信息，例如修改vpc描述信息，返回参数中也有描述信息字段，这个字段应该是调用openstack api查出来的，方能验证刚才的修改是否已生效。
5. 调用一些底层接口时，那些接口本身就是异步处理的，怎样获知已经处理完成？