package protocol

import (
	"encoding/binary"
	"github.com/google/gopacket"
)

type DataLayer struct {
	SequenceNumber uint32
	AckNumber uint32
	IdConnection uint16
	Flags uint16
	Payload []byte
}

var DataLayerType = gopacket.RegisterLayerType(
	2001,
	gopacket.LayerTypeMetadata {
		Name:    "DataLayerType",
		Decoder: gopacket.DecodeFunc(decodeDataLayer),
	},
)

func (d DataLayer) LayerType() gopacket.LayerType {
	return DataLayerType
}

func (d DataLayer) LayerContents() []byte {
	var sequenceNumberBytes = make([]byte, 4)
	var ackNumberBytes = make([]byte, 4)
	var idConnectionBytes = make([]byte, 2)
	var flagsBytes = make([]byte, 2)
	var payload = d.Payload

	binary.BigEndian.PutUint32(sequenceNumberBytes, d.SequenceNumber)
	binary.BigEndian.PutUint32(ackNumberBytes, d.AckNumber)
	binary.BigEndian.PutUint16(idConnectionBytes, d.IdConnection)
	binary.BigEndian.PutUint16(flagsBytes, d.Flags)

	var content []byte
	content = append(sequenceNumberBytes)
	content = append(ackNumberBytes)
	content = append(idConnectionBytes)
	content = append(flagsBytes)
	content = append(payload)

	return content
}

func (d DataLayer) LayerPayload() []byte {
	return d.Payload
}

func decodeDataLayer(data []byte, p gopacket.PacketBuilder) error {
	var sequenceNumber = binary.BigEndian.Uint32(data[0:4])
	var ackNumber = binary.BigEndian.Uint32(data[4:8])
	var idConnection = binary.BigEndian.Uint16(data[8:10])
	var flags = binary.BigEndian.Uint16(data[10:12])
	var payload []byte = nil

	if len(data) >= 12 {
		payload = data[12:]
	}

	p.AddLayer(&DataLayer{
		sequenceNumber,
		ackNumber,
		idConnection,
		flags,
		payload,
	})

	return p.NextDecoder(gopacket.LayerTypePayload)
}
