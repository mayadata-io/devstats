package devstats

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"

	lib "devstats"
	testlib "devstats/test"
)

// Copies Ctx structure
func copyContext(in *lib.Ctx) *lib.Ctx {
	out := lib.Ctx{
		Debug:               in.Debug,
		CmdDebug:            in.CmdDebug,
		MinGHAPIPoints:      in.MinGHAPIPoints,
		MaxGHAPIWaitSeconds: in.MaxGHAPIWaitSeconds,
		JSONOut:             in.JSONOut,
		DBOut:               in.DBOut,
		ST:                  in.ST,
		NCPUs:               in.NCPUs,
		PgHost:              in.PgHost,
		PgPort:              in.PgPort,
		PgDB:                in.PgDB,
		PgUser:              in.PgUser,
		PgPass:              in.PgPass,
		PgSSL:               in.PgSSL,
		Index:               in.Index,
		Table:               in.Table,
		Tools:               in.Tools,
		Mgetc:               in.Mgetc,
		IDBHost:             in.IDBHost,
		IDBPort:             in.IDBPort,
		IDBDB:               in.IDBDB,
		IDBUser:             in.IDBUser,
		IDBPass:             in.IDBPass,
		IDBMaxBatchPoints:   in.IDBMaxBatchPoints,
		QOut:                in.QOut,
		CtxOut:              in.CtxOut,
		DefaultStartDate:    in.DefaultStartDate,
		ForceStartDate:      in.ForceStartDate,
		LastSeries:          in.LastSeries,
		SkipIDB:             in.SkipIDB,
		SkipPDB:             in.SkipPDB,
		SkipGHAPI:           in.SkipGHAPI,
		SkipArtificailClean: in.SkipArtificailClean,
		SkipGetRepos:        in.SkipGetRepos,
		ResetIDB:            in.ResetIDB,
		ResetRanges:         in.ResetRanges,
		Explain:             in.Explain,
		OldFormat:           in.OldFormat,
		Exact:               in.Exact,
		LogToDB:             in.LogToDB,
		Local:               in.Local,
		IDBDrop:             in.IDBDrop,
		IDBDropProbN:        in.IDBDropProbN,
		MetricsYaml:         in.MetricsYaml,
		GapsYaml:            in.GapsYaml,
		TagsYaml:            in.TagsYaml,
		IVarsYaml:           in.IVarsYaml,
		PVarsYaml:           in.PVarsYaml,
		GitHubOAuth:         in.GitHubOAuth,
		ClearDBPeriod:       in.ClearDBPeriod,
		Trials:              in.Trials,
		LogTime:             in.LogTime,
		WebHookRoot:         in.WebHookRoot,
		WebHookPort:         in.WebHookPort,
		WebHookHost:         in.WebHookHost,
		CheckPayload:        in.CheckPayload,
		FullDeploy:          in.FullDeploy,
		DeployBranches:      in.DeployBranches,
		DeployStatuses:      in.DeployStatuses,
		DeployResults:       in.DeployResults,
		DeployTypes:         in.DeployTypes,
		ProjectRoot:         in.ProjectRoot,
		Project:             in.Project,
		TestsYaml:           in.TestsYaml,
		ReposDir:            in.ReposDir,
		ExecFatal:           in.ExecFatal,
		ExecQuiet:           in.ExecQuiet,
		ExecOutput:          in.ExecOutput,
		ProcessRepos:        in.ProcessRepos,
		ProcessCommits:      in.ProcessCommits,
		ExternalInfo:        in.ExternalInfo,
		ProjectsCommits:     in.ProjectsCommits,
		ProjectsYaml:        in.ProjectsYaml,
		ProjectsOverride:    in.ProjectsOverride,
		ExcludeRepos:        in.ExcludeRepos,
		InputDBs:            in.InputDBs,
		OutputDB:            in.OutputDB,
		TmOffset:            in.TmOffset,
		RecentRange:         in.RecentRange,
		OnlyIssues:          in.OnlyIssues,
		OnlyEvents:          in.OnlyEvents,
		CSVFile:             in.CSVFile,
		ComputeAll:          in.ComputeAll,
		ActorsFilter:        in.ActorsFilter,
		ActorsAllow:         in.ActorsAllow,
		ActorsForbid:        in.ActorsForbid,
		OnlyMetrics:         in.OnlyMetrics,
	}
	return &out
}

