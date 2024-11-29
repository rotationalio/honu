package object

// An object is the serial data that's written to the underlying storage and is composed
// of a package version, metadata, and the document data serialized in a format that
// can be easily unmarshaled and marshed without requiring copying of data into multiple
// byte slices.
type Object []byte
