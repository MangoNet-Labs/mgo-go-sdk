package response

type DelegatedStakeInfo struct {
	StakedMgoId       string `json:"stakedMgoId"`
	StakeRequestEpoch string `json:"stakeRequestEpoch"`
	StakeActiveEpoch  string `json:"stakeActiveEpoch"`
	Principal         string `json:"principal"`
	Status            string `json:"status"`
	EstimatedReward   string `json:"estimatedReward"`
}

type DelegatedStakesResponse struct {
	ValidatorAddress string               `json:"validatorAddress"`
	StakingPool      string               `json:"stakingPool"`
	Stakes           []DelegatedStakeInfo `json:"stakes"`
}

type MgoXGetCommitteeInfoResponse struct {
	Epoch      string     `json:"epoch"`
	Validators [][]string `json:"validators"`
}

type MgoAddress string

type MgoValidatorSummary struct {
	MgoAddress                   MgoAddress `json:"mgoAddress"`
	ProtocolPubkeyBytes          string     `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes           string     `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes            string     `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes       string     `json:"proofOfPossessionBytes"`
	OperationCapId               string     `json:"operationCapId"`
	Name                         string     `json:"name"`
	Description                  string     `json:"description"`
	ImageUrl                     string     `json:"imageUrl"`
	ProjectUrl                   string     `json:"projectUrl"`
	P2pAddress                   string     `json:"p2pAddress"`
	NetAddress                   string     `json:"netAddress"`
	PrimaryAddress               string     `json:"primaryAddress"`
	WorkerAddress                string     `json:"workerAddress"`
	NextEpochProtocolPubkeyBytes string     `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   string     `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  string     `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   string     `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          string     `json:"nextEpochNetAddress"`
	NextEpochP2pAddress          string     `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      string     `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       string     `json:"nextEpochWorkerAddress"`
	VotingPower                  string     `json:"votingPower"`
	GasPrice                     string     `json:"gasPrice"`
	CommissionRate               string     `json:"commissionRate"`
	NextEpochStake               string     `json:"nextEpochStake"`
	NextEpochGasPrice            string     `json:"nextEpochGasPrice"`
	NextEpochCommissionRate      string     `json:"nextEpochCommissionRate"`
	StakingPoolId                string     `json:"stakingPoolId"`
	StakingPoolActivationEpoch   string     `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch string     `json:"stakingPoolDeactivationEpoch"`
	StakingPoolMgoBalance        string     `json:"stakingPoolMgoBalance"`
	RewardsPool                  string     `json:"rewardsPool"`
	PoolTokenBalance             string     `json:"poolTokenBalance"`
	PendingStake                 string     `json:"pendingStake"`
	PendingPoolTokenWithdraw     string     `json:"pendingPoolTokenWithdraw"`
	PendingTotalMgoWithdraw      string     `json:"pendingTotalMgoWithdraw"`
	ExchangeRatesId              string     `json:"exchangeRatesId"`
	ExchangeRatesSize            string     `json:"exchangeRatesSize"`
}

type MgoSystemStateSummary struct {
	Epoch                                 string                `json:"epoch"`
	ProtocolVersion                       string                `json:"protocolVersion"`
	SystemStateVersion                    string                `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  string                `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       string                `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     string                `json:"referenceGasPrice"`
	SafeMode                              bool                  `json:"safeMode"`
	SafeModeStorageRewards                string                `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            string                `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                string                `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       string                `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 string                `json:"epochStartTimestampMs"`
	EpochDurationMs                       string                `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                string                `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     string                `json:"maxValidatorCount"`
	MinValidatorJoiningStake              string                `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            string                `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        string                `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          string                `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   string                `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       string                `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount string                `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              string                `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              int                   `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            string                `json:"totalStake"`
	ActiveValidators                      []MgoValidatorSummary `json:"activeValidators"`
	PendingActiveValidatorsId             string                `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           string                `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []string              `json:"pendingRemovals"`
	StakingPoolMappingsId                 string                `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               string                `json:"stakingPoolMappingsSize"`
	InactivePoolsId                       string                `json:"inactivePoolsId"`
	InactivePoolsSize                     string                `json:"inactivePoolsSize"`
	ValidatorCandidatesId                 string                `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               string                `json:"validatorCandidatesSize"`
	AtRiskValidators                      []string              `json:"atRiskValidators"`
	ValidatorReportRecords                [][]interface{}       `json:"validatorReportRecords"`
}

type Apy struct {
	Address string  `json:"address"`
	Apy     float64 `json:"apy"`
}

type ValidatorsApy struct {
	Apys  []Apy  `json:"apys"`
	Epoch string `json:"epoch"`
}
