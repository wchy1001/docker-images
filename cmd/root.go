package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wchy1001/docker-images/images"
	"log"
	"sync"
)

var rootCmd = &cobra.Command{
	Use:   "docker-image",
	Short: "docker-image used for pulling images and push images",
	Long: `docker-image read a json file that used for pull,tag and push images,`,
	//Args: cobra.MinimumNArgs(1),  # func args参数是用来获取第一个参数，而不是flag 比如 docker pull 中的pull就会在args列表中
	Run: func(cmd *cobra.Command,args []string) {
		Do(conf.Images)
	},
}


type Config struct {
	Images []string	`json:"images"`
}

var (
	cfgFile string
	dorg  	string
	dhost   string
	conf 	Config
	)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}


func init() {
	//run initConfig() before all command called
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/images.json)")
	rootCmd.PersistentFlags().StringVar(&dhost, "host", "127.0.0.1:4000", "destination image host and port")
	rootCmd.PersistentFlags().StringVar(&dorg, "org", "awesome", "destination organzion defalt is aswsome")
	//viper.BindPFlag("cfgFile", rootCmd.PersistentFlags().Lookup("config"))
}

//this func used for checking cfg-file
func initConfig(){
	if cfgFile != "" {
		// 直接使用flag提供的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil{
			log.Fatalln(err)
		}
		// 添加搜索路径并设置文件名字，不带后缀
		viper.AddConfigPath(home)
		viper.SetConfigName("images")
	}
	if  err := viper.ReadInConfig();err != nil{
		log.Fatalf("can not read the images file: %v", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("unable to decode into struct：  %s \n", err)
	}
}

func Do(config []string){
	wg := sync.WaitGroup{}
	wg.Add(len(config))
	for _, c := range config {
		go func(conf string){
			img := images.Newimage(conf)
			if err := img.Do(dhost,dorg); err != nil {
				log.Printf("%v image pull or push failed: %v",conf, err)
			}
			wg.Done()
		}(c)
	}
	wg.Wait()
}
