// +build linux
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
package vpncore

import (
	"net"
	"fmt"
	"errors"
)

func setUpHWAddr(ifce *Interface) (err error) {
	if ifce.IsTUN() {
		return fmt.Errorf("dont set hw addr to Tun device!")
	}

	cmd := fmt.Sprintf("ip link set dev %s address %s broadcast %s",
		ifce.Name(),
		fmt.Sprintf("%s%02x", DEFAULT_HWADDR_PREFIX, ifce.IP()[3]),
		DEFUALT_HWADDR_BRD)

	_, err = runCommand(cmd)
	return
}

func (ifce *Interface) SetupNetwork(ip net.IP, subnet net.IPNet, mtu int) (err error) {

	var cmd string

	err = ifce.changeMTU(mtu)
	if err != nil {
		return
	}

	if ifce.IsTUN() {
		peer_ip := generatePeerIP(ip)
		cmd = fmt.Sprintf("ip addr add dev %s %s peer %s", ifce.Name(), ip.String(), peer_ip.String())
	} else {
		cmd = fmt.Sprintf("ip addr add dev %s add %s", ifce.Name(), ip.String())
	}
	_, err = runCommand(cmd)
	if err != nil {
		return err
	} else {
		ifce.SetIP(ip, subnet)

	}

	if ifce.IsTAP() {
		err = setUpHWAddr(ifce)
	} else {
		_, err = runCommand(fmt.Sprintf("ip link set %s up", ifce.Name()))
	}

	if err != nil {
		return err
	}

	err = ifce.setupRoutes()

	return
}

func (ifce *Interface) SetupNATForServer() (err error) {

	subnet := ifce.Net()

	cmd1 := fmt.Sprintf("iptables -t nat -A POSTROUTING -o %s -s %s -j MASQUERADE", ifce.routes_m.default_nic, subnet.String())
	cmd2 := fmt.Sprintf("iptables -A FORWARD -d %s -i %s -o %s -j ACCEPT", subnet.String(), ifce.routes_m.default_nic, ifce.Name())
	cmd3 := fmt.Sprintf("iptables -A FORWARD -s %s -i %s -o %s -j ACCEPT", subnet.String(), ifce.Name(), ifce.routes_m.default_nic)
	cmd4 := "sysctl net.ipv4.ip_forward=1"

	_, err = runCommand(cmd1)
	if err != nil {
		return
	}
	_, err = runCommand(cmd2)
	if err != nil {
		return
	}
	_, err = runCommand(cmd3)
	if err != nil {
		return
	}

	_, err = runCommand(cmd4)
	if err != nil {
		return
	}

	return
}

func (ifce *Interface) changeMTU(mtu int) (err error) {

	cmd := fmt.Sprintf("ip link set dev %s up mtu %d qlen 100", ifce.Name(), mtu)
	_, err = runCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (ifce *Interface) setupRoutes() (err error) {

	if ifce.IP() == nil {
		return errors.New("Setup interface IP first!")
	}

	err = ifce.routes_m.AddRouteToNet(ifce.Name(), ifce.subnet, ifce.IP())
	return
}

