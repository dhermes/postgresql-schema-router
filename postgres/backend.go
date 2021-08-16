// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import "fmt"

// DescribeBackendMessage tries to determine the **type** of backend message
// based on the first few bytes of the TCP chunk.
func DescribeBackendMessage(chunk []byte) (string, error) {
	if len(chunk) < 1 {
		err := fmt.Errorf(
			"%w; message must contain at least 1 byte, has %d",
			ErrParsingServerMessage, len(chunk),
		)
		return "", err
	}

	messageType := chunk[0]
	switch messageType {
	case 'R':
		// - AuthenticationOk
		// - AuthenticationKerberosV5
		// - AuthenticationCleartextPassword
		// - AuthenticationMD5PasswordAuthenticationSCMCredential
		// - AuthenticationGSS
		// - AuthenticationSSPI
		// - AuthenticationGSSContinue
		// - AuthenticationSASL
		// - AuthenticationSASLContinue
		// - AuthenticationSASLFinal
		return "Authentication{*}", nil
	case 'K':
		return "BackendKeyData", nil
	case '2':
		return "BindComplete", nil
	case '3':
		return "CloseComplete", nil
	case 'C':
		return "CommandComplete", nil
	case 'd':
		// NOTE: This message type is both F&B
		return "CopyData", nil
	case 'c':
		// NOTE: This message type is both F&B
		return "CopyDone", nil
	case 'G':
		return "CopyInResponse", nil
	case 'H':
		return "CopyOutResponse", nil
	case 'W':
		return "CopyBothResponse", nil
	case 'D':
		return "DataRow", nil
	case 'I':
		return "EmptyQueryResponse", nil
	case 'E':
		return "ErrorResponse", nil
	case 'V':
		return "FunctionCallResponse", nil
	case 'v':
		return "NegotiateProtocolVersion", nil
	case 'n':
		return "NoData", nil
	case 'N':
		return "NoticeResponse", nil
	case 'A':
		return "NotificationResponse", nil
	case 't':
		return "ParameterDescription", nil
	case 'S':
		return "ParameterStatus", nil
	case '1':
		return "ParseComplete", nil
	case 's':
		return "PortalSuspended", nil
	case 'Z':
		return "ReadyForQuery", nil
	case 'T':
		return "RowDescription", nil

	}

	err := fmt.Errorf(
		"%w; unexpected message type %x",
		ErrParsingServerMessage, messageType,
	)
	return "", err
}
