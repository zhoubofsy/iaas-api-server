# 本地测试
#./main -method 'plCreate' -param '{"apikey": "12312", "tenant_id": "12312", "platform_userid": "12321", "peer_a_subnetid": "2c12437b-3f55-4cca-9149-11dbe3dc15f0", "peer_a_routerid": "66dbf5aa-43e3-4aa0-b8f4-8f755ed9d6da", "peer_b_routerid": "2fdcede4-dc07-43fe-b43d-18baa11967bc", "peer_b_subnetid": "db92f013-d0bc-46da-b522-c739f63a1059"}'
#./main -method 'plDelete' -param '{"apikey": "12312", "tenant_id": "12312", "platform_userid": "12321", "peer_a_subnetid": "2c12437b-3f55-4cca-9149-11dbe3dc15f0", "peer_a_routerid": "66dbf5aa-43e3-4aa0-b8f4-8f755ed9d6da", "peer_b_routerid": "2fdcede4-dc07-43fe-b43d-18baa11967bc", "peer_b_subnetid": "db92f013-d0bc-46da-b522-c739f63a1059"}'
#./main -method 'plGet' -param '{"apikey": "dMXFqy0H1w", "tenant_id": "t-0000001003", "platform_userid": "6bf806f67cdd4a2e98b05b2371c0c0bd", "peer_a_subnetid": "2c12437b-3f55-4cca-9149-11dbe3dc15f0", "peer_a_routerid": "66dbf5aa-43e3-4aa0-b8f4-8f755ed9d6da", "peer_b_routerid": "2fdcede4-dc07-43fe-b43d-18baa11967bc", "peer_b_subnetid": "db92f013-d0bc-46da-b522-c739f63a1059"}'

# security_group
#./main -method 'sgCreate' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "security_group_name": "mongia sg", "security_group_desc": "mongia test", "security_group_rule_sets": [{"rule_desc": "mongia rule 1", "direction": "ingress", "protocol": "tcp", "port_range_min": 8081, "port_range_max": 8086, "remote_ip_prefix": "10.4.0.0/24"}, { "rule_desc": "mongia rule 1", "direction": "egress", "protocol": "tcp", "port_range_min": 8081, "port_range_max": 8086, "remote_ip_prefix": "10.4.0.0/24"}]}'

# nat gateway

# firewall
#./main -method 'fCreate' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "firewall_name": "firewallA", "firewall_desc": "firewall desc", "firewall_ingress_policy_rules": [{"firewall_rule_desc": "ingress rule", "firewall_rule_action": "allow", "firewall_rule_protocol": "tcp", "firewall_rule_src_ip": "127.0.0.1", "firewall_rule_src_port": "80", "firewall_rule_dst_ip": "127.0.0.1", "firewall_rule_dst_port": "90"}], "firewall_egress_policy_rules": [{"firewall_rule_desc": "egress rule", "firewall_rule_action": "deny", "firewall_rule_protocol": "udp", "firewall_rule_src_ip": "127.0.0.1", "firewall_rule_src_port": "580", "firewall_rule_dst_ip": "127.0.0.1", "firewall_rule_dst_port": "900"}]}'
#./main -method 'fGet' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "firewall_id": "a03d3a21-dbf1-44bb-bddf-ce01f3699c50"}'
#./main -method 'fDelete' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "firewall_id": "a03d3a21-dbf1-44bb-bddf-ce01f3699c50"}'
./main -method 'fOperate' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "firewall_id": "e9c6120b-e97d-4176-a68a-f95eec4c39fd", "port_id": "24d68099-111d-4e10-bbb4-efaa8cc2d650", "ops_type": "detach"}'

# 金山云测试

