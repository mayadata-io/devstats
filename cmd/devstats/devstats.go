package main

import (
	lib "devstats"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Sync all projects from "projects.yaml", calling `gha2db_sync` for all of them
func syncAllProjects() bool {
	// Environment context parse
	var ctx lib.Ctx
	ctx.Init()

	// Set non-fatal exec mode, we want to run sync for next project(s) if current fails
	ctx.ExecFatal = false

	// Local or cron mode?
	cmdPrefix := ""
	dataPrefix := lib.DataDir
	if ctx.Local {
		cmdPrefix = "./"
		dataPrefix = "./"
	}

	// Read defined projects
	data, err := ioutil.ReadFile(dataPrefix + ctx.ProjectsYaml)
	lib.FatalOnError(err)

	var projects lib.AllProjects
	lib.FatalOnError(yaml.Unmarshal(data, &projects))

	// Create PID file (if not exists)
	// If PID file exists, exit
	pid := os.Getpid()
	pidFile := "/tmp/devstats.pid"
	f, err := os.OpenFile(pidFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0700)
	if err != nil {
		lib.Printf("Another `devstats` instance is running, PID file '%s' exists, exiting\n", pidFile)
		return false
	}
	fmt.Fprintf(f, "%d", pid)
	lib.FatalOnError(f.Close())

	// Schedule remove PID file when finished
	defer func() { lib.FatalOnError(os.Remove(pidFile)) }()

	// Sort projects by "order"
	orders := []int{}
	projectsMap := make(map[int]string)
	for name, proj := range projects.Projects {
		if lib.IsProjectDisabled(&ctx, name, proj.Disabled) {
			continue
		}
		orders = append(orders, proj.Order)
		projectsMap[proj.Order] = name
	}
	sort.Ints(orders)

	// Only run clone/pull part here
	// Remaining commit analysis in"gha2db_sync"
	// after new commits are fetched from GHA
	// So here we get repo files to the newest state
	// And the gha2db_sync takes Postgres DB commits to the newest state
	// after this it need to update commit files
	if !ctx.SkipGetRepos {
		lib.Printf("Updating git repos for all projects\n")
		dtStart := time.Now()
		_, res := lib.ExecCommand(
			&ctx,
			[]string{
				cmdPrefix + "get_repos",
			},
			map[string]string{
				"GHA2DB_PROCESS_REPOS": "1",
			},
		)
		dtEnd := time.Now()
		if res != nil {
			lib.Printf("Error updating git repos (took %v): %+v\n", dtEnd.Sub(dtStart), res)
			fmt.Fprintf(os.Stderr, "%v: Error updating git repos (took %v): %+v\n", dtEnd, dtEnd.Sub(dtStart), res)
			return false
		}
		lib.Printf("Updated git repos, took: %v\n", dtEnd.Sub(dtStart))
	}

	// Sync all projects
	for _, order := range orders {
		name := projectsMap[order]
		proj := projects.Projects[name]
		projEnv := map[string]string{
			"GHA2DB_PROJECT": name,
			"PG_DB":          proj.PDB,
			"IDB_DB":         proj.IDB,
		}
		// Apply eventual per project specific environment
		for envName, envValue := range proj.Env {
			projEnv[envName] = envValue
		}
		lib.Printf("Syncing #%d %s\n", order, name)
		dtStart := time.Now()
		_, res := lib.ExecCommand(
			&ctx,
			[]string{
				cmdPrefix + "gha2db_sync",
			},
			projEnv,
		)
		dtEnd := time.Now()
		if res != nil {
			lib.Printf("Error result for %s (took %v): %+v\n", name, dtEnd.Sub(dtStart), res)
			fmt.Fprintf(os.Stderr, "%v: Error result for %s (took %v): %+v\n", dtEnd, name, dtEnd.Sub(dtStart), res)
			continue
		}
		lib.Printf("Synced %s, took: %v\n", name, dtEnd.Sub(dtStart))
	}
	return true
}

func main() {
	dtStart := time.Now()
	synced := syncAllProjects()
	dtEnd := time.Now()
	if synced {
		lib.Printf("Synced all projects in: %v\n", dtEnd.Sub(dtStart))
	}
}
