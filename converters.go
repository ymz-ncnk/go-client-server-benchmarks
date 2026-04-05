package gcscb

import (
	"github.com/ymz-ncnk/go-client-server-communication-benchmarks/common"
	cstm_tj "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_json"
	cstm_tm "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_mus"
	cstm_tp "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/cmd-stream/tcp_protobuf"
	kthp_echo "github.com/ymz-ncnk/go-client-server-communication-benchmarks/projects/kitex/ttheader_protobuf/kitex_gen/echo"
)

func ToCstmTMDataSet(dataSet [][]common.Data) [][]cstm_tm.EchoCmd {
	cmdDataSet := make([][]cstm_tm.EchoCmd, len(dataSet))
	for i := range len(dataSet) {
		cmdDataSet[i] = make([]cstm_tm.EchoCmd, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			cmdDataSet[i][j] = cstm_tm.EchoCmd(dataSet[i][j])
		}
	}
	return cmdDataSet
}

func ToCstmTPDataSet(dataSet [][]common.Data) [][]*cstm_tp.EchoCmd {
	var (
		cmdDataSet   = make([][]*cstm_tp.EchoCmd, len(dataSet))
		protoDataSet = common.ToProtoData(dataSet)
	)
	for i := range len(protoDataSet) {
		cmdDataSet[i] = make([]*cstm_tp.EchoCmd, len(protoDataSet[i]))
		for j := range len(protoDataSet[i]) {
			cmdDataSet[i][j] = &cstm_tp.EchoCmd{ProtoData: protoDataSet[i][j]}
		}
	}
	return cmdDataSet
}

func ToCstmTJDataSet(dataSet [][]common.Data) [][]cstm_tj.EchoCmd {
	cmdDataSet := make([][]cstm_tj.EchoCmd, len(dataSet))
	for i := range len(dataSet) {
		cmdDataSet[i] = make([]cstm_tj.EchoCmd, len(dataSet[i]))
		for j := range len(dataSet[i]) {
			cmdDataSet[i][j] = cstm_tj.EchoCmd(dataSet[i][j])
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
				String:  dataSet[i][j].String,
				Float64: dataSet[i][j].Float64,
			}
		}
	}
	return kitexDataSet
}
