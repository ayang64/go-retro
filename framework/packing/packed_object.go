package packing

import (
	"crypto/sha256"
)

// Packed event represents a packed event in memory. The payload is a zlib
// deflated string of whatever the pack encoding scheme encodes to (json,
// msgpack, etc).
type PackedObject struct {
	payload []byte
	hash    Hash
}

func NewPackedObject(payloadStr string) PackedObject {
	return PackedObject{
		payload: []byte(payloadStr),
		hash: Hash{
			AlgoName: HashAlgoNameSHA256,
			Bytes:    sha256.New().Sum([]byte(payloadStr)),
		},
	}
}

func (po *PackedObject) Contents() []byte {
	return po.payload
}
func (po *PackedObject) Hash() Hash {
	return po.hash
}

type PackedEvent struct {
	PackedObject
}

func (pe *PackedEvent) TypeName() ObjectTypeName {
	return ObjectTypeEvent
}

type PackedAffix struct {
	PackedObject
}

func (pe *PackedAffix) TypeName() ObjectTypeName {
	return ObjectTypeAffix
}

type PackedCheckpoint struct {
	PackedObject
}

func (pc *PackedCheckpoint) TypeName() ObjectTypeName {
	return ObjectTypeCheckpoint
}
