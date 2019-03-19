package cmd

import "bytes"

type CmdIface interface {
	Run() (bytes.Buffer, error)
}
