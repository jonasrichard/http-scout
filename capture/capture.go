package capture

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Endpoint struct {
	IP   string
	Port uint16
}

type EndpointPair struct {
	Src Endpoint
	Dst Endpoint
}

type StreamFragment struct {
	Endpoints EndpointPair
	Payload   []byte
	FIN       bool
}

type Stream struct {
	Endpoints  EndpointPair
	SrcPayload []byte
	DstPayload []byte
	Timestamp  time.Time
	SrcFIN     bool
	DstFIN     bool
}

type Capture struct {
	Streams map[EndpointPair]*Stream
}

func (ep EndpointPair) Reverse() EndpointPair {
	return EndpointPair{
		Src: ep.Dst,
		Dst: ep.Src,
	}
}

func NewCapture() *Capture {
	return &Capture{
		Streams: make(map[EndpointPair]*Stream),
	}
}

func (c *Capture) Run() error {
	if handle, err := pcap.OpenLive("lo0", 1600, true, pcap.BlockForever); err != nil {
		return err
	} else if err := handle.SetBPFFilter("tcp and port 9000"); err != nil {
		return err
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for packet := range packetSource.Packets() {
			if sf := handlePacket(packet); sf != nil {
				c.AddStreamFragment(sf)
			}
		}
	}

	return nil
}

func (c *Capture) AddStreamFragment(sf *StreamFragment) {
	reversed := false

	stream, ok := c.Streams[sf.Endpoints]
	if !ok {
		reverseEndpoints := sf.Endpoints.Reverse()

		stream, ok = c.Streams[reverseEndpoints]
		if !ok {
			stream = &Stream{
				Endpoints:  sf.Endpoints,
				SrcPayload: make([]byte, 0),
				DstPayload: make([]byte, 0),
				Timestamp:  time.Now(),
				SrcFIN:     false,
				DstFIN:     false,
			}

			c.Streams[sf.Endpoints] = stream
		} else {
			reversed = true
		}
	}

	if reversed {
		stream.DstPayload = append(stream.DstPayload, sf.Payload...)
	} else {
		stream.SrcPayload = append(stream.SrcPayload, sf.Payload...)
	}

	alreadyFinished := stream.SrcFIN && stream.DstFIN

	if sf.FIN {
		if reversed {
			stream.DstFIN = true
		} else {
			stream.SrcFIN = true
		}
	}

	if !alreadyFinished && stream.SrcFIN && stream.DstFIN {
		fmt.Println(string(stream.SrcPayload))
		fmt.Println(string(stream.DstPayload))
	}
}

func handlePacket(packet gopacket.Packet) *StreamFragment {
	fmt.Printf("Packet %v\n", packet)

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		result := &StreamFragment{}

		tcp := tcpLayer.(*layers.TCP)

		src := Endpoint{
			Port: uint16(tcp.SrcPort),
		}
		dst := Endpoint{
			Port: uint16(tcp.DstPort),
		}

		result.Payload = tcp.Payload

		if tcp.FIN {
			result.FIN = true
		}

		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ip := ipLayer.(*layers.IPv4)

			src.IP = ip.SrcIP.To4().String()
			dst.IP = ip.DstIP.To4().String()

			result.Endpoints = EndpointPair{
				Src: src,
				Dst: dst,
			}
		}

		return result
	}

	return nil
}
