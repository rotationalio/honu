package object_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	mrand "math/rand"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/region"
	"go.rtnl.ai/honu/pkg/store/lamport"
	"go.rtnl.ai/honu/pkg/store/metadata"
	"go.rtnl.ai/honu/pkg/store/object"
	"go.rtnl.ai/ulid"
)

func TestObject(t *testing.T) {
	meta, data := loadObjectFixture(t)

	obj, err := object.Marshal(meta, data)
	require.NoError(t, err, "could not marshal object")

	require.Len(t, obj, 1264, "unexpected length of encoded object")
	require.Equal(t, object.StorageVersion, obj.StorageVersion())

	ometa, err := obj.Metadata()
	require.NoError(t, err, "could not decode metadata")
	require.Equal(t, meta, ometa, "metadata not correctly serialized")

	key, err := obj.Key()
	require.NoError(t, err, "could not decode key from object")
	require.True(t, bytes.Equal(ometa.Key()[:], key[:]), "object key did not match metadata key")

	odata, err := obj.Data()
	require.NoError(t, err, "could not decode data")
	require.Equal(t, data, odata, "data not correctly serialized")

	require.False(t, obj.Tombstone(), "object should not be a tombstone")
}

func TestTombstone(t *testing.T) {
	meta, _ := loadObjectFixture(t)

	obj, err := object.Marshal(meta, nil)
	require.NoError(t, err, "could not marshal object")

	require.True(t, obj.Tombstone(), "object should be a tombstone")

	odata, err := obj.Data()
	require.NoError(t, err, "could not decode data")
	require.Nil(t, odata, "data not correctly serialized")
}

func TestNil(t *testing.T) {
	obj := object.Object(nil)
	require.Equal(t, uint8(0), obj.StorageVersion())

	_, err := obj.Metadata()
	require.ErrorIs(t, err, object.ErrBadVersion)

	_, err = obj.Data()
	require.ErrorIs(t, err, object.ErrBadVersion)

	require.False(t, obj.Tombstone())
}

func TestMalformed(t *testing.T) {
	obj := object.Object([]byte{0x01})

	_, err := obj.Metadata()
	require.ErrorIs(t, err, object.ErrMalformed)

	_, err = obj.Data()
	require.ErrorIs(t, err, object.ErrMalformed)

	require.False(t, obj.Tombstone())
}

func loadObjectFixture(t *testing.T) (*metadata.Metadata, []byte) {
	var meta *metadata.Metadata
	loadFixture(t, "metadata.json", &meta)

	data, err := os.ReadFile("testdata/data.json")
	require.NoError(t, err, "could not read testdata/data.json")

	return meta, data
}

func loadFixture(t *testing.T, name string, v interface{}) {
	path := filepath.Join("testdata", name)
	f, err := os.Open(path)
	require.NoError(t, err, "could not open %s", path)
	defer f.Close()

	err = json.NewDecoder(f).Decode(v)
	require.NoError(t, err, "could not decode %s", path)
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
		Version:      randVersion(),
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

func randVersion() *metadata.Version {
	vers := &metadata.Version{
		Scalar:    lamport.Scalar{PID: mrand.Uint32(), VID: mrand.Uint64()},
		Region:    randRegion(),
		Tombstone: mrand.Float32() < 0.25,
		Created:   randTime(),
	}

	// 10% chance of nil parent
	if mrand.Float32() < 0.9 {
		vers.Parent = &lamport.Scalar{PID: mrand.Uint32(), VID: mrand.Uint64()}
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

var regions = region.List()

func randRegions() []region.Region {
	// 10% chance of nil
	if mrand.Float32() < 0.1 {
		return nil
	}

	regions := make([]region.Region, mrand.Intn(9)+1)
	for i := range regions {
		regions[i] = randRegion()
	}

	return regions
}

func randRegion() region.Region {
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
