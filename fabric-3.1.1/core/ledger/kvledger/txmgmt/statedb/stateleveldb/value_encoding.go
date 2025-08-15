/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package stateleveldb

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/ledger/internal/version"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv"
)

// encodeValue encodes the value, version, and metadata
func encodeValue(v *statedb.VersionedValue, ns, key string) ([]byte, error) {
	// EncryptValueMKV sẽ tự động đọc K1 từ file
	encryptedValue := mkv.EncryptValueMKV(v.Value)
	encryptedMetadata := mkv.EncryptValueMKV(v.Metadata)
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

	// DecryptValueMKV sẽ tự động đọc K1 từ file
	val := mkv.DecryptValueMKV(dbValue.Value)
	metadata := mkv.DecryptValueMKV(dbValue.Metadata)
	// protobuf always makes an empty byte array as nil
	if val == nil {
		val = []byte{}
	}
	return &statedb.VersionedValue{Version: ver, Value: val, Metadata: metadata}, nil
}
