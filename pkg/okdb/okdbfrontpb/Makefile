messages.pb.go: messages.proto
	protoc -I. --gogofast_out=\
	Mgoogle/empty.proto=github.com/gogo/protobuf/types,\
	Mgoogle/timestamp.proto=github.com/gogo/protobuf/types,\
	Mgogo.proto=github.com/gogo/protobuf/gogoproto,\
	Mthanos/storepb/types.proto=github.com/thanos-io/thanos/pkg/store/storepb,\
	Mthanos/labelpb/types.proto=github.com/thanos-io/thanos/pkg/store/labelpb,\
	plugins=grpc:. messages.proto

service.pb.go: service.proto messages.pb.go
	protoc -I. --gogofast_out=\
	Mgoogle/empty.proto=github.com/gogo/protobuf/types,\
	Mgoogle/timestamp.proto=github.com/gogo/protobuf/types,\
	Mgogo.proto=github.com/gogo/protobuf/gogoproto,\
	plugins=grpc:. service.proto
