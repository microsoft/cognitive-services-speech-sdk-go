// Copyright (c) Microsoft. All rights reserved.
// Licensed under the MIT license. See LICENSE.md file in the project root for full license information.

package helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
)

func PumpFileIntoStream(filename string, stream *audio.PushAudioInputStream) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1000)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			fmt.Println("Done reading file.")
			break
		}
		if err != nil {
			fmt.Println("Error reading file: ", err)
			break
		}
		err = stream.Write(buffer[0:n])
		if err != nil {
			fmt.Println("Error writing to the stream")
		}
	}
	stream.CloseStream()
}
