package keeper

import (
	"context"
	"errors"
	"fmt"

	"realfin/x/realestate/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateRate(ctx context.Context, msg *types.MsgCreateRate) (*types.MsgCreateRateResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.Rate.Has(ctx, msg.Symbol)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var rate = types.Rate{
		Creator:     msg.Creator,
		Symbol:      msg.Symbol,
		Rate:        msg.Rate,
		Name:        msg.Name,
		Description: msg.Description,
	}

	if err := k.Rate.Set(ctx, rate.Symbol, rate); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateRateResponse{}, nil
}

func (k msgServer) UpdateRate(ctx context.Context, msg *types.MsgUpdateRate) (*types.MsgUpdateRateResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Rate.Get(ctx, msg.Symbol)
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

	var rate = types.Rate{
		Creator:     msg.Creator,
		Symbol:      msg.Symbol,
		Rate:        msg.Rate,
		Name:        msg.Name,
		Description: msg.Description,
	}

	if err := k.Rate.Set(ctx, rate.Symbol, rate); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update rate")
	}

	return &types.MsgUpdateRateResponse{}, nil
}

func (k msgServer) DeleteRate(ctx context.Context, msg *types.MsgDeleteRate) (*types.MsgDeleteRateResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Rate.Get(ctx, msg.Symbol)
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

	if err := k.Rate.Remove(ctx, msg.Symbol); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove rate")
	}

	return &types.MsgDeleteRateResponse{}, nil
}
