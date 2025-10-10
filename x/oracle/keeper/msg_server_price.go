package keeper

import (
	"context"
	"errors"
	"fmt"

	"realfin/x/oracle/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreatePrice(ctx context.Context, msg *types.MsgCreatePrice) (*types.MsgCreatePriceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.Price.Has(ctx, msg.Symbol)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var price = types.Price{
		Creator:     msg.Creator,
		Symbol:      msg.Symbol,
		Rate:        msg.Rate,
		Name:        msg.Name,
		Description: msg.Description,
	}

	if err := k.Price.Set(ctx, price.Symbol, price); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreatePriceResponse{}, nil
}

func (k msgServer) UpdatePrice(ctx context.Context, msg *types.MsgUpdatePrice) (*types.MsgUpdatePriceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Price.Get(ctx, msg.Symbol)
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

	var price = types.Price{
		Creator:     msg.Creator,
		Symbol:      msg.Symbol,
		Rate:        msg.Rate,
		Name:        msg.Name,
		Description: msg.Description,
	}

	if err := k.Price.Set(ctx, price.Symbol, price); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update price")
	}

	return &types.MsgUpdatePriceResponse{}, nil
}

func (k msgServer) DeletePrice(ctx context.Context, msg *types.MsgDeletePrice) (*types.MsgDeletePriceResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Price.Get(ctx, msg.Symbol)
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

	if err := k.Price.Remove(ctx, msg.Symbol); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove price")
	}

	return &types.MsgDeletePriceResponse{}, nil
}
