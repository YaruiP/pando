// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/fox-one/pando/cmd/pando-server/config"
	"github.com/fox-one/pando/handler/api"
	"github.com/fox-one/pando/handler/rpc"
	"github.com/fox-one/pando/notifier"
	"github.com/fox-one/pando/server"
	asset2 "github.com/fox-one/pando/service/asset"
	user2 "github.com/fox-one/pando/service/user"
	"github.com/fox-one/pando/service/wallet"
	"github.com/fox-one/pando/session"
	"github.com/fox-one/pando/store/asset"
	"github.com/fox-one/pando/store/collateral"
	"github.com/fox-one/pando/store/flip"
	"github.com/fox-one/pando/store/message"
	"github.com/fox-one/pando/store/oracle"
	"github.com/fox-one/pando/store/transaction"
	"github.com/fox-one/pando/store/user"
	"github.com/fox-one/pando/store/vault"
)

// Injectors from wire.go:

func buildServer(cfg *config.Config) (*server.Server, error) {
	db, err := provideDatabase(cfg)
	if err != nil {
		return nil, err
	}
	userStore := user.New(db)
	client, err := provideMixinClient(cfg)
	if err != nil {
		return nil, err
	}
	userConfig := provideUserServiceConfig(cfg)
	userService := user2.New(client, userConfig)
	sessionConfig := provideSessionConfig(cfg)
	coreSession := session.New(userStore, userService, sessionConfig)
	assetStore := asset.New(db)
	vaultStore := vault.New(db)
	flipStore := flip.New(db)
	collateralStore := collateral.New(db)
	transactionStore := transaction.New(db)
	walletConfig := provideWalletServiceConfig(cfg)
	walletService := wallet.New(client, walletConfig)
	system := provideSystem(cfg)
	assetService := asset2.New(client)
	messageStore := message.New(db)
	localizer, err := provideLocalizer(cfg)
	if err != nil {
		return nil, err
	}
	coreNotifier := notifier.New(system, assetService, messageStore, vaultStore, collateralStore, userStore, flipStore, localizer)
	oracleStore := oracle.New(db)
	apiServer := api.New(coreSession, userService, assetStore, vaultStore, flipStore, collateralStore, transactionStore, walletService, coreNotifier, oracleStore, system)
	rpcServer := rpc.New(assetStore, vaultStore, flipStore, oracleStore, collateralStore, transactionStore)
	mux := provideRoute(apiServer, rpcServer, coreSession)
	serverServer := provideServer(mux)
	return serverServer, nil
}
