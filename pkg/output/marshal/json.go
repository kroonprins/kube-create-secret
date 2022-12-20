package marshal

import (
	"encoding/json"

	"github.com/kroonprins/kube-create-secret/pkg/constants"
	"github.com/kroonprins/kube-create-secret/pkg/core"
	"github.com/kroonprins/kube-create-secret/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubernetesResourceList[T interface{}] struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []T `json:"items"`
}

type JsonMarshaller struct {
}

func NewJsonMarshaller() *JsonMarshaller {
	return &JsonMarshaller{}
}

func (*JsonMarshaller) Marshal(config core.Config, items []interface{}) (bool, []byte, error) {
	if !(len(config.OutputFormats) == 0 && config.InputFormat == types.JSON) &&
		!(len(config.OutputFormats) == 1 && config.OutputFormats[0] == types.JSON) {
		return true, nil, nil
	}

	res, err := marshalJson(items)
	if err != nil {
		return false, nil, err
	}
	return false, res, nil
}

func marshalJson[T interface{}](items []T) ([]byte, error) {
	var toMarshal any
	if len(items) == 1 {
		toMarshal = items[0]
	} else {
		toMarshal = KubernetesResourceList[T]{
			TypeMeta: metav1.TypeMeta{
				Kind:       constants.SECRET_TEMPLATE_LIST_KIND,
				APIVersion: constants.SECRET_TEMPLATE_LIST_API_VERSION,
			},
			Items: items,
		}
	}

	bytes, err := json.MarshalIndent(toMarshal, "", "    ")
	if err != nil {
		return nil, err
	}
	bytes = append(bytes, []byte("\n")...)
	return bytes, nil
}
