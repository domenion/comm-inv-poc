package models

type GetProductResponse struct {
	ID                        string         `json:"id"`
	Version                   int64          `json:"version"`
	VersionModifiedAt         string         `json:"versionModifiedAt"`
	LastMessageSequenceNumber int64          `json:"lastMessageSequenceNumber"`
	CreatedAt                 string         `json:"createdAt"`
	LastModifiedAt            string         `json:"lastModifiedAt"`
	LastModifiedBy            LastModifiedBy `json:"lastModifiedBy"`
	CreatedBy                 CreatedBy      `json:"createdBy"`
	ProductType               ProductType    `json:"productType"`
	MasterData                MasterData     `json:"masterData"`
	Key                       string         `json:"key"`
	TaxCategory               ProductType    `json:"taxCategory"`
	PriceMode                 string         `json:"priceMode"`
	LastVariantID             int64          `json:"lastVariantId"`
}

type CreatedBy struct {
	IsPlatformClient bool        `json:"isPlatformClient"`
	User             ProductType `json:"user"`
}

type ProductType struct {
	TypeID string `json:"typeId"`
	ID     string `json:"id"`
}

type LastModifiedBy struct {
	IsPlatformClient bool `json:"isPlatformClient"`
}

type MasterData struct {
	Current          Current `json:"current"`
	Staged           Staged  `json:"staged"`
	Published        bool    `json:"published"`
	HasStagedChanges bool    `json:"hasStagedChanges"`
}

type Current struct {
	Name               MetaDescriptionClass `json:"name"`
	Description        MetaDescriptionClass `json:"description"`
	Categories         []ProductType        `json:"categories"`
	CategoryOrderHints CategoryOrderHints   `json:"categoryOrderHints"`
	Slug               MetaDescriptionClass `json:"slug"`
	MetaTitle          MetaDescriptionClass `json:"metaTitle"`
	MetaDescription    MetaDescriptionClass `json:"metaDescription"`
	MasterVariant      MasterVariant        `json:"masterVariant"`
	Variants           []Variant            `json:"variants"`
	SearchKeywords     CategoryOrderHints   `json:"searchKeywords"`
}

type CategoryOrderHints struct {
}

type MetaDescriptionClass struct {
	En string `json:"en"`
	Th string `json:"th"`
}

type MasterVariant struct {
	ID         int64         `json:"id"`
	Prices     []interface{} `json:"prices"`
	Images     []interface{} `json:"images"`
	Attributes []interface{} `json:"attributes"`
	Assets     []interface{} `json:"assets"`
}

type Staged struct {
	Name               MetaDescriptionClass `json:"name"`
	Description        StagedDescription    `json:"description"`
	Categories         []ProductType        `json:"categories"`
	CategoryOrderHints CategoryOrderHints   `json:"categoryOrderHints"`
	Slug               MetaDescriptionClass `json:"slug"`
	MetaTitle          MetaDescriptionClass `json:"metaTitle"`
	MetaDescription    MetaDescriptionClass `json:"metaDescription"`
	MasterVariant      Variant              `json:"masterVariant"`
	Variants           []Variant            `json:"variants"`
	SearchKeywords     CategoryOrderHints   `json:"searchKeywords"`
}

type StagedDescription struct {
	En string `json:"en"`
}

type Variant struct {
	ID           int64         `json:"id"`
	Sku          string        `json:"sku"`
	Key          string        `json:"key"`
	Prices       []Price       `json:"prices"`
	Images       []Image       `json:"images"`
	Attributes   []interface{} `json:"attributes"`
	Assets       []interface{} `json:"assets"`
	Availability Availability  `json:"availability"`
}

type Availability struct {
	IsOnStock         bool   `json:"isOnStock"`
	AvailableQuantity int64  `json:"availableQuantity"`
	Version           int64  `json:"version"`
	ID                string `json:"id"`
}

type Image struct {
	URL        string     `json:"url"`
	Label      string     `json:"label"`
	Dimensions Dimensions `json:"dimensions"`
}

type Dimensions struct {
	W int64 `json:"w"`
	H int64 `json:"h"`
}

type Price struct {
	ID            string      `json:"id"`
	Country       string      `json:"country"`
	CustomerGroup ProductType `json:"customerGroup"`
	Channel       ProductType `json:"channel"`
}
