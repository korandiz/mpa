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

import "io"

// A bitWriter wraps an io.Writer and allows it to be written bit-by-bit.
type bitWriter struct {
	output io.Writer
	buffer [4096]byte
	tail   int // index of the last byte in the buffer not fully written
	bits   int // number of bits already written to buffer[tail] (0 <= bits < 8)
}

// flush pads the stream with zeros until the next byte boundary and writes the
// contents of the buffer to the output.
func (wr *bitWriter) flush() error {
	if wr.bits > 0 {
		wr.buffer[wr.tail] <<= uint(8 - wr.bits)
		wr.tail++
		wr.bits = 0
	}

	n, err := wr.output.Write(wr.buffer[0:wr.tail])

	if n < wr.tail {
		if n < 0 {
			n = 0
		}
		copy(wr.buffer[0:], wr.buffer[n:wr.tail])
		wr.tail -= n
		if err == nil {
			err = io.ErrShortWrite
		}
		return err
	}

	wr.tail = 0
	return nil
}

// writeBits writes the n least significant bits of v to the stream.
func (wr *bitWriter) writeBits(v, n int) error {
	for n > 0 {
		if wr.tail == len(wr.buffer) {
			if err := wr.flush(); err != nil {
				return err
			}
		}

		k := 8 - wr.bits
		if n < k {
			k = n
		}
		wr.buffer[wr.tail] <<= uint(k)
		wr.buffer[wr.tail] |= byte((v >> uint(n-k)) &^ (^0 << uint(k)))
		wr.bits += k
		n -= k
		if wr.bits == 8 {
			wr.bits = 0
			wr.tail++
		}
	}
	return nil
}
