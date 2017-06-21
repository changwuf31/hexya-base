// Copyright 2016 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package base

import (
	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool"
)

func init() {
	currencyRate := pool.CurrencyRate().DeclareModel()
	currencyRate.AddDateTimeField("Name", models.SimpleFieldParams{String: "Date", Required: true, Index: true})
	currencyRate.AddFloatField("Rate", models.FloatFieldParams{Digits: types.Digits{Precision: 12, Scale: 6}, Help: "The rate of the currency to the currency of rate 1"})
	currencyRate.AddMany2OneField("Currency", models.ForeignKeyFieldParams{RelationModel: pool.Currency()})
	currencyRate.AddMany2OneField("Company", models.ForeignKeyFieldParams{RelationModel: pool.Company()})

	currency := pool.Currency().DeclareModel()
	currency.AddCharField("Name", models.StringFieldParams{String: "Currency", Help: "Currency Code [ISO 4217]", Size: 3, Unique: true})
	currency.AddCharField("Symbol", models.StringFieldParams{Help: "Currency sign, to be used when printing amounts", Size: 4})
	currency.AddFloatField("Rate", models.FloatFieldParams{String: "Current Rate", Help: "The rate of the currency to the currency of rate 1", Digits: types.Digits{Precision: 12, Scale: 6}}) //, Compute: "ComputeCurrentRate"})
	currency.AddOne2ManyField("Rates", models.ReverseFieldParams{RelationModel: pool.CurrencyRate(), ReverseFK: "Currency"})
	currency.AddFloatField("Rounding", models.FloatFieldParams{String: "Rounding Factor", Digits: types.Digits{Precision: 12, Scale: 6}})
	currency.AddIntegerField("DecimalPlaces", models.SimpleFieldParams{}) //Compute: "ComputeDecimalPlaces"})
	currency.AddBooleanField("Active", models.SimpleFieldParams{})
	currency.AddSelectionField("Position", models.SelectionFieldParams{Selection: types.Selection{"after": "After Amount", "before": "Before Amount"}, String: "Symbol Position", Help: "Determines where the currency symbol should be placed after or before the amount."})
}
