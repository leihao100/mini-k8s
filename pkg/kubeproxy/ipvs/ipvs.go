package ipvs

type IPSet struct {
	Name string
	// SetType specifies the ipset type.
	SetType string
	// HashFamily specifies the protocol family of the IP addresses to be stored in the set.
	// The default is inet, i.e IPv4.  If users want to use IPv6, they should specify inet6.
	HashFamily string
	// HashSize specifies the hash table size of ipset.
	HashSize int
	// MaxElem specifies the max element number of ipset.
	MaxElem int
	// PortRange specifies the port range of bitmap:port type ipset.
	PortRange string
	// comment message for ipset
	Comment string
	//Map maps the service and its ips
	Map map[string]string
}
