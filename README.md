# waybar-livestatus

Tiny waybar module to fetch new alerts from livestatus.
It use the go-livestatus module

# requirements

- fontsawesome
- golang

## Usage

You can use a toml config file, see `livestatus.toml.example` or use CLI params.

`waybar-livestatus -h`

## Configuration

`waybar-livestatus` will search config file in `$HOME/.config/waybar/livestatus.toml`

# waybar config

config:

```json
...
"custom/prod-status": {
    "exec": "~/path/to/waybar-livestatus",
    "return-type": "json",
    "markup": true,
    "on-click": "pkill -SIGUSR1 -x waybar-livestat",
},
...
```

style.css

```css
...
#custom-prod-status.pause {
  border-bottom: 3px solid @purple;
}
#custom-prod-status.warning {
  border-bottom: 3px solid @yellow;
}
#custom-prod-status.okcritical {
  border-top: 3px solid @green;
  border-bottom: 3px solid @red;
}
#custom-prod-status.okwarning {
  border-top: 3px solid @green;
  border-bottom: 3px solid @yellow;
}
#custom-prod-status.warningok {
  border-top: 3px solid @yellow;
  border-bottom: 3px solid @green;
}
#custom-prod-status.criticalok {
  border-top: 3px solid @red;
  border-bottom: 3px solid @green;
}
#custom-prod-status.warningwarning {
  border-top: 3px solid @yellow;
  border-bottom: 3px solid @yellow;
}
#custom-prod-status.warningcritical {
  border-top: 3px solid @yellow;
  border-bottom: 3px solid @red;
}
#custom-prod-status.criticalwarning {
  border-top: 3px solid @red;
  border-bottom: 3px solid @yellow;
}
#custom-prod-status.criticalcritical {
  border-top: 3px solid @red;
  border-bottom: 3px solid @red;
}
...
```

# pause

You can toggle pause of the poller with sending `SIGUSR1`.
