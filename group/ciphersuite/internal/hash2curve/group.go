// Package hash2curve wraps an hash-to-curve implementation and exposes functions for operations on points and scalars.
package hash2curve

import (
	"github.com/armfazh/h2c-go-ref"
	"github.com/bytemare/cryptotools/group"
)

// Hash2Curve implements the Group interface to Hash-to-Curve primitives.
type Hash2Curve struct {
	suite h2c.SuiteID
	h2c.HashToPoint
	*curve
	dst []byte
}

// New returns a pointer to a Hash2Curve structure instantiated for the given Hash-to-Curve identifier.
// and the domain separation tag.
func New(id h2c.SuiteID, dst []byte) *Hash2Curve {
	h, err := id.Get(dst)
	if err != nil {
		panic(err)
	}

	if !h.IsRandomOracle() {
		panic(errParamNotRandomOracle)
	}

	return &Hash2Curve{id, h, curves[id].New(h.GetCurve()), dst}
}

// NewScalar returns a new, empty, scalar.
func (h *Hash2Curve) NewScalar() group.Scalar {
	return scalar(h.GetHashToScalar().GetScalarField())
}

// ElementLength returns the byte size of an encoded element.
func (h *Hash2Curve) ElementLength() int {
	return pointLen(h.GetCurve().Field().BitLen())
}

// NewElement returns a new, empty, element.
func (h *Hash2Curve) NewElement() group.Element {
	return &Point{
		Hash2Curve: h,
		point:      h.base,
	}
}

// Identity returns the group's identity element.
func (h *Hash2Curve) Identity() group.Element {
	return &Point{
		Hash2Curve: h,
		point:      h.GetCurve().Identity(),
	}
}

// HashToGroup allows arbitrary input to be safely mapped to the curve of the Group.
func (h *Hash2Curve) HashToGroup(input, dst []byte) group.Element {
	var htp *Hash2Curve
	if dst == nil {
		htp = h
	} else {
		htp = New(h.suite, dst)
	}

	htp.checkDSTLen()

	return &Point{
		Hash2Curve: htp,
		point:      htp.Hash(input),
	}
}

// HashToScalar allows arbitrary input to be safely mapped to the field.
func (h *Hash2Curve) HashToScalar(input, dst []byte) group.Scalar {
	var htp *Hash2Curve
	if dst == nil {
		htp = h
	} else {
		htp = New(h.suite, dst)
	}

	htp.checkDSTLen()

	return &Scalar{
		s: htp.GetHashToScalar().Hash(input),
		f: htp.GetHashToScalar().GetScalarField(),
	}
}

// Base returns the group's base point a.k.a. canonical generator.
func (h *Hash2Curve) Base() group.Element {
	return h.NewElement()
}

// MultBytes allows []byte encodings of a scalar and an element of the Group to be multiplied.
func (h *Hash2Curve) MultBytes(scalar, element []byte) (group.Element, error) {
	s, err := h.NewScalar().Decode(scalar)
	if err != nil {
		return nil, err
	}

	e, err := h.NewElement().Decode(element)
	if err != nil {
		return nil, err
	}

	return e.Mult(s), nil
}

// DST returns the domain separation tag the group has been instantiated with.
func (h *Hash2Curve) DST() string {
	return string(h.dst)
}

func (h *Hash2Curve) checkDSTLen() {
	// todo bring this back after testing
	//if len(h.dst) < group.DstRecommendedMinLength {
	//	if len(h.dst) == group.DstMinLength {
	//		panic(errParamZeroLenDST)
	//	}
//
	//	panic(errParamShortDST)
	//}
}
