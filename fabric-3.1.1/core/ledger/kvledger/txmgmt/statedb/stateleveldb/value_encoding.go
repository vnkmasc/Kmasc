/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package stateleveldb

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/ledger/internal/version"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
)

// encodeValue encodes the value, version, and metadata
func encodeValue(v *statedb.VersionedValue, ns, key string) ([]byte, error) {
	// Mã hóa value và metadata trước khi lưu
	encryptedValue := statedb.EncryptValue(v.Value, ns, key)
	encryptedMetadata := statedb.EncryptValue(v.Metadata, ns, key)
	return proto.Marshal(
		&DBValue{
			Version:  v.Version.ToBytes(),
			Value:    encryptedValue,
			Metadata: encryptedMetadata,
		},
	)
}

// decodeValue decodes the statedb value bytes
func decodeValue(encodedValue []byte, ns, key string) (*statedb.VersionedValue, error) {
	dbValue := &DBValue{}
	err := proto.Unmarshal(encodedValue, dbValue)
	if err != nil {
		return nil, err
	}
	ver, _, err := version.NewHeightFromBytes(dbValue.Version)
	if err != nil {
		return nil, err
	}
	val := statedb.DecryptValue(dbValue.Value, ns, key)
	metadata := statedb.DecryptValue(dbValue.Metadata, ns, key)
	// protobuf always makes an empty byte array as nil
	if val == nil {
		val = []byte{}
	}
	return &statedb.VersionedValue{Version: ver, Value: val, Metadata: metadata}, nil
}
