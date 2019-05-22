# easypost
Integration with EasyPost API for tracking shipments

##### Create EasyPost client
`NewClient("[your_api_key]")`

##### Create shipment tracker
`c.GetTracker("[tracking_code]", "["carrier_name(optional)]")`
 
 it will create tracker in EasyPost and return pointer to Tracker and error. Error can be Payment required error, Unauthorized error or processing error

##### Create web hook handler
`NewWebHookHandler([username], [secret])` it returns `func(r *http.Request) (*Event, error)` which can be used in `http.HandleFunc`
 
##### Get result from WebHook Event
 ```
 r := http.Request{}
 handler := NewWebHookHandler("username", "password")
 e, err := handler(r)
 res, err := e.GetResult()
 switch res.(type) {
 case *Tracker:
 ....
 } 
 ```
