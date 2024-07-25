package rkafka

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
)

// KafkaConfig configuration parameters
type KafkaConfig struct {
	Brokers    string
	Properties struct {
		HeartbeatInterval      int `json:"heartbeat.interval.ms"`
		SessionTimeout         int `json:"session.timeout.ms"`
		RebalanceTimeout       int `json:"rebalance.timeout.ms"`
		RequestTimeoutOverhead int `json:"request.timeout.ms"`
	}
	ResetSaslRealm bool
	Security       map[string]string
	TLS            struct {
		Enable         bool
		CaCertFiles    string // CA cert.pem with which Kafka brokers certs be signed.  Leave empty for certificates trusted by the OS
		ClientCertFile string // Required for client authentication. It's client cert.pem.
		ClientKeyFile  string // Required if and only if ClientCertFile is present. It's client key.pem.

		TrustStoreLocation string // JKS format of CA certificate, used to extract CA cert.pem.
		TrustStorePassword string
		KeystoreLocation   string // JKS format of client certificate and key, used to extrace client cert.pem and key.pem.
		KeystorePassword   string
		EndpIdentAlgo      string
	}
	// simplified sarama.Config.Net.SASL to only support SASL/PLAIN and SASL/GSSAPI(Kerberos)
	Sasl struct {
		// Whether or not to use SASL authentication when connecting to the broker
		// (defaults to false).
		Enable bool
		// Mechanism is the name of the enabled SASL mechanism.
		// Possible values: PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, GSSAPI (defaults to PLAIN)
		Mechanism string
		// Username is the authentication identity (authcid) to present for
		// SASL/PLAIN or SASL/SCRAM authentication
		Username string
		// Password for SASL/PLAIN or SASL/SCRAM authentication
		Password string
		GSSAPI   struct {
			AuthType           int // 1. KRB5_USER_AUTH, 2. KRB5_KEYTAB_AUTH
			KeyTabPath         string
			KerberosConfigPath string
			ServiceName        string
			Username           string
			Password           string
			Realm              string
			DisablePAFXFAST    bool
		}
	}
}

func NewKafka(cfg *KafkaConfig) (*Kafka, error) {
	var err error
	if len(cfg.Brokers) == 0 {
		err = errors.New("invalid configuration, Kafka section is missing")
		return nil, err
	}
	cfg.convertKfkSecurity()
	if cfg.TLS.CaCertFiles == "" && cfg.TLS.TrustStoreLocation != "" {
		if cfg.TLS.CaCertFiles, _, err = JksToPem(cfg.TLS.TrustStoreLocation, cfg.TLS.TrustStorePassword, false); err != nil {
			return nil, err
		}
	}
	if cfg.TLS.ClientKeyFile == "" && cfg.TLS.KeystoreLocation != "" {
		if cfg.TLS.ClientCertFile, cfg.TLS.ClientKeyFile, err = JksToPem(cfg.TLS.KeystoreLocation, cfg.TLS.KeystorePassword, false); err != nil {
			return nil, err
		}
	}
	if cfg.Sasl.Enable {
		cfg.Sasl.Mechanism = strings.ToUpper(cfg.Sasl.Mechanism)
		switch cfg.Sasl.Mechanism {
		case "PLAIN", "SCRAM-SHA-256", "SCRAM-SHA-512", "GSSAPI":
		default:
			err = errors.New(fmt.Sprintf("rkafka SASL mechanism %s is unsupported", cfg.Sasl.Mechanism))
			return nil, err
		}

		if cfg.ResetSaslRealm {
			port := getKfkPort(cfg.Brokers)
			os.Setenv("DOMAIN_REALM", net.JoinHostPort("hadoop."+strings.ToLower(cfg.Sasl.GSSAPI.Realm), port))
		}
	}
	if cfg.Properties.HeartbeatInterval == 0 {
		cfg.Properties.HeartbeatInterval = defaultHeartbeatInterval
	}
	if cfg.Properties.RebalanceTimeout == 0 {
		cfg.Properties.RebalanceTimeout = defaultRebalanceTimeout
	}
	if cfg.Properties.RequestTimeoutOverhead == 0 {
		cfg.Properties.RequestTimeoutOverhead = defaultRequestTimeoutOverhead
	}
	if cfg.Properties.SessionTimeout == 0 {
		cfg.Properties.SessionTimeout = defaultSessionTimeout
	}

}

// convert java client style configuration into sinker
func (cfg *KafkaConfig) convertKfkSecurity() {
	protocol := cfg.Security["security.protocol"]
	if protocol == "" {
		return
	}

	if strings.Contains(protocol, "SSL") {
		TrySetValue(&cfg.TLS.Enable, true)
		TrySetValue(&cfg.TLS.EndpIdentAlgo, cfg.Security["ssl.endpoint.identification.algorithm"])
		TrySetValue(&cfg.TLS.TrustStoreLocation, cfg.Security["ssl.truststore.location"])
		TrySetValue(&cfg.TLS.TrustStorePassword, cfg.Security["ssl.truststore.password"])
		TrySetValue(&cfg.TLS.KeystoreLocation, cfg.Security["ssl.keystore.location"])
		TrySetValue(&cfg.TLS.KeystorePassword, cfg.Security["ssl.keystore.password"])
	}

	if strings.Contains(protocol, "SASL") {
		TrySetValue(&cfg.Sasl.Enable, true)
		TrySetValue(&cfg.Sasl.Mechanism, cfg.Security["sasl.mechanism"])
		if config, ok := cfg.Security["sasl.jaas.config"]; ok {
			configMap := readConfig(config)
			if strings.Contains(cfg.Sasl.Mechanism, "GSSAPI") {
				// GSSAPI
				if configMap["useKeyTab"] != "true" {
					//Username and password
					TrySetValue(&cfg.Sasl.GSSAPI.AuthType, 1)
					TrySetValue(&cfg.Sasl.GSSAPI.Username, configMap["username"])
					TrySetValue(&cfg.Sasl.GSSAPI.Password, configMap["password"])
				} else {
					//Keytab
					TrySetValue(&cfg.Sasl.GSSAPI.AuthType, 2)
					TrySetValue(&cfg.Sasl.GSSAPI.KeyTabPath, configMap["keyTab"])
					if principal, ok := configMap["principal"]; ok {
						prins := strings.Split(principal, "@")
						TrySetValue(&cfg.Sasl.GSSAPI.Username, prins[0])
						if len(prins) > 1 {
							TrySetValue(&cfg.Sasl.GSSAPI.Realm, prins[1])
						}
					}
					TrySetValue(&cfg.Sasl.GSSAPI.ServiceName, cfg.Security["sasl.kerberos.service.name"])
					TrySetValue(&cfg.Sasl.GSSAPI.KerberosConfigPath, defaultKerberosConfigPath)
				}
			} else {
				// PLAIN, SCRAM-SHA-256 or SCRAM-SHA-512
				TrySetValue(&cfg.Sasl.Username, configMap["username"])
				TrySetValue(&cfg.Sasl.Password, configMap["password"])
			}
		}
	}
}

