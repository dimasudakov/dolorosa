package nirvana_helper

import (
	"context"
	"nirvana/pkg/api/nirvana"
)

type NirvanaHelper struct {
	client nirvana.NirvanaClient
}

func NewNirvanaHelper(client nirvana.NirvanaClient) NirvanaHelper {
	return NirvanaHelper{
		client: client,
	}
}

func (n NirvanaHelper) CheckException(ctx context.Context, name string, attributes ExceptionAttributes) (bool, error) {
	found, err := n.client.CheckException(ctx, &nirvana.CheckExceptionRequest{
		Name: name,
		Attributes: &nirvana.ExceptionAttributes{
			ClientId: &attributes.ClientID,
			Amount:   &attributes.Amount,
		},
	})
	if err != nil {
		return false, err
	}
	return found.Found, err
}
