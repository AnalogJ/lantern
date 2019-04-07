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
	"github.com/chromedp/cdproto/security"
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
	switch dbType {

	case "requests":

		request := models.DbRequest{}
		orm.First(&request, "id = ?", dbId)

		message, err := generateRequestWillBeSentEvent(request)
		return message, request, err

	case "responses":
		response := models.DbResponse{}
		orm.First(&response, "id = ?", dbId)

		request := models.DbRequest{}
		orm.First(&request, "id = ?", response.RequestId)


		message, err := generateNetworkResponseReceived(request, response)
		return message, response, err

	default:
		return cdproto.Message{}, nil, errors.New("unknown DB event type")

	}
}


func generateRequestWillBeSentEvent(dbRequest models.DbRequest) (cdproto.Message, error){
	event := cdproto.Message{
		Method: cdproto.EventNetworkRequestWillBeSent,
	}

	timestamp := cdp.MonotonicTime(dbRequest.CreatedAt)
	walltime := cdp.TimeSinceEpoch(dbRequest.RequestedOn)
	initiator := network.Initiator{
		Type: network.InitiatorTypeOther,
	}

	payload := network.EventRequestWillBeSent{
		RequestID:   network.RequestID(fmt.Sprint(dbRequest.Id)),
		LoaderID:    cdp.LoaderID(""), // Loader identifier. Empty string if the request is fetched from worker.
		DocumentURL: dbRequest.Url,      // URL of the document this request is loaded for.
		Request: &network.Request{
			URL: dbRequest.Url, // Request URL (without fragment).
			//URLFragment      string                    `json:"urlFragment,omitempty"`      // Fragment of the requested URL starting with hash, if present.
			Method:      dbRequest.Method,                   // HTTP request method.
			Headers:     network.Headers(dbRequest.Headers), // HTTP request headers.
			PostData:    dbRequest.Body,                     // HTTP POST request data.
			HasPostData: len(dbRequest.Body) > 0,            // True when the request has POST data. Note that postData might still be omitted when this flag is true when the data is too long.
			//MixedContentType security.MixedContentType `json:"mixedContentType,omitempty"` // The mixed content type of the request.
			//InitialPriority  ResourcePriority          `json:"initialPriority"`            // Priority of the resource request at the time request is sent.
			//ReferrerPolicy: // The referrer policy of the request, as defined in https://www.w3.org/TR/referrer-policy/
		}, // Request data.
		Timestamp: &timestamp, // Timestamp.
		WallTime:  &walltime,  // Timestamp.
		Initiator: &initiator, // Request initiator.
		//RedirectResponse *Response           `json:"redirectResponse,omitempty"` // Redirect response data.
		Type: network.ResourceTypeOther,       //`json:"type,omitempty"`             // Type of this resource.
		//FrameID          cdp.FrameID         `json:"frameId,omitempty"`          // Frame identifier.
		//HasUserGesture   bool                `json:"hasUserGesture,omitempty"`   // Whether the request is initiated by a user gesture. Defaults to false.
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return event, err
	}

	event.Params = jsonBytes
	return event, nil
}

func generateNetworkResponseReceived(dbRequest models.DbRequest, dbResponse models.DbResponse) (cdproto.Message, error) {
	message := cdproto.Message {
		Method: cdproto.EventNetworkResponseReceived,
	}
	timestamp := cdp.MonotonicTime(dbResponse.CreatedAt)
	//walltime := cdp.TimeSinceEpoch(response.RespondedOn)
	//initiator := network.Initiator{
	//	Type: network.InitiatorTypeOther,
	//}

	//Type      ResourceType       `json:"type"`              // Resource type.
	//Response  *Response          `json:"response"`          // Response data.
	//FrameID   cdp.FrameID        `json:"frameId,omitempty"` // Frame identifier.


	//determine the correct resource type:
	resourceType := detectResourceType(dbResponse.MimeType)


	payload := network.EventResponseReceived{

		RequestID: network.RequestID(fmt.Sprint(dbResponse.RequestId)),
		LoaderID:  cdp.LoaderID(""),          // Loader identifier. Empty string if the request is fetched from worker.
		Timestamp: &timestamp,                // Timestamp.
		Type:      resourceType, // Resource type.
		Response: &network.Response{
			URL: 		dbRequest.Url,                          // Response URL. This URL can be different from CachedResource.url in case of redirect.
			Status:     int64(dbResponse.StatusCode),        // HTTP response status code.
			StatusText: dbResponse.Status,                   // HTTP response status text.
			Headers:    network.Headers(dbResponse.Headers), // HTTP response headers.
			MimeType:   dbResponse.MimeType,                 // Resource mimeType as determined by the browser.
			RequestHeaders: network.Headers(dbRequest.Headers),     // Refined HTTP request headers that were actually transmitted over the network.
			//RequestHeadersText string           `json:"requestHeadersText,omitempty"` // HTTP request headers text.
			ConnectionReused: false,
			//ConnectionID       float64          `json:"connectionId"`                 // Physical connection id that was actually used for this request.
			//RemoteIPAddress    string           `json:"remoteIPAddress,omitempty"`    // Remote IP address.
			//RemotePort         int64            `json:"remotePort,omitempty"`         // Remote port.
			FromDiskCache: false,     // Specifies that the request was served from the disk cache.
			FromServiceWorker: false, // Specifies that the request was served from the ServiceWorker.
			EncodedDataLength: float64(len(dbResponse.Body)), // Total number of bytes received for this request so far.
			//Timing:             &network.ResourceTiming{
			//	RequestTime: float64(dbRequest.RequestedOn.Unix()),
			//
			//},          // Timing information for the given request.
			Protocol: dbResponse.Protocol,           // Protocol used to fetch this request.
			SecurityState:      security.StateUnknown,                // Security state of the request resource.
			//SecurityDetails    *SecurityDetails `json:"securityDetails,omitempty"`    // Security details for the request.
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return message, err
	}

	message.Params = jsonBytes
	return message, err
}


func generateDataReceivedEvent(dbResponse models.DbResponse) (cdproto.Message, error){
	eventDataReceived := cdproto.Message{
		Method: cdproto.EventNetworkDataReceived,
	}
	timestamp := cdp.MonotonicTime(dbResponse.CreatedAt)
	dataRecievedPayload := network.EventDataReceived{
		RequestID:         network.RequestID(fmt.Sprint(dbResponse.RequestId)),
		Timestamp:         &timestamp,
		DataLength:        dbResponse.ContentLength,
		EncodedDataLength: int64(len(dbResponse.Body)),
	}
	dataRecievedJsonBytes, err := json.Marshal(dataRecievedPayload)
	if err != nil {
		return cdproto.Message{}, err
	}
	eventDataReceived.Params = dataRecievedJsonBytes
	return eventDataReceived, nil
}

func generateLoadingFinishedEvent(dbResponse models.DbResponse) (cdproto.Message, error) {
	eventLoadingFinished := cdproto.Message{
		Method: cdproto.EventNetworkLoadingFinished,
	}
	timestamp := cdp.MonotonicTime(dbResponse.CreatedAt)
	loadingFinishedPayload := network.EventLoadingFinished{
		RequestID:         network.RequestID(fmt.Sprint(dbResponse.RequestId)),
		Timestamp:         &timestamp,
		EncodedDataLength: float64(len(dbResponse.Body)),
	}
	loadingFinishedJsonBytes, err := json.Marshal(loadingFinishedPayload)
	if err != nil {
		return cdproto.Message{}, err
	}
	eventLoadingFinished.Params = loadingFinishedJsonBytes
	return eventLoadingFinished, nil
}