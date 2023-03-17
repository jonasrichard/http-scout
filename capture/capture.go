package capture

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func Run() error {
	if handle, err := pcap.OpenLive("lo0", 1600, true, pcap.BlockForever); err != nil {
		return err
	} else if err := handle.SetBPFFilter("tcp and port 9000"); err != nil {
		return err
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for packet := range packetSource.Packets() {
			handlePacket(packet)
		}
	}

	return nil
}

func handlePacket(packet gopacket.Packet) {
	fmt.Printf("%v\n", packet)
}
