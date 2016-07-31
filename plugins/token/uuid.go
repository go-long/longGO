package token

// This is a simple package for go that generates randomized base62 encoded uuids based on a single integer.
// It's ideal for shorturl services or for semi-secured randomized api primary keys.
//
// How it Works
//
// `UUID` is an alias for `uint64`.
// Its `UUID.Encode()` method interface returns a `Base62` encoded string based off of the number.
// Its implementation of the `json.Marshaler` interface encodes and decoded the `UUID` to and from the same
// `Base62` encoded string representation.
//
// Basically, the outside world will always address the uuid as its string equivolent and internally we can
// always be used as an `uint64` for fast, indexed, unique, lookups in various databases.
//
// **IMPORTANT:** Remember to always check for collisions when adding randomized uuids to a database

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
// Base62 is a string respresentation of every possible base62 character
	Base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// MaxUUIDLength is the largest possible character length of a uuid
	MaxUUIDLength = 10

// MinUUIDLength is the smallest possible character length of a uuid
	MinUUIDLength = 1

// DefaultUUIDLength is the default size of a uuid
	DefaultUUIDLength = MaxUUIDLength
)

var (
	base62Len = uint64(len(Base62))
)

// UUID is an alias of an int64 that is json marshalled into a base62 encoded uuid
type UUID uint64


// Encode encodes the uuid into a base62 string
func (t UUID) Encode() string {
	number := uint64(t)
	if number == 0 {
		return ""
	}

	var chars []byte
	for number > 0 {
		result := number / base62Len
		remainder := number % base62Len
		chars = append(chars, Base62[remainder])
		number = result
	}

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}

// UnmarshalJSON implements the `json.Marsheler` interface to decode the uuid from a base62 string back into an int64
func (t *UUID) UnmarshalJSON(data []byte) error {
	str := string(data)
	strLen := len(data)

	// return an error if the json was not a valid string representation of the uuid
	isString := (str[0] == '"' || str[0] == '\'') && (str[strLen-1] == '"' || str[strLen-1] == '\'')
	if !isString {
		return fmt.Errorf("attempted to parse a non-string uuid")
	}

	// decode the uuid
	decoded, err := Decode(str[1 : strLen-1])
	if err != nil {
		return err
	}

	// the uuid was successfully decoded
	*t = decoded
	return nil
}

// MarshalJSON implements the `json.Marsheler` interface to encode the uuid into a base62 string
func (t UUID) MarshalJSON() ([]byte, error) {
	uuid := t.Encode()
	if uuid == "" {
		return []byte{}, nil
	}
	return []byte("\"" + t.Encode() + "\""), nil
}

// New returns a `Base62` encoded `UUID` of *up to* `DefaultUUIDLength`
// if you pass in a `uuidLength` between `MinUUIDLength` and `MaxUUIDLength` this will return
// a `UUID` of *up to* that length instead if you pass in a `uuidLength` that is out of range it will panic
func NewUUID(uuidLength ...int) UUID {

	// calculate the max hash int based on the uuid length
	var max uint64
	if uuidLength != nil {
		isInRange := uuidLength[0] >= MinUUIDLength && uuidLength[0] <= MaxUUIDLength
		if isInRange {
			max = maxHashInt(uuidLength[0])
		} else {
			panic(fmt.Errorf("uuidLength ∉ [%d,%d]", MinUUIDLength, MaxUUIDLength))
			return UUID(0)
		}
	} else {
		max = maxHashInt(DefaultUUIDLength)
	}

	// generate a psuedo random uuid
	rand.Seed(time.Now().UTC().UnixNano())
	number := uint64(rand.Int63n(int64(max & math.MaxInt64)))

	return UUID(number)
}

// Decode returns a uuid from a 1-12 character base62 encoded string
func Decode(uuid string) (UUID, error) {

	number := uint64(0)
	idx := 0.0
	chars := []byte(Base62)

	charsLength := float64(len(chars))
	uuidLength := float64(len(uuid))

	if uuidLength > MaxUUIDLength {
		return UUID(0), fmt.Errorf("%d > MaxUUIDLength (%d)", int(uuidLength), MaxUUIDLength)
	} else if uuidLength < MinUUIDLength {
		return UUID(0), fmt.Errorf("%d < MinUUIDLength (%d)", int(uuidLength), MinUUIDLength)
	}

	for _, c := range []byte(uuid) {
		power := uuidLength - (idx + 1)
		index := bytes.IndexByte(chars, c)
		if index < 0 {
			return UUID(0), fmt.Errorf("%q is not present in %s", c, Base62)
		}
		number += uint64(index) * uint64(math.Pow(charsLength, power))
		idx++
	}

	return UUID(number), nil
}

// maxHashInt returns the largest possible int that will yeild a base62 encoded uuid of the specified length
func maxHashInt(length int) uint64 {
	return uint64(math.Max(0, math.Min(math.MaxUint64, math.Pow(float64(base62Len), float64(length)))))
}