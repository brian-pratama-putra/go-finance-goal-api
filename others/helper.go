package others

func GetStr(p_data map[string]interface{}, p_key string) string {
	if v_val, v_ok := p_data[p_key]; v_ok && v_val != nil {
		if v_str, v_ok2 := v_val.(string); v_ok2 {
			return v_str
		}
	}
	return ""
}
