package file

import "bytes"

// HandleFileName 去除文件名中不合法的部分
func HandleFileName(filename string) string {
	var buffer bytes.Buffer
	var i = 0
	for i = 0; i < len(filename); i++ {
		if filename[i] == '\\' || filename[i] == '/' || filename[i] == ':' || filename[i] == '*' || filename[i] == '?' || filename[i] == '"' || filename[i] == '<' || filename[i] == '>' || filename[i] == '|' {

		} else {
			buffer.WriteByte(filename[i])
		}

	}
	return buffer.String()
}
