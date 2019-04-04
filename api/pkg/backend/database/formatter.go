package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/analogj/lantern/common/models"
	"github.com/chromedp/cdproto"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/jinzhu/gorm"
)


func detectResourceType(mimeType string) network.ResourceType {
	switch mimeType {
	case "application/octet-stream":
		return network.ResourceTypeOther
	case "text/xml; charset=utf-8",
		"application/pdf",
		"application/postscript",
		"text/plain; charset=utf-16be",
		"text/plain; charset=utf-16le",
		"text/plain; charset=utf-8":
		return network.ResourceTypeDocument

	case "image/gif",
		"image/png",
		"image/jpeg",
		"image/bmp",
		"image/webp",
		"image/vnd.microsoft.icon":
		return network.ResourceTypeImage

	case "audio/aiff",
		"audio/basic",
		"audio/midi",
		"audio/mpeg",
		"application/ogg",
		"video/avi",
		"video/webm",
		"video/mp4":
		return network.ResourceTypeMedia

	case "application/vnd.ms-fontobject",
		"application/font-ttf",
		"application/font-off",
		"application/font-cff",
		"application/font-woff":
		return network.ResourceTypeFont
	case "application/x-rar-compressed",
		"application/zip",
		"application/x-gzip":
		return network.ResourceTypeOther

	case "text/html; charset=utf-8":
		return network.ResourceTypeDocument

	default:
		return network.ResourceTypeOther
	}

}


func processDatabaseEvent(orm *gorm.DB, dbType string, dbId int) (cdproto.Message, interface{}, error) {
	event := cdproto.Message{}

	switch dbType {

	case "requests":

		request := models.DbRequest{}
		orm.First(&request, "id = ?", dbId)

		event.Method = cdproto.EventNetworkRequestWillBeSent

		timestamp := cdp.MonotonicTime(request.CreatedAt)
		walltime := cdp.TimeSinceEpoch(request.RequestedOn)
		initiator := network.Initiator{
			Type: network.InitiatorTypeOther,
		}

		payload := network.EventRequestWillBeSent{
			RequestID:   network.RequestID(fmt.Sprint(request.Id)),
			LoaderID:    cdp.LoaderID(""), // Loader identifier. Empty string if the request is fetched from worker.
			DocumentURL: request.Url,      // URL of the document this request is loaded for.
			Request: &network.Request{
				URL: request.Url, // Request URL (without fragment).
				//URLFragment      string                    `json:"urlFragment,omitempty"`      // Fragment of the requested URL starting with hash, if present.
				Method:      request.Method,                   // HTTP request method.
				Headers:     network.Headers(request.Headers), // HTTP request headers.
				PostData:    request.Body,                     // HTTP POST request data.
				HasPostData: len(request.Body) > 0,            // True when the request has POST data. Note that postData might still be omitted when this flag is true when the data is too long.
				//MixedContentType security.MixedContentType `json:"mixedContentType,omitempty"` // The mixed content type of the request.
				//InitialPriority  ResourcePriority          `json:"initialPriority"`            // Priority of the resource request at the time request is sent.
				//ReferrerPolicy: // The referrer policy of the request, as defined in https://www.w3.org/TR/referrer-policy/
			}, // Request data.
			Timestamp: &timestamp, // Timestamp.
			WallTime:  &walltime,  // Timestamp.
			Initiator: &initiator, // Request initiator.
		}

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return event, nil, err
		}

		event.Params = jsonBytes
		return event, request, nil

	case "responses":
		response := models.DbResponse{}
		orm.First(&response, "id = ?", dbId)

		request := models.DbRequest{}
		orm.First(&request, "id = ?", response.RequestId)

		event.Method = cdproto.EventNetworkResponseReceived
		timestamp := cdp.MonotonicTime(response.CreatedAt)
		//walltime := cdp.TimeSinceEpoch(response.RespondedOn)
		//initiator := network.Initiator{
		//	Type: network.InitiatorTypeOther,
		//}

		//Type      ResourceType       `json:"type"`              // Resource type.
		//Response  *Response          `json:"response"`          // Response data.
		//FrameID   cdp.FrameID        `json:"frameId,omitempty"` // Frame identifier.


		//determine the correct resource type:
		resourceType := detectResourceType(response.MimeType)


		payload := network.EventResponseReceived{

			RequestID: network.RequestID(fmt.Sprint(response.RequestId)),
			LoaderID:  cdp.LoaderID(""),          // Loader identifier. Empty string if the request is fetched from worker.
			Timestamp: &timestamp,                // Timestamp.
			Type:      resourceType, // Resource type.
			Response: &network.Response{
				URL: 		request.Url,                          // Response URL. This URL can be different from CachedResource.url in case of redirect.
				Status:     int64(response.StatusCode),        // HTTP response status code.
				StatusText: response.Status,                   // HTTP response status text.
				Headers:    network.Headers(response.Headers), // HTTP response headers.
				MimeType:   response.MimeType,                 // Resource mimeType as determined by the browser.
				RequestHeaders: network.Headers(request.Headers),     // Refined HTTP request headers that were actually transmitted over the network.
				//RequestHeadersText string           `json:"requestHeadersText,omitempty"` // HTTP request headers text.
				ConnectionReused: false,
				//ConnectionID       float64          `json:"connectionId"`                 // Physical connection id that was actually used for this request.
				//RemoteIPAddress    string           `json:"remoteIPAddress,omitempty"`    // Remote IP address.
				//RemotePort         int64            `json:"remotePort,omitempty"`         // Remote port.
				//FromDiskCache      bool             `json:"fromDiskCache,omitempty"`      // Specifies that the request was served from the disk cache.
				//FromServiceWorker  bool             `json:"fromServiceWorker,omitempty"`  // Specifies that the request was served from the ServiceWorker.
				EncodedDataLength: float64(len(response.Body)), // Total number of bytes received for this request so far.
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
