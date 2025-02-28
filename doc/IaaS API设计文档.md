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
  string start_flavor_id = 4;
  int32 page_size = 5;
}

message ListFlavorsRes {
  int32 code = 1;
  string msg = 2;
  repeated Flavor flavors = 3;
  string next_flavor_id = 4;
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
  float image_size_in_g = 5;
}

message ListImagesReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string start_image_id = 4;
  int32 page_size = 5;
}

message ListImagesRes {
  int32 code = 1;
  string msg = 2;
  repeated Image images = 3;
  string next_image_id = 4;
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
    string guest_os_hostname = 13;
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
  string guest_os_hostname = 13;
  string root_pass = 14;
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
4. 系统盘（根卷），不先创建，创建server时通过field：block storage maping v2 那段设定，数据盘volume，也是过field：block storage maping v2 那段设定
5. 创建 server
6. 自定义主机名，自定义root密码，通过在user_data filed中注入脚本实现。

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
// 指定的当前proto语法的版本，有2和3
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
    string mount_point = 8;
    string created_time = 9;
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
  string region = 5;
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
  string region = 5;
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
  //列出云管平台用户在一个entrypoint下的bucket列表
  rpc ListBucketsInfo(ListBucketsInfoReq) returns(ListBucketsInfoRes);
  //扩容oss_user配额
  rpc SetOssUserQuota(SetOssUserQuotaReq) returns(SetOssUserQuotaRes);
  //找回key
  rpc RecoverKey(RecoverKeyReq) returns(RecoverKeyRes);
  //获取用户信息
  rpc GetUserInfo(GetUserInfoReq) returns(GetUserInfoRes);
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

message GetUserInfoReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string region = 4;
  string oss_uid = 5;
}

message GetUserInfoRes {
  int32 code = 1;
  string msg = 2;
  OssUser oss_user = 3;
}
```

### 创建Oss账号和bucket

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

### 列出云管平台用户在一个ceph集群下的bucket列表

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 根据传入的uid，遍历出其下所有bucker信息
4. 根据传入的page_number、page_size做分页处理

### 扩容oss_user配额

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 根据传入的uid和该用户的配额指标，修改该S3用户配额。
4. 根据传入的uid进行查询，返回当前该S3用户信息

### 找回key

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 根据传入的uid，返回key信息

### 获取oss用户信息

**要点**：

1. 鉴权，无需查表获取该租户的openstack连接参数;
2. 根据传入的region判断S3用户应建立在哪个ceph集群（每个region一个ceph集群）
3. 根据传入的uid进行查询，返回当前该S3用户信息

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
    message Subnet {
      string subnet = 1;
      string subnet_id = 2;
      string subnet_created_time = 3;
    }
    repeated Subnet subnet = 5;
    string vpc_status = 6;
    string vpc_created_time = 7;
    message Router {
      string router_id = 1;
      string router_name = 2;
      string router_created_time = 3;
      message Intf {
        string intf_id = 1;
        string intf_name = 2;
        string intf_ip = 3;
        string subnet_id = 4;
        string intf_created_time = 5;
      }
      repeated Intf Intfs = 4;
    }
    Router router = 8;
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
2. 根据传入的vpc名称、描述，调用network API v2中networks create接口，拿到network_id等信息
3. 根据拿到的network_id，和输入的subnet，调用network API v2中subnet create接口，创建subnet，子网网关IP设定为该sunnet的第1个IP，例如192.168.1.0/24，就使用192.168.1.1；并启用dhcp，dhcp池从第2个IP开始到该子网的最后一个（注意最后一个不一定是.254。要看子网的netmask是多少，子网地址范围，网上有专门计算这个函数，可以拿来用）创建完成后获得subnet_id等信息
4. 调用network API v2中router create接口，创建一个路由器，路由器的name命名规则：router-vpc名称；创建完成后拿到router_id等信息
5. 根据router_id, subnet_id，为上一步建好的路由器增加第一个接口，调用network API v2 add_router_interface方法创建接口，接口IP为第2步创建的子网网关IP，创建完成后获得interface_id等信息
6. 一个vpc下只创建一个router
7. 目前一个vpc下，仅支持创建一个subnet即可，vpc返回值中，subnet为数组，只是为了预留，也方便与openstack概念对应

### 查看vpc信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据传入的vpc_id，调用network API v2中show network detail接口；

### 修改vpc信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 调用network API v2中network update接口

## 路由相关服务

```jsx
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package route;

