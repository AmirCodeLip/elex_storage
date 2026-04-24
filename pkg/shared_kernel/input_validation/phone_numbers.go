package input_validation

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

var (
	space_reg       = regexp.MustCompile(`\s`)
	phonePrefix_reg = regexp.MustCompile(`^\+(\d{1,3})`)
	phone_reg       = regexp.MustCompile(`^\+?[0-9]{7,15}$`)
)

type Country byte

const (
	Iran Country = iota
	Germany
	France
	USACanada
	Cambodia
)

type PhoneRegion struct {
	Prefix  string
	Country Country
}

// PhoneRegions sorted by prefix length (longest first) to handle overlapping prefixes correctly
var PhoneRegions = func() []PhoneRegion {
	regions := []PhoneRegion{
		{Prefix: "+98", Country: Iran},
		{Prefix: "+49", Country: Germany},
		{Prefix: "+33", Country: France},
		{Prefix: "+1", Country: USACanada},
	}

	// Sort by prefix length descending, then alphabetically
	sort.Slice(regions, func(i, j int) bool {
		if len(regions[i].Prefix) != len(regions[j].Prefix) {
			return len(regions[i].Prefix) > len(regions[j].Prefix)
		}
		return regions[i].Prefix < regions[j].Prefix
	})

	return regions
}()

type PhoneInfo struct {
	Prefix         string
	Country        Country
	Number         string
	FixPhoneNumber string
}

// existPhone finds the matching phone region for a given phone number
func existPhone(phone string) *PhoneRegion {
	for _, region := range PhoneRegions {
		// Check if the phone starts with the region prefix
		if strings.HasPrefix(phone, region.Prefix) {
			return &region
		}
	}
	return nil
}

// CheckPhone validates and parses a phone number
func CheckPhone(phone string) (*PhoneInfo, error) {
	// Remove all whitespace
	phone = space_reg.ReplaceAllString(phone, "")

	// Basic validation
	if len(phone) < 7 {
		return nil, errors.New("invalid phone number: too short")
	}

	// Must start with '+'
	if !strings.HasPrefix(phone, "+") {
		return nil, errors.New("invalid phone number: must start with '+'")
	}

	// Find matching region
	phoneRegion := existPhone(phone)
	if phoneRegion == nil {
		return nil, errors.New("invalid phone number: unsupported country code")
	}

	// Extract the number part (remove the prefix)
	phoneNumber := strings.TrimPrefix(phone, phoneRegion.Prefix)

	// Validate the remaining number
	if len(phoneNumber) < 4 {
		return nil, errors.New("invalid phone number: local number too short")
	}

	// Ensure the number contains only digits
	for _, ch := range phoneNumber {
		if ch < '0' || ch > '9' {
			return nil, errors.New("invalid phone number: contains non-digit characters")
		}
	}

	return &PhoneInfo{
		Number:         phoneNumber,
		Prefix:         phoneRegion.Prefix,
		Country:        phoneRegion.Country,
		FixPhoneNumber: phoneRegion.Prefix + phoneNumber,
	}, nil
}

// Additional helper function to add more phone regions
func AddPhoneRegion(prefix string, country Country) {
	// Check if prefix already exists
	for i, region := range PhoneRegions {
		if region.Prefix == prefix {
			PhoneRegions[i].Country = country
			return
		}
	}

	// Add new region
	PhoneRegions = append(PhoneRegions, PhoneRegion{Prefix: prefix, Country: country})

	// Re-sort
	sort.Slice(PhoneRegions, func(i, j int) bool {
		if len(PhoneRegions[i].Prefix) != len(PhoneRegions[j].Prefix) {
			return len(PhoneRegions[i].Prefix) > len(PhoneRegions[j].Prefix)
		}
		return PhoneRegions[i].Prefix < PhoneRegions[j].Prefix
	})
}

func IsPhoneNumber(input string) error {
	if !phone_reg.MatchString(input) {
		return errors.New("Invalid phone number is provided")
	}
	return nil
}
