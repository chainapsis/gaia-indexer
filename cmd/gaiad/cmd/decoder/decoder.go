package decoder

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/v20/app/params"
	"github.com/cosmos/gogoproto/jsonpb"
	"google.golang.org/grpc"
)

type Decoder struct {
	EncodingConfig params.EncodingConfig
}

func (d *Decoder) ListenAndServe(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterCosmosDecoderServer(grpcServer, d)
	return grpcServer.Serve(lis)
}
func (d *Decoder) mustEmbedUnimplementedCosmosDecoderServer() {
	panic("Forward-compatibility: Unknown method!")
}
func (d *Decoder) Decode(ctx context.Context, request *DecodeRequest) (*DecodeResponse, error) {
	cosmosTx, err := d.EncodingConfig.TxConfig.TxDecoder()(request.TxByte)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	jsonpbMarshaller := jsonpb.Marshaler{}
	var msgs []*GeneralCosmosMsg
	codec := d.EncodingConfig.Marshaler
	for _, msg := range cosmosTx.GetMsgs() {
		msgString, err := jsonpbMarshaller.MarshalToString(msg)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		msgType := sdk.MsgTypeURL(msg)
		signersBytes, _, err := codec.GetMsgV1Signers(msg)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		signers := convertSignersToAddresses(signersBytes, codec)
		msgs = append(msgs, &GeneralCosmosMsg{
			Type:    msgType,
			Message: msgString,
			Signers: signers,
		})
	}
	resultString, err := d.EncodingConfig.TxConfig.TxJSONEncoder()(cosmosTx)
	fmt.Println(string(resultString))
	if err != nil {
		return nil, err
	}
	return &DecodeResponse{
		TxResult: string(resultString),
		Msgs:     msgs,
	}, nil
}

//	func getSliceFromAccAddress(addrs []sdk.AccAddress) []string {
//		var results []string
//		sMap := map[string]string{}
//		for _, s := range addrs {
//			if _, ok := sMap[s.String()]; ok {
//				continue
//			}
//			sMap[s.String()] = "exists"
//			results = append(results, s.String())
//		}
//		return results
//	}
func convertSignersToAddresses(signers [][]byte, cdc codec.Codec) []string {
	var results []string
	sMap := map[string]bool{}
	for _, bytes := range signers {
		signer, err := cdc.InterfaceRegistry().SigningContext().AddressCodec().BytesToString(bytes)
		if err != nil {
			// TODO: 일단 테스트보고 나중에 panic 제거
			panic(err)
			//fmt.Println(err)
			//return nil
		}
		if _, ok := sMap[signer]; ok {
			continue
		}
		sMap[signer] = true
		results = append(results, signer)
	}
	return results
}
