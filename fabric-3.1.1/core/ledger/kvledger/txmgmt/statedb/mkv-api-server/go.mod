module mkv-api-server

go 1.24.4

replace github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv => ../mkv

require github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/statedb/mkv v0.0.0-00010101000000-000000000000
