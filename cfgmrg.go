package main

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	out = flag.String("o", "config", "output file name")
)

func main() {
	flag.Parse()

	dir := os.Args[len(os.Args)-1]
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	res := api.NewConfig()
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "config.") {
			cfg, err := clientcmd.LoadFromFile(path.Join(dir, e.Name()))
			if err != nil {
				log.Printf("failed to load from %q: %v", e.Name(), err)
				continue
			}

			var currCtx string
			res.APIVersion = cfg.APIVersion
			res.Kind = cfg.Kind
			for name, cluster := range cfg.Clusters {
				res.Clusters[name] = cluster
			}
			for name, context := range cfg.Contexts {
				parts := strings.Split(name, "/")
				currCtx = parts[len(parts)-1]
				res.Contexts[currCtx] = context
			}
			for name, ext := range cfg.Extensions {
				res.Extensions[name] = ext
			}
			for name, auth := range cfg.AuthInfos {
				res.AuthInfos[name] = auth
			}
			res.Preferences = cfg.Preferences
			res.CurrentContext = currCtx
		}
	}

	if err := clientcmd.WriteToFile(*res, *out); err != nil {
		log.Fatal(err)
	}
}
