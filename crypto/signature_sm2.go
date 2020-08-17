// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// +build !nacl,!js,cgo,sm2

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"

	"github.com/tjfoc/gmsm/sm2"
)

func EcrecoverWithPub(hash, sig []byte) ([]byte, error) {
	if len(sig) != SignatureLength {
		// panic("size error , wang 98")
		return nil, fmt.Errorf("error ")
	}

	pub := sm2.Decompress(sig[65:])
	if pub.X == nil {
		// panic("DecompressPubkey error")
		return nil, fmt.Errorf("DecompressPubkey error")
	}
	// 使用标准压格式化算法
	// return FromECDSAPub(pub), nil
	return elliptic.Marshal(S256(), pub.X, pub.Y), nil
	// return sm2.MarshalSm2PublicKey(pub)
}

func SigToPubWithPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	if len(sig) != 98 {
		// panic("got error sig size ")
		return nil, fmt.Errorf("SigToPubWithPub: sig size error,want 98  but got %d", len(sig))

	}
	pubkey := sm2.Decompress(sig[65:])
	if pubkey.X == nil {
		// panic("DecompressPubkey error")
		return nil, fmt.Errorf("DecompressPubkey error")
	}

	key := &ecdsa.PublicKey{
		Curve: pubkey.Curve,
		X:     pubkey.X,
		Y:     pubkey.Y,
	}
	return key, nil
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given digest cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
// func Sign(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
// 	if len(digestHash) != DigestLength {
// 		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", DigestLength, len(digestHash))
// 	}
// 	seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
// 	defer zeroBytes(seckey)
// 	return secp256k1.Sign(digestHash, seckey)
// }

// The produced signature is in the [R || S || V] format where V is 0 or 1.
func signWithoutPub(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	if len(digestHash) != DigestLength {
		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", DigestLength, len(digestHash))
	}

	sm2Priv := &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: S256(),
			X:     prv.PublicKey.X,
			Y:     prv.PublicKey.Y,
		},
		D: prv.D,
	}

	var isok bool
	var r, s *big.Int
	for i := 0; i < 5; i++ {
		r, s, err = sm2.Sign(sm2Priv, digestHash)
		if err != nil {
			return nil, err
		}
		isok = sm2.Verify(&sm2Priv.PublicKey, digestHash, r, s)
		if isok {
			break
		}
	}
	if !isok {
		return nil, fmt.Errorf("sign failed")
	}

	sig = make([]byte, 65)
	copy(sig[:32], r.Bytes())
	copy(sig[32:64], s.Bytes())
	// sig[64] = 0 // invalid V
	return sig, nil
}

func SignWithPub(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	sig, err = signWithoutPub(digestHash, prv)
	if err != nil {
		// panic("sign failed: " + err.Error())
		return nil, err
	}

	pubKey := &sm2.PublicKey{
		Curve: S256(),
		X:     prv.PublicKey.X,
		Y:     prv.PublicKey.Y,
	}
	keyBuf := sm2.Compress(pubKey)
	sig = append(sig, keyBuf...)
	return sig, nil
}

// VerifySignature checks that the given public key created signature over digest.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should have the 64 byte [R || S] format.
func VerifySignature(pubkey, digestHash, signature []byte) bool {
	var sm2Pub *sm2.PublicKey
	if len(pubkey) == 33 {
		sm2Pub = sm2.Decompress(pubkey)
	} else {
		// UnmarshalPubkey(pubkey)
		x, y := elliptic.Unmarshal(S256(), pubkey)
		if x == nil {
			// panic("parse decompress pubkey failed:")
			return false
		}
		sm2Pub = &sm2.PublicKey{
			X:     x,
			Y:     y,
			Curve: S256(),
		}
	}

	if sm2Pub == nil || sm2Pub.X == nil {
		// panic("parse pubkey  failed: is nil ")
		return false
	}

	if len(signature) != 64 {
		// fmt.Println(len(signature))
		// panic("parse sig  failed: size error")
		return false
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	ok := sm2.Verify(sm2Pub, digestHash, r, s)
	if !ok {
		return false
		// fmt.Printf("%X\n", digestHash)
		// panic("verify signatrue failed")
	}
	return ok

}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	key, err := decompressPubkey(pubkey)
	// fmt.Println(err)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	// fmt.Println("y", key.Y)
	return key, nil
}

func decompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {

	defer func() {
		recover()
		// if data != nil {

		// }
	}()

	sm2Pub := sm2.Decompress(pubkey)
	if sm2Pub == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	if sm2Pub.X == nil {
		// panic("parse public err")
		return nil, fmt.Errorf("invalid public key")
	}

	return &ecdsa.PublicKey{X: sm2Pub.X, Y: sm2Pub.Y, Curve: S256()}, nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	sm2Pub := &sm2.PublicKey{
		Curve: S256(),
		X:     pubkey.X,
		Y:     pubkey.Y,
	}
	return sm2.Compress(sm2Pub)
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return sm2.P256Sm2()
}
