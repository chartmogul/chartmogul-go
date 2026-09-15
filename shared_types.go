package chartmogul

type Custom struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// AssociatedObjectIdentifier points a task or note at the customer or contact it belongs to.
type AssociatedObjectIdentifier struct {
	AssociatedObject string `json:"associated_object"` // AssociatedObjectCustomer or AssociatedObjectContact
	Method           string `json:"method"`            // AssociatedObjectIdentifierMethodUUID
	Value            string `json:"value"`
}

const (
	AssociatedObjectCustomer             = "customer"
	AssociatedObjectContact              = "contact"
	AssociatedObjectIdentifierMethodUUID = "uuid"
)
