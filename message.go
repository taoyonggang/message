// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
This go package is an encoder/decoder library for MQTT 3.1.1
(http://docs.oasis-open.org/mqtt/mqtt/v3.1.1) messages.

Create a new message and encode it:

	// Create new message.
	msg1 := NewConnectMessage()
	msg1.Username = []byte("gomqtt")
	msg1.Password = []byte("amazing!")

	// Allocate buffer.
	buf := make([]byte, msg1.Len())

	// Encode the message.
	if _, err := msg1.Encode(buf); err != nil {
		// there was an error while encoding
		panic(err)
	}

Decode bytes to a message:

	// Detect message.
	l, mt := DetectMessage(buf)

	// Check length
	if l == 0 {
		// buffer not complete yet
		return
	}

	// Create message.
	msg2, err := mt.New();
	if err != nil {
		// message type is invalid
		panic(err)
	}

	// Decode message.
	_, err = msg2.Decode(buf)
	if err != nil {
		// there was an error while decoding
		panic(err)
	}
*/
package message

const (
	// QoS 0: At most once delivery
	// The message is delivered according to the capabilities of the underlying network.
	// No response is sent by the receiver and no retry is performed by the sender. The
	// message arrives at the receiver either once or not at all.
	QosAtMostOnce byte = iota

	// QoS 1: At least once delivery
	// This quality of service ensures that the message arrives at the receiver at least once.
	// A QoS 1 PUBLISH Packet has a Packet Identifier in its variable header and is acknowledged
	// by a PUBACK Packet. Section 2.3.1 provides more information about Packet Identifiers.
	QosAtLeastOnce

	// QoS 2: Exactly once delivery
	// This is the highest quality of service, for use when neither loss nor duplication of
	// messages are acceptable. There is an increased overhead associated with this quality of
	// service.
	QosExactlyOnce

	// QosFailure is a return value for a subscription if there's a problem while subscribing
	// to a specific topic.
	QosFailure = 0x80
)

const (
	// MQTT protocol version 3.1.
	Version31 byte = 3

	// MQTT protocol version 3.1.1.
	Version311 byte = 4
)

// Message is an interface defined for all MQTT message types.
type Message interface {
	// Type return the messages message type.
	Type() MessageType

	// Len returns the byte length of the message.
	Len() int

	// Decode reads the bytes in the byte slice from the argument. It returns the
	// total number of bytes decoded, and whether there have been any errors during
	// the process. The byte slice MUST NOT be modified during the duration of this
	// message being available since the byte slice never gets copied.
	Decode([]byte) (int, error)

	// Encode writes the message bytes into the byte array from the argument. It
	// returns the number of bytes encoded and whether there's any errors along
	// the way. If there's any errors, then the byte slice and count should be
	// considered invalid.
	Encode([]byte) (int, error)

	// String returns a string representation of the message.
	String() string
}

/*
A basic fuzzing test that works with https://github.com/dvyukov/go-fuzz:

	$ go-fuzz-build github.com/gomqtt/message
	$ go-fuzz -bin=./message-fuzz.zip -workdir=./fuzz
*/
func Fuzz(data []byte) int {
	// check for zero length data
	if len(data) == 0 {
		return 1
	}

	// Detect message.
	l, mt := DetectMessage(data)

	// Check length
	if l == 0 {
		// for testing purposes we will not cancel
		// on incomplete buffers
	}

	// Create a new message
	msg, err := mt.New()
	if err != nil {
		return 0
	}

	// Decode it from the buffer.
	_, err = msg.Decode(data)
	if err != nil {
		return 0
	}

	// Prepare buffer with 255 as zero value
	buf := make([]byte, msg.Len())
	for i := 0; i < len(buf); i++ {
		buf[i] = 255
	}

	// Try encode the message again.
	_, err = msg.Encode(buf)
	if err != nil {
		// ignore protocol level errors
	}

	// Finally try to decode again.
	_, err = msg.Decode(buf)
	if err != nil {
		return 0
	}

	// Everything was ok!
	return 1
}
