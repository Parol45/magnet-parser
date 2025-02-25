package globals

import (
	"fmt"
)

const id = "5ba5c0ac0727f52516dba42f4f6bc8ad097deba6"
const version = "aaaaaaaa"

var transactionId = 0

func transIdToStr() string {
	transactionId++
	transactionId = transactionId % 10000000
	return fmt.Sprintf("%08d", transactionId)
}

// arguments
type Atype struct {
	Id          string  `json:"id,omitempty"`
	Target      *string `json:"target,omitempty"`
	InfoHash    *string `json:"infoHash,omitempty"`
	ImpliedPort *int    `json:"impliedPort,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Token       *string `json:"token,omitempty"`
}

// response
type Rtype struct {
	Id     string   `json:"id,omitempty"`
	Nodes  *string  `json:"nodes,omitempty"`
	Token  *string  `json:"token,omitempty"`
	Values []string `json:"values,omitempty"`
}

type PackageType struct {
	A *Atype `json:"a,omitempty"`
	R *Rtype `json:"r,omitempty"`
	Q string `json:"q,omitempty"`
	T string `json:"t,omitempty"`
	V string `json:"v,omitempty"`
	Y string `json:"y,omitempty"`
}

func NewPingRequest() PackageType {
	return PackageType{
		A: TakePointer(Atype{Id: id}),
		Q: "ping",
		T: transIdToStr(),
		V: version,
		Y: "q",
	}
}

func NewGetPeersRequest(infoHash string) PackageType {
	return PackageType{
		A: TakePointer(Atype{
			Id:       id,
			InfoHash: TakePointer(infoHash),
		}),
		Q: "get_peers",
		T: transIdToStr(),
		V: version,
		Y: "q",
	}
}

func NewFindNodeRequest(target string) PackageType {
	return PackageType{
		A: TakePointer(Atype{
			Id:     id,
			Target: TakePointer(target),
		}),
		Q: "find_node",
		T: transIdToStr(),
		V: version,
		Y: "q",
	}
}

func NewAnnouncePeerRequest(infoHash string, port int, token string) PackageType {
	return PackageType{
		A: TakePointer(Atype{
			Id:       id,
			InfoHash: TakePointer(infoHash),
			Port:     TakePointer(port),
			Token:    TakePointer(token),
		}),
		Q: "announce_peers",
		T: transIdToStr(),
		V: version,
		Y: "q",
	}
}
