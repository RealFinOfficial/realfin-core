package keeper

import (
	"context"
	"errors"
	"fmt"

	"realfin/x/insurance/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreatePolicy(ctx context.Context, msg *types.MsgCreatePolicy) (*types.MsgCreatePolicyResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.Policy.Has(ctx, msg.PolicyId)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var policy = types.Policy{
		Creator:            msg.Creator,
		PolicyId:           msg.PolicyId,
		AssetSymbol:        msg.AssetSymbol,
		Provider:           msg.Provider,
		CoverageType:       msg.CoverageType,
		CoveragePercentage: msg.CoveragePercentage,
	}

	if err := k.Policy.Set(ctx, policy.PolicyId, policy); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreatePolicyResponse{}, nil
}

func (k msgServer) UpdatePolicy(ctx context.Context, msg *types.MsgUpdatePolicy) (*types.MsgUpdatePolicyResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Policy.Get(ctx, msg.PolicyId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var policy = types.Policy{
		Creator:            msg.Creator,
		PolicyId:           msg.PolicyId,
		AssetSymbol:        msg.AssetSymbol,
		Provider:           msg.Provider,
		CoverageType:       msg.CoverageType,
		CoveragePercentage: msg.CoveragePercentage,
	}

	if err := k.Policy.Set(ctx, policy.PolicyId, policy); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update policy")
	}

	return &types.MsgUpdatePolicyResponse{}, nil
}

func (k msgServer) DeletePolicy(ctx context.Context, msg *types.MsgDeletePolicy) (*types.MsgDeletePolicyResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Policy.Get(ctx, msg.PolicyId)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Policy.Remove(ctx, msg.PolicyId); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove policy")
	}

	return &types.MsgDeletePolicyResponse{}, nil
}
