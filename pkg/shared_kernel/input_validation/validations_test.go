package input_validation

import (
	"fmt"
	"testing"
)

func TestCheckPhone(t *testing.T) {
	testPhones := []string{
		"+49 170 1234567",   // Germany
		"+33 6 12 34 56 78", // French (not in list)
		"+855 800 12345",    // Cambodia (not in list)
		"+1 212 555 1234",   // New York (not in list)
		"+1 415 555 2671",   // San Francisco (not in list)
		"+98123456789",      // Iran
		"+99 999 999999",    // Not Valid
	}

	// Add more regions for testing
	AddPhoneRegion("+855", 4) // Cambodia

	for _, phone := range testPhones {
		info, err := CheckPhone(phone)
		if err != nil {
			fmt.Printf("Error for %s: %v\n", phone, err)
		} else {
			fmt.Printf("Valid: %s -> Prefix: %s, Country: %v, Number: %s\n",
				phone, info.Prefix, info.Country, info.Number)
		}
	}
}
