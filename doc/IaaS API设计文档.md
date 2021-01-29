# IaaS API设计文档

## 租户相关服务

```jsx
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package tenant;

//租户相关服务
service TenantService{
  //创建租户
  rpc CreateTenant(CreateTenantReq) returns(CreateTenantRes);
}

message CreateTenantReq{
  string tenant_name = 1;
}

message CreateTenantRes{
  int32 code = 1;
  string msg = 2;
  string tenant_id = 3;
  string apikey = 4;
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
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package flavor;

//规格相关服务
service FlavorService{
  //获取规格列表
  rpc ListFlavors(ListFlavorsReq) returns(ListFlavorsRes);
  //获取规格信息
  rpc GetFlavor(GetFlavorReq) returns(GetFlavorRes);
}

message Flavor {
  string flavor_id = 1;
  string flavor_name = 2;
  string flavor_vcpus = 3;
  string flavor_ram = 4;
  string flavor_disk =5 ;
}

message ListFlavorsReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  int32 page_number = 4;
  int32 page_size = 5;
}

message ListFlavorsRes {
  int32 code = 1;
  string msg = 2;
  repeated Flavor flavors = 3;
}

message GetFlavorReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string flavor_id = 4;
}

message GetFlavorRes {
  int32 code = 1;
  string msg = 2;
  Flavor flavor = 3;
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
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package image;

//镜像相关服务
service ImageService{
  //获取镜像列表
  rpc ListImages(ListImagesReq) returns(ListImagesRes);
  //获取镜像信息
  rpc GetImage(GetImageReq) returns(GetImageRes);
}

message Image {
  string image_id = 1;
  string image_name = 2;
  string image_diskformat = 3;
  string image_containerformat = 4;
}

message ListImagesReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  int32 page_number = 4;
  int32 page_size = 5;
}

message ListImagesRes {
  int32 code = 1;
  string msg = 2;
  repeated Image images = 3;
}

message GetImageReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string image_id = 4;
}

message GetImageRes {
  int32 code = 1;
  string msg = 2;
  Image image = 3;
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
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package instance;

//云主机相关服务
service InstanceService{
  //创建云主机
  rpc CreateInstance(CreateInstanceReq) returns(InstanceRes);
  //获取云主机信息
  rpc GetInstance(GetInstanceReq) returns(InstanceRes);
  //修改云主机规格
  rpc UpdateInstanceFlavor(UpdateInstanceFlavorReq) returns(InstanceRes);
  //删除云主机
  rpc DeleteInstance(DeleteInstanceReq) returns(DeleteInstanceRes);
  //启动-停止-挂起-重启云主机
  rpc OperateInstance(OperateInstanceReq) returns(OperateInstanceRes);
}

message Flavor {
  string flavor_id = 1;
  string flavor_name = 2;
  string flavor_vcpus = 3;
  string flavor_ram = 4;
  string flavor_disk = 5;
}

message CloudDiskConf {
  string volume_type = 1;
  int32 size_in_g = 2;
}

message InstanceRes {
  int32 code = 1;
  string msg = 2;
  message Instance {
    string instance_id = 1;
    string instance_status = 2;
    string instance_addr = 3;
    string region = 4;
    string availability_zone = 5;
    Flavor flavor = 6;
    string image_ref = 7;
    CloudDiskConf system_disk = 8;
    repeated CloudDiskConf data_disks = 9;
    string network_uuid = 10;
    repeated string security_group_name = 11;
    string instance_name = 12;
    string hypervisor_hostname = 13;
    string created_time = 14;
    string updated_time = 15;
  }
  Instance instance = 3;
}

message CreateInstanceReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string region = 4;
  string availability_zone = 5;
  string flavor_ref = 6;
  string image_ref = 7;
  CloudDiskConf system_disk = 8;
  repeated CloudDiskConf data_disks = 9;
  string network_uuid = 10;
  repeated string security_group_name = 11;
  string instance_name = 12;
  string hypervisor_hostname = 13;
}

message GetInstanceReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string instance_id = 4;
}

message UpdateInstanceFlavorReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string instance_id = 4;
  string flavor_ref = 5;
}

message DeleteInstanceReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string instance_id = 4;
}

message DeleteInstanceRes {
  int32 code = 1;
  string msg = 2;
  string instance_id = 3;
  string deleted_time = 4;
}

message OperateInstanceReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string instance_id = 4;
  string operate_type = 5; //start/stop/softreboot/hardreboot/suspend
}

message OperateInstanceRes {
  int32 code = 1;
  string msg = 2;
  string instance_id = 3;
  string operated_time = 4;
  string operate_type = 5;
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
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package securitygroup;

//安全组相关服务
service SecurityGroupService{
  //创建安全组
  rpc CreateSecurityGroup(CreateSecurityGroupReq) returns(SecurityGroupRes);
  //获取安全组信息
  rpc GetSecurityGroup(GetSecurityGroupReq) returns(SecurityGroupRes);
  //修改安全组
  rpc UpdateSecurityGroup(UpdateSecurityGroupReq) returns(SecurityGroupRes);
  //删除安全组
  rpc DeleteSecurityGroup(DeleteSecurityGroupReq) returns(DeleteSecurityGroupRes);
  //安全组关联/取关云主机
  rpc OperateSecurityGroup(OperateSecurityGroupReq) returns(OperateSecurityGroupRes);
}

message SecurityGroupRes {
  int32 code = 1;
  string msg = 2;
  message SecurityGroup {
    string security_group_id = 1;
    string security_group_name = 2;
    string security_group_desc = 3;
    message SecurityGroupRule {
      string rule_id = 1;
      string rule_desc = 2;
      string direction = 3;
      string protocol = 4;
      int32  port_range_min = 5;
      int32 port_range_max = 6;
      string remote_ip_prefix = 7;
      string security_group_id = 8;
      string created_time = 9;
      string updated_time = 10;
    }
    repeated SecurityGroupRule security_group_rules = 4;
    string created_time = 5;
    string updated_time = 6;
  }
  SecurityGroup security_group =3;
}

message SecurityGroupRuleSet {
  string rule_desc = 1;
  string direction = 2;
  string protocol = 3;
  int32  port_range_min = 4;
  int32 port_range_max = 5;
  string remote_ip_prefix = 6;
}

message CreateSecurityGroupReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string security_group_name = 4;
  string security_group_desc = 5;
  repeated SecurityGroupRuleSet security_group_rule_sets = 6;
}

message GetSecurityGroupReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string security_group_id = 4;
}

message UpdateSecurityGroupReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string security_group_id = 4;
  string security_group_name = 5;
  string security_group_desc = 6;
  repeated SecurityGroupRuleSet security_group_rule_sets = 7;
}

message DeleteSecurityGroupReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string security_group_id = 4;
}

message DeleteSecurityGroupRes {
  int32 code = 1;
  string msg = 2;
  string security_group_id = 3;
  string deleted_time = 4;
}

message OperateSecurityGroupReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string security_group_id = 4;
  repeated string instance_ids = 5;
  string ops_type = 6; //attach or detach
}

message OperateSecurityGroupRes {
  int32 code = 1;
  string msg = 2;
  string security_group_id = 3;
  string ops_type = 4;
  string operateed_time = 5;
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
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package clouddisk;

//块存储相关服务
service CloudDiskService{
  //创建云硬盘
  rpc CreateCloudDisk(CreateCloudDiskReq) returns(CloudDiskRes);
  //获取云硬盘信息
  rpc GetCloudDisk(GetCloudDiskReq) returns(CloudDiskRes);
  //云硬盘扩容
  rpc ReqizeCloudDisk(ReqizeCloudDiskReq) returns(CloudDiskRes);
  //修改云硬盘信息
  rpc ModifyCloudDiskInfo(ModifyCloudDiskInfoReq) returns(CloudDiskRes);
  //云主机挂载/卸载云硬盘
  rpc OperateCloudDisk(OperateCloudDiskReq) returns(CloudDiskRes);
  //删除云硬盘
  rpc DeleteCloudDisk(DeleteCloudDiskReq) returns(DeleteCloudDiskRes);
}

message CloudDiskConf {
  string volume_type = 1;
  int32 size_in_g = 2;
}

message CloudDiskRes {
  int32 code = 1;
  string msg = 2;
  message CloudDisk {
    string volume_id = 1;
    string volume_name = 2;
    string volume_desc = 3;
    string region = 4;
    string availability_zone = 5;
    CloudDiskConf cloud_disk_conf = 6;
    string volume_status = 7;
    string created_time = 8;
    string updated_time = 9;
    string attach_instance_id = 10;
    string attach_instance_device = 11;
    string attached_time = 12;
  }
  CloudDisk cloud_disk = 3;
}

message CreateCloudDiskReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string volume_name = 4;
  string volume_desc = 5;
  string region = 6;
  string availability_zone = 7;
  CloudDiskConf cloud_disk_conf = 8;
}

message GetCloudDiskReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string volume_id = 4;
}

message ReqizeCloudDiskReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string volume_id = 4;
  CloudDiskConf cloud_disk_conf = 5;
}

message ModifyCloudDiskInfoReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string volume_id = 4;
  string volume_name = 5;
  string volume_desc = 6;
}

message OperateCloudDiskReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string volume_id = 4;
  string instance_id = 5;
  string ops_type = 6; //attach or detach
}

message DeleteCloudDiskReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string volume_id = 4;
}

message DeleteCloudDiskRes {
  int32 code = 1;
  string msg = 2;
  string volume_id = 3;
  string deleted_time = 4;
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
/ 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package nasdisk;

//NAS存储相关服务
service NasDiskService {
  //创建NAS盘
  rpc CreateNasDisk(CreateNasDiskReq) returns(CreateNasDiskRes);
  //删除NAS盘
  rpc DeleteNasDisk(DeleteNasDiskReq) returns(DeleteNasDiskRes);
  //查看挂载客户端
  rpc GetMountClients(GetMountClientsReq) returns(GetMountClientsRes);
}

message CreateNasDiskRes {
  int32 code = 1;
  string msg = 2;
  message NasDisk {
    string share_id = 1;
    string share_name = 2;
    string share_desc =3;
    string share_proto = 4;
    int32  share_size_in_g = 5;
    string region = 6;
    string  network_id = 7;
    string share_status = 8;
    string share_progress = 9;
    string share_server_id = 10;
    string share_server_host = 11;
    string share_network_id = 12;
    string created_time = 13;
    string updated_time = 14;
  }
  NasDisk nas_disk = 3;
}

message CreateNasDiskReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string share_name = 4;
  string share_desc = 5;
  string share_proto = 6;
  int32  share_size_in_g = 7;
  string region = 8;
  string network_id = 9;
}

message DeleteNasDiskReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string share_id = 4;
}

message DeleteNasDiskRes {
  int32 code = 1;
  string msg = 2;
  string share_id = 3;
  string deleted_time = 4;
}

message GetMountClientsReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string share_id = 4;
}

message GetMountClientsRes {
  int32 code = 1;
  string msg = 2;
  repeated string instance_id = 3;
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
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package oss;

//对象存储相关服务
service OSSService {
  //创建Oss账号和bucket
  rpc CreateUserAndBucket(CreateUserAndBucketReq) returns(CreateUserAndBucketRes);
  //查看具体一个bucket详情
  rpc GetBucketInfo(GetBucketInfoReq) returns(GetBucketInfoRes);
  //列出云管平台用户在一个entrypoint下的bucet列表
  rpc ListBucketsInfo(ListBucketsInfoReq) returns(ListBucketsInfoRes);
  //扩容oss_user配额
  rpc SetOssUserQuota(SetOssUserQuotaReq) returns(SetOssUserQuotaRes);
  //找回key
  rpc RecoverKey(RecoverKeyReq) returns(RecoverKeyRes);
}

message OssBucket {
  string bucket_name = 1;
  string bucket_policy = 2; //private/public-ro/public-rw
  int32 bucket_use_size_in_g = 3;
  int32 bucket_use_objects = 4;
  string belong_to_uid = 5;
  string bucket_created_time = 6;
}

message OssUser {
  string oss_uid = 1;
  string oss_user_created_time = 2;
  int32 user_max_size_in_g = 3;
  int32 user_max_objects = 4;
  int32 user_use_size_in_g = 5;
  int32 user_use_objects = 6;
  string user_created_time = 7;
  int32 total_buckets = 8;
}

message CreateUserAndBucketReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string region = 5;
  string bucket_name = 6;
  string storege_type = 7;
  int32 user_max_size_in_g = 8;
  int32 user_max_objects = 9;
  string bucket_policy = 10;
}

message CreateUserAndBucketRes {
  int32 code = 1;
  string msg = 2;
  string oss_endpoint = 3;
  string oss_access_key = 4;
  string oss_secret_key = 5;
  OssUser oss_user = 6;
  OssBucket oss_bucket = 7;
}

message GetBucketInfoReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string region = 4;
  string oss_uid = 5;
  string bucket_name = 6;
}

message GetBucketInfoRes {
  int32 code = 1;
  string msg = 2;
  OssBucket oss_bucket = 3;
}

message ListBucketsInfoReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string region = 4;
  string oss_uid = 5;
  int32 page_number = 6;
  int32 page_size = 7;
}

message ListBucketsInfoRes {
  int32 code = 1;
  string msg = 2;
  repeated OssBucket oss_buckets = 3;
}

message SetOssUserQuotaReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string region = 4;
  string oss_uid = 5;
  int32 user_max_size_in_g = 6;
  int32 user_max_objects = 7;
}

message SetOssUserQuotaRes {
  int32 code = 1;
  string msg = 2;
  OssUser oss_user = 3;
}

message RecoverKeyReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string region = 4;
  string oss_uid = 5;
}

message RecoverKeyRes {
  int32 code = 1;
  string msg = 2;
  string oss_endpoint = 3;
  string oss_access_key = 4;
  string oss_secret_key = 5;
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
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package vpc;

//专有网络相关服务
service VpcService {
  //创建vpc
  rpc CreateVpc(CreateVpcReq) returns(VpcRes);
  //查看vpc信息
  rpc GetVpcInfo(GetVpcInfoReq) returns(VpcRes);
  //修改vpc信息
  rpc SetVpcInfo(SetVpcInfoReq) returns(VpcRes);
}

message VpcRes {
  int32 code = 1;
  string msg = 2;
  message Vpc {
    string vpc_id = 1;
    string vpc_name = 2;
    string vpc_desc = 3;
    string region = 4;
    repeated string subnet = 5;
    string vpc_status = 6;
    string created_time = 7;
  }
  Vpc vpc = 3;
}


message CreateVpcReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string vpc_name = 4;
  string vpc_desc = 5;
  string region = 6;
  string subnet = 7;
}

message GetVpcInfoReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string vpc_id = 4;
}

message SetVpcInfoReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string vpc_id = 4;
  string vpc_name = 5;
  string vpc_desc = 6;
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
