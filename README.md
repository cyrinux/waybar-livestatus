# Waybar-livestatus

Tiny waybar module to fetch new alerts from livestatus.
It use the go-livestatus module

# Requirements

- fontsawesome

# Installation

- `make && make install`
- Archlinux: AUR package available https://aur.archlinux.org/packages/waybar-livestatus/

# requirements

- otf-font-awesome

## Configuration

`waybar-livestatus` will search config file in `$HOME/.config/waybar/livestatus.toml`

You can use a toml config file, see `livestatus.toml.example` or use CLI params.

`waybar-livestatus -h`

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

#custom-prod-status.warning {
  border-bottom: 3px solid @yellow;
}
#custom-prod-status.error {
  border-bottom: 3px solid @red;
}
#custom-prod-status.ok {
  border-bottom: 3px solid @green;
}
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

# Usage

## Daemon

You can toggle pause of the poller with sending `SIGUSR1`. Or clicking on the waybar icon if you use my configuration.
You will got the lists of alerts on mouse over the icon.

If `popup` are enable you will got a popup if alerts number raise.

## Client

Client mode was mainly done to be integrated with `dmenu`.

There is a client mode `waybar-livestatus -c`.
Without params this will display the list of alerts.
You can then pass `host` and `service` as params to get the notes_url.
