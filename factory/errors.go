package factory

import "errors"

// ErrNilEconomicsData signals that a nil economics data handler has been provided
var ErrNilEconomicsData = errors.New("nil economics data provided")

// ErrNilGenesisConfiguration signals that a nil genesis configuration has been provided
var ErrNilGenesisConfiguration = errors.New("nil genesis configuration provided")

// ErrNilCoreComponents signals that nil core components have been provided
var ErrNilCoreComponents = errors.New("nil core components provided")

// ErrNilTriesComponents signals that nil tries components have been provided
var ErrNilTriesComponents = errors.New("nil tries components provided")

// ErrNilShardCoordinator signals that nil core components have been provided
var ErrNilShardCoordinator = errors.New("nil shard coordinator provided")

// ErrNilPathManager signals that a nil path manager has been provided
var ErrNilPathManager = errors.New("nil path manager provided")

// ErrNilPath signals that a nil/empty path was provided
var ErrNilPath = errors.New("nil path provided")

// ErrNilMarshalizer signals that a nil marshalizer has been provided
var ErrNilMarshalizer = errors.New("nil marshalizer provided")

// ErrNilHasher signals that a nil hasher has been provided
var ErrNilHasher = errors.New("nil hasher provided")

// ErrNilEpochStartNotifier signals that a nil epoch start notifier has been provided
var ErrNilEpochStartNotifier = errors.New("nil epoch start notifier provided")

// ErrHasherCreation signals that the hasher cannot be created based on provided data
var ErrHasherCreation = errors.New("error creating hasher")

// ErrMarshalizerCreation signals that the marshalizer cannot be created based on provided data
var ErrMarshalizerCreation = errors.New("error creating marshalizer")

// ErrPubKeyConverterCreation signals that the public key converter cannot be created based on provided data
var ErrPubKeyConverterCreation = errors.New("error creating public key converter")

// ErrAccountsAdapterCreation signals that the accounts adapter cannot be created based on provided data
var ErrAccountsAdapterCreation = errors.New("error creating accounts adapter")

// ErrInitialBalancesCreation signals that the initial balances cannot be created based on provided data
var ErrInitialBalancesCreation = errors.New("error creating initial balances")

// ErrPublicKeyMismatch signals that the read public key mismatch the one read
var ErrPublicKeyMismatch = errors.New("public key mismatch between the computed and the one read from the file")

// ErrBlockchainCreation signals that the blockchain cannot be created
var ErrBlockchainCreation = errors.New("can not create blockchain")

// ErrDataStoreCreation signals that the data store cannot be created
var ErrDataStoreCreation = errors.New("can not create data store")

// ErrDataPoolCreation signals that the data pool cannot be created
var ErrDataPoolCreation = errors.New("can not create data pool")

// ErrInvalidConsensusConfig signals that an invalid consensus type is specified in the configuration file
var ErrInvalidConsensusConfig = errors.New("invalid consensus type provided in config file")

// ErrMultiSigHasherMissmatch signals that an invalid multisig hasher was provided
var ErrMultiSigHasherMissmatch = errors.New("wrong multisig hasher provided for bls consensus type")

// ErrMissingMultiHasherConfig signals that the multihasher type isn't specified in the configuration file
var ErrMissingMultiHasherConfig = errors.New("no multisig hasher provided in config file")

// ErrNilStatusHandler signals that a nil status handler has been provided
var ErrNilStatusHandler = errors.New("nil status handler provided")

// ErrWrongTypeAssertion signals that a wrong type assertion occurred
var ErrWrongTypeAssertion = errors.New("wrong type assertion")