//路由相关服务
service RouterService {
  //获取路由器信息
  rpc GetRouter(GetRouterReq) returns(GetRouterRes);
  //添加或删除路由表条目
  rpc SetRoutes(SetRoutesReq) returns(SetRoutesRes);
}

message GetRouterReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string router_id = 4;
}

message GetRouterRes {
  int32 code = 1;
  string msg = 2;
  message Router {
    string router_id = 1;
    string router_name = 2;
    string router_created_time = 3;
    message Intf {
      string intf_id = 1;
      string intf_name = 2;
      string intf_ip = 3;
      string subnet_id = 4;
      string intf_created_time = 5;
    }
    repeated Intf Intfs = 4;
    repeated Route current_routes = 5;
  }
  Router router = 3;
}

message Route {
  string destination = 1;
  string nexthop = 2;
}

message SetRoutesReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string router_id = 4;
  string set_type = 5; //add or remove
  repeated Route routes = 6;
}

message SetRoutesRes {
  int32 code = 1;
  string msg = 2;
  string router_id = 3;
  string set_type = 4;
  string set_time = 5;
  repeated Route current_routes = 6;
}
```

### 获取路由器信息

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据router_id调用network API v2中show router detail接口
3. 获得该路由器所有interfaces信息（实现方法待补充）

### 添加或删除路由表条目

**要点**：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据router_id和set_type，调用network API v2中Add extra routes to router/Remove extra routes to router接口
3. 返回值中的current_routes数组，是该路由器当前所有的路由条目（openstack api返回的正好也是这样的）

## NAT网关相关服务

```
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package natgateway;

//NAT网关相关服务
service NatGatewayService {
  //创建NAT网关
  rpc CreateNatGateway(CreateNatGatewayReq) returns(NatGatewayRes);
  //获取NAT网关信息
  rpc GetNatGateway(GetNatGatewayReq) returns(NatGatewayRes);
  //删除NAT网关
  rpc DeleteNatGateway(DeleteNatGatewayReq) returns(DeleteNatGatewayRes);
}

message NatGatewayRes {
  int32 code = 1;
  string msg = 2;
  message NatGateway {
    string gateway_id = 1;
    string router_id = 2;
    string external_network_id = 3;
    bool enable_snat = 4;
    string external_fixed_ip = 5;
    string created_time = 6;
  }
  NatGateway nat_gateway = 3;
}

message CreateNatGatewayReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string router_id = 4;
  string external_network_id = 5;
}

message GetNatGatewayReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string router_id = 4;
  string gateway_id = 5;
}

message DeleteNatGatewayReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string router_id = 4;
  string gateway_id = 5;
}

message DeleteNatGatewayRes {
  int32 code = 1;
  string msg = 2;
  string router_id = 3;
  string gateway_id = 4;
  string deleted_time = 5;
}
```

### 创建NAT网关

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据传入的routerID，调用network V2 sdk中routers.Update方法，为router增加Gatewayinfo field，Gatewayinfo中NetworkID取自external_network_id，EnableSNAT默认设置true；
3. 创建完成后，会返回exterbal_fixed_ip，取出这个fix ip结构体的subnet id作为返回参数中的gateway_id;

### 获取NAT网关信息

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据传入的router_id，gateway_id（即路由器连接外网的那个接口的ID），使用network V2 sdk中routers.Get相关方法，获取完整的router对象，从中解析出GatewayInfo filed，拼装返回参数

### 删除NAT网关

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据传入的router_id，gateway_id，使用network V2 sdk中routers.Update方法，将Gatewayinfo field置空



## 对等连接相关服务

```
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package peerlink;

