package srp

import (
	"errors"
)

// Handshake -
//
// Params:
//  A ([]byte) - a client's generated public key
//  v ([]byte) - a client's stored verifer
//
// Return:
//  []byte - the generated public key "B", to be sent to the client
//  []byte - the computed pre-session key "S", to be kept secret
//  []byte - the computed session key "K", to be kept secret
//  error
//
//  NOTE: Only the returned "B" value should be sent to the client. "S" and "K"
//		  are very secret.
//
func (g *group)Handshake(A, v []byte) ([]byte, []byte, []byte, error) {
	// "A" cannot be zero
	if isZero(A) {
		return nil, nil, nil, errors.New("Server found \"A\" to be zero. Aborting handshake")
	}

	// Create a random secret "b"
	b, err := randomBytes(32)
	if err != nil {
		return nil, nil, nil, err
	}

	// Calculate the SRP-6a version of the multiplier parameter "k"
	k := Hash(g.pad(g.N), g.pad(g.g))

	// Compute a value "B" based on "b"
	//   B = (v + g^b) % N
	B := g.add(g.mul(k, v), g.exp(g.g, b))

	// Calculate "u"
	u := Hash(g.pad(A), g.pad(B))

	// Compute the pseudo-session key, "S"
	//  S = (Av^u) ^ b
	S := g.exp(g.mul(A, g.exp(v, u)), b)

	// The actual session key is the hash of the pseudo-session key "S"
	K := Hash(S)

	return B, S, K, nil
}

// ServerProof -
//
// Params:
//  A ([]byte) - the client's session public key
//  ClientProof ([]byte) - the client's proof as computed with ClientProof()
//  K ([]byte) - the computed session secret
//
// Return:
//  []byte - the client's proof of knowing K
//  error
//
func ServerProof(A, ClientProof, K []byte) []byte {
	return Hash(A, ClientProof, K)
}
