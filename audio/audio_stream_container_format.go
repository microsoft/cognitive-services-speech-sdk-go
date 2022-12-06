// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package audio

// AudioStreamContainerFormat defines supported audio stream container format.
type AudioStreamContainerFormat int //nolint:revive

const (
	// OGGOPUS Stream ContainerFormat definition for OGG OPUS.
	OGGOPUS AudioStreamContainerFormat = 0x101

	// MP3 Stream ContainerFormat definition for MP3.
	MP3 AudioStreamContainerFormat = 0x102

	// FLAC Stream ContainerFormat definition for FLAC.
	FLAC AudioStreamContainerFormat = 0x103

	// ALAW Stream ContainerFormat definition for ALAW.
	ALAW AudioStreamContainerFormat = 0x104

	// MULAW Stream ContainerFormat definition for MULAW.
	MULAW AudioStreamContainerFormat = 0x105

	// AMRNB Stream ContainerFormat definition for AMRNB. Currently not supported.
	AMRNB AudioStreamContainerFormat = 0x106

	// AMRWB Stream ContainerFormat definition for AMRWB. Currently not supported.
	AMRWB AudioStreamContainerFormat = 0x107

	// ANY Stream ContainerFormat definition when the actual stream format is not known.
	ANY AudioStreamContainerFormat = 0x108
)

// AudioStreamWaveFormat represents the format specified inside WAV container which are sent directly as encoded to the speech service.
type AudioStreamWaveFormat int //nolint:revive

const (
	// AudioStreamWaveFormat definition for PCM (pulse-code modulated) data in integer format.
	WavePCM AudioStreamWaveFormat = 0x0001

	// AudioStreamWaveFormat definition A-law-encoded format.
	WaveALAW AudioStreamWaveFormat = 0x0006

	// AudioStreamWaveFormat definition for Mu-law-encoded format.
	WaveMULAW AudioStreamWaveFormat = 0x0007

)
