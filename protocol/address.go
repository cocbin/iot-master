package protocol

//Address 地址
//type Address struct {
//	Slave  uint8  `json:"slave,omitempty"`
//	Code   uint8  `json:"code"`
//	Block  uint16 `json:"block,omitempty"`
//	Offset uint16 `json:"offset"`
//}

type Address interface {
	String() string
	//Equal(addr Addr) bool
}

