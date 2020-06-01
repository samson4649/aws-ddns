package main

import (
  // "gitlab.monarch-ares.io/containers/aws-ddns"
  "context"
  "os/signal"
  "os"
  "log"
  "fmt"
  "net"
  "time"
  "github.com/alexflint/go-arg"
  "strconv"
  "gopkg.in/yaml.v2"
  // "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/route53"
)

var (
  appVersion string = "0.2.0"
  args Args
)

type ZoneValue struct {
  Record   string  `yaml:"record"`
  Target string  `yaml:"target"`
  TTL    int64   `yaml:"ttl"`
}

type ConfigFile struct {
  Configs []Config  `yaml:"domains"`
}

type Config struct {
  ZoneId     string       `yaml:"zone_id"`
  ZoneName   string       `yaml:"zone_name"`
  ZoneValues []ZoneValue  `yaml:"values"`
}

type Args struct {
  AccessToken     string  `arg:"--aws-access-key,env:AWS_ACCESS_KEY" help:"AWS API Access Token"`
  SecretKey       string  `arg:"--aws-secret-key,env:AWS_SECRET_KEY" help:"AWS API Secret Key"`
  Verbose         bool    `arg:"-v" help:"Provide information on the actions running"`
  ConfigFile      string  `arg:"-c,--config,env:DDNS_CONFIG_FILE"`
  UpdateInterval  int     `arg:"-i,--interval,env:DDNS_UPDATE_INTERVAL"`
}

func (Args) Version() string {
  return "Version: " + appVersion
}

func FloatToString(input_num float64) string {
    return strconv.FormatFloat(input_num, 'f', 2, 64)
}

func main()  {

  args.ConfigFile = "/etc/awsddns/awsddns.yml"
  args.UpdateInterval = 300
  p := arg.MustParse(&args)
  if args.AccessToken == "" || args.SecretKey == "" {
    p.Fail( "You must provide API authentication...")
  }

  ctx,cancel := context.WithCancel(context.Background())
  defer cancel()

  // capture sigint
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	defer signal.Stop(c)
	go func() {
		select {
		case <-c:
			log.Println("caught signal - initiating graceful shutdown, please wait...")
			cancel()
		// case <-scraper.Sigctx.Done():
    case <-ctx.Done():
		}
	}()

  // take in yaml and parse to ConfigFile
  f,err := os.Open( args.ConfigFile )
  if err != nil {
    log.Printf("Not able to access file '%s'", args.ConfigFile )
    return
  }
  log.Printf("Loading config: '%s'", args.ConfigFile )
  cc := ConfigFile{}
  yamlDecoder := yaml.NewDecoder(f)
  yamlDecoder.Decode(&cc)

  // fmt.Printf("%+v\n",cc)
  // r,err := yaml.Marshal(&cc)
  // fmt.Println( string(r) )

  log.Printf("Update interval set at: %ds", args.UpdateInterval )
  ticker := time.NewTicker( time.Second * time.Duration( args.UpdateInterval ) )

  main_loop:
  for {
    go Run(&cc, &ctx)
    log.Println("Waiting for next tick...")
    select {
    case <-ticker.C:
      continue
    case <-ctx.Done():
      break main_loop
    }
  }
}

func Run(cc *ConfigFile, ctx *context.Context ) {
  sess, err := session.NewSession()
  if err != nil {
    log.Println( err.Error() )
    return
  }
  svc := route53.New(sess)

  resolver := net.Resolver{
    PreferGo: true,
    Dial: OpenDNSDialer,
  }

  pipAddr, PI_err := resolver.LookupIPAddr( *ctx, "myip.opendns.com" )
  log.Printf("Public IP Resolved to: %s\n", pipAddr[0].IP)

  for _,zone := range cc.Configs {
    for _,data := range zone.ZoneValues {
      // log.Printf( "\tzone_id: %s\n\tRecord: %s.%s\n\tTarget: %s\n\tTTL: %d\n",  zone.ZoneId, data.Record, zone.ZoneName , data.Target, data.TTL)
      log.Printf( "Updating Record: %s.%s -> %s\n",  data.Record, zone.ZoneName , data.Target)
      if data.TTL == 0 {
        data.TTL = 300
      }
      var params *route53.ChangeResourceRecordSetsInput
      if data.Target == "_public" || data.Target == "" {
        if PI_err != nil {
          log.Printf( "Not able to update record '%s.%s' - not able to resolve public IP address", data.Record, zone.ZoneName )
          continue
        }
        params = buildCRRS( zone.ZoneId, fmt.Sprintf("%s.%s", data.Record, zone.ZoneName ), pipAddr[0].IP.String(), data.TTL )
      } else {
        params = buildCRRS( zone.ZoneId, fmt.Sprintf("%s.%s", data.Record, zone.ZoneName ), data.Target, data.TTL )
      }
      _,err := svc.ChangeResourceRecordSets(params)
      if err != nil {
        log.Println( err.Error() )
        continue
      }
      // log.Println("Change response:", resp )
      log.Printf( "%s.%s - update successful", data.Record, zone.ZoneName )
    }
  }

}

func OpenDNSDialer(ctx context.Context, network, address string) (net.Conn, error){
  d := net.Dialer{}
  return d.DialContext(ctx, "udp", "resolver1.opendns.com:53")
}

func buildCRRS( zoneId, name, target string, ttl int64 ) (*route53.ChangeResourceRecordSetsInput) {
  return &route53.ChangeResourceRecordSetsInput {
    ChangeBatch: &route53.ChangeBatch {
      Changes: []*route53.Change {
        {
          Action: aws.String("UPSERT"),
          ResourceRecordSet: &route53.ResourceRecordSet {
            Name: aws.String(name),
            Type: aws.String("A"),
            ResourceRecords: []*route53.ResourceRecord{
              {
                Value: aws.String(target),
              },
            },
            TTL: &ttl,
          },
        },
      },
      Comment: aws.String("updated now"),
    },
    HostedZoneId: aws.String(zoneId),
  }
}
