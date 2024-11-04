package chai

import (
	"syscall/js"

	"github.com/mhamedGd/chai/customtypes"
)

var audioContext js.Value

type AudioStream struct {
	m_AudioBuffer js.Value
}

func LoadAudioFile(_filePath string) AudioStream {
	m_AudioStream := AudioStream{}

	ch := make(chan js.Value)

	go func() {
		fetchPromise := js.Global().Get("fetch").Invoke(_filePath)
		fetchPromise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			response := args[0]
			return response.Call("arrayBuffer")
		})).Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			arrayBuffer := args[0]
			return audioContext.Call("decodeAudioData", arrayBuffer)
		})).Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			decodedAudio := args[0]
			ch <- decodedAudio
			return decodedAudio
		}))

	}()

	m_AudioStream.m_AudioBuffer = <-ch
	return m_AudioStream
}

func (a *AudioSourceComponent) Play(_audioName string, _async bool) {
	audioContext.Call("resume").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if a.audioSourceData.Get(_audioName).m_IsPlaying && !_async {
			a.audioSourceData.Get(_audioName).m_AudioSource.Call("stop", 0.0)
		}

		audioSourceData := a.audioSourceData.Get(_audioName)

		m_GainNode := audioContext.Call("createGain")
		m_GainNode.Get("gain").Set("value", audioSourceData.m_Volume)
		m_GainNode.Call("connect", audioContext.Get("destination"))
		audioSourceData.m_GainNode = m_GainNode

		audioSourceData.m_AudioSource = audioContext.Call("createBufferSource")
		audioSourceData.m_AudioSource.Get("playbackRate").Set("value", a.audioSourceData.Get(_audioName).m_Pitch)
		audioSourceData.m_AudioSource.Set("buffer", audioSourceData.m_AudioStream.m_AudioBuffer)
		audioSourceData.m_AudioSource.Call("connect", m_GainNode)

		audioSourceData.m_AudioSource.Set("onended", js.FuncOf(func(this js.Value, args []js.Value) any {
			audioSourceData.m_IsPlaying = false
			if audioSourceData.m_Loop {
				a.Play(_audioName, _async)
			}
			return nil
		}))

		audioSourceData.m_AudioSource.Call("start", 0.0)
		audioSourceData.m_IsPlaying = true

		a.audioSourceData.Set(_audioName, audioSourceData)
		return nil
	}))
}

func SuspendAudioContext() {
	audioContext.Call("suspend")
}
func ResumeAudioContext() {
	audioContext.Call("resume")
}

func triggerAudioContextPlaying() {
	if audioContext.Get("state").String() == "running" {
		SuspendAudioContext()
	} else {
		ResumeAudioContext()
	}
}

type audioJsSource = js.Value
type gainJsNode = js.Value

type AudioSourceData struct {
	m_AudioSource audioJsSource
	m_GainNode    gainJsNode
	m_AudioStream *AudioStream
	m_IsPlaying   bool
	m_Volume      float32
	m_Pitch       float32
	m_Loop        bool
}

type AudioSourceComponent struct {
	audioSourceData customtypes.Map[string, AudioSourceData]
}

func NewAudioSourceComponent() AudioSourceComponent {
	return AudioSourceComponent{
		audioSourceData: customtypes.NewMap[string, AudioSourceData](),
	}
}

func (a *AudioSourceComponent) SetVolume(_audioName string, _value float32) {
	audioSourceData := a.audioSourceData.Get(_audioName)

	audioSourceData.m_Volume = _value
	audioSourceData.m_GainNode.Get("gain").Set("m_Volume", audioSourceData.m_Volume)

	a.audioSourceData.Set(_audioName, audioSourceData)
}

func (a *AudioSourceComponent) GetVolume(_audioName string) float32 {
	return a.audioSourceData.Get(_audioName).m_Volume
}

func (a *AudioSourceComponent) SetPitch(_audioName string, _value float32) {
	s := a.audioSourceData.Get(_audioName)
	s.m_Pitch = _value
	a.audioSourceData.Set(_audioName, s)
}

func (a *AudioSourceComponent) SetLoop(_audioName string, _value bool) {
	audioS := a.audioSourceData.Get(_audioName)
	audioS.m_Loop = _value
	a.audioSourceData.Set(_audioName, audioS)
}

func (a *AudioSourceComponent) AddAudioSource(_audioName string, _audioStream AudioStream) {
	audioSourceData := a.audioSourceData.Get(_audioName)
	audioSourceData.m_Volume = 1.0
	audioSourceData.m_Loop = false

	m_GainNode := audioContext.Call("createGain")
	m_GainNode.Get("gain").Set("value", audioSourceData.m_Volume)
	m_GainNode.Call("connect", audioContext.Get("destination"))
	audioSourceData.m_GainNode = m_GainNode

	source := audioContext.Call("createBufferSource")
	source.Set("playbackRate", audioSourceData.m_Pitch)
	source.Set("buffer", _audioStream.m_AudioBuffer)
	source.Call("connect", m_GainNode)

	source.Set("onended", js.FuncOf(func(this js.Value, args []js.Value) any {
		audioSourceData.m_IsPlaying = false
		if audioSourceData.m_Loop {
			a.Play(_audioName, false)
		}
		return nil
	}))

	audioSourceData.m_AudioSource = source
	audioSourceData.m_AudioStream = &_audioStream
	audioSourceData.m_IsPlaying = false

	a.audioSourceData.Set(_audioName, audioSourceData)
}
