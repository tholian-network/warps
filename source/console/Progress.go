package console

import "os"
import "slices"
import "strings"

func toWords(message string) []string {

	var result []string

	chunk := ""

	for m := 0; m < len(message); m++ {

		chr := string(message[m])

		if chr == " " {

			if chunk != "" {
				result = append(result, chunk)
			}

			chunk = ""

		} else {
			chunk += chr
		}

	}

	if chunk != "" {
		result = append(result, chunk)
	}

	return result

}

func isSameProgress(old_message string, new_message string) bool {

	var result bool = false

	if strings.Contains(old_message, " of ") && strings.Contains(new_message, " of ") {

		old_words := toWords(old_message)
		new_words := toWords(new_message)

		old_index := slices.Index(old_words, "of")
		new_index := slices.Index(new_words, "of")

		if old_index > 0 && old_index < len(old_words)-1 && old_index == new_index && len(old_words) == len(new_words) {

			old_prefix := strings.Join(old_words[0:old_index - 1], " ")
			old_suffix := strings.Join(old_words[old_index + 2:], " ")
			new_prefix := strings.Join(new_words[0:new_index - 1], " ")
			new_suffix := strings.Join(new_words[new_index + 2:], " ")

			if old_prefix == new_prefix && old_suffix == new_suffix {
				result = true
			}

		}

	} else if strings.Contains(old_message, " of ") && !strings.Contains(new_message, " of ") {

		old_words := toWords(old_message)
		old_index := slices.Index(old_words, "of")

		if old_index > 0 && old_index < len(old_words)-1 {

			old_prefix := strings.Join(old_words[0:old_index - 1], " ")
			old_suffix := strings.Join(old_words[old_index + 2:], " ")

			if strings.HasPrefix(new_message, old_prefix) && strings.HasSuffix(new_message, old_suffix) {
				result = true
			}

		}

	}

	return result

}

func Progress(message string) {

	if features[FeatureProgress] == true {

		message = strings.ReplaceAll(message, "\n", "")
		message = sanitize(message)
		offset := toOffset()

		if len(MESSAGES) > 0 {

			last_method := MESSAGES[len(MESSAGES)-1].Method
			last_message := MESSAGES[len(MESSAGES)-1].Value

			if last_method == "Progress" && isSameProgress(last_message, message) {
				os.Stdout.WriteString("\033[A\033[2K\r")
				MESSAGES[len(MESSAGES)-1] = NewMessage("Progress", message)
			} else {
				MESSAGES = append(MESSAGES, NewMessage("Progress", message))
			}

		} else {
			MESSAGES = append(MESSAGES, NewMessage("Progress", message))
		}

		if COLORS == true {
			os.Stdout.WriteString("\u001b[40m" + offset + toSeparator(message) + message + "\u001b[K\u001b[0m\n")
		} else {
			os.Stdout.WriteString(offset + toSeparator(message) + message + "\n")
		}

	} else {

		message = strings.ReplaceAll(message, "\n", "")
		message = sanitize(message)

		if len(MESSAGES) > 0 {

			last_method := MESSAGES[len(MESSAGES)-1].Method
			last_message := MESSAGES[len(MESSAGES)-1].Value

			if last_method == "Progress" && isSameProgress(last_message, message) {
				MESSAGES[len(MESSAGES)-1] = NewMessage("Progress", message)
			} else {
				MESSAGES = append(MESSAGES, NewMessage("Progress", message))
			}

		} else {
			MESSAGES = append(MESSAGES, NewMessage("Progress", message))
		}

	}

}