# peerlink
#./main -address '120.48.27.208:58080' -method 'plCreate' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "peer_a_subnetid": "47580bdf-d2f6-4686-a7c6-a20f9f362182", "peer_a_routerid": "94d49e6c-51df-4617-a6c9-38a667b33032", "peer_b_routerid": "53b8872c-e5b2-4b54-91ec-0df4e0856048", "peer_b_subnetid": "0b7ea3ac-c076-466e-9187-4850f9622b88"}'
#./main -address '120.48.27.208:58080' -method 'plGet' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "peer_a_subnetid": "47580bdf-d2f6-4686-a7c6-a20f9f362182", "peer_a_routerid": "94d49e6c-51df-4617-a6c9-38a667b33032", "peer_b_routerid": "53b8872c-e5b2-4b54-91ec-0df4e0856048", "peer_b_subnetid": "0b7ea3ac-c076-466e-9187-4850f9622b88"}'
#./main -address '120.48.27.208:58080' -method 'plDelete' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "peer_a_subnetid": "47580bdf-d2f6-4686-a7c6-a20f9f362182", "peer_a_routerid": "94d49e6c-51df-4617-a6c9-38a667b33032", "peer_b_routerid": "53b8872c-e5b2-4b54-91ec-0df4e0856048", "peer_b_subnetid": "0b7ea3ac-c076-466e-9187-4850f9622b88"}'

# security_group
#./main -address '120.48.27.208:58080' -method 'sgCreate' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "security_group_name": "mongia sg", "security_group_desc": "mongia test", "security_group_rule_sets": [{"rule_desc": "mongia rule 1", "direction": "ingress", "protocol": "tcp", "port_range_min": 8081, "port_range_max": 8086, "remote_ip_prefix": "10.4.0.0/24"}, { "rule_desc": "mongia rule 1", "direction": "egress", "protocol": "tcp", "port_range_min": 8081, "port_range_max": 8086, "remote_ip_prefix": "10.4.0.0/24"}]}'
#./main -address '120.48.27.208:58080' -method 'sgDelete' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "security_group_id": "664588cb-45d9-47c9-a312-5c51ba492278"}'
#./main -address '120.48.27.208:58080' -method 'sgGet' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "security_group_id": "f3f22fde-84e8-4d84-82e5-5ab6791d5872"}'

# nat gateway
#./main -address '120.48.27.208:58080' -method 'ngDelete' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "router_id": "53b8872c-e5b2-4b54-91ec-0df4e0856048", "gateway_id": "983b7bcd-9a30-42ed-b9f9-3c1d96ffa466"}'
#./main -address '120.48.27.208:58080' -method 'ngDelete' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "router_id": "94d49e6c-51df-4617-a6c9-38a667b33032", "gateway_id": "8997e5a9-081b-4e9f-a8b3-f7847c1f1cb9"}'
#./main -address '120.48.27.208:58080' -method 'ngCreate' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "router_id": "94d49e6c-51df-4617-a6c9-38a667b33032", "external_network_id": "8321d0b2-4990-4047-a4c0-97a98b9cf63c"}'
#./main -address '120.48.27.208:58080' -method 'ngCreate' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "router_id": "53b8872c-e5b2-4b54-91ec-0df4e0856048", "external_network_id": "8321d0b2-4990-4047-a4c0-97a98b9cf63c"}'
#./main -address '120.48.27.208:58080' -method 'ngGet' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "router_id": "53b8872c-e5b2-4b54-91ec-0df4e0856048", "gateway_id": "983b7bcd-9a30-42ed-b9f9-3c1d96ffa466"}'
#./main -address '120.48.27.208:58080' -method 'ngGet' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "router_id": "53b8872c-e5b2-4b54-91ec-0df4e0856048", "gateway_id": "983b7bcd-9a30-42ed-b9f9-3c1d96ffa466"}'

# floadip
#./main -address '120.48.27.208:58080' -method 'sgCreate' -param '{"apikey": "RoYuYLWBdI", "tenant_id": "t-0000000019", "platform_userid": "xiaott14", "security_group_name": "mongia sg", "security_group_desc": "mongia test", "security_group_rule_sets": []}'


