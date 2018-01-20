package services

import (
	"fmt"
	"os"
	"crypto-users/pkg/database"
)

func VerifyToken(tokenId string, tokenValue string) (valid bool, message string) {
	token, err := database.FindTokenById(tokenId)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, "A token with the provided token id could not be found."
	}

	if token.Value != tokenValue {
		fmt.Fprintln(os.Stderr, "Provided token value and actual token value does not match. THIS IS SUSPICIOUS.")
		return false, "Provided token value and actual token value does not match."
	}

	if token.Invalidated != false {
		fmt.Fprintln(os.Stderr, "Token is already invalidated.")
		return false, "This token has previously been invalidated, please fetch a new token (login again)."
	}

	masterToken, err := database.GetMasterToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, "Could not locate current master token."
	}

	if token.MasterTokenValue != masterToken.Value {
		fmt.Fprintln(os.Stderr, err)
		return false, "Master token associated with provided token is invalid, please fetch a new token (login again)."
	}

	return true, "Provided token is valid."
}
