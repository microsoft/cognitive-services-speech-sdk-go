//
// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
//
package audio

// AudioStreamContainerFormat defines supported audio stream container format.
// Changed in version 1.4.0.
type AudioStreamContainerFormat int

const (
	// OGGOPUS Stream ContainerFormat definition for OGG OPUS.
	OGGOPUS AudioStreamContainerFormat = 0x101

	// MP3 Stream ContainerFormat definition for MP3.
	MP3 AudioStreamContainerFormat = 0x102

	// FLAC Stream ContainerFormat definition for FLAC. Added in version 1.7.0.
	FLAC AudioStreamContainerFormat = 0x103

	// ALAW Stream ContainerFormat definition for ALAW. Added in version 1.7.0.
	ALAW AudioStreamContainerFormat = 0x104

	// MULAW Stream ContainerFormat definition for MULAW. Added in version 1.7.0.>
	MULAW AudioStreamContainerFormat = 0x105

	// AMRNB Stream ContainerFormat definition for AMRNB. Currently not supported.
	AMRNB AudioStreamContainerFormat = 0x106

	// AMRWB Stream ContainerFormat definition for AMRWB. Currently not supported.
	AMRWB AudioStreamContainerFormat = 0x107
)
