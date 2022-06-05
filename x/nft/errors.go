package nft

import (
	errorsmod "cosmossdk.io/errors"
)

// x/nft module sentinel errors
var (
	ErrInvalidNFT     = errorsmod.Register(ModuleName, 2, "invalid nft")
	ErrClassExists    = errorsmod.Register(ModuleName, 3, "nft class already exist")
	ErrClassNotExists = errorsmod.Register(ModuleName, 4, "nft class does not exist")
	ErrNFTExists      = errorsmod.Register(ModuleName, 5, "nft already exist")
	ErrNFTNotExists   = errorsmod.Register(ModuleName, 6, "nft does not exist")
	ErrInvalidID      = errorsmod.Register(ModuleName, 7, "invalid id")
	ErrInvalidClassID = errorsmod.Register(ModuleName, 8, "invalid class id")
)
