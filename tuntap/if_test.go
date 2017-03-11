/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */
package tuntap

import (
	"testing"
	"net"
	"time"
	"fmt"
	"encoding/hex"
	"sync"
	"github.com/FTwOoO/vpncore/tuntap/tcpip"
	"github.com/FTwOoO/vpncore/tuntap/cmd"
	"github.com/FTwOoO/vpncore/tuntap/cmd/routes"
)

const BUFFERSIZE = 1522

func startRead(wg *sync.WaitGroup, ifce *Interface) bool {
	wg.Add(1)
	defer func () {
		fmt.Printf("startRead() end!")
		wg.Done()
	}()

	Reading:
	for {
		buffer := make([]byte, BUFFERSIZE)
		n, err := ifce.Read(buffer)
		if err == nil {
			fmt.Printf("Received a packet(%d bytes from %s)\n", n, ifce.Name())
			buffer = buffer[:n]

			var ipPacket tcpip.IPv4Packet

			if ifce.IsTAP() {
				ethertype := tcpip.MACPacket(buffer).MACEthertype()
				if ethertype != tcpip.IPv4 {
					fmt.Printf("Packet is not ipv4\n")
					continue Reading
				}
				if !tcpip.IsBroadcast(tcpip.MACPacket(buffer).MACDestination()) {
					fmt.Printf("Packet is Broadcast\n")
					continue Reading
				}

				ipPacket = tcpip.IPv4Packet(tcpip.MACPacket(buffer).MACPayload())
			} else {
				ipPacket = tcpip.IPv4Packet(buffer)
			}

			if !tcpip.IsIPv4(ipPacket) {
				fmt.Printf("Packet is not ipv4\n")
				continue Reading
			}

			if !ipPacket.SourceIP().Equal(ifce.IP()) {
				fmt.Printf("Packet source[%s] dont match [%s]\n", ipPacket.SourceIP().String(), ifce.IP().String())
				continue Reading
			}
			if ipPacket.Protocol() != tcpip.ICMP {
				fmt.Printf("Packet is not ICMP\n")
				continue Reading
			}
			fmt.Printf("Received ICMP frame: %#v\n", hex.EncodeToString(ipPacket))
			return true


		} else {
			fmt.Println(err)
			return false
		}
	}
}

func startPing(dst net.IP) {
	c := time.NewTicker(1 * time.Second)
	select {
	case <-c.C:
		c := fmt.Sprintf("ping -c 5 %s", dst.String())
		cmd.RunCommand(c)
	}
}

func ip4BroadcastAddr(subnet net.IPNet) (brdIp net.IP) {

	brdIp = net.IP{0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		brdIp[i] = subnet.IP[i] | (0xFF ^ subnet.Mask[i])
	}
	return

}

func testInterface(ifce *Interface, ip net.IP, subnet net.IPNet) {

	defer ifce.Close()


	err := ifce.SetupNetwork(ip, nil, subnet, 1400)
	if err != nil {
		panic(err)
	}

	router, err := routes.NewRoutesManager()
	if err != nil {
		panic(err)
	}
	defer router.Destroy()

	err = router.SetNewGateway(ifce.Name(), ifce.IP())
	if err != nil {
		panic(err)
	}

	go func() {
		if ifce.IsTUN() {
			startPing(ifce.PeerIP())
		} else {
			startPing(ip4BroadcastAddr(ifce.Net()))
		}
	}()

	wg := &sync.WaitGroup{}
	succed := startRead(wg, ifce)
	if succed == false {
		panic("Fail to read ICMP packets")
	}

	wg.Wait()
	fmt.Printf("Close the iterface %s\n", ifce.Name())

}

func TestAll(t *testing.T) {
	subnet := net.IPNet{IP:[]byte{192, 168, 77, 0}, Mask:net.IPv4Mask(255, 255, 255, 0)}
	ip := net.IP{192, 168, 77, 1}

	ifce, err := NewTUN("tun3")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("create %s\n", ifce.Name())
	testInterface(ifce, ip, subnet)
	ifce.Close()

	ifce2, err := NewTAP("tap1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("create %s\n", ifce2.Name())
	testInterface(ifce2, ip, subnet)
	ifce2.Close()

}
