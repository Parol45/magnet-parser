package requests

import (
	"fmt"
	"log/slog"
	"magnet-parser/bencode_obj"
	"magnet-parser/dht"
	"magnet-parser/globals"
)


const getPeersFmt = `{"a":{"id":"%s","info_hash":"%s"},"q":"get_peers","t":"%s","v":"%s","y":"q"}`
const findNodeFmt = `{"a":{"id":"%s","target":"%s"},"q":"find_node","t":"%s","v":"%s","y":"q"}`
const announcePeerFmt = `{"a":{"id":"%s","info_hash":"%s","port":%d,"token":"%s"},"q":"announce_peers","t":"%s","v":"%s","y":"q"}`


func Ping() []byte {
	req := globals.TakePointer(globals.NewPingRequest())
	req, err := dht.Compress(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while compressing bencode obj: %v\n", err))
		return nil
	}
	result, err := bencode_obj.Encode(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting json to bencode: %v\n", err))
		return nil
	}
	return result
}

func GetPeers(hash string) []byte {
	req := globals.TakePointer(globals.NewGetPeersRequest(hash))
	req, err := dht.Compress(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while compressing bencode obj: %v\n", err))
		return nil
	}
	result, err := bencode_obj.Encode(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting json to bencode: %v\n", err))
		return nil
	}
	return result
}

func FindNode(foreignId string) []byte {
	req := globals.TakePointer(globals.NewFindNodeRequest(foreignId))
	req, err := dht.Compress(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while compressing bencode obj: %v\n", err))
		return nil
	}
	result, err := bencode_obj.Encode(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting json to bencode: %v\n", err))
		return nil
	}
	return result
}

func AnnouncePeer(hash string, port int, token string) []byte {
	req := globals.TakePointer(globals.NewAnnouncePeerRequest(hash, port, token))
	req, err := dht.Compress(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while compressing bencode obj: %v\n", err))
		return nil
	}
	result, err := bencode_obj.Encode(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while converting json to bencode: %v\n", err))
		return nil
	}
	return result
}
