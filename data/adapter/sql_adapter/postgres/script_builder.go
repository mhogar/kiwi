package postgres

import (
	"fmt"
	"strings"

	"github.com/mhogar/kiwi/data/adapter"
	"github.com/mhogar/kiwi/data/query"
)

type ScriptBuilder struct{}

// Build a select query using the reflection model.
// Note that the model's name and fields are obtained using reflection and therefore sql injection is not possible.
func (s ScriptBuilder) BuildSelectQuery(model adapter.ReflectModel, where *query.WhereClause) string {
	script := `
		SELECT t1."%s"
			FROM "%s" t1
		WHERE %s
	`

	return fmt.Sprintf(
		script, strings.Join(model.Fields, `", t1."`), model.Name, s.buildWhereString(model, where),
	)
}

func (ScriptBuilder) buildWhereString(_ adapter.ReflectModel, _ *query.WhereClause) string {
	//TODO: implement where clause
	return ""
}

// Build an insert statement using the reflection model.
// Note that the model's name and fields are obtained using reflection and therefore sql injection is not possible.
func (s ScriptBuilder) BuildInsertStatement(model adapter.ReflectModel) string {
	script := `
		INSERT INTO "%s" ("%s")
			VALUES (%s)
	`

	return fmt.Sprintf(
		script, model.Name, strings.Join(model.Fields, `", "`), s.buildParametrizedString(model),
	)
}

func (ScriptBuilder) buildParametrizedString(model adapter.ReflectModel) string {
	params := []string{}
	for index := range model.Fields {
		params = append(params, fmt.Sprintf("$%d", index+1))
	}
	return strings.Join(params, ", ")
}

// Build an update statement using the reflection model.
// Note that the model's name and fields are obtained using reflection and therefore sql injection is not possible.
func (s ScriptBuilder) BuildUpdateStatement(model adapter.ReflectModel) string {
	script := `
		UPDATE "%s" SET
			%s
		WHERE "%s" = $1
	`

	return fmt.Sprintf(
		script, model.Name, s.buildSetString(model), model.UniqueField(),
	)
}

func (ScriptBuilder) buildSetString(model adapter.ReflectModel) string {
	setScripts := []string{}
	for index, field := range model.Fields {
		if index == 0 {
			continue
		}

		setScripts = append(setScripts,
			fmt.Sprintf(`"%s" = $%d`, field, index+1),
		)
	}

	return strings.Join(setScripts, `", "`)
}

// Build a delete statement using the reflection model.
// Note that the model's name and fields are obtained using reflection and therefore sql injection is not possible.
func (ScriptBuilder) BuildDeleteStatement(model adapter.ReflectModel) string {
	script := `
		DELETE FROM "%s" t1
			WHERE t1."%s" = $1
	`

	return fmt.Sprintf(
		script, model.Name, model.UniqueField(),
	)
}