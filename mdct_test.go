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
	"math"
	"math/rand"
	"testing"
)

func TestMdct(t *testing.T) {
	testMdct(t, new(mdct4Test))
	testMdct(t, new(mdct12Test))
	testMdct(t, new(mdct36Test))
}

func testMdct(t *testing.T, mt mdctTester) {
	rand.Seed(42) // make it repeatable
	x := mt.input()
	N := len(x)
	X1 := make([]float64, N/2)
	max := 0.0
	for i := 0; i < 1000; i++ {
		for n := range x {
			x[n] = 2*rand.Float32() - 1
		}
		directMdct(x, X1)
		X2 := mt.transform()
		for n := range X1 {
			max = math.Max(max, math.Abs(X1[n]-float64(X2[n])))
		}
	}

	t.Logf("N = %d, max. difference = %e", N, max)
	if max >= 1.0/(1<<16) {
		t.Fail()
	}
}

func directMdct(in []float32, out []float64) {
	N, Nf := len(in), float64(len(in))
	for k := 0; k < N/2; k++ {
		kf := float64(k)
		out[k] = 0
		for n := 0; n < N; n++ {
			nf := float64(n)
			in64 := float64(in[n])
			out[k] += in64 * math.Cos(math.Pi/(2*Nf)*(2*nf+1+Nf/2)*(2*kf+1))
		}
	}
}

type mdctTester interface {
	input() []float32
	transform() []float32
}

type mdct4Test struct {
	in  [4]float32
	out [2]float32
}

func (m *mdct4Test) input() []float32 {
	return m.in[:]
}

func (m *mdct4Test) transform() []float32 {
	mdct4(&m.in, &m.out)
	return m.out[:]
}

type mdct12Test struct {
	in  [12]float32
	out [6]float32
}

func (m *mdct12Test) input() []float32 {
	return m.in[:]
}

func (m *mdct12Test) transform() []float32 {
	mdct12(m.in[:], m.out[:])
	return m.out[:]
}

type mdct36Test struct {
	in  [36]float32
	out [18]float32
}

func (m *mdct36Test) input() []float32 {
	return m.in[:]
}

func (m *mdct36Test) transform() []float32 {
	mdct36(m.in[:], m.out[:])
	return m.out[:]
}
