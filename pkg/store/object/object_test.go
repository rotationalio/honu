package object_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	mrand "math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/honu/pkg/store/metadata"
	"github.com/rotationalio/honu/pkg/store/object"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {
	meta, data := loadFixture(t)

	obj, err := object.Marshal(meta, data)
	require.NoError(t, err, "could not marshal object")

	require.Len(t, obj, 1271, "unexpected length of encoded object")
	require.Equal(t, object.StorageVersion, obj.StorageVersion())

	ometa, err := obj.Metadata()
	require.NoError(t, err, "could not decode metadata")
	require.Equal(t, meta, ometa, "metadata not correctly serialized")

	odata, err := obj.Data()
	require.NoError(t, err, "could not decode data")
	require.Equal(t, data, odata, "data not correctly serialized")
}

func loadFixture(t *testing.T) (*metadata.Metadata, []byte) {
	var meta *metadata.Metadata
	f, err := os.Open("testdata/metadata.json")
	require.NoError(t, err, "could not open testdata/metadata.json")
	defer f.Close()

	err = json.NewDecoder(f).Decode(&meta)
	require.NoError(t, err, "could not decode metadata")

	data, err := os.ReadFile("testdata/data.json")
	require.NoError(t, err, "could not read testdata/data.json")

	return meta, data
}

//===========================================================================
// Benchmarks
//===========================================================================

type Size uint8

const (
	Small Size = iota
	Medium
	Large
	XLarge
)

func BenchmarkSerialization(b *testing.B) {

	makeHonuEncode := func(meta []*metadata.Metadata, data [][]byte) func(b *testing.B) {
		return func(b *testing.B) {
			b.StopTimer()
			for n := 0; n < b.N; n++ {
				m := meta[n%len(meta)]
				d := data[n%len(data)]

				b.StartTimer()
				data, err := object.Marshal(m, d)
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

				b.StartTimer()
				obj := object.Object(data)
				_, err := obj.Metadata()
				_, err2 := obj.Data()

				b.StopTimer()

				if err != nil || err2 != nil {
					b.FailNow()
				}
			}

		}
	}

	makeSizeBenchmark := func(size Size) func(b *testing.B) {
		return func(b *testing.B) {
			// Generate objects for testing
			meta := make([]*metadata.Metadata, 256)
			data := make([][]byte, 256)
			for i := range meta {
				meta[i], data[i] = generateRandomObject(size)
			}

			b.Run("Encode", makeHonuEncode(meta, data))

			hnd := make([][]byte, len(meta))
			for i, m := range meta {
				data, err := object.Marshal(m, data[i])
				if err != nil {
					b.FailNow()
				}
				hnd[i] = data
			}

			b.Run("Decode", makeHonuDecode(hnd))
		}
	}

	b.Run("Small", makeSizeBenchmark(Small))
	b.Run("Medium", makeSizeBenchmark(Medium))
	b.Run("Large", makeSizeBenchmark(Large))
	b.Run("XLarge", makeSizeBenchmark(XLarge))
}

//===========================================================================
// Generate Random Objects
//===========================================================================

func generateRandomObject(size Size) (*metadata.Metadata, []byte) {
	obj := &metadata.Metadata{
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

	data := make([]byte, nRandomBytes(size))
	if _, err := rand.Read(data); err != nil {
		panic(err)
	}

	return obj, data
}

func randVersion(isParent bool) *metadata.Version {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	vers := &metadata.Version{
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

func randSchema() *metadata.SchemaVersion {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	schema := &metadata.SchemaVersion{
		Name:  "RandomSchema",
		Major: mrand.Uint32(),
		Minor: mrand.Uint32(),
		Patch: mrand.Uint32(),
	}

	return schema
}

func randACL() []*metadata.AccessControl {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	acl := make([]*metadata.AccessControl, mrand.Intn(64)+1)
	for i := range acl {
		acl[i] = &metadata.AccessControl{
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

func randPublisher() *metadata.Publisher {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	return &metadata.Publisher{
		PublisherID: ulid.MustNew(ulid.Now(), rand.Reader),
		ClientID:    ulid.MustNew(ulid.Now(), rand.Reader),
		IPAddress:   net.IPv4(randUint8(), randUint8(), randUint8(), randUint8()),
		UserAgent:   "Random User Agent v1",
	}
}

func randEncryption() *metadata.Encryption {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	algs := []metadata.EncryptionAlgorithm{
		metadata.Plaintext, metadata.AES128_GCM, metadata.AES192_GCM, metadata.AES256_GCM,
	}

	enc := &metadata.Encryption{
		EncryptionAlgorithm: algs[mrand.Intn(len(algs))],
	}

	if enc.EncryptionAlgorithm == metadata.Plaintext {
		return enc
	}

	enc.SealingAlgorithm = metadata.RSA_OEAP_SHA512
	enc.SignatureAlgorithm = metadata.HMAC_SHA256
	enc.PublicKeyID = base64.RawStdEncoding.EncodeToString(randBytes(16))
	enc.EncryptionKey = randBytes(32)
	enc.HMACSecret = randBytes(32)
	enc.Signature = randBytes(256)

	return enc
}

func randCompression() *metadata.Compression {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	algs := []metadata.CompressionAlgorithm{
		metadata.None, metadata.GZIP, metadata.COMPRESS, metadata.DEFLATE, metadata.BROTLI,
	}

	cmp := &metadata.Compression{
		Algorithm: algs[mrand.Intn(len(algs))],
	}

	if cmp.Algorithm == metadata.GZIP || cmp.Algorithm == metadata.COMPRESS {
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
