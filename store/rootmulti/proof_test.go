package rootmulti

import (
	"encoding/hex"
	"fmt"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	ics23 "github.com/cosmos/ics23/go"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"
	"cosmossdk.io/store/iavl"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/types"
)

func TestVerifyIAVLStoreQueryProof(t *testing.T) {
	// Create main tree for testing.
	db := dbm.NewMemDB()
	iStore, err := iavl.LoadStore(db, log.NewNopLogger(), types.NewKVStoreKey("test"), types.CommitID{}, iavl.DefaultIAVLCacheSize, false, metrics.NewNoOpMetrics())
	store := iStore.(*iavl.Store)
	require.Nil(t, err)
	store.Set([]byte("MYKEY"), []byte("MYVALUE"))
	cid := store.Commit()

	// Get Proof
	res, err := store.Query(&types.RequestQuery{
		Path:  "/key", // required path to get key/value+proof
		Data:  []byte("MYKEY"),
		Prove: true,
	})
	require.NoError(t, err)
	require.NotNil(t, res.ProofOps)

	// Verify proof.
	prt := DefaultProofRuntime()
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/MYKEY", []byte("MYVALUE"))
	require.Nil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/MYKEY_NOT", []byte("MYVALUE"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/MYKEY/MYKEY", []byte("MYVALUE"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "MYKEY", []byte("MYVALUE"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/MYKEY", []byte("MYVALUE_NOT"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/MYKEY", []byte(nil))
	require.NotNil(t, err)
}

func TestVerifyMultiStoreQueryProof(t *testing.T) {
	// Create main tree for testing.
	db := dbm.NewMemDB()
	store := NewStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	iavlStoreKey := types.NewKVStoreKey("iavlStoreKey")

	store.MountStoreWithDB(iavlStoreKey, types.StoreTypeIAVL, nil)
	require.NoError(t, store.LoadVersion(0))

	iavlStore := store.GetCommitStore(iavlStoreKey).(*iavl.Store)

	for _, ikey := range []byte{0x11, 0x32, 0x50, 0x72, 0x99} {
		key := []byte{ikey}
		iavlStore.Set(key, key)
	}
	// Enter the key we want to proof
	key := []byte("prefix207-tendermint-42\x03\x00\x00\x00\x00\x00\x00\x00\x01")
	value, _ := hex.DecodeString("cf49bb81a77249af41ecbe7792d98ddf24b47b491a177ca5a8b1e82e2eaf011e")
	iavlStore.Set(key, value)

	cid := store.Commit()

	// Get Proof
	res, err := store.Query(&types.RequestQuery{
		Path:  "/iavlStoreKey/key",
		Data:  key,
		Prove: true,
	})
	require.NoError(t, err)
	require.NotNil(t, res.ProofOps)

	// Decode ics23 proof
	proofs := make([]*ics23.CommitmentProof, len(res.ProofOps.Ops))
	// spew.Dump(reqres.Response.ProofOps.Ops)
	for i, op := range res.ProofOps.Ops {
		var p ics23.CommitmentProof
		err = p.Unmarshal(op.Data)
		if err != nil || p.Proof == nil {
			panic(fmt.Sprintf("could not unmarshal proof op into CommitmentProof at index %d: %v", i, err))
		}
		proofs[i] = &p
	}
	spew.Config.DisableMethods = true
	// spew.Dump(proofs)
	fmt.Printf("ROOT %x\n", cid.Hash)
	for i, p := range proofs {
		fmt.Printf("KEY %x %q\n", p.GetExist().Key, p.GetExist().Key)
		fmt.Printf("VAL %x\n", p.GetExist().Value)
		fmt.Printf("%d LEAF PREFIX %x\n", i, p.GetExist().Leaf.Prefix)
		for j, op := range p.GetExist().Path {
			fmt.Printf("\t%d OP PREFIX %x\n", j, op.Prefix)
			fmt.Printf("\t%d OP SUFFIX %x\n", j, op.Suffix)
		}
	}

	// Verify proof.
	prt := DefaultProofRuntime()
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/iavlStoreKey/"+string(key), value)
	require.Nil(t, err)
	return

	// Verify proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/iavlStoreKey/MYKEY", []byte("MYVALUE"))
	require.Nil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/iavlStoreKey/MYKEY_NOT", []byte("MYVALUE"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/iavlStoreKey/MYKEY/MYKEY", []byte("MYVALUE"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "iavlStoreKey/MYKEY", []byte("MYVALUE"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/MYKEY", []byte("MYVALUE"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/iavlStoreKey/MYKEY", []byte("MYVALUE_NOT"))
	require.NotNil(t, err)

	// Verify (bad) proof.
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/iavlStoreKey/MYKEY", []byte(nil))
	require.NotNil(t, err)
}

func TestVerifyMultiStoreQueryProofAbsence(t *testing.T) {
	// Create main tree for testing.
	db := dbm.NewMemDB()
	store := NewStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	iavlStoreKey := types.NewKVStoreKey("iavlStoreKey")

	store.MountStoreWithDB(iavlStoreKey, types.StoreTypeIAVL, nil)
	err := store.LoadVersion(0)
	require.NoError(t, err)

	iavlStore := store.GetCommitStore(iavlStoreKey).(*iavl.Store)
	iavlStore.Set([]byte("MYKEY"), []byte("MYVALUE"))
	cid := store.Commit() // Commit with empty iavl store.

	// Get Proof
	res, err := store.Query(&types.RequestQuery{
		Path:  "/iavlStoreKey/key", // required path to get key/value+proof
		Data:  []byte("MYABSENTKEY"),
		Prove: true,
	})
	require.NoError(t, err)
	require.NotNil(t, res.ProofOps)

	// Verify proof.
	prt := DefaultProofRuntime()
	err = prt.VerifyAbsence(res.ProofOps, cid.Hash, "/iavlStoreKey/MYABSENTKEY")
	require.Nil(t, err)

	// Verify (bad) proof.
	prt = DefaultProofRuntime()
	err = prt.VerifyAbsence(res.ProofOps, cid.Hash, "/MYABSENTKEY")
	require.NotNil(t, err)

	// Verify (bad) proof.
	prt = DefaultProofRuntime()
	err = prt.VerifyValue(res.ProofOps, cid.Hash, "/iavlStoreKey/MYABSENTKEY", []byte(""))
	require.NotNil(t, err)
}
