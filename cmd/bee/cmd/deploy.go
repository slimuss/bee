// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/bee/pkg/settlement/swap/chequebook"
	"github.com/ethersphere/bee/pkg/settlement/swap/transaction"
	"github.com/ethersphere/bee/pkg/statestore/leveldb"
	mockinmem "github.com/ethersphere/bee/pkg/statestore/mock"
	"github.com/ethersphere/bee/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initBackend(ctx context.Context, logger logging.Logger, swapEndpoint string) (*ethclient.Client, error) {
	swapBackend, err := ethclient.Dial(swapEndpoint)
	if err != nil {
		return nil, err
	}
	maxDelay := 1 * time.Minute
	synced, err := transaction.IsSynced(ctx, swapBackend, maxDelay)
	if err != nil {
		return nil, err
	}
	if !synced {
		logger.Infof("waiting for ethereum backend to be synced.")
		err = transaction.WaitSynced(ctx, swapBackend, maxDelay)
		if err != nil {
			return nil, fmt.Errorf("could not wait for ethereum backend to sync: %w", err)
		}
	}
	return swapBackend, nil
}

func (c *command) initDeployCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy and fund the chequebook contract",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) > 0 {
				return cmd.Help()
			}

			var logger logging.Logger
			switch v := strings.ToLower(c.config.GetString(optionNameVerbosity)); v {
			case "0", "silent":
				logger = logging.New(ioutil.Discard, 0)
			case "1", "error":
				logger = logging.New(cmd.OutOrStdout(), logrus.ErrorLevel)
			case "2", "warn":
				logger = logging.New(cmd.OutOrStdout(), logrus.WarnLevel)
			case "3", "info":
				logger = logging.New(cmd.OutOrStdout(), logrus.InfoLevel)
			case "4", "debug":
				logger = logging.New(cmd.OutOrStdout(), logrus.DebugLevel)
			case "5", "trace":
				logger = logging.New(cmd.OutOrStdout(), logrus.TraceLevel)
			default:
				return fmt.Errorf("unknown verbosity level %q", v)
			}

			factoryAddressOption := c.config.GetString(optionNameSwapFactoryAddress)
			swapInitialDeposit := c.config.GetUint64(optionNameSwapInitialDeposit)
			swapEndpoint := c.config.GetString(optionNameSwapEndpoint)
			dataDir := c.config.GetString(optionNameDataDir)

			var stateStore storage.StateStorer
			if dataDir == "" {
				stateStore = mockinmem.NewStateStore()
				logger.Warning("using in-mem state store. no node state will be persisted")
			} else {
				stateStore, err = leveldb.NewStateStore(filepath.Join(dataDir, "statestore"))
				if err != nil {
					return fmt.Errorf("statestore: %w", err)
				}
			}

			signerConfig, err := c.configureSigner(cmd, logger)
			if err != nil {
				return err
			}

			swapBackend, err := initBackend(cmd.Context(), logger, swapEndpoint)
			if err != nil {
				return err
			}

			transactionService, err := transaction.NewService(logger, swapBackend, signerConfig.signer)
			if err != nil {
				return err
			}
			overlayEthAddress, err := signerConfig.signer.EthereumAddress()
			if err != nil {
				return err
			}

			chainID, err := swapBackend.ChainID(cmd.Context())
			if err != nil {
				logger.Infof("could not connect to backend at %v. In a swap-enabled network a working blockchain node (for goerli network in production) is required. Check your node or specify another node using --swap-endpoint.", swapEndpoint)
				return fmt.Errorf("could not get chain id from ethereum backend: %w", err)
			}

			var factoryAddress common.Address
			if factoryAddressOption == "" {
				var found bool
				factoryAddress, found = chequebook.DiscoverFactoryAddress(chainID.Int64())
				if !found {
					return errors.New("no known factory address for this network")
				}
				logger.Infof("using default factory address for chain id %d: %x", chainID, factoryAddress)
			} else if !common.IsHexAddress(factoryAddressOption) {
				return errors.New("malformed factory address")
			} else {
				factoryAddress = common.HexToAddress(factoryAddressOption)
				logger.Infof("using custom factory address: %x", factoryAddress)
			}

			chequebookFactory, err := chequebook.NewFactory(swapBackend, transactionService, factoryAddress, chequebook.NewSimpleSwapFactoryBindingFunc)
			if err != nil {
				return err
			}

			_, err = chequebook.Init(cmd.Context(),
				chequebookFactory,
				stateStore,
				logger,
				swapInitialDeposit,
				transactionService,
				swapBackend,
				chainID.Int64(),
				overlayEthAddress,
				nil,
				chequebook.NewSimpleSwapBindings,
				chequebook.NewERC20Bindings)
			if err != nil {
				return err
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return c.config.BindPFlags(cmd.Flags())
		},
	}

	c.setAllFlags(cmd)
	c.root.AddCommand(cmd)
	return nil
}
