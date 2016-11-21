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

package geoip

import (
	"testing"
	"time"
	"fmt"
	"math/rand"
	mtesting "github.com/FTwOoO/vpncore/testing"
)

func BenchmarkMaxMindDBCountry(b *testing.B) {
	start := time.Now()
	fmt.Printf("\nTest %d times:\n", b.N)

	db, err := Open("GeoLite2-Country.mmdb")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	r := rand.New(rand.NewSource(0))

	for i := 0; i < b.N; i++ {
		ip := mtesting.RandomIPv4Address(r)
		isCN, err := db.IsChineseIP(ip)
		country, _ := db.Country(ip)
		if err != nil {
			b.Fatal(err)
		}
		fmt.Printf("IP %v is CN IP ? : %v [%s:%s]\n", ip, isCN,
			country.Country.IsoCode, country.RegisteredCountry.IsoCode)
	}
	fmt.Printf("Time used: %s\n", time.Since(start))
}


func TestGeoIpQuery(t *testing.T) {
	if  GeoIpQuery("36.52.96.1") != "JP" {
		t.Failed()
	}
}