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

package message

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQosCodes(t *testing.T) {
	if QosAtMostOnce != 0 || QosAtLeastOnce != 1 || QosExactlyOnce != 2 {
		t.Errorf("QOS codes invalid")
	}
}

func TestConnackReturnCodes(t *testing.T) {
	require.Equal(t, ErrInvalidProtocolVersion.Error(), ConnackCode(1).Error(), "Incorrect ConnackCode error value.")
	require.Equal(t, ErrIdentifierRejected.Error(), ConnackCode(2).Error(), "Incorrect ConnackCode error value.")
	require.Equal(t, ErrServerUnavailable.Error(), ConnackCode(3).Error(), "Incorrect ConnackCode error value.")
	require.Equal(t, ErrBadUsernameOrPassword.Error(), ConnackCode(4).Error(), "Incorrect ConnackCode error value.")
	require.Equal(t, ErrNotAuthorized.Error(), ConnackCode(5).Error(), "Incorrect ConnackCode error value.")
}

func TestFixedHeaderFlags(t *testing.T) {
	type detail struct {
		name  string
		flags byte
	}

	details := map[MessageType]detail{
		RESERVED:    detail{"RESERVED", 0},
		CONNECT:     detail{"CONNECT", 0},
		CONNACK:     detail{"CONNACK", 0},
		PUBLISH:     detail{"PUBLISH", 0},
		PUBACK:      detail{"PUBACK", 0},
		PUBREC:      detail{"PUBREC", 0},
		PUBREL:      detail{"PUBREL", 2},
		PUBCOMP:     detail{"PUBCOMP", 0},
		SUBSCRIBE:   detail{"SUBSCRIBE", 2},
		SUBACK:      detail{"SUBACK", 0},
		UNSUBSCRIBE: detail{"UNSUBSCRIBE", 2},
		UNSUBACK:    detail{"UNSUBACK", 0},
		PINGREQ:     detail{"PINGREQ", 0},
		PINGRESP:    detail{"PINGRESP", 0},
		DISCONNECT:  detail{"DISCONNECT", 0},
		RESERVED2:   detail{"RESERVED2", 0},
	}

	for m, d := range details {
		if m.String() != d.name {
			t.Errorf("Name mismatch. Expecting %s, got %s.", d.name, m)
		}

		if m.defaultFlags() != d.flags {
			t.Errorf("Flag mismatch for %s. Expecting %d, got %d.", m, d.flags, m.defaultFlags())
		}
	}
}

func TestReadmeExample(t *testing.T) {
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

	// Detect message.
	l, mt := DetectMessage(buf)

	// Check length
	if l == 0 {
		// buffer not complete yet
		return
	}

	// Create message.
	msg2, err := mt.New()
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
}
