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

func TestIdct(t *testing.T) {
	testIdct(t, new(idct2Test))
	testIdct(t, new(idct4Test))
	testIdct(t, new(idct8Test))
	testIdct(t, new(idct16Test))
	testIdct(t, new(idct32Test))
}

func testIdct(t *testing.T, it idctTester) {
	rand.Seed(42) // make it repeatable
	max, X := 0.0, it.slice()
	for i := 0; i < 1000; i++ {
		for i := range X {
			X[i] = 2*rand.Float32() - 1
		}
		x := directIdct(X)
		it.transform()
		for i := range X {
			max = math.Max(max, math.Abs(x[i]-float64(X[i])))
		}
	}

	t.Logf("N = %d, max. difference = %e", len(X), max)
	if max >= 1.0/(1<<16) {
		t.Fail()
	}
}

func directIdct(X []float32) []float64 {
	N, N_fl := len(X), float64(len(X))
	x := make([]float64, N)
	for k := 0; k < N; k++ {
		k_fl := float64(k)
		for n := 0; n < N; n++ {
			n_fl := float64(n)
			x[k] += float64(X[n]) * math.Cos(math.Pi*(2*k_fl+1)*n_fl/(2*N_fl))
		}
	}
	return x
}

type idctTester interface {
	slice() []float32
	transform()
}

type idct2Test [2]float32

func (i *idct2Test) slice() []float32 {
	return i[:]
}

func (i *idct2Test) transform() {
	idct2((*[2]float32)(i))
}

type idct4Test [4]float32

func (i *idct4Test) slice() []float32 {
	return i[:]
}

func (i *idct4Test) transform() {
	idct4((*[4]float32)(i))
}

type idct8Test [8]float32

func (i *idct8Test) slice() []float32 {
	return i[:]
}

func (i *idct8Test) transform() {
	idct8((*[8]float32)(i))
}

type idct16Test [16]float32

func (i *idct16Test) slice() []float32 {
	return i[:]
}

func (i *idct16Test) transform() {
	idct16((*[16]float32)(i))
}

type idct32Test [32]float32

func (i *idct32Test) slice() []float32 {
	return i[:]
}

func (i *idct32Test) transform() {
	idct32(i[:])
}
