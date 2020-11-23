package accounts

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"time"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/indexer/types"
	"github.com/ElrondNetwork/elrond-go/data/esdt"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/marshal"
)

var log = logger.GetOrCreate("indexer/process/accounts")

const numDecimalsInFloatBalance = 10

// WrappedUserAccount is structure that is needed for ESDT accounts
type WrappedUserAccount struct {
	Account         state.UserAccountHandler
	TokenIdentifier string
}

// AlteredAccount is structure that holds information about an altered account
type AlteredAccount struct {
	IsESDTSender    bool
	IsESDTOperation bool
	TokenIdentifier string
}

// AccountsProcessor is structure responsible for processing accounts
type AccountsProcessor struct {
	dividerForDenomination float64
	balancePrecision       float64
	internalMarshalizer    marshal.Marshalizer
	addressPubkeyConverter core.PubkeyConverter
	accountsDB             state.AccountsAdapter
}

// NewAccountsProcessor will create a new instance of accounts processor
func NewAccountsProcessor(
	denomination int,
	marshalizer marshal.Marshalizer,
	addressPubkeyConverter core.PubkeyConverter,
	accountsDB state.AccountsAdapter,
) *AccountsProcessor {
	return &AccountsProcessor{
		internalMarshalizer:    marshalizer,
		addressPubkeyConverter: addressPubkeyConverter,
		balancePrecision:       math.Pow(10, float64(numDecimalsInFloatBalance)),
		dividerForDenomination: math.Pow(10, float64(core.MaxInt(denomination, 0))),
		accountsDB:             accountsDB,
	}
}

// GetAccounts will get accounts for egld operations and esdt operations
func (ap *AccountsProcessor) GetAccounts(alteredAccounts map[string]*AlteredAccount) ([]state.UserAccountHandler, []*WrappedUserAccount) {
	accountsToIndexEGLD := make([]state.UserAccountHandler, 0)
	accountsToIndexESDT := make([]*WrappedUserAccount, 0)
	for address, info := range alteredAccounts {
		addressBytes, err := ap.addressPubkeyConverter.Decode(address)
		if err != nil {
			log.Warn("cannot decode address", "address", address, "error", err)
			continue
		}

		account, err := ap.accountsDB.LoadAccount(addressBytes)
		if err != nil {
			log.Warn("cannot load account", "address bytes", addressBytes, "error", err)
			continue
		}

		userAccount, ok := account.(state.UserAccountHandler)
		if !ok {
			log.Warn("cannot cast AccountHandler to type UserAccountHandler")
			continue
		}

		if info.IsESDTOperation {
			accountsToIndexESDT = append(accountsToIndexESDT, &WrappedUserAccount{
				Account:         userAccount,
				TokenIdentifier: info.TokenIdentifier,
			})
		}

		if info.IsESDTSender {
			accountsToIndexEGLD = append(accountsToIndexEGLD, userAccount)
			continue
		}

		accountsToIndexEGLD = append(accountsToIndexEGLD, userAccount)
	}

	return accountsToIndexEGLD, accountsToIndexESDT
}

// PrepareAccountsMapEGLD will prepare a map of accounts with egld
func (ap *AccountsProcessor) PrepareAccountsMapEGLD(accounts []state.UserAccountHandler) map[string]*types.AccountInfo {
	accountsMap := make(map[string]*types.AccountInfo)
	for _, userAccount := range accounts {
		balanceAsFloat := ap.computeBalanceAsFloat(userAccount.GetBalance())
		acc := &types.AccountInfo{
			Nonce:      userAccount.GetNonce(),
			Balance:    userAccount.GetBalance().String(),
			BalanceNum: balanceAsFloat,
		}
		address := ap.addressPubkeyConverter.Encode(userAccount.AddressBytes())
		accountsMap[address] = acc
	}

	return accountsMap
}

// PrepareAccountsMapESDT will preaparare a map of accounts with ESDT tokens
func (ap *AccountsProcessor) PrepareAccountsMapESDT(accounts []*WrappedUserAccount) map[string]*types.AccountInfo {
	accountsESDTMap := make(map[string]*types.AccountInfo)
	for _, wrappedUserAccount := range accounts {
		address := ap.addressPubkeyConverter.Encode(wrappedUserAccount.Account.AddressBytes())

		balance, properties, err := ap.getESDTInfo(wrappedUserAccount)
		if err != nil {
			log.Warn("cannot get esdt info from account",
				"address", address,
				"error", err.Error())
			continue
		}

		acc := &types.AccountInfo{
			Address:         address,
			TokenIdentifier: wrappedUserAccount.TokenIdentifier,
			Balance:         balance.String(),
			BalanceNum:      ap.computeBalanceAsFloat(balance),
			Properties:      properties,
		}

		accountsESDTMap[address] = acc
	}

	return accountsESDTMap
}

// PrepareAccountsHistory will prepare a map of accounts history balance form a map of accounts
func (ap *AccountsProcessor) PrepareAccountsHistory(accounts map[string]*types.AccountInfo) map[string]*types.AccountBalanceHistory {
	currentTimestamp := time.Now().Unix()
	accountsMap := make(map[string]*types.AccountBalanceHistory)
	for address, userAccount := range accounts {
		acc := &types.AccountBalanceHistory{
			Address:         address,
			Balance:         userAccount.Balance,
			Timestamp:       currentTimestamp,
			TokenIdentifier: userAccount.TokenIdentifier,
		}
		addressKey := fmt.Sprintf("%s_%d", address, currentTimestamp)
		accountsMap[addressKey] = acc
	}

	return accountsMap
}

func (ap *AccountsProcessor) getESDTInfo(wrappedUserAccount *WrappedUserAccount) (*big.Int, string, error) {
	if wrappedUserAccount.TokenIdentifier == "" {
		return nil, "", nil
	}

	tokenKey := core.ElrondProtectedKeyPrefix + core.ESDTKeyIdentifier + wrappedUserAccount.TokenIdentifier
	valueBytes, err := wrappedUserAccount.Account.DataTrieTracker().RetrieveValue([]byte(tokenKey))
	if err != nil {
		return nil, "", nil
	}

	esdtToken := &esdt.ESDigitalToken{}
	err = ap.internalMarshalizer.Unmarshal(esdtToken, valueBytes)
	if err != nil {
		return nil, "", nil
	}

	return esdtToken.Value, hex.EncodeToString(esdtToken.Properties), nil
}

func (ap *AccountsProcessor) computeBalanceAsFloat(balance *big.Int) float64 {
	balanceBigFloat := big.NewFloat(0).SetInt(balance)
	balanceFloat64, _ := balanceBigFloat.Float64()

	bal := balanceFloat64 / ap.dividerForDenomination
	balanceFloatWithDecimals := math.Round(bal*ap.balancePrecision) / ap.balancePrecision

	return core.MaxFloat64(balanceFloatWithDecimals, 0)
}
