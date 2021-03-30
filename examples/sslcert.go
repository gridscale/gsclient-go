package main

import (
	"bufio"
	"context"
	"os"

	"github.com/gridscale/gsclient-go/v3"
	log "github.com/sirupsen/logrus"
)

var emptyCtx = context.Background()

// examplePrivateKey is an example of a SSL certificate private key, don't use it in production
const examplePrivateKey = "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC1AWnQPQpH3tNs\n+FpvmGRMS5I1wXuXSe3tpX65sdk+w9ZcatohQ0OAFYRE4cTBCWkIzjjymjdr1h2l\nsazjfBC8V67c2lk27iMo1q3rSYFFh4Z01AxpWadKEsE1IdbGc7O16TNSiWNM9gOA\njIHJqeXomj3BRRhBlqSLAPHPI2Y6nxtWQcWE89u3YLUEI9d3HEphhSzB83EdSp2c\ndEkjh4+LMFi98XwQNbdgOM3f+V2NnEK4kqGZG3U+JZEkHI/CZ4NSgLC1AkGltjdb\nzg5x5KKNM3caRugXluYGuQI/Ec0aR07IgxP2N3jlmQfT/85WbMH8OuBDlBY4tHyI\ntU6nqjJfAgMBAAECggEAZirmDyRlKSwdKuUEJvldo7MEVFNh74NLSVigrzAz77ma\nxY+KkDvnXeTHRBordMpa/x1oB4gEwFmbYmtnqv/ccnMLwJ1+vgKs1eBXSveygAx/\nWHJYjx6LzsPHSrZPBLVKOuPmlC/4XPiAAY9NswazPxfQw8a8akkdl1hxJPpWOb+h\nssonF5Gzt4DxHKd25Lplt1iDgRl0ilIvnLZk2Mkkl1OuVwLtngGKjWBm4d/Obtoy\nEUiNnIp7l86miPz5pXmpY/NJSy6/oOVceoDj+pR4eBNlI1QJNaCrGAO6JA8Hdh8I\nibNWWQBK7lQBFJr9JC8/NKbxEXZide+Fesxi7/KMAQKBgQDlIK/757GbaCLg5JUx\nTLsAKvi1RD4iOtdkwCFowgonVtDwlw7G5WjQAVC6n8SLGCAfspNWKtfS8n5iq9rb\nKV4dno/xPhI7n2kFRbM2IGrUAdLgS/KgI7E4Y4AuPxPJLcgeTk+g/LCjfoZFP3IJ\nladyPKh1eDrj9YLzEqc4hhRr6wKBgQDKO+jcbQgnAAornSwwRJ0Cj5XNrjpAdsb5\nM1A1UFdV0hnUz7wkjT4B14XUyQSnZUJUxV6/qI+lmfoEOuPKgTUkZATO+hfPcFfX\n9K8B+5qr6TD+9AYFcyoaM0kXC+Mjxwa6Yd6h1mFhezufT0Qfjfji/N84zGdz8ucp\nDrvxZTN6XQKBgFMURhtNyH10BemLmHkWvFt0OVfolartsPoMHFESwoG/HeWOsEH4\nHsgFIhN5KNfSeJtlsby1rioD2UXH0IRU/JY6zzCG9C+APqE1w6Rlnraerqq7fw8H\nwhOTKIAcSP1SR1SNypux5A50KxViyuOkyuFGE0L8xEWx2LhwVAfPvgnfAoGAffMx\n45ZELYXoz6DjlGwnHSEvuxl3Tg6rfShoG8wdmGVxkQiPtHQC2kLQJuXK8DYwSXti\ntxrT2985xsimdchiwHdKR12a1qaxDt5k4GdCvS5ORXrVBS/kWMz4CFJu9ClQF2Q8\ns65Al+WYDG/hjYVuLHAw1b73706oiPmUM5NDrEECgYBiak2SnpPXAfJZCtyc/aIA\nCHpH6iGrO04ERItcX8skDiyw+qoUvWmPZIY8KGyF2tF9uVhZoBE4cxn5q1x52RDg\nUV7Ax9IGPvQU9gGuzSIMjEzgitmsGRi1YRAQJ0UaBciTRCS0Jf12zFg/P7a0RyvU\naZ9BxDQPOR71IkWZMIF6EQ==\n-----END PRIVATE KEY-----"

