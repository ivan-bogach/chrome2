package chrome2

func scriptOpenURL(url string) string {
	return `
	var error = []
	try {
		location.href = '` + url + `';
	} catch(err) {
		error
	}
	`
}

func scriptGetStringsSlice(jsString string) string {
	return `
		var result = [];
		try {
			` + jsString + `
		} catch(err) {
  			result.push(err.stack)
			result
		}
	`
}
func scriptGetString(jsString string) string {
	return `
		try {
			` + jsString + `
		} catch(err) {
  			err.stack
		}
	`
}
