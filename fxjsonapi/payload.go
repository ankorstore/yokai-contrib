package fxjsonapi

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/google/jsonapi"
)

type MarshallParams struct {
	Payload         any
	WithoutIncluded bool
	Metadata        map[string]interface{}
}

func Marshall(params MarshallParams) ([]byte, error) {
	mp, err := jsonapi.Marshal(params.Payload)
	if err != nil {
		return nil, err
	}

	cast := false
	buf := bytes.Buffer{}

	if omp, ok := mp.(*jsonapi.OnePayload); ok {
		cast = true

		if params.WithoutIncluded {
			omp.Included = []*jsonapi.Node{}
		}

		if len(params.Metadata) > 0 {
			var meta jsonapi.Meta = params.Metadata

			omp.Meta = &meta
		}

		err = json.NewEncoder(&buf).Encode(omp)
		if err != nil {
			return nil, err
		}
	}

	if mmp, ok := mp.(*jsonapi.ManyPayload); ok {
		cast = true

		if params.WithoutIncluded {
			mmp.Included = []*jsonapi.Node{}
		}

		if len(params.Metadata) > 0 {
			var meta jsonapi.Meta = params.Metadata

			mmp.Meta = &meta
		}

		err = json.NewEncoder(&buf).Encode(mmp)
		if err != nil {
			return nil, err
		}
	}

	if !cast {
		return nil, errors.New("error casting marshalled payload")
	}

	return buf.Bytes(), nil
}
