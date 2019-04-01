package database

import (
	"reflect"
	"time"
	"github.com/chromedp/cdproto"
	"github.com/analogj/lantern/common/models"
	"github.com/mitchellh/mapstructure"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"fmt"
	"encoding/json"
	"errors"
)

func stringToDateTimeHook(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {
	if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
		return time.Parse(time.RFC3339, data.(string))
	}

	return data, nil
}


func processDatabaseEvent(dbType string, dbEvent map[string]interface{}) (cdproto.Message, interface{}, error) {
	event := cdproto.Message{}


	switch dbType {

	case "requests":
		request := models.DbRequest{}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: stringToDateTimeHook,
			Result: &request,

		})
		if err != nil {
			return event, nil,  err
		}
		err = decoder.Decode(dbEvent)
		if err != nil {
			return event, nil, err
		}

		event.Method = cdproto.EventNetworkRequestWillBeSent

		timestamp := cdp.MonotonicTime(request.CreatedAt)
		walltime := cdp.TimeSinceEpoch(request.RequestedOn)
		initiator := network.Initiator{
			Type: network.InitiatorTypeOther,
		}

		payload := network.EventRequestWillBeSent {
			RequestID: network.RequestID(fmt.Sprint(request.Id)),
			LoaderID:        cdp.LoaderID(""), // Loader identifier. Empty string if the request is fetched from worker.
			DocumentURL:     request.Url, // URL of the document this request is loaded for.
			Request: &network.Request {
				URL:              request.Url, // Request URL (without fragment).
				//URLFragment      string                    `json:"urlFragment,omitempty"`      // Fragment of the requested URL starting with hash, if present.
				Method: request.Method, // HTTP request method.
				Headers: network.Headers(request.Headers), // HTTP request headers.
				PostData: request.Body, // HTTP POST request data.
				HasPostData: len(request.Body) > 0, // True when the request has POST data. Note that postData might still be omitted when this flag is true when the data is too long.
				//MixedContentType security.MixedContentType `json:"mixedContentType,omitempty"` // The mixed content type of the request.
				//InitialPriority  ResourcePriority          `json:"initialPriority"`            // Priority of the resource request at the time request is sent.
				//ReferrerPolicy: // The referrer policy of the request, as defined in https://www.w3.org/TR/referrer-policy/
			}, // Request data.
			Timestamp:        &timestamp, // Timestamp.
			WallTime:         &walltime, // Timestamp.
			Initiator:       &initiator, // Request initiator.
		}

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return event, nil, err
		}

		event.Params = jsonBytes
		return event, request, nil



	case "responses":
		response := models.DbResponse{}
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: stringToDateTimeHook,
			Result: &response,

		})
		if err != nil {
			return event, nil, err
		}
		err = decoder.Decode(dbEvent)
		if err != nil {
			return event, nil, err
		}

		event.Method = cdproto.EventNetworkResponseReceived


		timestamp := cdp.MonotonicTime(response.CreatedAt)
		//walltime := cdp.TimeSinceEpoch(response.RespondedOn)
		//initiator := network.Initiator{
		//	Type: network.InitiatorTypeOther,
		//}

		//Type      ResourceType       `json:"type"`              // Resource type.
		//Response  *Response          `json:"response"`          // Response data.
		//FrameID   cdp.FrameID        `json:"frameId,omitempty"` // Frame identifier.


		payload := network.EventResponseReceived {


			RequestID: network.RequestID(fmt.Sprint(response.RequestId)),
			LoaderID:        cdp.LoaderID(""), // Loader identifier. Empty string if the request is fetched from worker.
			Timestamp:        &timestamp, // Timestamp.
			Type:         network.ResourceTypeOther, // Resource type.
			Response:       &network.Response{
				//URL: response.           `json:"url"`                          // Response URL. This URL can be different from CachedResource.url in case of redirect.
				Status: int64(response.StatusCode), // HTTP response status code.
				StatusText: response.Status, // HTTP response status text.
				Headers: network.Headers(response.Headers), // HTTP response headers.
				MimeType: response.MimeType,  // Resource mimeType as determined by the browser.
				//RequestHeaders     Headers          `json:"requestHeaders,omitempty"`     // Refined HTTP request headers that were actually transmitted over the network.
				//RequestHeadersText string           `json:"requestHeadersText,omitempty"` // HTTP request headers text.
				ConnectionReused: false,
				//ConnectionID       float64          `json:"connectionId"`                 // Physical connection id that was actually used for this request.
				//RemoteIPAddress    string           `json:"remoteIPAddress,omitempty"`    // Remote IP address.
				//RemotePort         int64            `json:"remotePort,omitempty"`         // Remote port.
				//FromDiskCache      bool             `json:"fromDiskCache,omitempty"`      // Specifies that the request was served from the disk cache.
				//FromServiceWorker  bool             `json:"fromServiceWorker,omitempty"`  // Specifies that the request was served from the ServiceWorker.
				EncodedDataLength:  float64(len(response.Body)), // Total number of bytes received for this request so far.
				//Timing:             *ResourceTiming  `json:"timing,omitempty"`             // Timing information for the given request.
				//Protocol           string           `json:"protocol,omitempty"`           // Protocol used to fetch this request.
				//SecurityState      security.State   `json:"securityState"`                // Security state of the request resource.
				//SecurityDetails    *SecurityDetails `json:"securityDetails,omitempty"`    // Security details for the request.
			},
		}

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return event, nil, err
		}

		event.Params = jsonBytes
		return event, response, nil


	default:
		return event, nil, errors.New("unknown DB event type")

	}
}
