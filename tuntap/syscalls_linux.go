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
package tuntap

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
	"fmt"
)

const (
	IFNAMSIZ = 16

)

type ifReq struct {
	Name  [IFNAMSIZ]byte
	Flags uint16
	pad   [0x20 - IFNAMSIZ - 2]byte
}

func newTAP(ifName string) (ifce *Interface, err error) {
	if !checkTapName(ifName) {
		return nil, fmt.Errorf("Error name:%s", ifName)
	}

	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	name, err := createInterface(file.Fd(), ifName, syscall.IFF_TAP | syscall.IFF_NO_PI)
	if err != nil {
		return nil, err
	}

	ifce = &Interface{isTAP: true,
		ReadWriteCloser: file,
		name: name,
	}
	return
}

func newTUN(ifName string) (ifce *Interface, err error) {

	if !checkTunName(ifName) {
		return nil, fmt.Errorf("Error name:%s", ifName)
	}

	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	name, err := createInterface(file.Fd(), ifName, syscall.IFF_TUN | syscall.IFF_NO_PI)
	if err != nil {
		fmt.Println("err %v", err)

		return nil, err
	}

	ifce = &Interface{isTAP: false,
		ReadWriteCloser: file,
		name: name,
	}
	return
}

func createInterface(fd uintptr, ifName string, flags uint16) (createdIFName string, err error) {
	var req ifReq
	req.Flags = flags
	copy(req.Name[:], ifName)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		err = fmt.Errorf("ioctl: %v", errno)
		return
	}
	createdIFName = strings.Trim(string(req.Name[:]), "\x00")

	return
}
