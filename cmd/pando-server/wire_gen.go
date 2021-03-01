// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/handler/api"
	"github.com/fox-one/pando/handler/rpc"
	"github.com/fox-one/pando/server"
	"github.com/fox-one/pando/service/user"
	"github.com/fox-one/pando/store/asset"
	"github.com/fox-one/pando/store/collateral"
	"github.com/fox-one/pando/store/transaction"
	"github.com/fox-one/pando/store/vault"
)

// Injectors from wire.go:

func buildServer(cfg *config.Config) (*server.Server, error) {
	client, err := provideMixinClient(cfg)
	if err != nil {
		return nil, err
	}
	userService := user.New(client)
	session := provideSessions(userService, cfg)
	db, err := provideDatabase(cfg)
	if err != nil {
		return nil, err
	}
	assetStore := asset.New(db)
	vaultStore := vault.New(db)
	collateralStore := Collateral.New(db)
	transactionStore := transaction.New(db)
	system := provideSystem(cfg)
	walletService := provideWalletService(client, cfg, system)
	apiServer := api.New(session, assetStore, vaultStore, collateralStore, transactionStore, walletService, system)
	rpcServer := rpc.New(assetStore, vaultStore, collateralStore, transactionStore)
	mainHealthHandler := provideHealth(system)
	mux := provideRoute(apiServer, rpcServer, session, mainHealthHandler)
	serverServer := provideServer(mux)
	return serverServer, nil
}
