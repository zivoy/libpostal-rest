package main

import (
	"log/slog"

	expand "github.com/openvenues/gopostal/expand"
	// parser "github.com/openvenues/gopostal/parser"

	"github.com/go-fuego/fuego"
)

type Expansion struct {
	Address    string   `json:"address"`
	Expansions []string `json:"expansions"`
}

type ExpansionRequest = []string
type ExpansionResponse = []Expansion

type ExpandOptionsRequest struct {
	Options   ExpandOptions `json:"Options"`
	Addresses []string      `json:"addresses"`
}

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "LibPostal rest wrapper", nil
	}, fuego.OptionSummary("Welcome page"))

	defailtLibpostalOptions := expand.GetDefaultExpansionOptions()

	expand := fuego.Group(s, "/expand")
	fuego.Post(expand, "", func(c fuego.ContextWithBody[ExpansionRequest]) (ExpansionResponse, error) {
		addresses, err := c.Body()
		if err != nil {
			return nil, err
		}

		return expandAddresses(addresses, defailtLibpostalOptions), nil
	}, fuego.OptionSummary("Expand many addresses"), fuego.OptionDescription("Expand many addresses using the libpostal expand function"))

	fuego.Post(expand, "/advanced", func(c fuego.ContextWithBody[ExpandOptionsRequest]) (ExpansionResponse, error) {
		request, err := c.Body()
		if err != nil {
			return nil, err
		}

		return expandAddresses(request.Addresses, getOptions(request.Options)), nil
	}, fuego.OptionSummary("Expand many addresses with options"), fuego.OptionDescription("Expand many addresses using the libpostal expand function, you can also specify options"))

	fuego.Get(expand, "/default", func(c fuego.ContextNoBody) (ExpandOptions, error) {
		return getExteriorOptions(defailtLibpostalOptions), nil
	}, fuego.OptionSummary("Get default options"), fuego.OptionDescription("Get the default options used by the expand function"))

	s.Run()
}

func expandAddresses(addresses []string, options expand.ExpandOptions) ExpansionResponse {
	slog.Debug("expanding addresses", "addresses", addresses, "options", options)
	expansions := make([]Expansion, len(addresses))

	for i, str := range addresses {
		expanded := expand.ExpandAddressOptions(str, options)
		expansions[i] = Expansion{Address: str, Expansions: expanded}
		slog.Debug("expanded", "expansions", expansions[i])
	}

	return expansions
}

func getOptions(options ExpandOptions) expand.ExpandOptions {
	return expand.ExpandOptions{
		Languages:              options.Languages,
		AddressComponents:      options.AddressComponents,
		LatinAscii:             options.LatinAscii,
		Transliterate:          options.Transliterate,
		StripAccents:           options.StripAccents,
		Decompose:              options.Decompose,
		Lowercase:              options.Lowercase,
		TrimString:             options.TrimString,
		ReplaceWordHyphens:     options.ReplaceWordHyphens,
		DeleteWordHyphens:      options.DeleteWordHyphens,
		ReplaceNumericHyphens:  options.ReplaceNumericHyphens,
		DeleteNumericHyphens:   options.DeleteNumericHyphens,
		SplitAlphaFromNumeric:  options.SplitAlphaFromNumeric,
		DeleteFinalPeriods:     options.DeleteFinalPeriods,
		DeleteAcronymPeriods:   options.DeleteAcronymPeriods,
		DropEnglishPossessives: options.DropEnglishPossessives,
		DeleteApostrophes:      options.DeleteApostrophes,
		ExpandNumex:            options.ExpandNumex,
		RomanNumerals:          options.RomanNumerals,
	}
}

func getExteriorOptions(options expand.ExpandOptions) ExpandOptions {
	return ExpandOptions{
		Languages:              options.Languages,
		AddressComponents:      options.AddressComponents,
		LatinAscii:             options.LatinAscii,
		Transliterate:          options.Transliterate,
		StripAccents:           options.StripAccents,
		Decompose:              options.Decompose,
		Lowercase:              options.Lowercase,
		TrimString:             options.TrimString,
		ReplaceWordHyphens:     options.ReplaceWordHyphens,
		DeleteWordHyphens:      options.DeleteWordHyphens,
		ReplaceNumericHyphens:  options.ReplaceNumericHyphens,
		DeleteNumericHyphens:   options.DeleteNumericHyphens,
		SplitAlphaFromNumeric:  options.SplitAlphaFromNumeric,
		DeleteFinalPeriods:     options.DeleteFinalPeriods,
		DeleteAcronymPeriods:   options.DeleteAcronymPeriods,
		DropEnglishPossessives: options.DropEnglishPossessives,
		DeleteApostrophes:      options.DeleteApostrophes,
		ExpandNumex:            options.ExpandNumex,
		RomanNumerals:          options.RomanNumerals,
	}
}

type ExpandOptions struct {
	Languages              []string `json:"languages,omitempty"`
	AddressComponents      uint16   `json:"address_components,omitempty"`
	LatinAscii             bool     `json:"latin_ascii,omitempty"`
	Transliterate          bool     `json:"transliterate,omitempty"`
	StripAccents           bool     `json:"strip_accents,omitempty"`
	Decompose              bool     `json:"decompose,omitempty"`
	Lowercase              bool     `json:"lowercase,omitempty"`
	TrimString             bool     `json:"trim_string,omitempty"`
	ReplaceWordHyphens     bool     `json:"replace_word_hyphens,omitempty"`
	DeleteWordHyphens      bool     `json:"delete_word_hyphens,omitempty"`
	ReplaceNumericHyphens  bool     `json:"replace_numeric_hyphens,omitempty"`
	DeleteNumericHyphens   bool     `json:"delete_numeric_hyphens,omitempty"`
	SplitAlphaFromNumeric  bool     `json:"split_alpha_from_numeric,omitempty"`
	DeleteFinalPeriods     bool     `json:"delete_final_periods,omitempty"`
	DeleteAcronymPeriods   bool     `json:"delete_acronym_periods,omitempty"`
	DropEnglishPossessives bool     `json:"drop_english_possessives,omitempty"`
	DeleteApostrophes      bool     `json:"delete_apostrophes,omitempty"`
	ExpandNumex            bool     `json:"expand_numex,omitempty"`
	RomanNumerals          bool     `json:"roman_numerals,omitempty"`
}
