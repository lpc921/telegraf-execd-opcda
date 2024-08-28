# telegraf-execd-opcda

telegraf-execd-opcda is a [telegraf](https://github.com/influxdata/telegraf) external input plugin to gather data from OPC DA using the [konimarti/opc](https://github.com/konimarti/opc) library.

## Usage

* Install OPC DA Automation Wrapper 2.02 by following the installation instruction from [konimarti/opc](https://github.com/konimarti/opc) library.

* Download the [latest release package](https://github.com/lpc921/telegraf-execd-opcda/releases/latest) for your platform.

* Edit opcda.conf as needed. Example:

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

* Edit execd plugin configuration as needed:

```toml
[[inputs.execd]]
   command = ["/path/to/opcda", "-config", "/path/to/opcda.conf"]
   signal = "none"
   restart_delay = "10s"
   data_format = "influx"
```

* Restart or reload Telegraf.

## Development

* Install and register OPC DA Automation Wrapper 2.02. For 32-bit OPC servers on 64-bit systems follow these steps:
  * copy the x86 version `gbda_aut.dll` to `C:\Windows\SysWOW64` folder;
  * register this module - enter `%systemroot%\SysWoW64\regsvr32.exe gbda_aut.dll`
  * Set go architecture - `$ENV:GOARCH=386` (powershell), or `SET GOARCH=386` (batch)

* Build the executable `go build -trimpath -o opcda.exe .\cmd\opcda\main.go`

* Edit opcda.conf as needed.

From here, you can already test the plugin with your config file.

```ps1
opcda -config opcda.conf
```

If everything is ok, you should see something like this

```text
sim PORT.PLC.int=3i 1633624798000000000
sim,equipment=maker1,input=temperature randomFloat=10.511474609375 1633624798000000000
```
