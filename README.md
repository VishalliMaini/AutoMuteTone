# Audio Mixer with GStreamer in Go

This project implements an audio mixer using GStreamer in Go. The application takes audio from a microphone and a test tone generator, adjusting the output based on the microphone's audio level.

## Table of Contents

- [Setup](#setup)
- [Approach](#approach)
- [Challenges Faced](#challenges-faced)
- [License](#license)

## Setup

### Prerequisites

1. **Go**: Install Go on your machine. Download it from [golang.org](https://golang.org/dl/). Install Go using the following commands:

   **For Ubuntu/Debian:**
   ```bash
   sudo apt update
   sudo apt install golang-go
   ```
   
   **For macOS (using Homebrew):**
   ```bash
   brew install go
   ```
   
   **For Windows:** Download the installer from [golang.org](https://golang.org/dl/) and follow the installation instructions.
   

2. **GStreamer**: Install GStreamer on your system. You can find installation instructions for various platforms on the [GStreamer website](https://gstreamer.freedesktop.org/documentation/installing/index.html). Here are some common installation commands:

   **For Ubuntu/Debian:**
   ```bash
   sudo apt update
   sudo apt install gstreamer1.0-tools gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav
   ```
   
   **For macOS (using Homebrew):**
   ```bash
   brew install gstreamer gst-plugins-base gst-plugins-good gst-plugins-bad gst-plugins-ugly gst-libav
   ```
   
   **For Windows:** Download the GStreamer installer from the [GStreamer website](https://gstreamer.freedesktop.org/) and follow the installation instructions.

3. **Go GStreamer Bindings**: This project uses the `go-gst` bindings. You can install them using the following command:
   ```bash
   go get github.com/tinyzimmer/go-gst/gst
   ```

### Running the Application

1. Create a new directory for the project:
   ```bash
   mkdir audio-mute
   ```
2. Change into the project directory:
   ```bash
   cd audio-mute
   ```
3. Initialize a new Go module:
   ```bash
   go mod init audio-mute
   ```
   This will creat a new go.mod

4. Create the main Go file:
   ```bash
   nano main.go
   ```
   Write your application code in `main.go`, then save and exit the editor.

5. Tidy up the module dependencies and add module requirements and sums:
   ```bash
   go mod tidy
   ```
6. Build and run the application:
   ```bash
   go run main.go
   ```

## Approach

The application is structured as follows:

1. **Initialization**: The GStreamer library is initialized, and a new pipeline is created.


2. **Element Creation**: Various GStreamer elements are created:

   - `audiotestsrc`: Generates a test tone.
   
   - `autoaudiosrc`: Captures audio from the default microphone.
   
   - `level`: Monitors the audio level of the microphone input.
   
   - `audiomixer`: Mixes the audio from the test tone and the microphone.
   
   - `autoaudiosink`: Outputs the mixed audio to the default audio output device.


3. **Linking Elements**: The elements are linked together to form a complete audio processing pipeline.


4. **Adding a Probe**: A probe is added to the `level` element to monitor the RMS (Root Mean Square) level of the microphone input. Based on the RMS value, the application decides whether to play or stop the test tone.


5. **Pipeline Execution**: The pipeline is set to the `PLAYING` state, and the application enters a loop to wait for events (End of Stream or errors) or termination signals.


6. **Cleanup**: Upon termination, the pipeline is set to the `NULL` state to clean up resources.


## Challenges Faced

1. **Element Linking**: Ensuring that all elements were correctly linked was a challenge. GStreamer requires that elements be compatible in terms of their pads, and debugging link failures can be tricky.


2. **RMS Value Interpretation**: Understanding how to interpret the RMS value from the `level` element was initially confusing. It required some experimentation to determine appropriate thresholds for silence detection.


3. **Cross-Platform Compatibility**: Ensuring that the application works across different operating systems (Linux, macOS, Windows) required testing and adjustments, particularly in audio source configurations.


4. **Background Noise**: Managing background noise from the microphone input posed a challenge. The application needed to differentiate between actual speech and ambient noise, which sometimes led to false positives in triggering the test tone playback.

