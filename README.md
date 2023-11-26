# ATAT-Telemetry
Telemetry Server for collecting telemetry data

## Code Generation with Proto

- Ensure that `GOPATH`, `GOROOT`, `GOBIN`, and `PATH` are set appropriately 

- Use `protoc` to compile from within the `Protobuf` folder 
`protoc --go_out=paths=source_relative:gen -I. telemetry.proto`

- Generate c code for Arduino using Nanopb
- Make sure you have protobuf installed `pip3 install protobuf`
`python ~/Documents/nanopb/generator/nanopb_generator.py telemetry.proto`


## The only video I took of this working 🥲

https://github.com/ParthSareen/ATAT-Telemetry/assets/29360864/315c2dfe-0c0f-4c6e-9197-ad439d5720ed

