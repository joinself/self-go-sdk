// Copyright 2020 Self Group Ltd. All Rights Reserved.

package siggraph

import (
	"time"

	"golang.org/x/crypto/ed25519"
)

// Node node
type Node struct {
	kid      string            // id of the key
	did      string            // id of the device
	typ      string            // type of key
	seq      int               // the sequence of the operatin the key was added
	ca       int64             // when the key was created at/valid from
	ra       int64             // when the key was revoked at
	pk       ed25519.PublicKey // the public key
	incoming []*Node           // all keys that signed this key when it was created
	outgoing []*Node           // all keys that this key has signed
}

// TODO : this might exceed stack size on large graphs, we should use a slice as a stack to avoid this issue

// collect collects all descendents of this node
func (n *Node) collect() []*Node {
	var nodes []*Node

	nodes = append(nodes, n.outgoing...)

	for _, c := range n.outgoing {
		nodes = append(nodes, c.collect()...)
	}

	return nodes
}

func (n *Node) createdAt() time.Time {
	if n.ca == 0 {
		return time.Time{}
	}

	// is this greater than the maximum unix timestamp (seconds)?
	if n.ca > 1<<32-1 {
		return time.UnixMilli(n.ca)
	}

	return time.Unix(n.ca, 0)
}

func (n *Node) revokedAt() time.Time {
	if n.ra == 0 {
		return time.Time{}
	}

	if n.ra > 1<<32-1 {
		return time.UnixMilli(n.ra)
	}

	return time.Unix(n.ra, 0)
}
