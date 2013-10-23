package config_test

import (
	. "github.com/cloudfoundry/hm9000/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Config", func() {
	configJSON := `
    {
        "heartbeat_period_in_seconds": 10,
        "heartbeat_ttl_in_heartbeats": 3,
        "actual_freshness_ttl_in_heartbeats": 3,
        "grace_period_in_heartbeats": 3,
        "desired_state_ttl_in_heartbeats": 60,
        "desired_freshness_ttl_in_heartbeats": 12,
        "desired_state_batch_size": 500,
        "fetcher_network_timeout_in_seconds": 10,
        "actual_freshness_key": "/actual-fresh",
        "desired_freshness_key": "/desired-fresh",
        "cc_auth_user": "mcat",
        "cc_auth_password": "testing",
        "cc_base_url": "http://127.0.0.1:6001",
        "store_urls": ["http://127.0.0.1:4001"],
        "store_max_concurrent_requests": 30,
        "sender_nats_start_subject": "hm9000.start",
        "sender_nats_stop_subject": "hm9000.stop",
        "sender_message_limit": 30,
        "sender_polling_interval_in_heartbeats": 1,
        "sender_timeout_in_heartbeats": 10,
        "fetcher_polling_interval_in_heartbeats": 6,
        "fetcher_timeout_in_heartbeats": 60,        
        "analyzer_polling_interval_in_heartbeats": 1,
        "analyzer_timeout_in_heartbeats": 10,
        "number_of_crashes_before_backoff_begins": 3,
        "starting_backoff_delay_in_heartbeats": 3,
        "maximum_backoff_delay_in_heartbeats": 96,
        "metrics_server_port": 7879,
        "metrics_server_user": "metrics_server_user",
        "metrics_server_password": "canHazMetrics?",
        "api_server_port": 1812,
        "api_server_user": "magnet",
        "api_server_password": "orangutan4sale",
        "nats": {
            "host": "127.0.0.1",
            "port": 4222,
            "user": "",
            "password": ""
        }
    }
    `

	Context("when passed valid JSON", func() {
		It("deserializes", func() {
			config, err := FromJSON([]byte(configJSON))
			Ω(err).ShouldNot(HaveOccured())
			Ω(config.HeartbeatPeriod).Should(BeNumerically("==", 10))
			Ω(config.HeartbeatTTL()).Should(BeNumerically("==", 30))
			Ω(config.ActualFreshnessTTL()).Should(BeNumerically("==", 30))
			Ω(config.GracePeriod()).Should(BeNumerically("==", 30))
			Ω(config.DesiredStateTTL()).Should(BeNumerically("==", 600))
			Ω(config.DesiredFreshnessTTL()).Should(BeNumerically("==", 120))

			Ω(config.SenderPollingInterval().Seconds()).Should(BeNumerically("==", 10))
			Ω(config.SenderTimeout().Seconds()).Should(BeNumerically("==", 100))
			Ω(config.FetcherPollingInterval().Seconds()).Should(BeNumerically("==", 60))
			Ω(config.FetcherTimeout().Seconds()).Should(BeNumerically("==", 600))
			Ω(config.AnalyzerPollingInterval().Seconds()).Should(BeNumerically("==", 10))
			Ω(config.AnalyzerTimeout().Seconds()).Should(BeNumerically("==", 100))

			Ω(config.NumberOfCrashesBeforeBackoffBegins).Should(BeNumerically("==", 3))
			Ω(config.StartingBackoffDelay().Seconds()).Should(BeNumerically("==", 30))
			Ω(config.MaximumBackoffDelay().Seconds()).Should(BeNumerically("==", 960))

			Ω(config.DesiredStateBatchSize).Should(BeNumerically("==", 500))
			Ω(config.FetcherNetworkTimeout().Seconds()).Should(BeNumerically("==", 10))
			Ω(config.ActualFreshnessKey).Should(Equal("/actual-fresh"))
			Ω(config.DesiredFreshnessKey).Should(Equal("/desired-fresh"))
			Ω(config.CCAuthUser).Should(Equal("mcat"))
			Ω(config.CCAuthPassword).Should(Equal("testing"))
			Ω(config.CCBaseURL).Should(Equal("http://127.0.0.1:6001"))
			Ω(config.StoreURLs).Should(Equal([]string{"http://127.0.0.1:4001"}))
			Ω(config.StoreMaxConcurrentRequests).Should(Equal(30))
			Ω(config.SenderNatsStartSubject).Should(Equal("hm9000.start"))
			Ω(config.SenderNatsStopSubject).Should(Equal("hm9000.stop"))
			Ω(config.SenderMessageLimit).Should(Equal(30))

			Ω(config.MetricsServerPort).Should(Equal(7879))
			Ω(config.MetricsServerUser).Should(Equal("metrics_server_user"))
			Ω(config.MetricsServerPassword).Should(Equal("canHazMetrics?"))

			Ω(config.NATS.Host).Should(Equal("127.0.0.1"))
			Ω(config.NATS.Port).Should(Equal(4222))
			Ω(config.NATS.User).Should(Equal(""))
			Ω(config.NATS.Password).Should(Equal(""))
		})
	})

	Describe("loading up the default config", func() {
		It("should load up the JSON in default_config.json", func() {
			config, err := DefaultConfig()
			Ω(err).ShouldNot(HaveOccured())

			expectedConfig, _ := FromJSON([]byte(configJSON))
			Ω(config).Should(Equal(expectedConfig))
		})
	})

	Describe("loading from a file", func() {
		It("should load up the JSON in default_config.json", func() {
			ioutil.WriteFile("/tmp/_test_config.json", []byte(configJSON), 0777)

			config, err := FromFile("/tmp/_test_config.json")
			Ω(err).ShouldNot(HaveOccured())

			expectedConfig, _ := FromJSON([]byte(configJSON))
			Ω(config).Should(Equal(expectedConfig))
		})
	})

	Context("when passed valid JSON", func() {
		It("deserializes", func() {
			config, err := FromJSON([]byte("¥"))
			Ω(err).Should(HaveOccured())
			Ω(config).Should(BeZero())
		})
	})
})