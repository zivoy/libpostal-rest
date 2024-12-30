package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-fuego/fuego"
	"github.com/rs/cors"
)

type Expansion struct {
	Address    string   `json:"address"`
	Expansions []string `json:"expansions"`
}

type Parse struct {
	Address string          `json:"address"`
	Parse   ParsedComponent `json:"parse"`
}

type ExpansionResponse = []Expansion
type ParseResponse = []Parse
type AddressRequest = []string

type ExpandOptionsRequest struct {
	Options   ExpandOptions `json:"options"`
	Addresses []string      `json:"addresses"`
}

type ParseOptionsRequest struct {
	Options   ParserOptions `json:"options"`
	Addresses []string      `json:"addresses"`
}

var Version string

func main() {
	s := fuego.NewServer(
		fuego.WithAddr("0.0.0.0:8724"),
		fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			DisableLocalSave: true,
			DisableSwaggerUI: false,
			DisableSwagger:   false,

			JsonUrl:          "/docs/openapi.json",
			SwaggerUrl:       "/docs",
			PrettyFormatJson: true,
		}),
		fuego.WithCorsMiddleware(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST"},
		}).Handler),
	)
	fuego.Use(s, chiMiddleware.RealIP, chiMiddleware.Recoverer)

	creds := map[string]string{
		"admin": "admin",
	}

	fuego.Get(s, "/", func(c fuego.ContextNoBody) (string, error) {
		return fmt.Sprintf("LibPostal rest wrapper (%s)", Version), nil
	}, fuego.OptionSummary("Welcome page"))

	defaultLibpostalExpandOptions := expand.GetDefaultExpansionOptions()
	defaultLibpostalParseOptions := parser.ParserOptions{ // the get function is not exposed for some reason
		Language: "",
		Country:  "",
	}

	expand := fuego.Group(s, "/expand")
	fuego.Use(expand, chiMiddleware.Logger, chiMiddleware.BasicAuth("libpostal", creds))
	fuego.Post(expand, "", func(c fuego.ContextWithBody[AddressRequest]) (ExpansionResponse, error) {
		addresses, err := parseAddressList(c.Request())
		if err != nil {
			return nil, err
		}

		return expandAddresses(addresses, defaultLibpostalExpandOptions), nil
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
	fuego.Use(parse, chiMiddleware.Logger, chiMiddleware.BasicAuth("libpostal", creds))
	fuego.Post(parse, "", func(c fuego.ContextWithBody[AddressRequest]) (ParseResponse, error) {
		addresses, err := parseAddressList(c.Request())
		if err != nil {
			return nil, err
		}

		return parseAddresses(addresses, defaultLibpostalParseOptions), nil
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

	err := s.Run()
	if err != nil {
		slog.Error("Error running server", "error", err)
	}
}

func parseAddressList(r *http.Request) (AddressRequest, error) {
	var addresses []string
	err := json.NewDecoder(r.Body).Decode(&addresses)
	if err != nil {
		return nil, err
	}

	return addresses, nil
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
		parses[i] = Parse{Address: str, Parse: getParsedComponents(parsed)}
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
	Languages              []string `json:"languages"`
	AddressComponents      uint16   `json:"address_components"`
	LatinAscii             bool     `json:"latin_ascii"`
	Transliterate          bool     `json:"transliterate"`
	StripAccents           bool     `json:"strip_accents"`
	Decompose              bool     `json:"decompose"`
	Lowercase              bool     `json:"lowercase"`
	TrimString             bool     `json:"trim_string"`
	ReplaceWordHyphens     bool     `json:"replace_word_hyphens"`
	DeleteWordHyphens      bool     `json:"delete_word_hyphens"`
	ReplaceNumericHyphens  bool     `json:"replace_numeric_hyphens"`
	DeleteNumericHyphens   bool     `json:"delete_numeric_hyphens"`
	SplitAlphaFromNumeric  bool     `json:"split_alpha_from_numeric"`
	DeleteFinalPeriods     bool     `json:"delete_final_periods"`
	DeleteAcronymPeriods   bool     `json:"delete_acronym_periods"`
	DropEnglishPossessives bool     `json:"drop_english_possessives"`
	DeleteApostrophes      bool     `json:"delete_apostrophes"`
	ExpandNumex            bool     `json:"expand_numex"`
	RomanNumerals          bool     `json:"roman_numerals"`
}

type ParserOptions struct {
	Language string `json:"language"`
	Country  string `json:"country"`
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

type ParsedComponent struct {
	House         string `json:"house,omitempty"`
	Category      string `json:"category,omitempty"`
	Near          string `json:"near,omitempty"`
	HouseNumber   string `json:"house_number,omitempty"`
	Road          string `json:"road,omitempty"`
	Unit          string `json:"unit,omitempty"`
	Level         string `json:"level,omitempty"`
	Staircase     string `json:"staircase,omitempty"`
	Entrance      string `json:"entrance,omitempty"`
	PoBox         string `json:"po_box,omitempty"`
	Postcode      string `json:"postcode,omitempty"`
	Suburb        string `json:"suburb,omitempty"`
	CityDistrict  string `json:"city_district,omitempty"`
	City          string `json:"city,omitempty"`
	Island        string `json:"island,omitempty"`
	StateDistrict string `json:"state_district,omitempty"`
	State         string `json:"state,omitempty"`
	CountryRegion string `json:"country_region,omitempty"`
	Country       string `json:"country,omitempty"`
	WorldRegion   string `json:"world_region,omitempty"`
}

func getParsedComponents(parsedComponents []parser.ParsedComponent) ParsedComponent {
	compoonent := ParsedComponent{}

	for _, component := range parsedComponents {
		switch component.Label {
		case ParserHouse:
			compoonent.House = component.Value
		case ParserCategory:
			compoonent.Category = component.Value
		case ParserNear:
			compoonent.Near = component.Value
		case ParserHouse_number:
			compoonent.HouseNumber = component.Value
		case ParserRoad:
			compoonent.Road = component.Value
		case ParserUnit:
			compoonent.Unit = component.Value
		case ParserLevel:
			compoonent.Level = component.Value
		case ParserStaircase:
			compoonent.Staircase = component.Value
		case ParserEntrance:
			compoonent.Entrance = component.Value
		case ParserPo_box:
			compoonent.PoBox = component.Value
		case ParserPostcode:
			compoonent.Postcode = component.Value
		case ParserSuburb:
			compoonent.Suburb = component.Value
		case ParserCity_district:
			compoonent.CityDistrict = component.Value
		case ParserCity:
			compoonent.City = component.Value
		case ParserIsland:
			compoonent.Island = component.Value
		case ParserState_district:
			compoonent.StateDistrict = component.Value
		case ParserState:
			compoonent.State = component.Value
		case ParserCountry_region:
			compoonent.CountryRegion = component.Value
		case ParserCountry:
			compoonent.Country = component.Value
		case ParserWorld_region:
			compoonent.WorldRegion = component.Value
		default:
			slog.Warn("unknown component", "component", component)
		}
	}

	return compoonent
}
