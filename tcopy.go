package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

var (
	t             *template.Template
	valuesFile    string
	templatesFile string
	destFile      string
)

type K8s struct {
	Name           string
	Port           string
	Path           string
	Host           string
	Ns	       string
	Svc_name       string
	Svc_port       string
	Tag_port       string
	Port_name      string
	Replicas       string
	Image          string
	Img_policy     string
	Env_name01     string
	Env_value01    string
	Volume_name01  string
	Volume_value01 string
}

func initArgs() {
	flag.StringVar(&templatesFile, "s", "svc", "-s 指定template源文件")
	flag.StringVar(&valuesFile, "v", "./values", "-v 指定values配置文件")
	flag.StringVar(&destFile, "d", "./k8s/app.yaml", "-d 指定生成的目标配置文件")
	flag.Parse()
}

func initTemplate(templatesFile string) (err error) {
	switch {
	case templatesFile == "svc":
		templatesFile = "./svc.yaml"
	case templatesFile == "ingress":
		templatesFile = "./ingress.yaml"
	case templatesFile == "deploy":
		templatesFile = "./deployment.yaml"
	default:
		templatesFile = "./template.yaml"
	}
	t, err = template.ParseFiles(templatesFile)
	if err != nil {
		fmt.Println("parse file err:", err)
		return
	}
	return
}

func parseYaml(v *viper.Viper) (yamlObj *K8s) {
	if err := v.Unmarshal(&yamlObj); err != nil {
		fmt.Printf("err:%s", err)
	}
	return yamlObj
}
func main() {
	initArgs()
	initTemplate(templatesFile)
	//读取yaml文件
	vv := viper.New()
	//设置读取的配置文件
	vv.SetConfigName(valuesFile)
	//添加读取的配置文件路径
	vv.AddConfigPath("./")
	//windows环境下为%GOPATH，linux环境下为$GOPATH
	vv.AddConfigPath("%GOPATH/src/")
	//设置配置文件类型
	vv.SetConfigType("yaml")

	if err := vv.ReadInConfig(); err != nil {
		fmt.Printf("err1:%s\n", err)
	}
	yaml := parseYaml(vv)

	file, err := os.OpenFile(destFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}
	if err := t.Execute(file, yaml); err != nil {
		fmt.Println("There was an error:", err.Error())
	}
}
