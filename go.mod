module gitlab.com/thorchain/tss/go-tss

go 1.14

require (
	github.com/binance-chain/tss-lib v0.0.0-20201118045712-70b2cb4bf916
	github.com/blang/semver v3.5.1+incompatible
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/cosmos/cosmos-sdk v0.41.0
	github.com/deckarep/golang-set v1.7.1
	github.com/decred/dcrd/dcrec/secp256k1 v1.0.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/ipfs/go-log v1.0.5
	github.com/libp2p/go-libp2p v0.18.1
	github.com/libp2p/go-libp2p-core v0.14.0
	github.com/libp2p/go-libp2p-discovery v0.5.0
	github.com/libp2p/go-libp2p-kad-dht v0.10.0
	github.com/libp2p/go-libp2p-peerstore v0.6.0
	github.com/libp2p/go-libp2p-testing v0.9.2
	github.com/magiconair/properties v1.8.4
	github.com/multiformats/go-multiaddr v0.5.0
	github.com/olekukonko/tablewriter v0.0.0-20170122224234-a0225b3f23b5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/client_model v0.2.0
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/tendermint v0.34.3
	gitlab.com/thorchain/binance-sdk v1.2.3-0.20210117202539-d569b6b9ba5d
	go.uber.org/atomic v1.9.0
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e
	golang.org/x/text v0.3.7
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)

replace (
	github.com/binance-chain/tss-lib => gitlab.com/thorchain/tss/tss-lib v0.0.0-20201118045712-70b2cb4bf916
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
