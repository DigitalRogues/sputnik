package requesthandling

import (
	"fmt"
	"github.com/q231950/sputnik/keymanager"
	"net/http"
	"strings"
	"time"
)

type RequestManager interface {
	PingRequest() (*http.Request, error)
}

type CloudkitRequestManager struct {
	KeyManager keymanager.KeyManager
}

func (r *CloudkitRequestManager) PingRequest() (*http.Request, error) {

	request, err := http.NewRequest("Get", "https://elbedev.com", nil)

	keyId := r.KeyManager.KeyId()
	request.Header.Add("X-Apple-CloudKit-Request-KeyID", keyId)

	timeString := r.formattedTime(time.Now())
	request.Header.Add("X-Apple-CloudKit-Request-ISO8601Date", timeString)

	signature := r.signatureForParameters(timeString, "", "")
	request.Header.Add("X-Apple-CloudKit-Request-SignatureV1", signature)

	publicKey := r.KeyManager.PublicKey()
	fmt.Println(publicKey)

	privateKey := r.KeyManager.PrivateKey()
	fmt.Println(privateKey)

	return request, err
}

func (r *CloudkitRequestManager) signatureForParameters(date string, body string, subpath string) string {
	parameters := []string{date, body, subpath}
	signature := strings.Join(parameters, ":")
	return signature
}

// [path]/database/[version]/[container]/[environment]/[operation-specific subpath]
// https://api.apple-cloudkit.com/database/1/[container ID]/development/public/users/lookup/email
func (r *CloudkitRequestManager) subpath() string {
	version := "1"
	containerId := "iCloud.com.elbedev.shelve.dev"
	subpath := "public/users/lookup/email"

	components := []string{"database", version, containerId, "development", subpath}
	return strings.Join(components, "/")
}

// https://developer.apple.com/library/content/documentation/DataManagement/Conceptual/CloutKitWebServicesReference/SettingUpWebServices/SettingUpWebServices.html#//apple_ref/doc/uid/TP40015240-CH24-SW4
func (r *CloudkitRequestManager) url() string {
	path := "https://api.apple-cloudkit.com"
	subpath := r.subpath()
	return strings.Join([]string{path, subpath}, "/")
}

func (r *CloudkitRequestManager) formattedTime(time time.Time) string {
	//2006-01-02T15:04:05MST-0700
	// timeString := time.Format("Mon Jan 2 15:04:05 -0700 MST 2006")
	return time.Format("2006-01-02T15:04:05MST-0700")
}

func (r *CloudkitRequestManager) payload(date string, body string, service string) string {
	components := []string{date, body, service}
	return strings.Join(components, ":")
}
