package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
)

var (
	_, _, _, _, _ sdk.Msg = &MsgUpdateParams{}, &MsgCreateAccount{}, &MsgCreateProposal{}, &MsgVote{}, &MsgExecuteProposal{}

	_ codectypes.UnpackInterfacesMessage = &MsgCreateProposal{}
)

// ValidateBasic implements the sdk.Msg interface.
func (m MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	return m.Params.ValidateBasic()
}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgCreateAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if len(m.Members) == 0 {
		return ErrMissingMembers
	}
	if m.Threshold <= 0 {
		return ErrZeroThreshold
	}
	membersMap := map[string]struct{}{} // to check for duplicates
	for i := range m.Members {
		if _, ok := membersMap[m.Members[i].Address]; ok {
			return ErrDuplicateMember
		}

		membersMap[m.Members[i].Address] = struct{}{}

		if m.Members[i].Weight == 0 {
			return ErrZeroMemberWeight
		}
		_, err := sdk.AccAddressFromBech32(m.Members[i].Address)
		if err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid member address: %s", err)
		}
	}
	return nil
}

// GetMsgs unpacks m.Messages Any's into sdk.Msg's
func (m *MsgCreateProposal) GetMsgs() ([]sdk.Msg, error) {
	return sdktx.GetMsgs(m.Messages, "sdk.MsgProposal")
}

// SetMsgs packs sdk.Msg's into m.Messages Any's
// NOTE: this will overwrite any existing messages
func (m *MsgCreateProposal) SetMsgs(msgs []sdk.Msg) error {
	anys, err := sdktx.SetMsgs(msgs)
	if err != nil {
		return err
	}

	m.Messages = anys
	return nil
}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgCreateProposal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid sender address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(m.AccountAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid multisig account address: %s", err)
	}
	// TODO assert max length
	if m.Title == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("proposal title cannot be empty") //nolint:staticcheck
	}
	// TODO assert max length
	if m.Summary == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("proposal summary cannot be empty") //nolint:staticcheck
	}
	if len(m.Messages) == 0 {
		// TODO allow no messages for text proposals?
		return sdkerrors.ErrInvalidRequest.Wrap("Proposal.Messages length must be non-nil") //nolint:staticcheck
	}
	msgs, err := m.GetMsgs()
	if err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("unable to read proposal messages: %v", err) //nolint:staticcheck
	}
	for i, msg := range msgs {
		// perform a basic validation of the message
		if m, ok := msg.(sdk.HasValidateBasic); ok {
			if err := m.ValidateBasic(); err != nil {
				return sdkerrors.ErrInvalidRequest.Wrapf("validation fail for proposal message %d: %v", i, err) //nolint:staticcheck
			}
		}
	}
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgCreateProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return sdktx.UnpackInterfaces(unpacker, m.Messages)
}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgVote) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Voter); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid voter address: %s", err)
	}
	if m.Vote == VoteOption_VOTE_OPTION_UNSPECIFIED {
		return ErrInvalidVote
	}
	return nil
}

// ValidateBasic implements the sdk.Msg interface.
func (m MsgExecuteProposal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Executor); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid executor address: %s", err)
	}
	return nil
}