// exampleLeafCert is an example of a SSL certificate's leaf certificate, don't use it in production
const exampleLeafCert = "-----BEGIN CERTIFICATE-----\nMIIDZDCCAkwCCQDAfYZ8/ZCtYDANBgkqhkiG9w0BAQsFADB0MQswCQYDVQQGEwJH\nRTEMMAoGA1UECAwDQkVSMQwwCgYDVQQHDANCRVIxDDAKBgNVBAoMA29nczEMMAoG\nA1UECwwDb2dzMRQwEgYDVQQDDAt3d3cub2dzLmNvbTEXMBUGCSqGSIb3DQEJARYI\naWlAbC5jb20wHhcNMjEwMTIwMTE1NzAyWhcNMjIwMTIwMTE1NzAyWjB0MQswCQYD\nVQQGEwJHRTEMMAoGA1UECAwDQkVSMQwwCgYDVQQHDANCRVIxDDAKBgNVBAoMA29n\nczEMMAoGA1UECwwDb2dzMRQwEgYDVQQDDAt3d3cub2dzLmNvbTEXMBUGCSqGSIb3\nDQEJARYIaWlAbC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC1\nAWnQPQpH3tNs+FpvmGRMS5I1wXuXSe3tpX65sdk+w9ZcatohQ0OAFYRE4cTBCWkI\nzjjymjdr1h2lsazjfBC8V67c2lk27iMo1q3rSYFFh4Z01AxpWadKEsE1IdbGc7O1\n6TNSiWNM9gOAjIHJqeXomj3BRRhBlqSLAPHPI2Y6nxtWQcWE89u3YLUEI9d3HEph\nhSzB83EdSp2cdEkjh4+LMFi98XwQNbdgOM3f+V2NnEK4kqGZG3U+JZEkHI/CZ4NS\ngLC1AkGltjdbzg5x5KKNM3caRugXluYGuQI/Ec0aR07IgxP2N3jlmQfT/85WbMH8\nOuBDlBY4tHyItU6nqjJfAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAH3VDEuIKeIt\nDMzvs4gxowrKyUKP2OzIc447QA34RgiDiroFivV5G133yXmoTJ9YJnUnSbDFtZIY\n/zcY0JhLmBnFDfg4Uim+x7TA2+S3JOdZKV6rg3/CMJ2WYSqczj17c8MI5XQDOZGd\nHvCc5+XrOfWKWbY2toxCw1xpB325r9ufw3hS/NGaVRsvBsJ8A9b5TGLkoUTGR5pl\nPXfgoi4n3fvotDo6Ew1hp0Xlcxriqf/cmMX110cpsXbA4hm9OCv+vKWliOF7EWUO\nZfaaZlI4LuEw1ukFPVgGVeFBUsDAwbxLuVx3u5DVGM/3oj86bgmiBSYySLW5zoYc\nz72bskKMaeg=\n-----END CERTIFICATE-----"

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.DefaultConfiguration(uuid, token)
	client := gsclient.NewClient(config)
	log.Info("gridscale client configured")

	log.Info("Create SSL certificate: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cSSLCert, err := client.CreateSSLCertificate(
		emptyCtx,
		gsclient.SSLCertificateCreateRequest{
			Name:            "go-client-ssl-cert",
			PrivateKey:      examplePrivateKey,
			LeafCertificate: exampleLeafCert,
			Labels:          []string{"test"},
		})
	if err != nil {
		log.Error("Create SSL certificate has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"sslcert_uuid": cSSLCert.ObjectUUID,
	}).Info("SSL certificate is successfully created")
	defer func() {
		err := client.DeleteSSLCertificate(emptyCtx, cSSLCert.ObjectUUID)
		if err != nil {
			log.Error("Delete SSL certificate has failed with error", err)
			return
		}
		log.Info("SSL certificate has been successfully deleted")
	}()

	log.Info("Get a SSL certificate: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	sslCert, err := client.GetSSLCertificate(emptyCtx, cSSLCert.ObjectUUID)
	if err != nil {
		log.Error("Get SSL certificate has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"sslcert": sslCert,
	}).Info("SSL certificate is successfully retrieved")

	log.Info("Get a list of SSL certificates: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	sslCertList, err := client.GetSSLCertificateList(emptyCtx)
	if err != nil {
		log.Error("Get SSL certificate list has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"sslcerts": sslCertList,
	}).Info("SSL certificate list is successfully retrieved")

	log.Info("Delete SSL certificate: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
