package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"taylz.io/db"
	"taylz.io/db/mysql"
	"taylz.io/db/patch"
	"taylz.io/env"
	"taylz.io/log"
	"taylz.io/types"
)

// HelpMessage is printed when you use arg "-help" or -"h"
var HelpMessage = `
	db-patch runs database migrations
	internally uses (or can create) table "patch" to manage revision number

	--- options
	[name]			[default]			[comment]

	-help, -h		false				print this help page and then quit

	-PATCH_DIR		"./"				directory to load patch files from

	-DB_USER		(required)			username to use when connecting to database

	-DB_PASSWORD		(required)			password to use when conencting to database

	-DB_HOST		(required)			database host ip address

	-DB_PORT		(required)			port to open database host ip with mysql

	-DB_NAME		(required)			database name to connect to within database server
`

func main() {

	env := ENV().Merge("DB_", db.ENV()).ParseDefault()

	logger := log.StdOutService(log.LevelDebug, log.DefaultFormatWithColor())
	logger.With(types.Dict{
		"DB_NAME":   env["DB_NAME"],
		"PATCH_DIR": env["PATCH_DIR"],
	}).Debug("db-patch")

	if types.Bool(env["help"]) || types.Bool(env["h"]) {
		fmt.Print(HelpMessage)
		return
	}

	conn, err := mysql.Open(db.ParseDSN(env.Match("DB_")))
	if conn == nil {
		logger.Add("Error", err).Error("failed to open db")
		return
	}
	logger.With(types.Dict{
		"HOST": env["DB_HOST"],
		"NAME": env["DB_NAME"],
	}).Info("opened connection")

	// get current patch info
	pid, err := patch.Get(conn)
	if err == db.ErrPatchTable {
		logger.Warn(err.Error())
		ansbuf := "?"
		for ansbuf != "y" && ansbuf != "" && ansbuf != "n" {
			fmt.Print(`patch table does not exist, create patch table now? (y/n): `)
			fmt.Scanln(&ansbuf)
			ansbuf = strings.Trim(ansbuf, " \t")
		}
		if ansbuf == "n" {
			logger.Info("exit without creating patch table")
			return
		}
		if err := mysql.CreatePatchTable(conn); err != nil {
			logger.Add("Error", err).Error("failed to create patch table")
			return
		}
		logger.Info("created patch table")
		pid = 0 // reset pid=-1 from the error state
	} else if err != nil {
		logger.Add("Error", err).Error("failed to identify patch number")
		return
	} else {
		logger.Info("found patch#", pid)
	}

	patches := patch.GetFiles(env["PATCH_DIR"])
	if len(patches) < 1 {
		logger.Error("no patches found")
		return
	}

	for i := pid + 1; patches[i] != ""; i++ {
		fmt.Println("queue patch#", i, " 	file:", patches[i])
	}

	// ask about patches
	ansbuf := "?"
	for ansbuf != "y" && ansbuf != "" && ansbuf != "n" {
		fmt.Print("Apply patches? [y/n] (default 'y'): ")
		fmt.Scanln(&ansbuf)
		ansbuf = strings.Trim(ansbuf, " \t\n")
	}
	if ansbuf == "n" {
		logger.Info("not applying patches")
		return
	}

	// apply patches
	for pid++; patches[pid] != ""; pid++ {
		pf := patches[pid]
		log := logger.With(types.Dict{
			"PatchID":   pid,
			"PatchFile": pf,
		})
		tStart := time.Now()
		sql, err := ioutil.ReadFile(pf)

		if err = db.ExecTx(conn, string(sql)); err != nil {
			log.Add("Error", err).Error("failed to patch")
			return
		} else if _, err = conn.Exec("UPDATE patch SET patch=?", pid); err != nil {
			log.Add("Error", err).Error("failed to update patch number")
			return
		}
		log.Add("Time", time.Now().Sub(tStart)).Info("applied patch")
	}

	logger.Add("Patch", pid-1).Info("done")
}

func ENV() env.Service {
	return env.Service{
		"PATCH_DIR": "",
	}
}
