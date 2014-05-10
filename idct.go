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

// All functions in this file use Lee's algorithm to compute the inverse
// discrete cosine transform. For details, see:
//
//   B. G. Lee, "A New Algorithm to Compute The Discrete Cosine Transform",
//   IEEE Transactions on Acoustics, Speech and Signal Processing, Vol. 32,
//   N. 6, pp. 1243-1245, Dec. 1984
//
// The Nth-order IDCT is defined by the following formula:
//
//   x[k] = Sum(n=0...N-1) { X[n] * cos(π * (2*k + 1) * n / (2*N)) }
//

// idct32 computes the 32nd-order inverse discrete cosine transform.
func idct32(data []float32) {
	var even, odd [16]float32
	for i := 0; i < 16; i++ {
		even[i] = data[2*i]
		odd[i] = data[2*i+1]
	}

	for i := 15; i >= 1; i-- {
		odd[i] += odd[i-1]
	}

	idct16(&even)
	idct16(&odd)

	for i := 0; i < 16; i++ {
		tmp := odd[i] * dct32c[i]
		data[i] = even[i] + tmp
		data[31-i] = even[i] - tmp
	}
}

// idct16 computes the 16th-order inverse discrete cosine transform.
func idct16(data *[16]float32) {
	var even, odd [8]float32
	for i := 0; i < 8; i++ {
		even[i] = data[2*i]
		odd[i] = data[2*i+1]
	}

	for i := 7; i >= 1; i-- {
		odd[i] += odd[i-1]
	}

	idct8(&even)
	idct8(&odd)

	for i := 0; i < 8; i++ {
		tmp := odd[i] * dct16c[i]
		data[i] = even[i] + tmp
		data[15-i] = even[i] - tmp
	}
}

// idct8 computes the 8th-order inverse discrete cosine transform.
func idct8(data *[8]float32) {
	var even, odd [4]float32
	for i := 0; i < 4; i++ {
		even[i] = data[2*i]
		odd[i] = data[2*i+1]
	}

	for i := 3; i >= 1; i-- {
		odd[i] += odd[i-1]
	}

	idct4(&even)
	idct4(&odd)

	for i := 0; i < 4; i++ {
		tmp := odd[i] * dct8c[i]
		data[i] = even[i] + tmp
		data[7-i] = even[i] - tmp
	}
}

// idct4 computes the 4th-order inverse discrete cosine transform.
func idct4(data *[4]float32) {
	var even, odd [2]float32
	for i := 0; i < 2; i++ {
		even[i] = data[2*i]
		odd[i] = data[2*i+1]
	}

	odd[1] += odd[0]

	idct2(&even)
	idct2(&odd)

	for i := 0; i < 2; i++ {
		tmp := odd[i] * dct4c[i]
		data[i] = even[i] + tmp
		data[3-i] = even[i] - tmp
	}
}

// idct2 computes the 2nd-order inverse discrete cosine transform.
func idct2(data *[2]float32) {
	even := data[0]
	odd := data[1]

	tmp := dct2c * odd
	data[0] = even + tmp
	data[1] = even - tmp
}
