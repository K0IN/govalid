// Package rules implements validation rules for fields in structs.
package rules

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/codegen"

	"github.com/sivchari/govalid/internal/markers"
	"github.com/sivchari/govalid/internal/validator"
	"github.com/sivchari/govalid/internal/validator/registry"
)

type patternValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	pattern    string
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*patternValidator)(nil)

const patternKey = "%s-pattern"

func (v *patternValidator) Validate() string {
	fieldName := v.FieldName()
	// Use external helper function for regex pattern matching
	return fmt.Sprintf("!validationhelper.MatchPattern(%q, t.%s)", v.pattern, fieldName)
}

func (v *patternValidator) FieldName() string {
	return v.field.Names[0].Name
}

func (v *patternValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(v.structName, v.parentPath, v.FieldName())
}

func (v *patternValidator) Err() string {
	key := fmt.Sprintf(patternKey, v.FieldPath().CleanedPath())

	if validator.GeneratorMemory[key] {
		return ""
	}

	validator.GeneratorMemory[key] = true

	const deprecationNoticeTemplate = `
		// Deprecated: Use [@ERRVARIABLE]
		//
		// [@LEGACYERRVAR] is deprecated and is kept for compatibility purpose.
		[@LEGACYERRVAR] = [@ERRVARIABLE]
	`

	const errTemplate = `
		// [@ERRVARIABLE] is returned when the field does not match pattern [@PATTERN].
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason:"field [@FIELD] must match pattern [@PATTERN]",Path:"[@PATH]",Type:"[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sPatternValidation", v.structName, v.FieldName())
	currentErrVarName := v.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", v.FieldName(),
		"[@PATH]", v.FieldPath().String(),
		"[@PATTERN]", v.pattern,
		"[@TYPE]", v.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (v *patternValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]PatternValidation", "[@PATH]", v.FieldPath().CleanedPath())
}

func (v *patternValidator) Imports() []string {
	return []string{"github.com/sivchari/govalid/validation/validationhelper"}
}

// ValidatePattern creates a new pattern validator for string types.
func ValidatePattern(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Only works on string types
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	pattern, ok := input.Expressions[markers.GoValidMarkerPattern]
	if !ok || pattern == "" {
		return nil
	}

	return &patternValidator{
		pass:       input.Pass,
		field:      input.Field,
		pattern:    pattern,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}
