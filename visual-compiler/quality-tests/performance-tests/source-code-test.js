import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
      {duration: '20s', target: 50},
      {duration: '20s', target: 50},
      {duration: '20s', target: 0},
    ],
    thresholds: {
      http_req_duration: ["p(95)<1000"],
      http_req_failed: ["rate<0.01"],
    },
}

export default function () {
    const url = "http://localhost:8080/api/lexing/code";

    const source_code_data = JSON.stringify({
          source_code: `func CreateTokens(source string, rules []TypeRegex) ([]TypeValue, []string, error) {

    tokens := []TypeValue{}
    tokens_unidentified := []string{}

    if source == "" {
        return nil, nil, fmt.Errorf("source code is empty")
    }

    if len(rules) == 0 {
        return nil, nil, fmt.Errorf("no rules specified")
    }

    var builder strings.Builder

    for i := 0; i < len(source); i++ {

        r := rune(source[i])

        if unicode.IsLetter(r) || r == '_' || unicode.IsDigit(r) || r == '.' || unicode.IsSpace(r) || r == '-' {

            builder.WriteRune(r)

        } else {

            builder.WriteRune(' ')
            builder.WriteRune(r)

            if i+1 < len(source) {

                i++
                r = rune(source[i])

                if !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || r == ';' || r == '_' || r == '.' || r == '-' || r == ')') {

                    builder.WriteRune(r)

                } else {
                    i--
                }
            }

            builder.WriteRune(' ')
        }
    }

    var words = strings.Fields(builder.String())

    for _, word := range words {

        found := false

        for _, rule := range rules {

            re := regexp.MustCompile("^" + rule.Regex + "$")

            if re.MatchString(word) {
                found = true
                tokens = append(tokens, TypeValue{Type: rule.Type, Value: word})
                break
            }
        }

        if !found {
            tokens_unidentified = append(tokens_unidentified, word)
        }
    }

    for current_token := 0; current_token < len(tokens_unidentified); current_token++ {
        for other_token := current_token + 1; other_token < len(tokens_unidentified); other_token++ {
            if tokens_unidentified[current_token] == tokens_unidentified[other_token] {
                tokens_unidentified = append(tokens_unidentified[:other_token], tokens_unidentified[other_token+1:]...)
                other_token--
            }
    }
}

    return tokens, tokens_unidentified, nil
}`,
          project_name: "guest_project1"
    });

    const params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let res = http.post(url, source_code_data, params);

    check(res, {
      "status is 200": (r) => r.status === 200,
    });
    
    sleep(1);
}