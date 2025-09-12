package errors

import (
	"fmt"
	"strings"
)

//===========================================================================
// Configuration Validation Methods
//===========================================================================

func ConfigError(err error, errs ...*InvalidConfiguration) ConfigurationErrors {
	var cerrs ConfigurationErrors
	if err == nil {
		cerrs = make(ConfigurationErrors, 0, len(errs))
	} else {
		var ok bool
		if cerrs, ok = err.(ConfigurationErrors); !ok {
			cerrs = make(ConfigurationErrors, 0, len(errs)+1)
			cerrs = append(cerrs, &InvalidConfiguration{issue: err.Error()})
		}
	}

	for _, e := range errs {
		if e != nil {
			cerrs = append(cerrs, e)
		}
	}

	if len(cerrs) == 0 {
		return nil
	}
	return cerrs
}

func RequiredConfig(conf, field string) *InvalidConfiguration {
	return &InvalidConfiguration{
		conf:  conf,
		field: field,
		issue: "is required but not set",
	}
}

func InvalidConfig(conf, field, issue string, args ...any) *InvalidConfiguration {
	return &InvalidConfiguration{
		conf:  conf,
		field: field,
		issue: fmt.Sprintf(issue, args...),
	}
}

func ConfigParseError(conf, field string, err error) *InvalidConfiguration {
	return &InvalidConfiguration{
		conf:  conf,
		field: field,
		issue: fmt.Sprintf("could not parse %s: %s", field, err.Error()),
	}
}

//===========================================================================
// Configuration Validation Errors
//===========================================================================

type InvalidConfiguration struct {
	conf  string
	field string
	issue string
}

func (e *InvalidConfiguration) Error() string {
	field := e.field
	if e.conf != "" {
		field = e.conf + "." + e.field
	}

	return fmt.Sprintf("invalid configuration: %s %s", field, e.issue)
}

func (e *InvalidConfiguration) Field() string {
	return e.field
}

type ConfigurationErrors []*InvalidConfiguration

func (e ConfigurationErrors) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}

	errs := make([]string, 0, len(e))
	for _, err := range e {
		errs = append(errs, err.Error())
	}

	return fmt.Sprintf("%d configuration errors occurred:\n  %s", len(e), strings.Join(errs, "\n  "))
}

func (e ConfigurationErrors) Errors() []*InvalidConfiguration {
	return e
}
