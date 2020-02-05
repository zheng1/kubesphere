package meteringconfig

import "time"

// TODO: use CRD

type MeteringStatus struct {
	AccessSystemId string
	CatalogId      string
	Products       map[string]Product
}

func (m MeteringStatus) GetProduct(productName string) Product {
	return m.Products[productName]
}

type Product struct {
	ResourceName string // same as resource name
	Attributes   []Attribute

	ProductId string
}

type Attribute struct {
	Name            string
	Period          time.Duration
	IsMeteringValue bool
	MetricName      string

	AttributeId string
}

var CurrentStatus = MeteringStatus{
	AccessSystemId: "sys_37AAB8JNq27M",
	CatalogId:      "cata_99jov8BGwvY6",
	Products: map[string]Product{
		"pod": {
			ResourceName: "pod",
			ProductId:    "prd_VGVgk9PBDq9p",
			Attributes: []Attribute{
				{
					Name:            "pod_cpu_usage",
					IsMeteringValue: true,
					Period:          1 * time.Minute,
					MetricName:      "pod_cpu_usage",

					AttributeId: "attr_jz9Lk032poG2",
				},
			},
		},
		"deployment": {
			ResourceName: "deployment",
			ProductId:    "prd_N15k2rM8oRJO",
			Attributes: []Attribute{
				{
					Name:            "pod_cpu_usage",
					IsMeteringValue: true,
					Period:          1 * time.Minute,
					MetricName:      "pod_cpu_usage",

					AttributeId: "attr_14R9PEZQy9JV",
				},
			},
		},
	},
}