func readConfig(config string) map[string]string {
	configMap := make(map[string]string)
	config = strings.TrimSuffix(config, ";")
	fields := strings.Split(config, " ")
	for _, field := range fields {
		if strings.Contains(field, "=") {
			key := strings.Split(field, "=")[0]
			value := strings.Split(field, "=")[1]
			value = strings.Trim(value, "\"")
			configMap[key] = value
		}
	}
	return configMap
}

func getKfkPort(brokers string) string {
	hosts := strings.Split(brokers, ",")
	var port string
	for _, host := range hosts {
		_, p, err := net.SplitHostPort(host)
		if err != nil {
			port = p
			break
		}
	}
	return port
}

// set v2 to v1, if v1 didn't bind any value
// FIXME: how about v1 bind default value?
func TrySetValue(v1, v2 interface{}) bool {
	var ok bool
	rt := reflect.TypeOf(v1)
	rv := reflect.ValueOf(v1)

	if rt.Kind() != reflect.Ptr {
		return ok
	}
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	if rv.IsValid() && rv.IsZero() {
		ok = true
		switch rt.Kind() {
		case reflect.Uint:
			v, _ := v2.(uint)
			rv.SetUint(uint64(v))
		case reflect.Uint8:
			v, _ := v2.(uint8)
			rv.SetUint(uint64(v))
		case reflect.Uint16:
			v, _ := v2.(uint16)
			rv.SetUint(uint64(v))
		case reflect.Uint32:
			v, _ := v2.(uint32)
			rv.SetUint(uint64(v))
		case reflect.Uint64:
			v, _ := v2.(uint64)
			rv.SetUint(uint64(v))
		case reflect.Int:
			v, _ := v2.(int)
			rv.SetInt(int64(v))
		case reflect.Int8:
			v, _ := v2.(int8)
			rv.SetInt(int64(v))
		case reflect.Int16:
			v, _ := v2.(int16)
			rv.SetInt(int64(v))
		case reflect.Int32:
			v, _ := v2.(int32)
			rv.SetInt(int64(v))
		case reflect.Int64:
			v, _ := v2.(int64)
			rv.SetInt(int64(v))
		case reflect.Float32:
			v, _ := v2.(float32)
			rv.SetFloat(float64(v))
		case reflect.Float64:
			v, _ := v2.(float64)
			rv.SetFloat(float64(v))
		case reflect.String:
			rv.SetString(v2.(string))
		case reflect.Bool:
			rv.SetBool(v2.(bool))
		default:
			ok = false
		}
	}
	return ok
}

// JksToPem converts JKS to PEM
// Refers to:
// https://serverfault.com/questions/715827/how-to-generate-key-and-crt-file-from-jks-file-for-httpd-apache-server
func JksToPem(jksPath, jksPassword string, overwrite bool) (certPemPath, keyPemPath string, err error) {
	dir, fn := filepath.Split(jksPath)
	certPemPath = filepath.Join(dir, fn+".cert.pem")
	keyPemPath = filepath.Join(dir, fn+".key.pem")
	pkcs12Path := filepath.Join(dir, fn+".p12")
	if overwrite {
		for _, fp := range []string{certPemPath, keyPemPath, pkcs12Path} {
			if err = os.RemoveAll(fp); err != nil {
				return
			}
		}
	} else {
		for _, fp := range []string{certPemPath, keyPemPath, pkcs12Path} {
			if _, err = os.Stat(fp); err == nil {
				return
			}
		}
	}
	cmds := [][]string{
		{"keytool", "-importkeystore", "-srckeystore", jksPath, "-destkeystore", pkcs12Path, "-deststoretype", "PKCS12"},
		{"openssl", "pkcs12", "-in", pkcs12Path, "-nokeys", "-out", certPemPath, "-passin", "env:password"},
		{"openssl", "pkcs12", "-in", pkcs12Path, "-nodes", "-nocerts", "-out", keyPemPath, "-passin", "env:password"},
	}
	for _, cmd := range cmds {
		exe := exec.Command(cmd[0], cmd[1:]...)
		if cmd[0] == "keytool" {
			exe.Stdin = bytes.NewReader([]byte(jksPassword + "\n" + jksPassword + "\n" + jksPassword))
		} else if cmd[0] == "openssl" {
			exe.Env = []string{fmt.Sprintf("password=%s", jksPassword)}
		}
		_, err = exe.CombinedOutput()
		if err != nil {
			return
		}
	}
	return
}
