//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("not an email", func(t *testing.T) {
		res, err := GetDomainStat(bytes.NewBufferString("{\"Id\":2,\"Name\":\"Jesse Vasquez\",\"Username\":\"qRichardson\",\"Email\":\"mLynchbroWsecat.com\",\"Phone\":\"9-373-949-64-00\",\"Password\":\"SiZLeNSGn\",\"Address\":\"Fulton Hill 80\"}"), "ru")
		require.NoError(t, err)
		require.Len(t, res, 0)
	})

	t.Run("error a JSON", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString("{\"\":,\"Name\":\"Howard Mendoza\",\"Username\":\"0Oliver\",\"Email\":\"aliquid_qui_ea@Browsedrive.gov\",\"Phone\":\"6-866-899-36-79\",\"Password\":\"InAQJvsq\",\"Address\":\"Blackbird Place 25\"}"), "gov")
		require.Error(t, err)
		require.Nil(t, result)
	})

}

func BenchmarkGetDomainStatLong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			b.StopTimer()
			r, err := zip.OpenReader("testdata/users.dat.zip")
			require.NoError(b, err)
			defer r.Close()

			data, err := r.File[0].Open()
			require.NoError(b, err)

			b.StartTimer()
			stat, err := GetDomainStat(data, "biz")
			b.StopTimer()
			require.Len(b, stat, 383)
			require.NoError(b, err)
		}()
	}
}

var shortData string = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

func BenchmarkGetDomainStatShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		stat, err := GetDomainStat(bytes.NewBufferString(shortData), "linktype")
		b.StopTimer()
		require.Len(b, stat, 0)
		require.NoError(b, err)
	}
}
