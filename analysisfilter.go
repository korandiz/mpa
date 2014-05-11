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

// analysisFilter implements the analysis subband filter (Figure 3-C.1).
type analysisFilter [512]float32

// filter feeds 32 PCM samples to the filterbank and computes the next sample
// for every subband.
func (f *analysisFilter) filter(x []float32) {
	// Shifting
	copy(f[32:], f[0:480])

	// Input samples
	for i := 0; i < 32; i++ {
		f[i] = x[31-i]
	}

	// Windowing
	var y [64]float32
	for i := 0; i < 64; i++ {
		var t float32
		for p := i; p < 512; p += 64 {
			t += analysisWindow[p] * f[p]
		}
		y[i] = t
	}

	// Matrixing can be carried out efficiently using the inverse discrete
	// cosine transform. For details, see:
	//
	//   K. Konstantinides, "Fast Subband Filtering in MPEG Audio Coding",
	//   IEEE Signal Processing Letters, Vol. 1, N. 2, pp. 26-28, Feb. 1994
	//
	x[0] = y[16]
	for k := 1; k <= 16; k++ {
		x[k] = y[k+16] + y[16-k]
	}
	for k := 17; k <= 31; k++ {
		x[k] = y[k+16] - y[80-k]
	}
	idct32(x)
}
