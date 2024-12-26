package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"

	"github.com/rahulsanju/go_grpc_implementation/invoicer"
	"google.golang.org/grpc"
)

type MyInvoicerServer struct {
	invoicer.UnimplementedInvoicerServer
}

func (m MyInvoicerServer) Create(ctx context.Context, req *invoicer.CreateRequest) (*invoicer.CreateResponse, error) {
	statements := make([]string, 0)
	statements = append(statements, "Amount "+req.Amount.String()+" "+req.Amount.Currency+" transferred from "+req.From+" to "+req.To)
	var response invoicer.CreateResponse
	pdfData, err := json.Marshal(&statements)
	response.Pdf = pdfData
	if err != nil {
		return nil, errors.New("error while creating PDF")
	}

	docxData, err := json.Marshal(&response)
	response.Docx = docxData

	if err != nil {
		return nil, errors.New("error while creating DOCX")
	}

	return &response, nil

}

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("Cannot create listener : %s", err)
	}

	serverRegsitrar := grpc.NewServer()
	service := &MyInvoicerServer{}

	invoicer.RegisterInvoicerServer(serverRegsitrar, service)
	err = serverRegsitrar.Serve(lis)
	if err != nil {
		log.Fatalf("Impossible to server : %s", err)
	}

}
