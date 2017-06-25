// Copyright 2016 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package base

import (
	"fmt"

	"github.com/hexya-erp/hexya/hexya/models"
	"github.com/hexya-erp/hexya/hexya/models/types"
	"github.com/hexya-erp/hexya/pool"
)

func init() {
	partnerTitle := pool.PartnerTitle().DeclareModel()
	partnerTitle.AddCharField("Name", models.StringFieldParams{String: "Title", Required: true, Translate: true, Unique: true})
	partnerTitle.AddCharField("Shortcut", models.StringFieldParams{String: "Abbreviation", Translate: true})

	partnerCategory := pool.PartnerCategory().DeclareModel()
	partnerCategory.AddCharField("Name", models.StringFieldParams{String: "Category Name", Required: true, Translate: true})
	partnerCategory.AddIntegerField("Color", models.SimpleFieldParams{String: "Color Index"})
	partnerCategory.AddMany2OneField("Parent", models.ForeignKeyFieldParams{RelationModel: pool.PartnerCategory(),
		String: "Parent Tag", Index: true, OnDelete: models.Cascade})
	partnerCategory.AddCharField("CompleteName", models.StringFieldParams{String: "Full Name", Compute: "ComputeCompleteName"})
	partnerCategory.AddOne2ManyField("Children", models.ReverseFieldParams{RelationModel: pool.PartnerCategory(),
		ReverseFK: "Parent", String: "Children Tags"})
	partnerCategory.AddMany2ManyField("Partners", models.Many2ManyFieldParams{RelationModel: pool.Partner()})

	partnerCategory.AddMethod("ComputeCompleteName",
		`ComputeCompleteName returns the complete name of the tag with all the parents`,
		func(s pool.PartnerCategorySet) (*pool.PartnerCategoryData, []models.FieldNamer) {
			completeName := s.Name()
			for rs := s; !rs.Parent().IsEmpty(); rs = rs.Parent() {
				completeName = fmt.Sprintf("%s/%s", rs.Parent().Name(), completeName)
			}
			res := &pool.PartnerCategoryData{
				CompleteName: completeName,
			}
			return res, []models.FieldNamer{pool.PartnerCategory().CompleteName()}
		})

	partner := pool.Partner().DeclareModel()
	partner.AddCharField("Name", models.StringFieldParams{Required: true, Index: true, NoCopy: true})
	partner.AddDateField("Date", models.SimpleFieldParams{})
	partner.AddMany2OneField("Title", models.ForeignKeyFieldParams{RelationModel: pool.PartnerTitle()})
	partner.AddMany2OneField("Parent", models.ForeignKeyFieldParams{RelationModel: pool.Partner()})
	partner.AddOne2ManyField("Children", models.ReverseFieldParams{RelationModel: pool.Partner(), ReverseFK: "Parent"})
	partner.AddCharField("Ref", models.StringFieldParams{})
	partner.AddCharField("Lang", models.StringFieldParams{})
	partner.AddCharField("TZ", models.StringFieldParams{})
	partner.AddCharField("TZOffset", models.StringFieldParams{})
	partner.AddMany2OneField("User", models.ForeignKeyFieldParams{RelationModel: pool.User()})
	partner.AddCharField("VAT", models.StringFieldParams{})
	//Banks            []*PartnerBank
	partner.AddCharField("Website", models.StringFieldParams{})
	partner.AddCharField("Comment", models.StringFieldParams{})
	partner.AddMany2ManyField("Categories", models.Many2ManyFieldParams{RelationModel: pool.PartnerCategory()})
	partner.AddFloatField("CreditLimit", models.FloatFieldParams{})
	partner.AddCharField("EAN13", models.StringFieldParams{})
	partner.AddBooleanField("Active", models.SimpleFieldParams{})
	partner.AddBooleanField("Customer", models.SimpleFieldParams{})
	partner.AddBooleanField("Supplier", models.SimpleFieldParams{})
	partner.AddBooleanField("Employee", models.SimpleFieldParams{})
	partner.AddCharField("Function", models.StringFieldParams{})
	partner.AddSelectionField("Type", models.SelectionFieldParams{Selection: types.Selection{
		"contact": "Contact", "invoice": "Invoice Address", "delivery": "Shipping Address", "other": "Other Address"},
		Help:    "Used to select automatically the right address according to the context in sales and purchases documents.",
		Default: models.DefaultValue("contact"),
	})
	partner.AddCharField("Street", models.StringFieldParams{})
	partner.AddCharField("Street2", models.StringFieldParams{})
	partner.AddCharField("Zip", models.StringFieldParams{})
	partner.AddCharField("City", models.StringFieldParams{})
	partner.AddMany2OneField("State", models.ForeignKeyFieldParams{RelationModel: pool.CountryState(),
		Filter: pool.CountryState().Country().EqualsEval("country_id")})
	partner.AddMany2OneField("Country", models.ForeignKeyFieldParams{RelationModel: pool.Country()})
	partner.AddCharField("Email", models.StringFieldParams{})
	partner.AddCharField("Phone", models.StringFieldParams{})
	partner.AddCharField("Fax", models.StringFieldParams{})
	partner.AddCharField("Mobile", models.StringFieldParams{})
	partner.AddDateField("Birthdate", models.SimpleFieldParams{})
	partner.AddBooleanField("IsCompany", models.SimpleFieldParams{Compute: "ComputeIsCompany", Stored: true, Depends: []string{"CompanyType"}})
	partner.AddBooleanField("UseParentAddress", models.SimpleFieldParams{})
	partner.AddBinaryField("Image", models.SimpleFieldParams{})
	partner.AddBinaryField("ImageMedium", models.SimpleFieldParams{})
	partner.AddSelectionField("CompanyType", models.SelectionFieldParams{Selection: types.Selection{"person": "Individual", "company": "Company"},
		OnChange: "ComputeIsCompany", Default: models.DefaultValue("person")})
	partner.AddMany2OneField("Company", models.ForeignKeyFieldParams{RelationModel: pool.Company()})
	partner.AddIntegerField("Color", models.SimpleFieldParams{})
	partner.AddOne2ManyField("Users", models.ReverseFieldParams{RelationModel: pool.User(), ReverseFK: "Partner"})

	partner.Methods().ComputeIsCompany().DeclareMethod(
		`ComputeIsCompany computes the IsCompany field from the selected CompanyType`,
		func(s pool.PartnerSet) (*pool.PartnerData, []models.FieldNamer) {
			var res pool.PartnerData
			res.IsCompany = s.CompanyType() == "company"
			return &res, []models.FieldNamer{pool.Partner().IsCompany()}
		})

	partner.AddSQLConstraint("check_name", "CHECK( (type='contact' AND name IS NOT NULL) or (type != 'contact') )", "Contacts require a name.")
	//'has_image': fields.function(_has_image, type="boolean"),
	//'contact_address': fields.function(_address_display,  type='char', string='Complete Address'),
	//'commercial_partner_id': fields.function(_commercial_partner_id, type='many2one', relation='res.partner', string='Commercial Entity', store=_commercial_partner_store_triggers)

}