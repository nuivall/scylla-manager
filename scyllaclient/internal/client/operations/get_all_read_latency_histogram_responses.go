// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/scylladb/mermaid/scyllaclient/internal/models"
)

// GetAllReadLatencyHistogramReader is a Reader for the GetAllReadLatencyHistogram structure.
type GetAllReadLatencyHistogramReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAllReadLatencyHistogramReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetAllReadLatencyHistogramOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetAllReadLatencyHistogramOK creates a GetAllReadLatencyHistogramOK with default headers values
func NewGetAllReadLatencyHistogramOK() *GetAllReadLatencyHistogramOK {
	return &GetAllReadLatencyHistogramOK{}
}

/*GetAllReadLatencyHistogramOK handles this case with default header values.

GetAllReadLatencyHistogramOK get all read latency histogram o k
*/
type GetAllReadLatencyHistogramOK struct {
	Payload []*models.RateMovingAverageAndHistogram
}

func (o *GetAllReadLatencyHistogramOK) Error() string {
	return fmt.Sprintf("[GET /column_family/metrics/read_latency/moving_average_histogram/][%d] getAllReadLatencyHistogramOK  %+v", 200, o.Payload)
}

func (o *GetAllReadLatencyHistogramOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
