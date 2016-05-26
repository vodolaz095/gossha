/*
Package inflection pluralizes and singularizes English nouns.

		inflection.Plural("person") => "people"
		inflection.Plural("Person") => "People"
		inflection.Plural("PERSON") => "PEOPLE"

		inflection.Singularize("people") => "person"
		inflection.Singularize("People") => "Person"
		inflection.Singularize("PEOPLE") => "PERSON"

		inflection.Plural("FancyPerson") => "FancydPeople"
		inflection.Singularize("FancyPeople") => "FancydPerson"

Standard rules are from Rails's ActiveSupport (https://github.com/rails/rails/blob/master/activesupport/lib/active_support/inflections.rb)

If you want to register more rules, follow:

		inflection.AddUncountable("fish")
		inflection.AddIrregular("person", "people")
		inflection.AddPlural("(bu)s$", "${1}ses") # "bus" => "buses" / "BUS" => "BUSES" / "Bus" => "Buses"
		inflection.AddSingular("(bus)(es)?$", "${1}") # "buses" => "bus" / "Buses" => "Bus" / "BUSES" => "BUS"
*/
package inflection

import (
	"regexp"
	"strings"
)

var pluralInflections = [][]string{
	{"([a-z])$", "${1}s"},
	{"s$", "s"},
	{"^(ax|test)is$", "${1}es"},
	{"(octop|vir)us$", "${1}i"},
	{"(octop|vir)i$", "${1}i"},
	{"(alias|status)$", "${1}es"},
	{"(bu)s$", "${1}ses"},
	{"(buffal|tomat)o$", "${1}oes"},
	{"([ti])um$", "${1}a"},
	{"([ti])a$", "${1}a"},
	{"sis$", "ses"},
	{"(?:([^f])fe|([lr])f)$", "${1}${2}ves"},
	{"(hive)$", "${1}s"},
	{"([^aeiouy]|qu)y$", "${1}ies"},
	{"(x|ch|ss|sh)$", "${1}es"},
	{"(matr|vert|ind)(?:ix|ex)$", "${1}ices"},
	{"^(m|l)ouse$", "${1}ice"},
	{"^(m|l)ice$", "${1}ice"},
	{"^(ox)$", "${1}en"},
	{"^(oxen)$", "${1}"},
	{"(quiz)$", "${1}zes"},
}

var singularInflections = [][]string{
	{"s$", ""},
	{"(ss)$", "${1}"},
	{"(n)ews$", "${1}ews"},
	{"([ti])a$", "${1}um"},
	{"((a)naly|(b)a|(d)iagno|(p)arenthe|(p)rogno|(s)ynop|(t)he)(sis|ses)$", "${1}sis"},
	{"(^analy)(sis|ses)$", "${1}sis"},
	{"([^f])ves$", "${1}fe"},
	{"(hive)s$", "${1}"},
	{"(tive)s$", "${1}"},
	{"([lr])ves$", "${1}f"},
	{"([^aeiouy]|qu)ies$", "${1}y"},
	{"(s)eries$", "${1}eries"},
	{"(m)ovies$", "${1}ovie"},
	{"(x|ch|ss|sh)es$", "${1}"},
	{"^(m|l)ice$", "${1}ouse"},
	{"(bus)(es)?$", "${1}"},
	{"(o)es$", "${1}"},
	{"(shoe)s$", "${1}"},
	{"(cris|test)(is|es)$", "${1}is"},
	{"^(a)x[ie]s$", "${1}xis"},
	{"(octop|vir)(us|i)$", "${1}us"},
	{"(alias|status)(es)?$", "${1}"},
	{"^(ox)en", "${1}"},
	{"(vert|ind)ices$", "${1}ex"},
	{"(matr)ices$", "${1}ix"},
	{"(quiz)zes$", "${1}"},
	{"(database)s$", "${1}"},
}

var irregularInflections = [][]string{
	{"person", "people"},
	{"man", "men"},
	{"child", "children"},
	{"sex", "sexes"},
	{"move", "moves"},
	{"mombie", "mombies"},
}

var uncountableInflections = []string{"equipment", "information", "rice", "money", "species", "series", "fish", "sheep", "jeans", "police"}

type inflection struct {
	regexp  *regexp.Regexp
	replace string
}

var compiledPluralMaps []inflection
var compiledSingularMaps []inflection

func compile() {
	compiledPluralMaps = []inflection{}
	compiledSingularMaps = []inflection{}
	for _, uncountable := range uncountableInflections {
		inf := inflection{
			regexp:  regexp.MustCompile("^(?i)(" + uncountable + ")$"),
			replace: "${1}",
		}
		compiledPluralMaps = append(compiledPluralMaps, inf)
		compiledSingularMaps = append(compiledSingularMaps, inf)
	}

	for _, value := range irregularInflections {
		infs := []inflection{
			{regexp: regexp.MustCompile(strings.ToUpper(value[0]) + "$"), replace: strings.ToUpper(value[1])},
			{regexp: regexp.MustCompile(strings.Title(value[0]) + "$"), replace: strings.Title(value[1])},
			{regexp: regexp.MustCompile(value[0] + "$"), replace: value[1]},
		}
		compiledPluralMaps = append(compiledPluralMaps, infs...)
	}

	for _, value := range irregularInflections {
		infs := []inflection{
			{regexp: regexp.MustCompile(strings.ToUpper(value[1]) + "$"), replace: strings.ToUpper(value[0])},
			{regexp: regexp.MustCompile(strings.Title(value[1]) + "$"), replace: strings.Title(value[0])},
			{regexp: regexp.MustCompile(value[1] + "$"), replace: value[0]},
		}
		compiledSingularMaps = append(compiledSingularMaps, infs...)
	}

	for i := len(pluralInflections) - 1; i >= 0; i-- {
		value := pluralInflections[i]
		infs := []inflection{
			{regexp: regexp.MustCompile(strings.ToUpper(value[0])), replace: strings.ToUpper(value[1])},
			{regexp: regexp.MustCompile(value[0]), replace: value[1]},
			{regexp: regexp.MustCompile("(?i)" + value[0]), replace: value[1]},
		}
		compiledPluralMaps = append(compiledPluralMaps, infs...)
	}

	for i := len(singularInflections) - 1; i >= 0; i-- {
		value := singularInflections[i]
		infs := []inflection{
			{regexp: regexp.MustCompile(strings.ToUpper(value[0])), replace: strings.ToUpper(value[1])},
			{regexp: regexp.MustCompile(value[0]), replace: value[1]},
			{regexp: regexp.MustCompile("(?i)" + value[0]), replace: value[1]},
		}
		compiledSingularMaps = append(compiledSingularMaps, infs...)
	}
}

func init() {
	compile()
}

func AddPlural(key, value string) {
	pluralInflections = append(pluralInflections, []string{key, value})
	compile()
}

func AddSingular(key, value string) {
	singularInflections = append(singularInflections, []string{key, value})
	compile()
}

func AddIrregular(key, value string) {
	irregularInflections = append(irregularInflections, []string{key, value})
	compile()
}

func AddUncountable(value string) {
	uncountableInflections = append(uncountableInflections, value)
	compile()
}

func Plural(str string) string {
	for _, inflection := range compiledPluralMaps {
		if inflection.regexp.MatchString(str) {
			return inflection.regexp.ReplaceAllString(str, inflection.replace)
		}
	}
	return str
}

func Singular(str string) string {
	for _, inflection := range compiledSingularMaps {
		if inflection.regexp.MatchString(str) {
			return inflection.regexp.ReplaceAllString(str, inflection.replace)
		}
	}
	return str
}
