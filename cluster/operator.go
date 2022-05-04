// Copyright © 2022 Obol Labs Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of  MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along with
// this program.  If not, see <http://www.gnu.org/licenses/>.

package cluster

import (
	ssz "github.com/ferranbt/fastssz"

	"github.com/obolnetwork/charon/app/errors"
)

// Operator identifies a charon node and its operator.
type Operator struct {
	// Address is the Ethereum address identifying the operator.
	Address string `json:"address"`

	// ENR identifies the charon node.
	ENR string `json:"enr"`

	// Nonce is incremented each time the ENR is signed.
	Nonce int `json:"nonce"`

	// ENRSignature is a EIP712 signature of the ENR by the Address, authorising the charon node to act on behalf of the operator in the cluster.
	ENRSignature []byte `json:"enr_signature"`
}

// VerifySignature returns an error if the ENR signature doesn't match the address and enr fields.
func (o Operator) VerifySignature() error {
	digest, err := digestEIP712(o.Address, []byte(o.ENR), o.Nonce)
	if err != nil {
		return err
	}

	if ok, err := verifySig(o.Address, digest[:], o.ENRSignature); err != nil {
		return err
	} else if !ok {
		return errors.New("invalid operator enr signature")
	}

	return nil
}

// HashTreeRoot ssz hashes the Definition object.
func (o Operator) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(o) //nolint:wrapcheck
}

// HashTreeRootWith ssz hashes the Operator object with a hasher.
func (o Operator) HashTreeRootWith(hh *ssz.Hasher) error {
	indx := hh.Index()

	// Field (0) 'Address'
	hh.PutBytes([]byte(o.Address))

	// Field (1) 'ENR'
	hh.PutBytes([]byte(o.ENR))

	// Field (2) 'Nonce'
	hh.PutUint64(uint64(o.Nonce))

	// Field (3) 'ENRSignature'
	hh.PutBytes(o.ENRSignature)

	hh.Merkleize(indx)

	return nil
}