// Dynamically sets Ctx fields (uses map of field names into their new values)
func dynamicSetFields(t *testing.T, ctx *lib.Ctx, fields map[string]interface{}) *lib.Ctx {
	// Prepare mapping field name -> index
	valueOf := reflect.Indirect(reflect.ValueOf(*ctx))
	nFields := valueOf.Type().NumField()
	namesToIndex := make(map[string]int)
	for i := 0; i < nFields; i++ {
		namesToIndex[valueOf.Type().Field(i).Name] = i
	}

	// Iterate map of interface{} and set values
	elem := reflect.ValueOf(ctx).Elem()
	for fieldName, fieldValue := range fields {
		// Check if structure actually  contains this field
		fieldIndex, ok := namesToIndex[fieldName]
		if !ok {
			t.Errorf("context has no field: \"%s\"", fieldName)
			return ctx
		}
		field := elem.Field(fieldIndex)
		fieldKind := field.Kind()
		// Switch type that comes from interface
		switch interfaceValue := fieldValue.(type) {
		case int:
			// Check if types match
			if fieldKind != reflect.Int {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.SetInt(int64(interfaceValue))
		case bool:
			// Check if types match
			if fieldKind != reflect.Bool {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.SetBool(interfaceValue)
		case string:
			// Check if types match
			if fieldKind != reflect.String {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.SetString(interfaceValue)
		case time.Time:
			// Check if types match
			fieldType := field.Type()
			if fieldType != reflect.TypeOf(time.Now()) {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.Set(reflect.ValueOf(fieldValue))
		case []int:
			// Check if types match
			fieldType := field.Type()
			if fieldType != reflect.TypeOf([]int{}) {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.Set(reflect.ValueOf(fieldValue))
		case []int64:
			// Check if types match
			fieldType := field.Type()
			if fieldType != reflect.TypeOf([]int64{}) {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.Set(reflect.ValueOf(fieldValue))
		case []string:
			// Check if types match
			fieldType := field.Type()
			if fieldType != reflect.TypeOf([]string{}) {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.Set(reflect.ValueOf(fieldValue))
		case map[string]bool:
			// Check if types match
			fieldType := field.Type()
			if fieldType != reflect.TypeOf(map[string]bool{}) {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.Set(reflect.ValueOf(fieldValue))
		case *regexp.Regexp:
			// Check if types match
			fieldType := field.Type()
			if fieldType != reflect.TypeOf(regexp.MustCompile("a")) {
				t.Errorf("trying to set value %v, type %T for field \"%s\", type %v", interfaceValue, interfaceValue, fieldName, fieldKind)
				return ctx
			}
			field.Set(reflect.ValueOf(fieldValue))
		default:
			// Unknown type provided
			t.Errorf("unknown type %T for field \"%s\"", interfaceValue, fieldName)
		}
	}

	// Return dynamically updated structure
	return ctx
}

func TestInit(t *testing.T) {
	// This is the expected default struct state
	defaultContext := lib.Ctx{
		Debug:               0,
		CmdDebug:            0,
		MinGHAPIPoints:      1,
		MaxGHAPIWaitSeconds: 1,
		JSONOut:             false,
		DBOut:               true,
		ST:                  false,
		NCPUs:               0,
		PgHost:              "localhost",
		PgPort:              "5432",
		PgDB:                "gha",
		PgUser:              "gha_admin",
		PgPass:              "password",
		PgSSL:               "disable",
		Index:               false,
		Table:               true,
		Tools:               true,
		Mgetc:               "",
		IDBHost:             "http://localhost",
		IDBPort:             "8086",
		IDBDB:               "gha",
		IDBUser:             "gha_admin",
		IDBPass:             "password",
		IDBMaxBatchPoints:   10240,
		QOut:                false,
		CtxOut:              false,
		DefaultStartDate:    time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC),
		ForceStartDate:      false,
		LastSeries:          "events_h",
		SkipIDB:             false,
		SkipPDB:             false,
		SkipGHAPI:           false,
		SkipArtificailClean: false,
		SkipGetRepos:        false,
		ResetIDB:            false,
		ResetRanges:         false,
		Explain:             false,
		OldFormat:           false,
		Exact:               false,
		LogToDB:             true,
		Local:               false,
		IDBDrop:             false,
		IDBDropProbN:        20,
		MetricsYaml:         "metrics/metrics.yaml",
		GapsYaml:            "metrics/gaps.yaml",
		TagsYaml:            "metrics/idb_tags.yaml",
		IVarsYaml:           "metrics/idb_vars.yaml",
		PVarsYaml:           "metrics/pdb_vars.yaml",
		GitHubOAuth:         "/etc/github/oauth",
		ClearDBPeriod:       "1 week",
		Trials:              []int{10, 30, 60, 120, 300, 600},
		LogTime:             true,
		WebHookRoot:         "/hook",
		WebHookPort:         ":1982",
		WebHookHost:         "127.0.0.1",
		CheckPayload:        true,
		FullDeploy:          true,
		DeployBranches:      []string{"master"},
		DeployStatuses:      []string{"Passed", "Fixed"},
		DeployResults:       []int{0},
		DeployTypes:         []string{"push"},
		ProjectRoot:         "",
		Project:             "",
		TestsYaml:           "tests.yaml",
		ReposDir:            os.Getenv("HOME") + "/devstats_repos/",
		ExecFatal:           true,
		ExecQuiet:           false,
		ExecOutput:          false,
		ProcessRepos:        false,
		ProcessCommits:      false,
		ExternalInfo:        false,
		ProjectsCommits:     "",
		ProjectsYaml:        "projects.yaml",
		ProjectsOverride:    map[string]bool{},
		ExcludeRepos:        map[string]bool{},
		InputDBs:            []string{},
		OutputDB:            "",
		TmOffset:            0,
		RecentRange:         "2 hours",
		OnlyIssues:          []int64{},
		OnlyEvents:          []int64{},
		CSVFile:             "",
		ComputeAll:          false,
		ActorsFilter:        false,
		ActorsAllow:         nil,
		ActorsForbid:        nil,
		OnlyMetrics:         map[string]bool{},
	}

	var nilRegexp *regexp.Regexp

	// Test cases
	var testCases = []struct {
		name            string
		environment     map[string]string
		expectedContext *lib.Ctx
	}{
		{
			"Default values",
			map[string]string{},
			&defaultContext,
		},
		{
			"Setting debug level",
			map[string]string{"GHA2DB_DEBUG": "2"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"Debug": 2},
			),
		},
		{
			"Setting negative debug level",
			map[string]string{"GHA2DB_DEBUG": "-1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"Debug": -1},
			),
		},
		{
			"Setting command debug level",
			map[string]string{"GHA2DB_CMDDEBUG": "3"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"CmdDebug": 3},
			),
		},
		{
			"Setting GitHub API Points 1",
			map[string]string{"GHA2DB_MIN_GHAPI_POINTS": "0"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"MinGHAPIPoints": 0},
			),
		},
		{
			"Setting GitHub API Points 2",
			map[string]string{"GHA2DB_MIN_GHAPI_POINTS": "-1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"MinGHAPIPoints": 1},
			),
		},
		{
			"Setting GitHub API Points 3",
			map[string]string{"GHA2DB_MIN_GHAPI_POINTS": "1000"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"MinGHAPIPoints": 1000},
			),
		},
		{
			"Setting GitHub API Wait 1",
			map[string]string{"GHA2DB_MAX_GHAPI_WAIT": "0"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"MaxGHAPIWaitSeconds": 0},
			),
		},
		{
			"Setting GitHub API Wait 2",
			map[string]string{"GHA2DB_MAX_GHAPI_WAIT": "-1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"MaxGHAPIWaitSeconds": 1},
			),
		},
		{
			"Setting GitHub API Wait 3",
			map[string]string{"GHA2DB_MAX_GHAPI_WAIT": "1000"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"MaxGHAPIWaitSeconds": 1000},
			),
		},
		{
			"Setting JSON out and disabling DB out",
			map[string]string{"GHA2DB_JSON": "set", "GHA2DB_NODB": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"JSONOut": true, "DBOut": false},
			),
		},
		{
			"Setting ST (singlethreading) and NCPUs",
			map[string]string{"GHA2DB_ST": "1", "GHA2DB_NCPUS": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"ST": true, "NCPUs": 1},
			),
		},
		{
			"Setting TmOffset",
			map[string]string{"GHA2DB_TMOFFSET": "5"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"TmOffset": 5},
			),
		},
		{
			"Setting Postgres parameters",
			map[string]string{
				"PG_HOST": "example.com",
				"PG_PORT": "1234",
				"PG_DB":   "test",
				"PG_USER": "pgadm",
				"PG_PASS": "123!@#",
				"PG_SSL":  "enable",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"PgHost": "example.com",
					"PgPort": "1234",
					"PgDB":   "test",
					"PgUser": "pgadm",
					"PgPass": "123!@#",
					"PgSSL":  "enable",
				},
			),
		},
		{
			"Setting index, table, tools",
			map[string]string{
				"GHA2DB_INDEX":     "1",
				"GHA2DB_SKIPTABLE": "yes",
				"GHA2DB_SKIPTOOLS": "Y",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"Index": true,
					"Table": false,
					"Tools": false,
				},
			),
		},
		{
			"Setting skip log time",
			map[string]string{
				"GHA2DB_SKIPTIME": "Y",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"LogTime": false,
				},
			),
		},
		{
			"Setting getchar default to string longer than 1 character",
			map[string]string{"GHA2DB_MGETC": "yes"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"Mgetc": "y"},
			),
		},
		{
			"Setting InfluxDB parameters",
			map[string]string{
				"IDB_HOST": "example.com",
				"IDB_PORT": "1234",
				"IDB_DB":   "test",
				"IDB_USER": "pgadm",
				"IDB_PASS": "123!@#",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"IDBHost": "http://example.com",
					"IDBPort": "1234",
					"IDBDB":   "test",
					"IDBUser": "pgadm",
					"IDBPass": "123!@#",
				},
			),
		},
		{
			"Setting IDBMaxBatchPoints",
			map[string]string{"IDB_MAXBATCHPOINTS": "1000000"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"IDBMaxBatchPoints": 1000000},
			),
		},
		{
			"Setting query out & context out",
			map[string]string{"GHA2DB_QOUT": "1", "GHA2DB_CTXOUT": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"QOut": true, "CtxOut": true},
			),
		},
		{
			"Setting skip IDB, reset IDB, reset quick ranges",
			map[string]string{
				"GHA2DB_SKIPIDB":     "1",
				"GHA2DB_RESETIDB":    "yes",
				"GHA2DB_RESETRANGES": "yeah",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"SkipIDB":      true,
					"ResetIDB":     true,
					"ResetRanges":  true,
					"IDBDrop":      false,
					"IDBDropProbN": 0,
				},
			),
		},
		{
			"Setting skip PDB",
			map[string]string{"GHA2DB_SKIPPDB": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"SkipPDB": true},
			),
		},
		{
			"Setting skip GHAPI and GetRepos",
			map[string]string{
				"GHA2DB_GETREPOSSKIP": "1",
				"GHA2DB_GHAPISKIP":    "1",
				"GHA2DB_AECLEANSKIP":  "1",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"SkipGHAPI":           true,
					"SkipGetRepos":        true,
					"SkipArtificailClean": true,
				},
			),
		},
		{
			"Setting explain query mode",
			map[string]string{"GHA2DB_EXPLAIN": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"Explain": true},
			),
		},
		{
			"Setting last series",
			map[string]string{"GHA2DB_LASTSERIES": "reviewers_q"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"LastSeries": "reviewers_q"},
			),
		},
		{
			"Setting default start date to 2017",
			map[string]string{"GHA2DB_STARTDT": "2017"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"DefaultStartDate": time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			),
		},
		{
			"Setting default start date to 1982-07-16 10:15:45",
			map[string]string{"GHA2DB_STARTDT": "1982-07-16 10:15:45"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"DefaultStartDate": time.Date(1982, 7, 16, 10, 15, 45, 0, time.UTC),
				},
			),
		},
		{
			"Setting force start date",
			map[string]string{"GHA2DB_STARTDT_FORCE": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ForceStartDate": true,
				},
			),
		},
		{
			"Setting Old pre 2015 GHA JSONs format",
			map[string]string{"GHA2DB_OLDFMT": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"OldFormat": true},
			),
		},
		{
			"Setting exact repository names mode",
			map[string]string{"GHA2DB_EXACT": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"Exact": true},
			),
		},
		{
			"Setting skip DB log mode mode",
			map[string]string{"GHA2DB_SKIPLOG": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"LogToDB": false},
			),
		},
		{
			"Setting local mode",
			map[string]string{"GHA2DB_LOCAL": "yeah"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"Local": true},
			),
		},
		{
			"Setting IDB drop series mode",
			map[string]string{"GHA2DB_IDB_DROP_SERIES": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"IDBDrop": true},
			),
		},
		{
			"Setting IDB drop series prob n 0",
			map[string]string{"GHA2DB_IDB_DROP_PROB_N": "0"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"IDBDropProbN": 0},
			),
		},
		{
			"Setting IDB drop series prob n -10",
			map[string]string{"GHA2DB_IDB_DROP_PROB_N": "-10"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"IDBDropProbN": 0},
			),
		},
		{
			"Setting IDB drop series prob n 100",
			map[string]string{"GHA2DB_IDB_DROP_PROB_N": "100"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"IDBDropProbN": 100},
			),
		},
		{
			"Setting non standard YAML files",
			map[string]string{
				"GHA2DB_METRICS_YAML": "met.YAML",
				"GHA2DB_GAPS_YAML":    "/gapz.yml",
				"GHA2DB_TAGS_YAML":    "/t/g/s.yml",
				"GHA2DB_IVARS_YAML":   "/vari.yml",
				"GHA2DB_PVARS_YAML":   "/varp.yml",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"MetricsYaml": "met.YAML",
					"GapsYaml":    "/gapz.yml",
					"TagsYaml":    "/t/g/s.yml",
					"IVarsYaml":   "/vari.yml",
					"PVarsYaml":   "/varp.yml",
				},
			),
		},
		{
			"Setting GitHub OAUth key",
			map[string]string{
				"GHA2DB_GITHUB_OAUTH": "1234567890123456789012345678901234567890",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"GitHubOAuth": "1234567890123456789012345678901234567890",
				},
			),
		},
		{
			"Setting GitHub OAUth file",
			map[string]string{
				"GHA2DB_GITHUB_OAUTH": "/home/keogh/gh.key",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"GitHubOAuth": "/home/keogh/gh.key",
				},
			),
		},
		{
			"Setting clear DB logs period",
			map[string]string{"GHA2DB_MAXLOGAGE": "3 days"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"ClearDBPeriod": "3 days"},
			),
		},
		{
			"Setting webhook data",
			map[string]string{
				"GHA2DB_WHROOT": "/root",
				"GHA2DB_WHPORT": ":1666",
				"GHA2DB_WHHOST": "0.0.0.0",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"WebHookRoot": "/root",
					"WebHookPort": ":1666",
					"WebHookHost": "0.0.0.0",
				},
			),
		},
		{
			"Setting webhook data missing ':'",
			map[string]string{"GHA2DB_WHPORT": "1986"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"WebHookPort": ":1986"},
			),
		},
		{
			"Setting skip check webhook payload",
			map[string]string{"GHA2DB_SKIP_VERIFY_PAYLOAD": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"CheckPayload": false},
			),
		},
		{
			"Setting skip full deploy",
			map[string]string{"GHA2DB_SKIP_FULL_DEPLOY": "1"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"FullDeploy": false},
			),
		},
		{
			"Setting trials",
			map[string]string{"GHA2DB_TRIALS": "1,2,3,4"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"Trials": []int{1, 2, 3, 4}},
			),
		},
		{
			"Setting webhook params",
			map[string]string{
				"GHA2DB_DEPLOY_BRANCHES": "master,staging,production",
				"GHA2DB_DEPLOY_STATUSES": "ok,passed,fixed",
				"GHA2DB_DEPLOY_RESULTS":  "-1,0,1",
				"GHA2DB_DEPLOY_TYPES":    "push,pull_request",
				"GHA2DB_PROJECT_ROOT":    "/home/lukaszgryglicki/dev/go/src/gha2db",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"DeployBranches": []string{"master", "staging", "production"},
					"DeployStatuses": []string{"ok", "passed", "fixed"},
					"DeployResults":  []int{-1, 0, 1},
					"DeployTypes":    []string{"push", "pull_request"},
					"ProjectRoot":    "/home/lukaszgryglicki/dev/go/src/gha2db",
				},
			),
		},
		{
			"Setting project",
			map[string]string{"GHA2DB_PROJECT": "prometheus"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"Project":     "prometheus",
					"MetricsYaml": "metrics/prometheus/metrics.yaml",
					"GapsYaml":    "metrics/prometheus/gaps.yaml",
					"TagsYaml":    "metrics/prometheus/idb_tags.yaml",
					"IVarsYaml":   "metrics/prometheus/idb_vars.yaml",
					"PVarsYaml":   "metrics/prometheus/pdb_vars.yaml",
				},
			),
		},
		{
			"Setting project and non standard yaml",
			map[string]string{
				"GHA2DB_PROJECT":   "prometheus",
				"GHA2DB_GAPS_YAML": "/gapz.yml",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"Project":     "prometheus",
					"MetricsYaml": "metrics/prometheus/metrics.yaml",
					"GapsYaml":    "/gapz.yml",
					"TagsYaml":    "metrics/prometheus/idb_tags.yaml",
					"IVarsYaml":   "metrics/prometheus/idb_vars.yaml",
					"PVarsYaml":   "metrics/prometheus/pdb_vars.yaml",
				},
			),
		},
		{
			"Setting tests.yaml",
			map[string]string{
				"GHA2DB_TESTS_YAML": "foobar.yml",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"TestsYaml": "foobar.yml",
				},
			),
		},
		{
			"Setting projects.yaml",
			map[string]string{
				"GHA2DB_PROJECTS_YAML": "baz.yml",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ProjectsYaml": "baz.yml",
				},
			),
		},
		{
			"Setting repos dir without ending '/'",
			map[string]string{
				"GHA2DB_REPOS_DIR": "/abc",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ReposDir": "/abc/",
				},
			),
		},
		{
			"Setting repos dir with ending '/'",
			map[string]string{
				"GHA2DB_REPOS_DIR": "~/temp/",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ReposDir": "~/temp/",
				},
			),
		},
		{
			"Setting recent range",
			map[string]string{
				"GHA2DB_RECENT_RANGE": "1 year",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"RecentRange": "1 year",
				},
			),
		},
		{
			"Setting CSV output",
			map[string]string{
				"GHA2DB_CSVOUT": "report.csv",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"CSVFile": "report.csv",
				},
			),
		},
		{
			"Set process repos & commits",
			map[string]string{
				"GHA2DB_PROCESS_REPOS":   "1",
				"GHA2DB_PROCESS_COMMITS": "1",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ProcessRepos":   true,
					"ProcessCommits": true,
				},
			),
		},
		{
			"Set get_repos external info for cncf/gitdm",
			map[string]string{
				"GHA2DB_EXTERNAL_INFO": "1",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ExternalInfo": true,
				},
			),
		},
		{
			"Set compute all periods mode",
			map[string]string{
				"GHA2DB_COMPUTE_ALL": "1",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ComputeAll": true,
				},
			),
		},
		{
			"Set actors filter",
			map[string]string{
				"GHA2DB_ACTORS_FILTER": "1",
				"GHA2DB_ACTORS_ALLOW":  `lukasz\s+gryglicki`,
				"GHA2DB_ACTORS_FORBID": `linus`,
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ActorsFilter": true,
					"ActorsAllow":  regexp.MustCompile(`lukasz\s+gryglicki`),
					"ActorsForbid": regexp.MustCompile(`linus`),
				},
			),
		},
		{
			"Incorrectly set actors filter",
			map[string]string{
				"GHA2DB_ACTORS_FILTER": "",
				"GHA2DB_ACTORS_ALLOW":  `lukasz\s+gryglicki`,
				"GHA2DB_ACTORS_FORBID": `linus`,
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ActorsFilter": false,
					"ActorsAllow":  nilRegexp,
					"ActorsForbid": nilRegexp,
				},
			),
		},
		{
			"Set actors filter allow",
			map[string]string{
				"GHA2DB_ACTORS_FILTER": "1",
				"GHA2DB_ACTORS_ALLOW":  `lukasz\s+gryglicki`,
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ActorsFilter": true,
					"ActorsAllow":  regexp.MustCompile(`lukasz\s+gryglicki`),
					"ActorsForbid": nilRegexp,
				},
			),
		},
		{
			"Set actors filter forbid",
			map[string]string{
				"GHA2DB_ACTORS_FILTER": "yes",
				"GHA2DB_ACTORS_FORBID": `lukasz\s+gryglicki`,
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ActorsFilter": true,
					"ActorsAllow":  nilRegexp,
					"ActorsForbid": regexp.MustCompile(`lukasz\s+gryglicki`),
				},
			),
		},
		{
			"Setting projects commits",
			map[string]string{
				"GHA2DB_PROJECTS_COMMITS": "a,b,c",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ProjectsCommits": "a,b,c",
				},
			),
		},
		{
			"Setting projects override",
			map[string]string{
				"GHA2DB_PROJECTS_OVERRIDE": "a,,c,-,+,,",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ProjectsOverride": map[string]bool{},
				},
			),
		},
		{
			"Setting projects override",
			map[string]string{
				"GHA2DB_PROJECTS_OVERRIDE": "nothing",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ProjectsOverride": map[string]bool{},
				},
			),
		},
		{
			"Setting projects override",
			map[string]string{
				"GHA2DB_PROJECTS_OVERRIDE": "+pro1",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ProjectsOverride": map[string]bool{"pro1": true},
				},
			),
		},
		{
			"Setting projects override",
			map[string]string{
				"GHA2DB_PROJECTS_OVERRIDE": ",+pro1,-pro2,,pro3,,+-pro4,-+pro5,",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"ProjectsOverride": map[string]bool{
						"pro1":  true,
						"pro2":  false,
						"-pro4": true,
						"+pro5": false,
					},
				},
			),
		},
		{
			"Setting exclude repos",
			map[string]string{"GHA2DB_EXCLUDE_REPOS": "repo1,org1/repo2,,abc"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"ExcludeRepos": map[string]bool{
					"repo1":      true,
					"org1/repo2": true,
					"abc":        true,
				},
				},
			),
		},
		{
			"Setting only metrics mode",
			map[string]string{"GHA2DB_ONLY_METRICS": "metric1,metric2,,metric3"},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{"OnlyMetrics": map[string]bool{
					"metric1": true,
					"metric2": true,
					"metric3": true,
				},
				},
			),
		},
		{
			"Setting input & output DBs for 'merge_pdbs' tool",
			map[string]string{
				"GHA2DB_INPUT_DBS": "db1,db2,db3",
				"GHA2DB_OUTPUT_DB": "db4",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"InputDBs": []string{"db1", "db2", "db3"},
					"OutputDB": "db4",
				},
			),
		},
		{
			"Setting debug issues mode on ghapi2db",
			map[string]string{
				"GHA2DB_ONLY_ISSUES": "1,2000,3000000,4000000000,5000000000000,6000000000000000",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"OnlyIssues": []int64{
						1,
						2000,
						3000000,
						4000000000,
						5000000000000,
						6000000000000000,
					},
				},
			),
		},
		{
			"Setting debug events mode on ghapi2db",
			map[string]string{
				"GHA2DB_ONLY_EVENTS": "1,2000,3000000,4000000000,5000000000000,6000000000000000",
			},
			dynamicSetFields(
				t,
				copyContext(&defaultContext),
				map[string]interface{}{
					"OnlyEvents": []int64{
						1,
						2000,
						3000000,
						4000000000,
						5000000000000,
						6000000000000000,
					},
				},
			),
		},
	}

	// Context Init() is verbose when called with CtxDebug
	// For this case we want to discard its STDOUT
	stdout := os.Stdout

	// Execute test cases
	for index, test := range testCases {
		var gotContext lib.Ctx

		// Remember initial environment
		currEnv := make(map[string]string)
		for key := range test.environment {
			currEnv[key] = os.Getenv(key)
		}

		// Set new environment
		for key, value := range test.environment {
			err := os.Setenv(key, value)
			if err != nil {
				t.Errorf(err.Error())
			}
		}

		// When CTXOUT is set, Ctx.Init() writes debug data to STDOUT
		// We don't want to see it while running tests
		if test.environment["GHA2DB_CTXOUT"] != "" {
			fd, err := os.Open(os.DevNull)
			if err != nil {
				t.Errorf(err.Error())
			}
			os.Stdout = fd
		}

		// Initialize context while new environment is set
		gotContext.Init()
		if test.environment["GHA2DB_CTXOUT"] != "" {
			os.Stdout = stdout
		}

		// Restore original environment
		for key := range test.environment {
			err := os.Setenv(key, currEnv[key])
			if err != nil {
				t.Errorf(err.Error())
			}
		}

		// Maps are not directly compareable (due to unknown key order) - need to transorm them
		testlib.MakeComparableMap(&gotContext.ProjectsOverride)
		testlib.MakeComparableMap(&test.expectedContext.ProjectsOverride)
		testlib.MakeComparableMap(&gotContext.ExcludeRepos)
		testlib.MakeComparableMap(&test.expectedContext.ExcludeRepos)
		testlib.MakeComparableMap(&gotContext.OnlyMetrics)
		testlib.MakeComparableMap(&test.expectedContext.OnlyMetrics)

		// Check if we got expected context
		got := fmt.Sprintf("%+v", gotContext)
		expected := fmt.Sprintf("%+v", *test.expectedContext)
		if got != expected {
			t.Errorf(
				"Test case number %d \"%s\"\nExpected:\n%+v\nGot:\n%+v\n",
				index+1, test.name, expected, got,
			)
		}
	}
}
