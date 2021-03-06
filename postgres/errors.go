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

import (
	"errors"
)

var (
	// ErrNotImplemented indicates a method or message type is not implemented.
	ErrNotImplemented = errors.New("not implemented")
	// ErrParsingClientMessage indicates a failure occurred when parsing a TCP
	// packet as a PostgreSQL client message.
	ErrParsingClientMessage = errors.New("failed to parse TCP packet as PostgesSQL client message")
	// ErrParsingServerMessage indicates a failure occurred when parsing a TCP
	// packet as a PostgreSQL server message.
	ErrParsingServerMessage = errors.New("failed to parse TCP packet as PostgesSQL server message")
)
