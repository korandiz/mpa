// mpa, an MPEG-1 Audio library
// Copyright (C) 2014 KORÁNDI Zoltán <korandi.z@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License, version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// Please note that, being hungarian, my last name comes before my first
// name. That's why it's in all caps, and not because I like to shout my
// name. So please don't start your emails with "Hi Korandi" or "Dear Mr.
// Zoltan", because it annoys the hell out of me. Thanks.

package mpa

// mdctFilter performs the windowing and MDCT steps of Layer III encoding.
type mdctFilter [18]float32

// filter transforms the 18 subband samples in the input array, and writes the
// result in the output array.
func (f *mdctFilter) filter(input, output []float32, typ int) {
	var tmp [36]float32
	copy(tmp[0:], f[:])
	copy(tmp[18:], input)
	copy(f[:], input)

	if typ != 2 {
		for t := 0; t < 36; t++ {
			tmp[t] *= mdctWindows[typ][t]
		}
		mdct36(tmp[:], output)
	} else {
		for w := 0; w < 3; w++ {
			var wInput [12]float32
			for t := 0; t < 12; t++ {
				wInput[t] = tmp[6+6*w+t] * mdctWindows[2][t]
			}
			mdct12(wInput[:], output[6*w:])
		}
	}
}
