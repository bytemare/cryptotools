// Package ristretto allows simple and abstracted operations in the Ristretto255 group
package ristretto

import (
	"github.com/gtank/ristretto255"

	"github.com/bytemare/cryptotools/group"
	"github.com/bytemare/cryptotools/group/ciphersuite/internal/ristretto/h2r"
	"github.com/bytemare/cryptotools/hash"
)

const ristrettoInputLength = 64

// Ristretto implements the group interface.
type Ristretto struct {
	h2r *h2r.HashToRistretto
}

// New returns a structure interfacing operations in the Ristretto255 group
// The protocol and version arguments are used for the HashToGroup domain separation,
// and should represent values from the calling protocol.
func New(hashID hash.Identifier, dst []byte) *Ristretto {
	return &Ristretto{h2r.New(dst, hashID)}
}

// NewScalar returns a new, empty, scalar.
func (r *Ristretto) NewScalar() group.Scalar {
	return &Scalar{Scalar: ristretto255.NewScalar()}
}

// ElementLength returns the byte size of an encoded element.
func (r *Ristretto) ElementLength() int {
	return canonicalEncodingLength
}

// NewElement returns a new, empty, element.
func (r *Ristretto) NewElement() group.Element {
	return &Element{HashToRistretto: r.h2r, element: ristretto255.NewElement()}
}

// Identity returns the group's identity element.
func (r *Ristretto) Identity() group.Element {
	return &Element{
		HashToRistretto: r.h2r,
		element:         ristretto255.NewElement().Zero(),
	}
}

// HashToGroup allows arbitrary input to be safely mapped to the curve of the group.
func (r *Ristretto) HashToGroup(input, dst []byte) group.Element {
	var h *h2r.HashToRistretto
	if dst == nil {
		h = r.h2r
	} else {
		h = h2r.New(dst, r.h2r.Identifier())
	}

	uniform := h.Expand(input, ristrettoInputLength)

	return &Element{
		HashToRistretto: h,
		element:         ristretto255.NewElement().FromUniformBytes(uniform),
	}
}

// HashToScalar allows arbitrary input to be safely mapped to the field.
func (r *Ristretto) HashToScalar(input, dst []byte) group.Scalar {
	var h *h2r.HashToRistretto
	if dst == nil {
		h = r.h2r
	} else {
		h = h2r.New(dst, r.h2r.Identifier())
	}

	hashed := h.Expand(input, ristrettoInputLength)
	return &Scalar{Scalar: ristretto255.NewScalar().FromUniformBytes(hashed)}
}

// Base returns Ristretto255's base point a.k.a. canonical generator.
func (r *Ristretto) Base() group.Element {
	return &Element{r.h2r, ristretto255.NewElement().Base()}
}

// MultBytes allows []byte encodings of a scalar and an element of the group to be multiplied.
func (r *Ristretto) MultBytes(s, e []byte) (group.Element, error) {
	sc, err := r.NewScalar().Decode(s)
	if err != nil {
		return nil, err
	}

	el, err := r.NewElement().Decode(e)
	if err != nil {
		return nil, err
	}

	return el.Mult(sc), nil
}

// DST returns the domain separation tag the group has been instantiated with.
func (r *Ristretto) DST() string {
	return r.h2r.GetOriginalDST()
}
