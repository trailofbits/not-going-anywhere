```go
if len(s) > 1 && (s[0] == 'E' || s[0] == 'e') {
parsed, err := strconv.ParseInt(string(s[1:]), 10, 64)
	if err != nil {
		return 0, 0, DecimalExponent, false
	}
	return 10, int32(parsed), DecimalExponent, true
}
```
