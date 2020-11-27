package execute

import (
	"fmt"
	log2 "log"
	"strings"

	"github.com/powerqueue/fitque-users-api/models"
	"github.com/powerqueue/fitque-users-api/routes"
	"github.com/powerqueue/fitque-users-api/services"

	cobra "github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

var viper = v.New()
var rootCmd = &cobra.Command{
	Use:   "int-address-hierarchychy",
	Short: "int-address-hierarchy Interface",
	Long:  "int-address-hierarchy Interface is a NextGen integration service that functions as a REST API integration layer for CaaS to have access to retrieve Address Hierarchies used to determine correspondence destination.",
	Run: func(cmd *cobra.Command, args []string) {

		// initialize logger
		// log.InitializeLogger(nil, "int-address-hierarchy")
		// closer, err := trace.InitTracing("", "int-address-hierarchy")
		// if err != nil {
		// 	log2.Fatalf("cannot initialize tracing. Error:%s", err)
		// }
		// defer closer.Close()

		//initialize and setup rest
		// var jwtMethod jwt.SigningMethod

		port := viper.GetString("port")
		metricsPort := viper.GetString("metrics-port")
		// jwtSecret := viper.GetString("jwt-secret")
		// jwtMethodStr := viper.GetString("jwt-method")
		// switch jwtMethodStr {
		// case "SigningMethodHS256":
		// 	jwtMethod = jwt.SigningMethodHS256
		// case "SigningMethodHS384":
		// 	jwtMethod = jwt.SigningMethodHS384
		// case "SigningMethodHS512":
		// 	jwtMethod = jwt.SigningMethodHS512
		// default:
		// 	// log.Fatalf("Invalid jwt signing method %s", jwtMethodStr)
		// 	fmt.Println("Invalid jwt signing method %s", jwtMethodStr)
		// }

		// connect to mongo and migrate schema
		dbConfigs := &models.DBConfigs1{
			Host:     viper.GetString("db-host"),
			Port:     viper.GetInt64("db-port"),
			User:     viper.GetString("db-user"),
			Password: viper.GetString("db-pwd"),
		}
		dbConn, dbErr := models.ConnectAndMigrate(dbConfigs, "fitqueue-db")
		if dbErr != nil {
			// log.Panicf("Cannot connect to DB, %s", dbErr)
			fmt.Println("Cannot connect to DB, %s", dbErr)
		}

		// create mongo client
		mongoClient := models.NewMongoClient(dbConn)

		// initialize validator
		// validatorInstance := core.NewValidator(validator.New())

		// initialize repository
		loginRepo := models.NewLoginRepository(mongoClient)

		// initialize manager
		loginService := services.LoginService{
			LoginRepo: loginRepo,
		}

		// start web server
		// rest.InitSecurity(jwtMethod, jwtSecret)
		routes.InitServer(port, metricsPort, loginService)
	},
}

func init() {
	var (
		webPort        string
		metricsWebPort string
		jwtSecret      string
		jwtMethod      string
		mongoPort      int64
		mongoHost      string
		mongoPw        string
		mongoUser      string
	)

	// core.InitViperCfg(viper, "int-address-hierarchy")
	v.SetEnvPrefix("fitqueue-login-api")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// core.BindFlagWithCmd(rootCmd, log.InitObsWithViper)
	/*dont think i need this
	for _, init := range inits {
		v, f := init()
		rootCmd.Flags().AddFlagSet(f)
		_ = v.BindPFlags(rootCmd.Flags())
	}
	*/

	rootCmd.PersistentFlags().StringVar(&webPort, "port", "8085", "web server port")
	_ = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	rootCmd.PersistentFlags().StringVar(&metricsWebPort, "metrics-port", "9090", "web server port for metrics")
	_ = viper.BindPFlag("metrics-port", rootCmd.PersistentFlags().Lookup("metrics-port"))
	rootCmd.PersistentFlags().StringVar(&jwtSecret, "jwt-secret", "m$hp#sec4org*1234567891011121314151617181920212223242526272829303132", "JWT secret that was used to sign tokens")
	_ = viper.BindPFlag("jwt-secret", rootCmd.PersistentFlags().Lookup("jwt-secret"))
	rootCmd.PersistentFlags().StringVar(&jwtMethod, "jwt-method", "SigningMethodHS512", "JWT algorithm (SigningMethodHS256, SigningMethodHS384, SigningMethodHS512)")
	_ = viper.BindPFlag("jwt-method", rootCmd.PersistentFlags().Lookup("jwt-method"))

	// rootCmd.PersistentFlags().StringVar(&mongoHost, "db-host", "localhost", "MongoDB Host")
	rootCmd.PersistentFlags().StringVar(&mongoHost, "db-host", "ln004prd", "MongoDB Host")
	_ = viper.BindPFlag("db-host", rootCmd.PersistentFlags().Lookup("db-host"))
	rootCmd.PersistentFlags().Int64Var(&mongoPort, "db-port", 27017, "MongoDB port")
	_ = viper.BindPFlag("db-port", rootCmd.PersistentFlags().Lookup("db-port"))
	rootCmd.PersistentFlags().StringVar(&mongoPw, "db-pwd", "", "MongoDB Password")
	_ = viper.BindPFlag("db-pwd", rootCmd.PersistentFlags().Lookup("db-pwd"))
	rootCmd.PersistentFlags().StringVar(&mongoUser, "db-user", "", "MongoDB Username")
	_ = viper.BindPFlag("db-user", rootCmd.PersistentFlags().Lookup("db-user"))

}

// Execute -- execute command parser.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log2.Fatal(err)
	}

}
