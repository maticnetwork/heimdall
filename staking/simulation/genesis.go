package simulation

// RandomizedGenState generates a random GenesisState for staking
// func RandomGenesisValidatorsData(r *rand.Rand, n int) ([]authTypes.GenesisAccount, []hmTypes.Validator, *hmTypes.ValidatorSet) {
// 	accounts := RandomAccounts(r, n)

// 	validators := make([]*hmTypes.Validator, n)
// 	dividendAccounts := make([]hmTypes.DividendAccount, n)

// 	for i = 0; i < n ; i++ {
// 			validators[i] = hmTypes.NewValidator(
// 				hmTypes.NewValidatorID(uint64(int64(i))),
// 				0,
// 				0,
// 				RandIntBetween(r, 10, 100), // power
// 				accounts[i].PubKey,
// 				accounts[i].Address(),
// 			)

// 			// create dividend account for validator
// 			dividendAccounts[i] = hmTypes.NewDividendAccount(
// 				hmTypes.NewDividendAccountID(uint64(validators[i].ID)),
// 				big.NewInt(0).String(),
// 				big.NewInt(0).String(),
// 			)
// 		}
// 	}

// 	// validator set
// 	validatorSet := hmTypes.NewValidatorSet(validators)
// 	return validators, validatorSet, accounts
// }
