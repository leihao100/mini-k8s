package ipvsManager

const defaultProtocol = 6      //tcp
const defaultAddressFamily = 2 //ipv4
const defaultScheduler = "rr"  //round-robin
/*
协议类型	值
描述
Ggp	3
网关到网关协议。
Icmp	1
Internet 控制消息协议。
IcmpV6	58
IPv6 的 Internet 控制消息协议。
Idp	22
Internet 数据报协议。
Igmp	2
Internet 组管理协议。
IP	0
Internet 协议。
IPSecAuthenticationHeader	51
IPv6 身份验证标头。 有关详细信息，请参阅https://www.ietf.org 上的 RFC 2292，第 2.2.1 节。
IPSecEncapsulatingSecurityPayload	50
IPv6 封装安全负载标头。
IPv4	4
Internet 协议版本 4。
IPv6	41
Internet 协议版本 6 (IPv6)。
IPv6DestinationOptions	60
IPv6 目标选项标头。
IPv6FragmentHeader	44
IPv6 片段标头。
IPv6HopByHopOptions	0
IPv6 逐跳选项标头。
IPv6NoNextHeader	59
IPv6 无下一个标头。
IPv6RoutingHeader	43
IPv6 路由标头。
Ipx	1000
Internet 数据包交换协议。
ND	77
网络磁盘协议（非正式）。
Pup	12
PARC 通用数据包协议。
Raw	255
原始 IP 数据包协议。
Spx	1256
顺序包交换协议。
SpxII	1257
顺序包交换版本 2 协议。
Tcp	6
传输控制协议。
Udp	17
用户数据报协议。
Unknown	-1
未知的协议。

Unspecified	0
未指定的协议。
*/

/*

地址类型
值
描述
AppleTalk	16
AppleTalk 地址。
Atm	22
本机 ATM 服务地址。
Banyan	21
Banyan 地址。
Ccitt	10
CCITT 协议（如 X.25）的地址。
Chaos	5
MIT CHAOS 协议的地址。
Cluster	24
Microsoft 群集产品的地址。
DataKit	9
Datakit 协议的地址。
DataLink	13
直接数据链接接口地址。
DecNet	12
DECnet 地址。
Ecma	8
欧洲计算机制造商协会 (ECMA) 地址。
FireFox	19
FireFox 地址。
HyperChannel	15
NSC Hyperchannel 地址。
Ieee12844	25
IEEE 1284.4 工作组地址。
ImpLink	3
ARPANET IMP 地址。
InterNetwork	2
IP 版本 4 的地址。
InterNetworkV6	23
IP 版本 6 的地址。
Ipx	6
IPX 或 SPX 地址。
Irda	26
IrDA 地址。
Iso	7
ISO 协议的地址。
Lat	14
LAT 地址。
Max	29
MAX 地址。
NetBios	17
NetBios 地址。
NetworkDesigners	28
支持网络设计器 OSI 网关的协议的地址。
NS	6
Xerox NS 协议的地址。
Osi	7
OSI 协议的地址。
Pup	4
PUP 协议的地址。
Sna	11
IBM SNA 地址。
Unix	1
Unix 本地到主机地址。
Unknown	-1
未知的地址族。
Unspecified	0
未指定的地址族。
VoiceView	18
VoiceView 地址。
*/
