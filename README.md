# irtc

[![GitHub](https://img.shields.io/github/license/coopersong/irtc)](https://github.com/coopersong/irtc/blob/master/LICENSE)

This package is used for converting IP range to CIDR(s). 

1 IP range can have several CIDRs. For example, the IP range [192.168.1.0, 192.168.1.255] can be converted to [192.168.1.0/24]. 
[192.168.1.0, 192.168.1.254] can be converted to [192.168.1.0/25 192.168.1.128/26 192.168.1.192/27 192.168.1.224/28 192.168.1.240/29 192.168.1.248/30 192.168.1.252/31 192.168.1.254/32].

## CIDR

Classless Inter-Domain Routing(CIDR) is a method of assigning IP addresses that improves the efficiency of address distribution and replaces the previous system based on Class A, Class B and Class C networks.

For example, CIDR 192.168.1.0/24 can represent the IP net whose IP range is from 192.168.1.0 to 192.168.1.255.

You can refer to [this page](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) for more information.

## Contributing

- Please create an issue in [issue list](https://github.com/coopersong/irtc/issues).
- Following the golang coding standards.

## License

The project is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
