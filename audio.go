package chai

import (
	"syscall/js"
)

var audioContext js.Value

type AudioStream struct {
	audioBuffer js.Value
}

func LoadAudioFile(_file_path string) AudioStream {
	audioStream := AudioStream{}

	ch := make(chan js.Value)

	go func() {
		fetchPromise := js.Global().Get("fetch").Invoke(_file_path)
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

	audioStream.audioBuffer = <-ch
	return audioStream
}

func (a *AudioSourceComponent) Play(_audio_name string, _async bool) {
	audioContext.Call("resume").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// source := audioContext.Call("createBufferSource")
		// source.Set("buffer", a.audioSources[_audio_name].audioBuffer)
		// source.Call("connect", audioContext.Get("destination"))
		// source.Call("start", 0.0)
		if a.audioSourceData[_audio_name].isPlaying && !_async {
			a.audioSourceData[_audio_name].audioSource.Call("stop", 0.0)
		}

		audioSourceData := a.audioSourceData[_audio_name]

		gainNode := audioContext.Call("createGain")
		gainNode.Get("gain").Set("value", audioSourceData.volume)
		gainNode.Call("connect", audioContext.Get("destination"))
		audioSourceData.gainNode = gainNode

		audioSourceData.audioSource = audioContext.Call("createBufferSource")
		audioSourceData.audioSource.Set("buffer", audioSourceData.audioStream.audioBuffer)
		audioSourceData.audioSource.Call("connect", gainNode)

		audioSourceData.audioSource.Set("onended", js.FuncOf(func(this js.Value, args []js.Value) any {
			audioSourceData.isPlaying = false
			if audioSourceData.loop {
				a.Play(_audio_name, _async)
			}
			return nil
		}))

		audioSourceData.audioSource.Call("start", 0.0)
		audioSourceData.isPlaying = true

		a.audioSourceData[_audio_name] = audioSourceData
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
	audioSource audioJsSource
	gainNode    gainJsNode
	audioStream *AudioStream
	isPlaying   bool
	volume      float32
	loop        bool
}

type AudioSourceComponent struct {
	audioSourceData map[string]AudioSourceData
}

func NewAudioSourceComponent() AudioSourceComponent {
	return AudioSourceComponent{
		audioSourceData: make(map[string]AudioSourceData),
	}
}

func (a *AudioSourceComponent) SetVolume(_audio_name string, _value float32) {
	audioSourceData := a.audioSourceData[_audio_name]

	audioSourceData.volume = _value
	audioSourceData.gainNode.Get("gain").Set("volume", audioSourceData.volume)

	a.audioSourceData[_audio_name] = audioSourceData
}

func (a *AudioSourceComponent) GetVolume(_audio_name string) float32 {
	return a.audioSourceData[_audio_name].volume
}

func (a *AudioSourceComponent) SetLoop(_audio_name string, _value bool) {
	audioS := a.audioSourceData[_audio_name]
	audioS.loop = _value
	a.audioSourceData[_audio_name] = audioS
}

func (a *AudioSourceComponent) AddAudioSource(_audio_name string, _audio_stream AudioStream) {
	audioSourceData := a.audioSourceData[_audio_name]
	audioSourceData.volume = 1.0
	audioSourceData.loop = false

	gainNode := audioContext.Call("createGain")
	gainNode.Get("gain").Set("value", audioSourceData.volume)
	gainNode.Call("connect", audioContext.Get("destination"))
	audioSourceData.gainNode = gainNode

	source := audioContext.Call("createBufferSource")
	source.Set("buffer", _audio_stream.audioBuffer)
	source.Call("connect", gainNode)

	source.Set("onended", js.FuncOf(func(this js.Value, args []js.Value) any {
		audioSourceData.isPlaying = false
		if audioSourceData.loop {
			a.Play(_audio_name, false)
		}
		return nil
	}))

	audioSourceData.audioSource = source
	audioSourceData.audioStream = &_audio_stream
	audioSourceData.isPlaying = false

	a.audioSourceData[_audio_name] = audioSourceData
}

// type AudioPlaySystem struct {
// 	EcsSystem
// }

// func (a *AudioPlaySystem) Update(dt float32) {

// }
