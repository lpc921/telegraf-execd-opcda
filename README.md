# telegraf-execd-opcda

[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/lpc921/telegraf-execd-opcda/blob/master/LICENSE)

telegraf-execd-opcda is a [telegraf](https://github.com/influxdata/telegraf) external input plugin to gather data from OPC DA using the [konimarti/opc](https://github.com/konimarti/opc) library.

## Installation

Install the Graybox OPC Automation Wrapper. You can get the Graybox DA Automation Wrapper [here](http://gray-box.net/download_daawrapper.php?lang=en). Follow the installation instruction for this wrapper.

For 32-bit OPC servers on 64-bit systems:

- copy the x86 version `gbda_aut.dll` to `C:\Windows\SysWOW64` folder;
- register this module - enter `%systemroot%\SysWoW64\regsvr32.exe gbda_aut.dll`
- Set go architecture - `$ENV:GOARCH=386` (powershell), or `SET GOARCH=386` (batch)
- Build the executable `go build -o opcda.exe .\cmd\opcda\main.go`

Create your plugin.config file

```toml
[[inputs.opcda]]
  # Measurement name
  name = "sim"

  # OPC DA server
  server = "OI.SIM.1"

  # First successful node in the list is used
  nodes = ["localhost"]

  # Node ID configuration
  # item        - OPC DA Item Name or ItemID
  # name        - field name to use in the output (optional)
  # tags        - extra tags to be added to the output metric (optional)
  # Example:
  items = [
    { item = "PORT.PLC.int" },
    { item = "PORT.PLC.float", name = "randomFloat", tags = [
      [
        "equipment",
        "maker1",
      ],
      [
        "input",
        "temperature",
      ],
    ] },
  ]
```

From here, you can already test the plugin with your config file.

```ps1
opcda -config opcda.conf
```

If everything is ok, you should see something like this

```text
sim PORT.PLC.int=3i 1633624798000000000
sim,equipment=maker1,input=temperature randomFloat=10.511474609375 1633624798000000000
```

## Telegraf configuration

To use the plugin with telegraf, add this configuration to your main telegraf.conf file. telegraf-execd-opcda is an external plugin using [shim](https://github.com/influxdata/telegraf/blob/master/plugins/common/shim/README.md) and [execd](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd). Go see their doc for more information.

```toml
[[inputs.execd]]
   command = ["/path/to/opcda", "-config", "/path/to/opcda.conf"]
   signal = "none"
   restart_delay = "10s"
   data_format = "influx"
```
