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

// mdct4 performs a direct computation of the MDCT, while mdct12 and mdct36 use
// the method described in the following paper:
//
//   H. Shu, X. Bao, Ch. Toumoulin, L. Luo, "Radix-3 Algorithm for the Fast
//   Computation of Forward and Inverse MDCT", IEEE Signal Processing Letters,
//   Vol. 14, N. 2, pp. 93-96, Feb. 2007
//
// The MDCT with a window size of N is defined as
//
//   X[k] = Sum(n=0...N-1) { x[n] * cos(π/(2*N) * (2*n + 1 + N/2) * (2*k + 1)) }
//

// mdct36 computes the MDCT for N = 36.
func mdct36(in, out []float32) {
	var (
		inA, inB, inC    [12]float32
		outA, outB, outC [6]float32
	)

	for n := 0; n < 12; n++ {
		a, b, c := in[n], in[12+n], in[24+n]
		ad, bd, cd := in[11-n], in[23-n], in[35-n]

		inA[n] = ad - bd + cd
		inB[n] = (2*cd-ad+bd)*imdct36s[n] - (ad+bd)*imdct36c3[n]
		inC[n] = (2*a+b-c)*imdct36c[n] + (b+c)*imdct36s3[n]
	}

	mdct12(inA[:], outA[:])
	mdct12(inB[:], outB[:])
	mdct12(inC[:], outC[:])

	for k := 1; k < 6; k += 2 {
		outC[k] *= -1
	}

	for k := 0; k < 6; k++ {
		out[3*k+1] = outA[k]
		out[3*k] = (outB[k] + outC[k]) / 2
		out[3*k+2] = (outB[k] - outC[k]) / 2
	}
}

// mdct12 computes the MDCT for N = 12.
func mdct12(in, out []float32) {
	var (
		inA, inB, inC    [4]float32
		outA, outB, outC [2]float32
	)

	for n := 0; n < 4; n++ {
		a, b, c := in[n], in[4+n], in[8+n]
		ad, bd, cd := in[3-n], in[7-n], in[11-n]

		inA[n] = ad - bd + cd
		inB[n] = (2*cd-ad+bd)*imdct12s[n] - (ad+bd)*imdct12c3[n]
		inC[n] = (2*a+b-c)*imdct12c[n] + (b+c)*imdct12s3[n]
	}

	mdct4(&inA, &outA)
	mdct4(&inB, &outB)
	mdct4(&inC, &outC)

	outC[1] *= -1

	for k := 0; k < 2; k++ {
		out[3*k+1] = outA[k]
		out[3*k] = (outB[k] + outC[k]) / 2
		out[3*k+2] = (outB[k] - outC[k]) / 2
	}
}

// mdct4 computes the MDCT for N = 4.
func mdct4(in *[4]float32, out *[2]float32) {
	x1, x2 := in[0]-in[1], in[2]+in[3]
	out[0] = x1*imdct4c0 + x2*imdct4c1
	out[1] = x1*imdct4c1 - x2*imdct4c0
}
