package openkitgo

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const DEFAULT_SEND_INTERVAL = time.Duration(2) * time.Minute
const DEFAULT_MAX_BEACON_SIZE = 30 * 1024 // default: max 30KB (in B) to send in one beacon
const DEFAULT_CAPTURE = true
const DEFAULT_CAPTURE_ERRORS = true
const DEFAULT_CAPTURE_CRASHES = true

const DEFAULT_DATA_COLLECTION_LEVEL = 2
const DEFAULT_CRASH_REPORTING_LEVEL = 2

const DEFAULT_MAX_RECORD_AGE_IN_MILLIS = 6300000 // 1hour and 45 minutes in ms
const DEFAULT_UPPER_MEMORY_BOUNDARY_IN_BYTES = 100 * 1024 * 1024 // 100MiB
const DEFAULT_LOWER_MEMORY_BOUNDARY_IN_BYTES = 80 * 1024 * 1024 // 80 MiB

type Configuration struct {
	openKitType openKitType

	applicationName string
	applicationID   string
	deviceID        string
	endpointURL     string

	capture        bool
	sendInterval   time.Duration
	maxBeaconSize  int
	captureErrors  bool
	captureCrashes bool

	device *Device

	applicationVersion string

	httpClientConfiguration *HTTPClientConfiguration
	beaconConfiguration     *BeaconConfiguration

	serverConfigurationSet bool

	lock sync.Mutex

	BeaconCacheConfiguration *BeaconCacheConfiguration
}

func NewConfiguration(endpointURL string,
	applicationName string,
	applicationID string,
	applicationVersion string,
	deviceID int,
	operatingSystem string,
	manufacturer string,
	modelID string,
	proxyAddress string,
	verifyCertificates bool) *Configuration {

	// TODO - Implement getTrustManager

	c := new(Configuration)
	c.endpointURL = endpointURL
	c.applicationName = applicationName
	c.applicationID = applicationID
	c.applicationVersion = applicationVersion
	c.deviceID = strconv.Itoa(deviceID)

	c.maxBeaconSize = DEFAULT_MAX_BEACON_SIZE

	d := &Device{
		operatingSystem: operatingSystem,
		manufacturer:    manufacturer,
		modelID:         modelID,
	}

	c.httpClientConfiguration = &HTTPClientConfiguration{
		serverID:           1,
		applicationID:      applicationID,
		baseURL:            endpointURL,
		proxyAddress:       proxyAddress,
		verifyCertificates: verifyCertificates,
	}

	c.beaconConfiguration = &BeaconConfiguration{
		multiplicity:        1,
		dataCollectionLevel: DEFAULT_DATA_COLLECTION_LEVEL,
		crashReportingLevel: DEFAULT_CRASH_REPORTING_LEVEL,
		deviceID:			 deviceID,
	}

	// TODO - Implement BeaconCacheConfiguration
	c.BeaconCacheConfiguration = &BeaconCacheConfiguration {
		maxRecordAge:			DEFAULT_MAX_RECORD_AGE_IN_MILLIS,
		cacheSizeLowerBound:	DEFAULT_LOWER_MEMORY_BOUNDARY_IN_BYTES,
		cacheSizeUpperBound:	DEFAULT_UPPER_MEMORY_BOUNDARY_IN_BYTES,
	}

	c.device = d

	return c
}

func (c *Configuration) createSessionNumber() int {
	return rand.Intn(2147483647)
}

func (c *Configuration) makeTimestamp() int {
	return TimeToMillis(time.Now())
}

func (c *Configuration) updateSettings(statusResponse *StatusResponse) {
	statusResponse.log.Debugf("Registering new config properties %+v", statusResponse)

	c.capture = statusResponse.capture

	newServerID := statusResponse.serverID
	if newServerID == -1 {
		newServerID = 1
	}

	if c.httpClientConfiguration.serverID != newServerID {
		c.httpClientConfiguration = &HTTPClientConfiguration{
			serverID:      newServerID,
			applicationID: c.applicationID,
			baseURL:       c.endpointURL,
		}
	}

	newSendInterval := statusResponse.sendInterval
	if newSendInterval == -1 {
		newSendInterval = DEFAULT_SEND_INTERVAL
	}

	if c.sendInterval != newSendInterval {
		c.sendInterval = newSendInterval
	}

	newMaxBeaconSize := statusResponse.maxBeaconSize
	if newMaxBeaconSize == -1 {
		newMaxBeaconSize = DEFAULT_MAX_BEACON_SIZE
	}

	if c.maxBeaconSize != newMaxBeaconSize {
		c.maxBeaconSize = newMaxBeaconSize
	}

}

func (c *Configuration) isServerConfigurationSet() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.serverConfigurationSet
}

type HTTPClientConfiguration struct {
	baseURL            string
	applicationID      string
	serverID           int
	proxyAddress       string
	verifyCertificates bool
}

type BeaconConfiguration struct {
	multiplicity        int
	dataCollectionLevel int
	crashReportingLevel int
	deviceID			uint32
}

type ServerConfiguration struct {
}

type BeaconCacheConfiguration struct {
	maxRecordAge		uint64
    cacheSizeLowerBound	uint64
    cacheSizeUpperBound	uint64
}	