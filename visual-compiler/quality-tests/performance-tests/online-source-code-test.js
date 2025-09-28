import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
      {duration: '20s', target: 10},
      {duration: '20s', target: 10},
      {duration: '20s', target: 0},
    ],
    thresholds: {
      http_req_duration: ["p(95)<2000"],
      http_req_failed: ["rate<0.01"],
    },
}

export default function () {
     let project_name = "guest_project_" + Math.floor(Math.random() * 100000);
    
        const project_url = "https://www.visual-compiler.co.za/api/users/save";
        const project_data =  JSON.stringify({
            users_id: "68d32088d29390ec2c897f35",
            project_name: project_name
        });
        const project_params = {
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer guestuser"
            },
        };
    
        let project_res = http.post(project_url, project_data, project_params);
    
        check(project_res, {
          "status is 200": (r) => r.status === 200,
        });
    
        // SOURCE CODE
    const url = "https://www.visual-compiler.co.za/api/lexing/code";

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
    
        const delete_project_url = "https://www.visual-compiler.co.za/api/users/deleteProject";
                const delete_project_data =  JSON.stringify({
                    users_id: "68d32088d29390ec2c897f35",
                    project_name: project_name
                });
                const delete_project_params = {
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": "Bearer guestuser"
                    },
                };
            
                let delete_project_res = http.request("DELETE", delete_project_url, delete_project_data, delete_project_params);
            
                check(delete_project_res, {
                  "status is 200": (r) => r.status === 200,
                });
        
        sleep(1);
}