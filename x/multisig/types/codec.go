package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAccount{}, &MsgCreateProposal{}, &MsgUpdateParams{}, &MsgVote{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateAccount{}, "atomone/multisig/v1/MsgCreateAccount")
	legacy.RegisterAminoMsg(cdc, &MsgCreateProposal{}, "atomone/multisig/v1/MsgCreateProposal")
	legacy.RegisterAminoMsg(cdc, &MsgVote{}, "atomone/multisig/v1/MsgVote")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "atomone/x/multisig/v1/MsgUpdateParams")
	cdc.RegisterConcrete(&Params{}, "atomone/multisig/v1/Params", nil)
}

/* TODO REMOVE
var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Need to add registration in the atomone multisig amino for all modules
	// because of the MsgCreateProposal which can embed any other messages.
	banktypes.RegisterLegacyAminoCodec(amino)
	govtypes.RegisterLegacyAminoCodec(amino)
	consensustypes.RegisterLegacyAminoCodec(amino)
	crisistypes.RegisterLegacyAminoCodec(amino)
	distributiontypes.RegisterLegacyAminoCodec(amino)
	evidencetypes.RegisterLegacyAminoCodec(amino)
	minttypes.RegisterLegacyAminoCodec(amino)
	slashingtypes.RegisterLegacyAminoCodec(amino)
	stakingtypes.RegisterLegacyAminoCodec(amino)
	upgradetypes.RegisterLegacyAminoCodec(amino)
}
*/
