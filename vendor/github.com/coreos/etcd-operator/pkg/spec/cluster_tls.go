// Copyright 2017 The etcd-operator Authors
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

package spec

import "errors"

// TLSPolicy defines the TLS policy of an etcd cluster
type TLSPolicy struct {
	// StaticTLS enables user to generate static x509 certificates and keys,
	// put them into Kubernetes secrets, and specify them into here.
	Static *StaticTLS `json:"static"`
}

type StaticTLS struct {
	// Member contains secrets containing TLS certs used by each etcd member pod.
	Member *MemberSecret `json:"member"`
	// OperatorSecret is the secret containing TLS certs used by operator to
	// talk securely to this cluster.
	OperatorSecret string `json:"operatorSecret"`
}

type MemberSecret struct {
	// PeerSecret is the secret containing TLS certs used by each etcd member pod
	// for the communication between etcd peers.
	PeerSecret string `json:"peerSecret"`
	// ClientSecret is the secret containing TLS certs used by each etcd member pod
	// for the communication between etcd and its clients.
	ClientSecret string `json:"clientSecret"`
}

func (tp *TLSPolicy) Validate() error {
	if tp.Static == nil {
		return nil
	}
	st := tp.Static

	if len(st.OperatorSecret) != 0 {
		if len(st.Member.ClientSecret) == 0 {
			return errors.New("operator secret set but member clientSecret not set")
		}
	} else if st.Member != nil && len(st.Member.ClientSecret) != 0 {
		return errors.New("member clientSecret set but operator secret not set")
	}
	return nil
}

func (tp *TLSPolicy) IsSecureClient() bool {
	if tp == nil || tp.Static == nil {
		return false
	}
	return len(tp.Static.OperatorSecret) != 0
}

func (tp *TLSPolicy) IsSecurePeer() bool {
	if tp == nil || tp.Static == nil || tp.Static.Member == nil {
		return false
	}
	return len(tp.Static.Member.PeerSecret) != 0
}
