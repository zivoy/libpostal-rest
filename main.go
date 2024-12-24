package main

import (
	"log/slog"

	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"

	"github.com/go-fuego/fuego"
)

type Expansion struct {
	Address    string   `json:"address"`
	Expansions []string `json:"expansions"`
}

type Parse struct {
	Address string                   `json:"address"`
	Parse   []parser.ParsedComponent `json:"parse"`
}

type ExpansionResponse = []Expansion
type ParseResponse = []Parse

type AddressRequest struct {
	Addresses []string `json:"addresses"`
}

type ExpandOptionsRequest struct {
	Options   ExpandOptions `json:"options"`
	Addresses []string      `json:"addresses"`
}

type ParseOptionsRequest struct {
	Options   ParserOptions `json:"options"`
	Addresses []string      `json:"addresses"`
}

func main() {
	s := fuego.NewServer()

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return "LibPostal rest wrapper", nil
	}, fuego.OptionSummary("Welcome page"))

	defaultLibpostalExpandOptions := expand.GetDefaultExpansionOptions()
	defaultLibpostalParseOptions := parser.ParserOptions{ // the get function is not exposed for some reason
		Language: "",
		Country:  "",
	}

	expand := fuego.Group(s, "/expand")
	fuego.Post(expand, "", func(c fuego.ContextWithBody[AddressRequest]) (ExpansionResponse, error) {
		request, err := c.Body()
		if err != nil {
			return nil, err
		}

		return expandAddresses(request.Addresses, defaultLibpostalExpandOptions), nil
	}, fuego.OptionSummary("Expand many addresses"), fuego.OptionDescription("Expand many addresses using the libpostal expand function"))

	fuego.Post(expand, "/advanced", func(c fuego.ContextWithBody[ExpandOptionsRequest]) (ExpansionResponse, error) {
		request, err := c.Body()
		if err != nil {
			return nil, err
		}

		return expandAddresses(request.Addresses, importExpandOptions(request.Options)), nil
	}, fuego.OptionSummary("Expand many addresses with options"), fuego.OptionDescription("Expand many addresses using the libpostal expand function, you can also specify options"))

	fuego.Get(expand, "/default", func(c fuego.ContextNoBody) (ExpandOptions, error) {
		return exportExpandOptions(defaultLibpostalExpandOptions), nil
	}, fuego.OptionSummary("Get default options"), fuego.OptionDescription("Get the default options used by the expand function"))

	parse := fuego.Group(s, "/parse")
	fuego.Post(parse, "", func(c fuego.ContextWithBody[AddressRequest]) (ParseResponse, error) {
		request, err := c.Body()
		if err != nil {
			return nil, err
		}

		return parseAddresses(request.Addresses, defaultLibpostalParseOptions), nil
	}, fuego.OptionSummary("Parse many addresses"), fuego.OptionDescription("Parse many addresses using the libpostal parse function"))

	fuego.Post(parse, "/advanced", func(c fuego.ContextWithBody[ParseOptionsRequest]) (ParseResponse, error) {
		request, err := c.Body()
		if err != nil {
			return nil, err
		}

		return parseAddresses(request.Addresses, importParseOptions(request.Options)), nil
	}, fuego.OptionSummary("Parse many addresses with options"), fuego.OptionDescription("Parse many addresses using the libpostal parse function, you can also specify options"))

	fuego.Get(parse, "/default", func(c fuego.ContextNoBody) (ParserOptions, error) {
		return exportParseOptions(defaultLibpostalParseOptions), nil
	}, fuego.OptionSummary("Get default options"), fuego.OptionDescription("Get the default options used by the parse function"))

	s.Run()
}

func expandAddresses(addresses []string, options expand.ExpandOptions) ExpansionResponse {
	slog.Debug("expanding addresses", "addresses", addresses, "options", options)
	expansions := make([]Expansion, len(addresses))

	for i, str := range addresses {
		var expanded []string
		expanded = expand.ExpandAddressOptions(str, options)
		expansions[i] = Expansion{Address: str, Expansions: expanded}
		slog.Debug("expanded", "expansions", expansions[i])
	}

	return expansions
}

func parseAddresses(addresses []string, options parser.ParserOptions) ParseResponse {
	slog.Debug("parsing addresses", "addresses", addresses, "options", options)
	parses := make([]Parse, len(addresses))

	for i, str := range addresses {
		parsed := parser.ParseAddressOptions(str, options)
		parses[i] = Parse{Address: str, Parse: parsed}
		slog.Debug("parsed", "parses", parses[i])
	}

	return parses
}

func importExpandOptions(options ExpandOptions) expand.ExpandOptions {
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

func exportExpandOptions(options expand.ExpandOptions) ExpandOptions {
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

func importParseOptions(options ParserOptions) parser.ParserOptions {
	return parser.ParserOptions{
		Language: options.Language,
		Country:  options.Country,
	}
}

func exportParseOptions(options parser.ParserOptions) ParserOptions {
	return ParserOptions{
		Language: options.Language,
		Country:  options.Country,
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

type ParserOptions struct {
	Language string `json:"language,omitempty"`
	Country  string `json:"country,omitempty"`
}

// parser labels
const (
	ParserHouse          = "house"
	ParserCategory       = "category"
	ParserNear           = "near"
	ParserHouse_number   = "house_number"
	ParserRoad           = "road"
	ParserUnit           = "unit"
	ParserLevel          = "level"
	ParserStaircase      = "staircase"
	ParserEntrance       = "entrance"
	ParserPo_box         = "po_box"
	ParserPostcode       = "postcode"
	ParserSuburb         = "suburb"
	ParserCity_district  = "city_district"
	ParserCity           = "city"
	ParserIsland         = "island"
	ParserState_district = "state_district"
	ParserState          = "state"
	ParserCountry_region = "country_region"
	ParserCountry        = "country"
	ParserWorld_region   = "world_region"
)

var ParserLabels = [...]string{
	ParserHouse,
	ParserCategory,
	ParserNear,
	ParserHouse_number,
	ParserRoad,
	ParserUnit,
	ParserLevel,
	ParserStaircase,
	ParserEntrance,
	ParserPo_box,
	ParserPostcode,
	ParserSuburb,
	ParserCity_district,
	ParserCity,
	ParserIsland,
	ParserState_district,
	ParserState,
	ParserCountry_region,
	ParserCountry,
	ParserWorld_region,
}
