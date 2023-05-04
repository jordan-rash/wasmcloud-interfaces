// Provide helpful metadata used by code generators
metadata package = [{
    namespace: "org.wasmcloud.interface.display",
    crate: "wasmcloud_interface_display",    
    doc: "Interface with small LCD displays",
}]

namespace org.wasmcloud.interface.display

use org.wasmcloud.model#wasmbus
use org.wasmcloud.model#U8

@wasmbus(
    contractId: "wasmcloud:display",
    providerReceive: true )
service Display {
  version: "0.1.0",
  operations: [
    DisplayLine
    Clear
  ]
}

operation DisplayLine {
  input: Line
  output: Boolean
}

operation Clear {
  output: Boolean
}

structure Line {
  text: String,
  lineNumber: U8
}


