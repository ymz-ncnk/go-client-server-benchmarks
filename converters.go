package gcscb

import (
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	cstj_cmds "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_json/cmds"
	cstm_cmds "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus/cmds"
	cstp_cmds "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf/cmds"
	kthp_echo "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo"
)

func ToCstmDataSet(dataSet [][]common.Data) [][]cstm_cmds.EchoCmd {
	cmdDataSet := make([][]cstm_cmds.EchoCmd, len(dataSet))
	for i := range len(dataSet) {
		cmdDataSet[i] = make([]cstm_cmds.EchoCmd, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			cmdDataSet[i][j] = cstm_cmds.EchoCmd(dataSet[i][j])
		}
	}
	return cmdDataSet
}

func ToCstpDataSet(dataSet [][]common.Data) [][]*cstp_cmds.EchoCmd {
	var (
		cmdDataSet   = make([][]*cstp_cmds.EchoCmd, len(dataSet))
		protoDataSet = common.ToProtoData(dataSet)
	)
	for i := range len(protoDataSet) {
		cmdDataSet[i] = make([]*cstp_cmds.EchoCmd, len(protoDataSet[i]))
		for j := range len(protoDataSet[i]) {
			cmdDataSet[i][j] = &cstp_cmds.EchoCmd{ProtoData: protoDataSet[i][j]}
		}
	}
	return cmdDataSet
}

func ToCstjDataSet(dataSet [][]common.Data) [][]cstj_cmds.EchoCmd {
	cmdDataSet := make([][]cstj_cmds.EchoCmd, len(dataSet))
	for i := range len(dataSet) {
		cmdDataSet[i] = make([]cstj_cmds.EchoCmd, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			cmdDataSet[i][j] = cstj_cmds.EchoCmd(dataSet[i][j])
		}
	}
	return cmdDataSet
}

func ToKthpDataSet(dataSet [][]common.Data) [][]*kthp_echo.KitexData {
	kitexDataSet := make([][]*kthp_echo.KitexData, len(dataSet))
	for i := range len(dataSet) {
		kitexDataSet[i] = make([]*kthp_echo.KitexData, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			kitexDataSet[i][j] = &kthp_echo.KitexData{
				Bool:    dataSet[i][j].Bool,
				Int64:   dataSet[i][j].Int64,
				String_: dataSet[i][j].String,
				Float64: dataSet[i][j].Float64,
			}
		}
	}
	return kitexDataSet
}
