package store_test

import (
	"crypto/rand"
	"encoding/base64"
	mrand "math/rand"
	"net"
	"testing"
	"time"

	"github.com/oklog/ulid"
	"github.com/rotationalio/honu/pkg/object/v1"
	"github.com/rotationalio/honu/pkg/store"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestObjectSerialization(t *testing.T) {
	obj := generateRandomObject(Small)
	data, err := store.Marshal(obj)
	require.NoError(t, err, "could not marshal object")
	require.Greater(t, len(data), 512, "expected object to be greater than the minimum data size")

	cmp := &store.Object{}
	err = store.Unmarshal(data, cmp)
	require.NoError(t, err, "could not unmarshal object")
	require.Equal(t, obj, cmp, "deserialized object does not match original")
}

func BenchmarkSerialization(b *testing.B) {

	makeHonuEncode := func(objs []*store.Object) func(b *testing.B) {
		return func(b *testing.B) {
			b.StopTimer()
			for n := 0; n < b.N; n++ {
				obj := objs[n%len(objs)]

				b.StartTimer()
				data, err := store.Marshal(obj)
				b.StopTimer()

				if err != nil {
					b.FailNow()
				}

				b.ReportMetric(float64(len(data)), "bytes")
			}
		}
	}

	makeHonuDecode := func(hnd [][]byte) func(b *testing.B) {
		return func(b *testing.B) {
			b.StopTimer()
			for n := 0; n < b.N; n++ {
				data := hnd[n%len(hnd)]
				obj := &store.Object{}

				b.StartTimer()
				err := store.Unmarshal(data, obj)
				b.StopTimer()

				if err != nil {
					b.FailNow()
				}
			}

		}
	}

	makeProtobufEncode := func(opbs []*object.Object) func(b *testing.B) {
		return func(b *testing.B) {
			b.StopTimer()
			for n := 0; n < b.N; n++ {
				obj := opbs[n%len(opbs)]

				b.StartTimer()
				data, err := proto.Marshal(obj)
				b.StopTimer()

				if err != nil {
					b.FailNow()
				}

				b.ReportMetric(float64(len(data)), "bytes")
			}
		}
	}

	makeProtobufDecode := func(pbs [][]byte) func(b *testing.B) {
		return func(b *testing.B) {
			b.StopTimer()
			for n := 0; n < b.N; n++ {
				data := pbs[n%len(pbs)]
				obj := &object.Object{}

				b.StartTimer()
				err := proto.Unmarshal(data, obj)
				b.StopTimer()

				if err != nil {
					b.FailNow()
				}
			}
		}
	}

	makeSizeBenchmark := func(size Size) func(b *testing.B) {
		return func(b *testing.B) {
			// Generate objects for testing
			objs := make([]*store.Object, 256)
			for i := range objs {
				objs[i] = generateRandomObject(size)
			}

			opbs := make([]*object.Object, len(objs))
			for i, obj := range objs {
				opbs[i] = convertObject(obj)
			}

			b.Run("Encode", func(b *testing.B) {
				b.Run("Honu", makeHonuEncode(objs))
				b.Run("Protobuf", makeProtobufEncode(opbs))
			})

			hnd := make([][]byte, len(objs))
			for i, obj := range objs {
				data, err := store.Marshal(obj)
				if err != nil {
					b.FailNow()
				}
				hnd[i] = data
			}

			pbs := make([][]byte, len(opbs))
			for i, obj := range opbs {
				var err error
				pbs[i], err = proto.Marshal(obj)
				if err != nil {
					b.FailNow()
				}
			}

			b.Run("Decode", func(b *testing.B) {
				b.Run("Honu", makeHonuDecode(hnd))
				b.Run("Protobuf", makeProtobufDecode(pbs))
			})
		}
	}

	b.Run("Small", makeSizeBenchmark(Small))
	b.Run("Medium", makeSizeBenchmark(Medium))
	b.Run("Large", makeSizeBenchmark(Large))
	b.Run("XLarge", makeSizeBenchmark(XLarge))
}

func convertObject(o *store.Object) *object.Object {
	p := &object.Object{
		Version:      convertVersion(o.Version),
		Schema:       convertSchemaVersion(o.Schema),
		Mimetype:     o.MIME,
		Owner:        o.Owner[:],
		Group:        o.Group[:],
		Permissions:  []byte{o.Permissions},
		Acl:          convertACL(o.ACL),
		WriteRegions: o.WriteRegions,
		Publisher:    convertPublisher(o.Publisher),
		Encryption:   convertEncryption(o.Encryption),
		Compression:  convertCompression(o.Compression),
		Flags:        []byte{o.Flags},
		Created:      timestamppb.New(o.Created),
		Modified:     timestamppb.New(o.Modified),
		Data:         o.Data,
	}

	return p
}

func convertVersion(v *store.Version) *object.Version {
	if v == nil {
		return nil
	}

	vers := &object.Version{
		Pid:       v.PID,
		Version:   v.Version,
		Region:    v.Region,
		Tombstone: v.Tombstone,
		Created:   timestamppb.New(v.Created),
	}

	if v.Parent != nil {
		vers.Parent = convertVersion(v.Parent)
	}

	return vers
}

func convertSchemaVersion(o *store.SchemaVersion) *object.SchemaVersion {
	if o == nil {
		return nil
	}

	return &object.SchemaVersion{
		Name:         o.Name,
		MajorVersion: o.Major,
		MinorVersion: o.Minor,
		PatchVersion: o.Patch,
	}
}

func convertACL(a []*store.AccessControl) []*object.ACL {
	if len(a) == 0 {
		return nil
	}

	acl := make([]*object.ACL, len(a))
	for i, c := range a {
		acl[i] = &object.ACL{
			ClientId:    c.ClientID[:],
			Permissions: []byte{c.Permissions},
		}
	}

	return acl
}

func convertPublisher(o *store.Publisher) *object.Publisher {
	if o == nil {
		return nil
	}

	return &object.Publisher{
		PublisherId: o.PublisherID[:],
		ClientId:    o.ClientID[:],
		IpAddr:      o.IPAddress.String(),
		UserAgent:   o.UserAgent,
	}
}

func convertEncryption(o *store.Encryption) *object.Encryption {
	if o == nil {
		return nil
	}

	var encmap = map[store.EncryptionAlgorithm]object.Encryption_Algorithm{
		store.Plaintext:       object.Encryption_PLAINTEXT,
		store.AES256_GCM:      object.Encryption_AES256_GCM,
		store.AES192_GCM:      object.Encryption_AES192_GCM,
		store.AES128_GCM:      object.Encryption_AES128_GCM,
		store.HMAC_SHA256:     object.Encryption_HMAC_SHA256,
		store.RSA_OEAP_SHA512: object.Encryption_RSA_OAEP_SHA512,
	}

	return &object.Encryption{
		PublicKeyId:         o.PublicKeyID,
		EncryptionKey:       o.EncryptionKey,
		HmacSecret:          o.HMACSecret,
		Signature:           o.Signature,
		SealingAlgorithm:    encmap[o.SealingAlgorithm],
		EncryptionAlgorithm: encmap[o.EncryptionAlgorithm],
		SignatureAlgorithm:  encmap[o.SignatureAlgorithm],
	}
}

func convertCompression(o *store.Compression) *object.Compression {
	if o == nil {
		return nil
	}
	return &object.Compression{
		Algorithm: object.Compression_Algorithm(int32(o.Algorithm)),
		Level:     o.Level,
	}
}

func generateRandomObject(size Size) *store.Object {
	obj := &store.Object{
		Version:      randVersion(false),
		Schema:       randSchema(),
		MIME:         "application/random",
		Owner:        ulid.MustNew(ulid.Now(), rand.Reader),
		Group:        ulid.MustNew(ulid.Now(), rand.Reader),
		Permissions:  randUint8(),
		ACL:          randACL(),
		WriteRegions: randRegions(),
		Publisher:    randPublisher(),
		Encryption:   randEncryption(),
		Compression:  randCompression(),
		Flags:        randUint8(),
		Created:      randTime(),
		Modified:     randTime(),
	}

	obj.Data = make([]byte, nRandomBytes(size))
	if _, err := rand.Read(obj.Data); err != nil {
		panic(err)
	}

	return obj
}

func randVersion(isParent bool) *store.Version {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	vers := &store.Version{
		PID:       mrand.Uint64(),
		Version:   mrand.Uint64(),
		Region:    randRegion(),
		Tombstone: mrand.Float32() < 0.25,
		Created:   randTime(),
	}

	if !isParent {
		vers.Parent = randVersion(true)
	}

	return vers
}

func randSchema() *store.SchemaVersion {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	schema := &store.SchemaVersion{
		Name:  "RandomSchema",
		Major: mrand.Uint32(),
		Minor: mrand.Uint32(),
		Patch: mrand.Uint32(),
	}

	return schema
}

func randACL() []*store.AccessControl {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	acl := make([]*store.AccessControl, mrand.Intn(64)+1)
	for i := range acl {
		acl[i] = &store.AccessControl{
			ClientID:    ulid.MustNew(ulid.Now(), rand.Reader),
			Permissions: randUint8(),
		}
	}

	return acl
}

var regions = []string{
	"southafrica-north-a", "southafrica-north-b", "southafrica-north-c",
	"asia-east1-a", "asia-east1-b", "asia-east1-c",
	"asia-east2-a", "asia-east2-b", "asia-east2-c",
	"asia-southeast2-a", "asia-southeast2-b", "asia-southeast2-c",
	"asia-southeast1-a", "asia-southeast1-b", "asia-southeast1-c",
	"asia-south1-a", "asia-south1-b", "asia-south1-c",
	"asia-northeast2-a", "asia-northeast2-b", "asia-northeast2-c",
	"asia-northeast3-a", "asia-northeast3-b", "asia-northeast3-c",
	"australia-southeast1-a", "australia-southeast1-b", "australia-southeast1-c",
	"asia-northeast1-a", "asia-northeast1-b", "asia-northeast1-c",
	"europe-west3-a", "europe-west3-b", "europe-west3-c",
	"europe-north1-a", "europe-north1-b", "europe-north1-c",
	"europe-west2-a", "europe-west2-b", "europe-west2-c",
	"europe-southwest1-a", "europe-southwest1-b", "europe-southwest1-c",
	"europe-west8-a", "europe-west8-b", "europe-west8-c",
	"europe-west9-a", "europe-west9-b", "europe-west9-c",
	"europe-west1-a", "europe-west1-b", "europe-west1-c",
	"europe-central2-a", "europe-central2-b", "europe-central2-c",
	"europe-west6-a", "europe-west6-b", "europe-west6-c",
	"me-central1-a", "me-central1-b", "me-central1-c",
	"me-west1-a", "me-west1-b", "me-west1-c",
	"northamerica-northeast1-a", "northamerica-northeast1-b", "northamerica-northeast1-c",
	"northamerica-northeast2-a", "northamerica-northeast2-b", "northamerica-northeast2-c",
	"us-central1-a", "us-central1-b", "us-central1-c", "us-central1-f",
	"us-west4-a", "us-west4-b", "us-west4-c",
	"us-west2-a", "us-west2-b", "us-west2-c",
	"us-east1-b", "us-east1-c", "us-east1-d",
	"us-east4-a", "us-east4-b", "us-east4-c",
	"us-west3-a", "us-west3-b", "us-west3-c",
	"us-west1-a", "us-west1-b", "us-west1-c",
	"us-east4-a", "us-east4-b", "us-east4-c",
	"us-south1-a", "us-south1-b", "us-south1-c",
	"us-west1-a", "us-west1-b", "us-west1-c",
	"southamerica-west1-a", "southamerica-west1-b", "southamerica-west1-c",
	"southamerica-east1-a", "southamerica-east1-b", "southamerica-east1-c",
	"northamerica-northeast1-a", "northamerica-northeast1-b", "northamerica-northeast1-c",
	"us-east5-a", "us-east5-b", "us-east5-c",
	"europe-central2-a", "europe-central2-b", "europe-central2-c",
	"southamerica-west1-a", "southamerica-west1-b", "southamerica-west1-c",
}

func randRegions() []string {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	regions := make([]string, mrand.Intn(9)+1)
	for i := range regions {
		regions[i] = randRegion()
	}

	return regions
}

func randRegion() string {
	return regions[mrand.Intn(len(regions))]
}

func randPublisher() *store.Publisher {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	return &store.Publisher{
		PublisherID: ulid.MustNew(ulid.Now(), rand.Reader),
		ClientID:    ulid.MustNew(ulid.Now(), rand.Reader),
		IPAddress:   net.IPv4(randUint8(), randUint8(), randUint8(), randUint8()),
		UserAgent:   "Random User Agent v1",
	}
}

func randEncryption() *store.Encryption {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	algs := []store.EncryptionAlgorithm{
		store.Plaintext, store.AES128_GCM, store.AES192_GCM, store.AES256_GCM,
	}

	enc := &store.Encryption{
		EncryptionAlgorithm: algs[mrand.Intn(len(algs))],
	}

	if enc.EncryptionAlgorithm == store.Plaintext {
		return enc
	}

	enc.SealingAlgorithm = store.RSA_OEAP_SHA512
	enc.SignatureAlgorithm = store.HMAC_SHA256
	enc.PublicKeyID = base64.RawStdEncoding.EncodeToString(randBytes(16))
	enc.EncryptionKey = randBytes(32)
	enc.HMACSecret = randBytes(32)
	enc.Signature = randBytes(256)

	return enc
}

func randCompression() *store.Compression {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	algs := []store.CompressionAlgorithm{
		store.None, store.GZIP, store.COMPRESS, store.DEFLATE, store.BROTLI,
	}

	cmp := &store.Compression{
		Algorithm: algs[mrand.Intn(len(algs))],
	}

	if cmp.Algorithm == store.GZIP || cmp.Algorithm == store.COMPRESS {
		cmp.Level = mrand.Int63n(9) + 1
	}

	return cmp
}

func randUint8() uint8 {
	return uint8(mrand.Int31n(255))
}

func randTime() time.Time {
	td := mrand.Int63n(3.154e+16)
	if mrand.Float32() < 0.5 {
		td = td * -1
	}
	return time.Now().Add(time.Duration(td)).Truncate(time.Nanosecond)
}

func randBytes(n int) []byte {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return buf
}

func nRandomBytes(size Size) int64 {
	switch size {
	case Small:
		return mrand.Int63n(4096) + 512
	case Medium:
		return mrand.Int63n(32768) + 8192
	case Large:
		return mrand.Int63n(262144) + 65536
	case XLarge:
		return mrand.Int63n(4194304) + 1048576
	default:
		panic("unknown size")
	}
}

type Size uint8

const (
	Small Size = iota
	Medium
	Large
	XLarge
)
