package entities

import (
	"fmt"
	"github-clone/src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Attrs = map[string]*dynamodb.AttributeValue

type Entity interface {
	ToItem() Attrs
}

//func (k key) keyToItemMap() Attrs {
//	keyAttrs := make(Attrs)
//	err := mapstructure.Decode(k, &keyAttrs)
//	if err != nil {
//		panic(err)
//	}
//
//	return keyAttrs
//}

//func newAttrDefinition() attrDefinition {
//
//}

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

func (attrDefinition *attrDefinition) withSecondaryIndexKey(index int, partitionKeyValue, sortKeyValue string) *attrDefinition {
	secondaryPartitionKey := fmt.Sprintf("GS%dPK", index)
	secondarySortKey := fmt.Sprintf("GS%dSK", index)

	secondaryIndexKey := map[string]*dynamodb.AttributeValue{
		secondaryPartitionKey: {S: aws.String(partitionKeyValue)},
		secondarySortKey:      {S: aws.String(sortKeyValue)},
	}
	attrDefinition.extraAttrs = util.MergeMaps(attrDefinition.extraAttrs, secondaryIndexKey)
	return attrDefinition
}

func (attrDefinition *attrDefinition) allAttributes() Attrs {
	primaryKeyAttrs := attrDefinition.getPrimaryKey()
	typeAttr := attrDefinition.getType()

	allAttrs := util.MergeMaps(primaryKeyAttrs, typeAttr, attrDefinition.extraAttrs)
	return allAttrs
}
