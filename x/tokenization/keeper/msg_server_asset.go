package keeper

import (
	"context"
	"errors"
	"fmt"

	"realfin/x/tokenization/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateAsset(ctx context.Context, msg *types.MsgCreateAsset) (*types.MsgCreateAssetResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if the value already exists
	ok, err := k.Asset.Has(ctx, msg.Symbol)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var asset = types.Asset{
		Creator:     msg.Creator,
		Symbol:      msg.Symbol,
		Name:        msg.Name,
		Description: msg.Description,
		AssetType:   msg.AssetType,
		Metadata:    msg.Metadata,
	}

	if err := k.Asset.Set(ctx, asset.Symbol, asset); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateAssetResponse{}, nil
}

func (k msgServer) UpdateAsset(ctx context.Context, msg *types.MsgUpdateAsset) (*types.MsgUpdateAssetResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Asset.Get(ctx, msg.Symbol)
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

	var asset = types.Asset{
		Creator:     msg.Creator,
		Symbol:      msg.Symbol,
		Name:        msg.Name,
		Description: msg.Description,
		AssetType:   msg.AssetType,
		Metadata:    msg.Metadata,
	}

	if err := k.Asset.Set(ctx, asset.Symbol, asset); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update asset")
	}

	return &types.MsgUpdateAssetResponse{}, nil
}

func (k msgServer) DeleteAsset(ctx context.Context, msg *types.MsgDeleteAsset) (*types.MsgDeleteAssetResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	// Check if the value exists
	val, err := k.Asset.Get(ctx, msg.Symbol)
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

	if err := k.Asset.Remove(ctx, msg.Symbol); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove asset")
	}

	return &types.MsgDeleteAssetResponse{}, nil
}
