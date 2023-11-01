package packet

import (
	"errors"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/nttcom/fluvia/pkg/bpf"
)

const MAX_SEGMENTLIST_ENTRIES = 10

type ProbeData struct {
	H_source     string
	H_dest       string
	V6Srcaddr    string
	V6Dstaddr    string
	NextHdr      uint8
	HdrExtLen    uint8
	RoutingType  uint8
	SegmentsLeft uint8
	LastEntry    uint8
	Flags        uint8
	Tag          uint16
	Segments     [MAX_SEGMENTLIST_ENTRIES]string
}

func Parse(data []byte) (*ProbeData, error) {
	var pd ProbeData
	packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)

	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	eth, ok := ethLayer.(*layers.Ethernet)
	if !ok {
		return nil, errors.New("Could not parse a packet with Ethernet")
	}

	pd.H_dest = eth.DstMAC.String()
	pd.H_source = eth.SrcMAC.String()

	ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
	ipv6, ok := ipv6Layer.(*layers.IPv6)
	if !ok {
		return nil, errors.New("Could not parse a packet with IPv6")
	}

	pd.V6Srcaddr = ipv6.SrcIP.String()
	pd.V6Dstaddr = ipv6.DstIP.String()

	srv6Layer := packet.Layer(bpf.Srv6LayerType)
	srv6, ok := srv6Layer.(*bpf.Srv6Layer)
	if !ok {
		return nil, errors.New("Could not parse a packet with SRv6")
	}

	pd.NextHdr = srv6.NextHeader
	pd.HdrExtLen = srv6.HdrExtLen
	pd.RoutingType = srv6.RoutingType
	pd.SegmentsLeft = srv6.SegmentsLeft
	pd.LastEntry = srv6.LastEntry
	pd.Flags = srv6.Flags
	pd.Tag = srv6.Tag

	for idx := 0; idx < MAX_SEGMENTLIST_ENTRIES; idx++ {
		if idx >= len(srv6.Segments) {
			break
		}
		pd.Segments[idx] = srv6.Segments[idx].String()
	}

	return &pd, nil
}
