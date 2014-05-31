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

import (
	"bytes"
	"testing"
)

func TestBitwriter_writeBits(t *testing.T) {
	var buf bytes.Buffer
	w := &bitWriter{output: &buf}

	for i := 0; i < 1234; i++ {
		if err := w.writeBits(i, 0); err != nil {
			t.Error("Unexpected error:", err)
			return
		}
	}

	for i := 0; i < 4096; i++ {
		for j := 0; j < 16; j++ {
			if err := w.writeBits(j+(i<<uint(j+1)), j+1); err != nil {
				t.Error("Unexpected error:", err)
				return
			}
			if err := w.writeBits(i*j, 0); err != nil {
				t.Error("Unexpected error:", err)
				return
			}
		}
	}
	w.flush()

	testData := []byte{
		0x28, 0xc8, 0x28, 0x60, 0x70, 0x40, 0x12, 0x02,
		0x80, 0x2c, 0x01, 0x80, 0x06, 0x80, 0x0e, 0x00,
		0x0f,
	}
	out := buf.Bytes()

	if len(out) != 4096*len(testData) {
		t.Errorf("Output length: %d, expected %d", len(out), 4096*len(testData))
		return
	}

	for k, v := range out {
		i := k / len(testData)
		j := k % len(testData)
		e := testData[j]
		if v != e {
			t.Errorf("i = %d, j = %d, k = %d, v = %d, e = %d", i, j, k, v, e)
			return
		}
	}
}