//对等连接相关服务
service PeerLinkService {
  //创建对等连接
  rpc CreatePeerLink(PeerLinkReq) returns(PeerLinkRes);
  //获取对等连接信息
  rpc GetPeerLink(PeerLinkReq) returns(PeerLinkRes);
  //删除对等连接
  rpc DeletePeerLink(PeerLinkReq) returns(DeletePeerLinkRes);
}

message PeerLinkRes {
  int32 code = 1;
  string msg = 2;
  message LinkConf {
    string intf_id = 1;
    string intf_ip = 2;
    message Route {
      string destination = 1;
      string nexthop = 2;
    }
    Route route_to_peer = 3;
    string created_time = 4;
  }
  LinkConf link_conf_on_peer_a = 3;
  LinkConf link_conf_on_peer_b = 4;
}

message PeerLinkReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string peer_a_subnetid = 4;
  string peer_a_routerid = 5;
  string peer_b_subnetid = 6;
  string peer_b_routerid = 7;
}

message DeletePeerLinkRes {
  int32 code = 1;
  string msg = 2;
  string peer_a_subnetid = 3;
  string peer_a_routerid = 4;
  string peer_b_subnetid = 5;
  string peer_b_routerid = 6;
  string deleted_time = 7;
}
```

### 创建对等连接

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 先决条件：管理员事先创建一个share类型的网络，该网络的subnetid已知，将此share网络的subnetid作为配置项放入配置文件中
3. 为peer_a_routerid这个路由器添加一个接口，接口连接的subnetid是share网络的subnetid，接口要指定ip，ip从该share网络的subnetpool中获取；同理，为peer_b_routerid这个路由器添加一个接口，接口连接的subnetid是share网络的subnetid，接口要指定ip，ip从该share网络的subnetpool中获取；
4. 为peer_a_routerid这个路由器添加一条路由，路由的destination目的网段，是peer_b_subnetid这个子网的cidr，路由的nexthop，是前面第3步，peer_b_routerid这个路由器新加接口的IP；同理，为peer_b_routerid这个路由器添加一条路由，路由的destination目的网段，是peer_a_subnetid这个子网的cidr，路由的nexthop，是前面第3步，peer_a_routerid这个路由器新加接口的IP;

### 获取对等连接信息

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 列出peer_a_routerid这个路由器所有port信息（使用ports2.List方法，listOpts中可以指明deviceid），遍历查找，如果一个port的FixedIPs[0].SubnetID==share网络的subnetid，取出该port的id和FixedIPs[0].IPAddress作为返回值link_conf_on_peer_a中的intf_id和intf_ip；同理，列出peer_b_routerid这个该路由器所有port信息（使用ports2.List方法，listOpts中可以指明deviceid），遍历查找，如果一个port的FixedIPs[0].SubnetID==share网络的subnetid，取出该port的id和FixedIPs[0].IPAddress作为返回值link_conf_on_peer_b中的intf_id和intf_ip；
3. 列出peer_a_routerid这个路由器所有路由表条目，遍历找出destination目的网段是peer_b_subnetid对应的cidr的那条路由，作为返回值link_conf_on_peer_a中的route_to_peer；同理，列出peer_a_routerid这个路由器所有路由表条目，遍历找出destination目的网段是peer_a_subnetid对应的cidr的那条路由，作为返回值link_conf_on_peer_b中的route_to_peer；

### 删除对等连接

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 针对peer_a_routerid这个路由器：调用routers.RemoveInterface，删除subnetid为share网络id的那个接口，然后调用原生API的remove_extraroutes，删除peer_b_subnetid所对应cidr的路由条目；同理，针对peer_b_routerid这个路由器：调用routers.RemoveInterface，删除subnetid为share网络id的那个接口，调用原生API的remove_extraroutes，删除peer_a_subnetid所对应cidr的路由条目；

## 浮动IP相关服务

```
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package floatip;

//浮动IP相关服务
service FloatIpService{
  //创建浮动ip并绑定到instance
  rpc BindFloatIpToInstance(BindFloatIpToInstanceReq) returns(BindFloatIpToInstanceRes);
  //解绑浮动ip并回收
  rpc RevokeFloatIpFromInstance(RevokeFloatIpFromInstanceReq) returns(RevokeFloatIpFromInstanceRes);
}

message BindFloatIpToInstanceReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string instance_id = 4;
  int32 vpc_router_id = 5;
}

message BindFloatIpToInstanceRes {
  int32 code = 1;
  string msg = 2;
  string float_ip = 3;
  string binded_time = 4;
}

message RevokeFloatIpFromInstanceReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string instance_id = 4;
  string float_ip = 5;
}

message RevokeFloatIpFromInstanceRes {
  int32 code = 1;
  string msg = 2;
  string revoked_time = 3;
}
```

### 创建浮动ip并绑定到instance

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据传入的vpc_router_id（该instance所属的vpc的路由器的id），检查该路由器是否有Gatewayinfo，如果没有说明该vpc还没有外部网关，没有与外部public网络相连，给出报错提示：该实例所在vpc没有外部网关，返回不再往下执行
3. 上述两步验证通过后，调用networking/v2/extensions/layer3/floatingips.Create方法产生一个新floating ip
4. 调用compute/v2/extensions/floatingips.AssociateInstance方法将上一步产生的floating ip关联到传入的instance id

### 解绑浮动ip并回收

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 调用compute/v2/extensions/floatingips..DisassociateInstance方法取消与instance id的关联
3. 根据传入的floating ip，调用networking/v2/extensions/layer3/floatingips.list方法，找到该floating ip对象的id
4. 根据上一步获得的floatingip的id，调用networking/v2/extensions/layer3/floatingips.Delete方法，删除该floating ip

## 防火墙相关服务

```
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";

// 指定文件生成出来的package
package firewall;

//防火墙相关服务
service FirewallService{
  //创建防火墙
  rpc CreateFirewall(CreateFirewallReq) returns(FirewallRes);
  //获取防火墙信息
  rpc GetFirewall(GetFirewallReq) returns(FirewallRes);
  //修改防火墙
  rpc UpdateFirewall(UpdateFirewallReq) returns(FirewallRes);
  //删除防火墙
  rpc DeleteFirewall(DeleteFirewallReq) returns(DeleteFirewallRes);
  //防火墙绑定/取消绑定路由器接口
  rpc OperateFirewall(OperateFirewallReq) returns(OperateFirewallRes);
}

message FirewallRule {
  string filewall_rule_id = 1;
  string filewall_rule_name = 2;
  string filewall_rule_desc = 3;
  string filewall_rule_action = 4;
  string filewall_rule_protocol = 5;
  string filewall_rule_src_ip = 6;
  string filewall_rule_src_port = 7;
  string filewall_rule_dst_ip = 8;
  string filewall_rule_dst_port = 9;
}

message FirewallPolicy {
  string filewall_policy_id = 1;
  string filewall_policy_name = 2;
  string filewall_policy_desc = 3;
  repeated FirewallRule filewall_policy_rules = 4;
}

//openstack原生支持一个firewall group对多个ports，本api限定firewall group与port为一一对应关系，故firewall_attached_port_id为单个string
message Firewall {
  string firewall_id = 1;
  string firewall_name = 2;
  string firewall_desc = 3;
  string firewall_attached_port_id = 4;
  string firewall_status = 5;
  FirewallPolicy firewall_ingress_policy = 6;
  FirewallPolicy firewall_egress_policy = 7;
  string created_time = 8;
  string updated_time = 9;
}

message FirewallRes {
  int32 code = 1;
  string msg = 2;
  Firewall firewall = 3;
}

message FirewallRuleSet {
  string filewall_rule_desc = 1;
  string filewall_rule_action = 2;
  string filewall_rule_protocol = 3;
  string filewall_rule_src_ip = 4;
  string filewall_rule_src_port = 5;
  string filewall_rule_dst_ip = 6;
  string filewall_rule_dst_port = 7;
}
message CreateFirewallReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string filewall_name = 4;
  string firewall_desc = 5;
  repeated FirewallRuleSet firewall_ingress_policy_rules = 6;
  repeated FirewallRuleSet firewall_egress_policy_rules = 7;
}

message GetFirewallReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string firewall_id = 4;
}

message UpdateFirewallReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string firewall_id = 4;
  string firewall_name = 5;
  string firewall_desc = 6;
  repeated FirewallRuleSet firewall_ingress_policy_rules = 7;
  repeated FirewallRuleSet firewall_egress_policy_rules = 8;
}

message DeleteFirewallReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string firewall_id = 4;
}

message DeleteFirewallRes {
  int32 code = 1;
  string msg = 2;
  string firewall_id = 3;
  string deleted_time = 4;
}

message OperateFirewallReq {
  string apikey = 1;
  string tenant_id = 2;
  string platform_userid = 3;
  string firewall_id = 4;
  string port_id = 5;
  string ops_type = 6; //attach or detach
}

message OperateFirewallRes {
  int32 code = 1;
  string msg = 2;
  string firewall_id = 3;
  string firewall_attached_port_id = 4;
  string ops_type = 5;
  string operated_time = 6;
}
```

### 创建防火墙

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 创建一个空firewall group，name，describe根据入参
3. 创建两个firewall policy，一个命名为<firewall_group_id>_ingress_policy，另一个命名为<firewall_group_id>_egress_policy
4. 根据入参中的firewall_ingress_policy_rules、firewall_egress_policy_rules，分别创建入向firewall rule集，和出向firewall rule集，并分别装入刚才创建的firewall policy
5. 更新firewall group，指定egress_policy_id，ingress_policy_id
6. firewall_attached_port_id返回空串

### 获取防火墙信息

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 根据入参的firewall_id，依次查询firewall group，firewall policy，firewall rules信息
3. created_time，updated_timed，返回空。

### 修改防火墙

要点：

1. 由于firewall rule集中可能既有permit规则，也有deny规则，并且有严格的顺序，如采用原子化的增删差集条目的方式操作，会造成条目顺序错误，故更新时采取整体替换方式（注意这一点和安全组规则不同）。
2. 鉴权 && 查表获取该租户的openstack连接参数；
3. 如果该firewall group已经绑定了路由器接口，先解绑
4. 创建两个firewall policy，一个命名为<firewall_group_id>_ingress_policy，另一个命名为<firewall_group_id>_egress_policy
5. 根据入参中的firewall_ingress_policy_rules、firewall_egress_policy_rules，分别创建入向firewall rule集，和出向firewall rule集，并分别装入刚才创建的firewall policy
6. 更新firewall group，egress_policy_id，ingress_policy_id指向刚才新建的两个policy
7. 如果原先firewall group已经绑定了路由器接口，重新绑定
8. 根据两个旧的policy_id，先删除其下的firewall rules，再删除policy
9. 回滚原则：第4，5步出现错误，清理垃圾数据即可；第6步出错，firewall group重新指向旧policyid

### 删除防火墙

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 如果该firewall group已经绑定了路由器接口，不允许删除，提示解绑路由器端口后方可删除
3. 根据传入的firewall_id，查出ingress_policy_id，egress_policy_id以及其下的各rule id，依次删除，删除顺序，firewall rule-->firewall policy-->firewall group

### 防火墙绑定/取消绑定路由器接口

要点：

1. 鉴权 && 查表获取该租户的openstack连接参数；
2. 如果该firewall group已经绑定了路由器接口，不允许再次绑定，提示解绑路由器端口后方可绑定
3. 根据传入的firewall_id，port_id，ops_type，更新firewall group



## 统一注意事项

1. 调用openstack API，ceph API，有的响应比较慢，应有超时处理
2. 收到底层API抛出error，都触发向上（云管平台）抛出error，error信息可以隐藏细节信息
3. 尽量使用gophercloud sdk，减少原生openstack API调用
4. 返回参数值一定是从openstack查出来的真实信息，例如修改vpc描述信息，返回参数中也有描述信息字段，这个字段应该是调用openstack api查出来的，方能验证刚才的修改是否已生效。
5. 调用一些底层接口时，那些接口本身就是异步处理的，怎样获知已经处理完成？
