package rules

import "merge-dsl/pkg/cursor"

var (
	RulesDefaultKey  = "rules"
	RulesDefaultName = "name"
)

func ConvertToRule(cursor cursor.Cursor[cursor.SchemaData], factory map[string]RuleFactory) cursor.Cursor[[]Rule] {
	if cursor == nil {
		return nil
	}
	return &Cursor{
		Cursor:  cursor,
		Factory: factory,
	}
}

func ConvertToRules(cursors []cursor.Cursor[cursor.SchemaData], factory map[string]RuleFactory) []cursor.Cursor[[]Rule] {
	result := make([]cursor.Cursor[[]Rule], len(cursors))
	for i, cursor := range cursors {
		result[i] = ConvertToRule(cursor, factory)
	}
	return result
}

type Cursor struct {
	Cursor  cursor.Cursor[cursor.SchemaData]
	Factory map[string]RuleFactory
}

func (c *Cursor) HasChildren() bool {
	return c.Cursor.HasChildren()
}

func (c *Cursor) Value() []Rule {
	rules := []Rule{}
	value := c.Cursor.Value()
	if rawRules, ok := value[RulesDefaultKey].([]interface{}); ok {
		for _, rawRule := range rawRules {
			if mapRule, ok := rawRule.(map[string]interface{}); ok {
				if name, ok := mapRule[RulesDefaultName].(string); ok {
					if ruleFactory, ok := c.Factory[name]; ok {
						rules = append(rules, ruleFactory(mapRule))
					}
				}
			}
		}
	}
	return rules
}

func (c *Cursor) GetKey(key string) cursor.Cursor[[]Rule] {
	return ConvertToRule(c.Cursor.GetKey(key), c.Factory)
}

func (c *Cursor) GetKeys() []string {
	return c.Cursor.GetKeys()
}

func (c *Cursor) GetItems() []cursor.Cursor[[]Rule] {
	return ConvertToRules(c.Cursor.GetItems(), c.Factory)
}

func (c *Cursor) GetDefault() cursor.Cursor[[]Rule] {
	return ConvertToRule(c.Cursor.GetDefault(), c.Factory)
}
