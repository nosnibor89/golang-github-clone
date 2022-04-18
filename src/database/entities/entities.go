package entities

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
)

type Attrs = map[string]*dynamodb.AttributeValue

type Entity struct {
	CreatedAt, UpdatedAt time.Time
}

type attrDefinition struct {
	pk, sk, typeLabel string
	extraAttrs        Attrs
}

func (attrDefinition *attrDefinition) getPrimaryKey() Attrs {
	return Attrs{
		"PK": {S: aws.String(attrDefinition.pk)},
		"SK": {S: aws.String(attrDefinition.sk)},
	}
}

func (attrDefinition *attrDefinition) getType() Attrs {
	return Attrs{
		"Type": {
			S: aws.String(attrDefinition.typeLabel),
		},
	}
}

func (attrDefinition *attrDefinition) withStringAttribute(name, value string) *attrDefinition {
	if len(attrDefinition.extraAttrs) == 0 {
		attrDefinition.extraAttrs = make(Attrs)
	}

	attrDefinition.extraAttrs[name] = &dynamodb.AttributeValue{
		S: aws.String(value),
	}
	return attrDefinition
}

func (attrDefinition *attrDefinition) withBoolAttribute(name string, value bool) *attrDefinition {
	if len(attrDefinition.extraAttrs) == 0 {
		attrDefinition.extraAttrs = make(Attrs)
	}

	attrDefinition.extraAttrs[name] = &dynamodb.AttributeValue{
		BOOL: aws.Bool(value),
	}
	return attrDefinition
}

func (attrDefinition *attrDefinition) withIntAttribute(name string, value string) *attrDefinition {
	if len(attrDefinition.extraAttrs) == 0 {
		attrDefinition.extraAttrs = make(Attrs)
	}

	attrDefinition.extraAttrs[name] = &dynamodb.AttributeValue{
		N: aws.String(value),
	}
	return attrDefinition
}

func (attrDefinition *attrDefinition) withSecondaryIndexKey(index int, partitionKeyValue, sortKeyValue string) *attrDefinition {
	secondaryPartitionKey := fmt.Sprintf("GS%dPK", index)
	secondarySortKey := fmt.Sprintf("GS%dSK", index)

	secondaryIndexKey := map[string]*dynamodb.AttributeValue{
		secondaryPartitionKey: {S: aws.String(partitionKeyValue)},
		secondarySortKey:      {S: aws.String(sortKeyValue)},
	}
	attrDefinition.extraAttrs = mergeAttrs(attrDefinition.extraAttrs, secondaryIndexKey)
	return attrDefinition
}

func (attrDefinition *attrDefinition) allAttributes() Attrs {
	primaryKeyAttrs := attrDefinition.getPrimaryKey()
	typeAttr := attrDefinition.getType()

	allAttrs := mergeAttrs(primaryKeyAttrs, typeAttr, attrDefinition.extraAttrs)
	return allAttrs
}

func parseTimeAttr(datetime string) time.Time {
	parsed, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		fmt.Printf("Count not parse time for value %v\n", datetime)
		return time.Date(1997, time.January, 1, 1, 1, 1, 1, time.UTC)
	}

	return parsed
}

func parseTimeItem(datetime time.Time) string {
	return datetime.Format(time.RFC3339)
}

func mergeAttrs[K string | int64, V any](maps ...map[K]V) map[K]V {
	newMap := make(map[K]V)

	for _, currentMap := range maps {
		for key, value := range currentMap {
			newMap[key] = value
		}
	}

	return newMap
}
