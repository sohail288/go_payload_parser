# go_payload_parser

This allows you to easily parse request payloads

```go
type QueryStringOptional struct {
  value string
} // ?

type ListThingsRequest struct {
  search  string  `query:"q,required" json:"searchString"`
  authToken string `cookie:"authToken,required" query:"token"`
}

func (ltr *ListThingsRequest) ValidateSearch() error {
  if ltr.search is empty and structtag is required {
     // return error
  }
  if q == "forbidden term" {
    return error
  }
  func cleanUp(term string) string {

  }
  ltr.search = cleanUp(ltr.search)
}
```

And then using it:

```go
error = payloadParser.ParsePayload(&ListThingsRequest{}, req)
if error != nil {
  // lol
}

// the request is fully parsed, body / header / query strings / cookies?
```

## API Reference

### Structtag format

#### Query Strings

`query:"name,[required],[default],[serialization],[description]"`

#### Headers

`header:"name,required"

#### Cookies
`cookie:"name,required,[default],[serialization],[description]"`

## TODO:

- Create MVP
- Allow multiple validators
- Allow parsing of arbitrarily deep structs
