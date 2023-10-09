package common

import data.abbey.functions

allow[msg] {
    functions.expire_after("{{ .TimeExpiry }}")
    msg := "granting access for {{ .TimeExpiry }}"
}
