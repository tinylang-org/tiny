Prism.languages.tiny = Prism.languages.extend('c', {
	'keyword': /\b(break|case|const|continue|default|else|for|fun|if|import|int16|int32|int64|int8|namespace|struct|switch|uint16|uint32|uint64|uint8|var)\b/,
	'builtin': /\b(f)\b/,
	'constant': /\b(DIGITAL_MESSAGE|FIRMATA_STRING|ANALOG_MESSAGE|REPORT_DIGITAL|REPORT_ANALOG|INPUT_PULLUP|SET_PIN_MODE|INTERNAL2V56|SYSTEM_RESET|LED_BUILTIN|INTERNAL1V1|SYSEX_START|INTERNAL|EXTERNAL|DEFAULT|OUTPUT|INPUT|HIGH|LOW)\b/
});