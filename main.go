package main

import (
   
    "os"
    "os/signal"
    "math"
    "syscall"
    "time"
    "unsafe"

    "github.com/tinyzimmer/go-gst/gst"
)

// Function to create a GStreamer element safely
func createElement(factory string) *gst.Element {
    element, err := gst.NewElement(factory)
    if err != nil {
        log.Fatalf("Failed to create element %s: %v", factory, err)
    }
    return element
}

// Function to calculate the RMS (Root Mean Square) value of the buffer
func calculateRMS(buffer *gst.Buffer) float64 {
    memory := buffer.GetMemory(0)
    mapInfo := memory.Map(gst.MapRead)
    if mapInfo == nil {
        log.Println("Failed to map buffer memory")
        return 0.0
    }
    defer memory.Unmap()

    data := (*[1 << 30]int16)(unsafe.Pointer(mapInfo.Data()))[:mapInfo.Size()/2]

    var sumSquares float64
    for i := 0; i < len(data); i++ {
        sample := float64(data[i])
        sumSquares += sample * sample
    }
    return math.Sqrt(sumSquares / float64(len(data)))
}

// Detect microphone input and mute/unmute the audio stream based on volume level (RMS)
func addLevelProbe(level *gst.Element, playbin *gst.Element) {
    level.GetStaticPad("src").AddProbe(gst.PadProbeTypeBuffer, func(pad *gst.Pad, info *gst.PadProbeInfo) gst.PadProbeReturn {
        buffer := info.GetBuffer()
        if buffer != nil {
            rmsLevel := calculateRMS(buffer)

            // Log the RMS level for debugging
            log.Printf("📢 RMS Level: %f", rmsLevel)

            // Define a threshold for speech detection
            speechThreshold := 10000.0 // Adjust this threshold as needed

            // If RMS value exceeds the threshold, mute the audio stream (speech detected)
            if rmsLevel > speechThreshold {
                log.Println("🎤 User is speaking, muting audio stream...")
                playbin.SetProperty("volume", 0.0) // Mute audio stream
            } else {
                log.Println("🤫 No significant noise detected, playing audio stream...")
                playbin.SetProperty("volume", 1.0) // Unmute audio stream
            }
        }
        return gst.PadProbeOK
    })
}

func main() {
    // Initialize GStreamer
    gst.Init(nil)

    // Create a new pipeline for audio streaming
    pipeline, err := gst.NewPipeline("audio-streamer")
    if err != nil {
        log.Fatalf("Failed to create pipeline: %v", err)
    }

    // Create the necessary GStreamer elements for the pipeline
    playbin := createElement("playbin")         // Playbin element for streaming
    microphone := createElement("pulsesrc")     // Microphone input (audio source)
    level := createElement("level")             // Level element to measure volume
    audioConvert := createElement("audioconvert")  // Convert audio formats if needed
    audioResample := createElement("audioresample") // Resample audio if needed
    audioSink := createElement("autoaudiosink")    // Sink to output audio to speakers
    noiseGate := createElement("audioamplify")  // Audio amplify element for noise gating

    // Set the URI for the stream (replace this URL with the audio stream URL)
    url := "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3" // Example URL
    playbin.SetProperty("uri", url)

    // Add all the elements to the pipeline so they can be used during processing
    pipeline.Add(playbin)         // Add the playbin element for audio streaming from a URL
    pipeline.Add(microphone)      // Add the microphone input element (pulsesrc for capturing audio)
    pipeline.Add(noiseGate)       // Add the noise gate element (audioamplify for controlling volume levels)
    pipeline.Add(level)           // Add the level element to measure the audio signal's volume (RMS)
    pipeline.Add(audioConvert)    // Add the audio converter to convert audio formats if necessary
    pipeline.Add(audioResample)   // Add the audio resampler to ensure the audio is in the correct sample rate
    pipeline.Add(audioSink)       // Add the audio sink to output the audio to the speakers or another destination

    // Link elements in the pipeline
    if err := microphone.Link(noiseGate); err != nil {
        log.Fatalf("Failed to link microphone to noise gate: %v", err)
    }
    if err := noiseGate.Link(level); err != nil {
        log.Fatalf("Failed to link noise gate to level: %v", err)
    }
    if err := level.Link(audioConvert); err != nil {
        log.Fatalf("Failed to link level to audioconvert: %v", err)
    }
    if err := audioConvert.Link(audioResample); err != nil {
        log.Fatalf("Failed to link audioconvert to audioresample: %v", err)
    }
    if err := audioResample.Link(audioSink); err != nil {
        log.Fatalf("Failed to link audioresample to audiosink: %v", err)
    }

    // Link playbin to the rest of the pipeline
    playbin.Link(audioConvert)

    // Add a probe to monitor volume levels and mute/unmute based on speech detection
    addLevelProbe(level, playbin)

    // Start playing the pipeline
    pipeline.SetState(gst.StatePlaying)
    log.Println("🔊 Streaming audio from URL!")

    // Stop pipeline after 5 seconds (adjustable)
    time.AfterFunc(5*time.Second, func() {
        log.Println("🛑 Stopping pipeline after 5 seconds...")
        pipeline.SetState(gst.StateNull) // Stop the pipeline
        log.Println("✅ Pipeline stopped. Exiting.")
        os.Exit(0)
    })

    // Wait for SIGINT or SIGTERM to cleanly exit
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
    <-sigc

    // Clean up: set the pipeline to Null state (stopped)
    pipeline.SetState(gst.StateNull)
}

