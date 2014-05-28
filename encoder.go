package mpa

import "io"

type InvalidParameter string

func (err InvalidParameter) Error() string {
	return string(err)
}

type Encoder struct {
	Samples [2][1152]float32
	Output  io.Writer

	// Header fields
	Layer             int
	Bitrate           int
	SamplingFrequency int
	Mode              int
	Copyrighted       bool
	Original          bool
	Emphasis          int

	// All layers
	bitrateIndex int
	samplingFreq int
	nChannels    int
	anal         [2]analysisFilter
	samples      [2][1152]float32
}

// All layers

func (e *Encoder) EncodeFrame() error {
	var err error

	if err = e.checkParameters(); err != nil {
		return err
	}

	e.analyzeInput()

	switch e.Layer {
	case 1:
		err = e.encodeFrame1()
	case 2:
		err = e.encodeFrame2()
	case 3:
		err = e.encodeFrame3()
	}

	return err
}

func (e *Encoder) Flush() error {
	// TODO
	return nil
}

func (e *Encoder) checkParameters() error {
	if e.Layer < 1 || e.Layer > 3 {
		return InvalidParameter("Layer invalid")
	}

	if e.Bitrate < bitrateBps[e.Layer][1] {
		return InvalidParameter("Bitrate too low")
	} else if e.Bitrate > bitrateBps[e.Layer][14] {
		return InvalidParameter("Bitrate too high")
	}
	e.bitrateIndex = 0 // FreeFormat
	for i := 1; i < 15; i++ {
		if e.Bitrate == bitrateBps[e.Layer][i] {
			e.bitrateIndex = i
			break
		}
	}

	switch e.SamplingFrequency {
	case samplingFreqHz[0]:
		e.samplingFreq = 0
	case samplingFreqHz[1]:
		e.samplingFreq = 1
	case samplingFreqHz[2]:
		e.samplingFreq = 2
	default:
		return InvalidParameter("SamplingFrequency invalid")
	}

	switch e.Mode {
	case ModeStereo:
		fallthrough
	case ModeJointStereo:
		fallthrough
	case ModeDualChannel:
		e.nChannels = 2
	case ModeMono:
		e.nChannels = 1
	default:
		return InvalidParameter("Mode invalid")
	}

	switch e.Emphasis {
	case EmphNone:
	case Emph5015:
	case EmphCCITT:
	default:
		return InvalidParameter("Emphasis invalid")
	}

	return nil
}

func (e *Encoder) analyzeInput() {
	samplesPerSubband := 36
	if e.Layer == 1 {
		samplesPerSubband = 12
	}

	for ch := 0; ch < e.nChannels; ch++ {
		for s := 0; s < samplesPerSubband; s++ {
			in := e.Samples[ch][32*s : 32*s+32]
			out := e.samples[ch][32*s : 32*s+32]
			for i := 0; i < 32; i++ {
				sample := in[i]
				if sample < -1 {
					sample = -1
				} else if sample > 1 {
					sample = 1
				} else if sample != sample { // NaN
					sample = 0
				}
				out[i] = sample
			}
			e.anal[ch].filter(out)
		}
	}
}

// Layer I

func (e *Encoder) encodeFrame1() error {
	// TODO
	return nil
}

// Layer II

func (e *Encoder) encodeFrame2() error {
	// TODO
	return nil
}

// Layer III

func (e *Encoder) encodeFrame3() error {
	// TODO
	return nil
}
