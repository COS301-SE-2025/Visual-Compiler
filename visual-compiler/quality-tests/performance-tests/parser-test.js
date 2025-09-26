import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
      {duration: '20s', target: 10},
      {duration: '20s', target: 10},
      {duration: '20s', target: 0},
    ],
}

export default function () {

    let project_name = "guest_project_" + Math.floor(Math.random() * 100000);

    const project_url = "http://localhost:8080/api/users/save";
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
    const source_code_url = "http://localhost:8080/api/lexing/code";

    const source_code_data = JSON.stringify({
            source_code: `func main() {

            var red int
            red = 13

            var orange int
            orange = 5

            var yellow int
            yellow = 55

            var green int
            green = 5

            var blue int
            blue = 5

            var purple int
            purple = 5

            var pink int
            pink = 5

            var red int
            red = 13

            var orange int
            orange = 5

            var yellow int
            yellow = 55

            var green int
            green = 5

            var blue int
            blue = 5

            var purple int
            purple = 5

            var pink int
            pink = 5

            var red int
            red = 13

            var orange int
            orange = 5

            var yellow int
            yellow = 55

            var green int
            green = 5

            var blue int
            blue = 5

            var purple int
            purple = 5

            var pink int
            pink = 5

            }`,
            project_name: project_name
    });

    const source_code_params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let source_code_res = http.post(source_code_url, source_code_data, source_code_params);

    check(source_code_res, {
      "status is 200": (r) => r.status === 200,
    });

    sleep(1);

    // LEXER

    const lexer_rules_url = "http://localhost:8080/api/lexing/rules";

    const lexer_rules_data = JSON.stringify({
        pairs: [
            { Type: "KEYWORD", Regex: "int|var|func" },
            { Type: "IDENTIFIER", Regex: "[a-zA-Z_]+" },
            { Type: "OPERATOR", Regex: "=" },
            { Type: "INTEGER", Regex: "[0-9]+" },
            { Type: "OPEN_SCOPE", Regex: "\\{" },
            { Type: "CLOSE_SCOPE", Regex: "\\}" },
            { Type: "OPEN_BRACKET", Regex: "\\(" },
            { Type: "CLOSE_BRACKET", Regex: "\\)" }
        ],
        project_name: project_name
    });

    const lexer_rules_params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let lexer_rules_res = http.post(lexer_rules_url, lexer_rules_data, lexer_rules_params);
    check(lexer_rules_res, {
      "status is 200": (r) => r.status === 200,
    });
    
    sleep(1);

    const lexer_url = "http://localhost:8080/api/lexing/lexer";

    const lexer_data = JSON.stringify({
          project_name: project_name
    });

    const lexer_params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let lexer_res = http.post(lexer_url, lexer_data, lexer_params);
    check(lexer_res, {
      "status is 200": (r) => r.status === 200,
    });
    //console.log(lexer_res)
    
    sleep(1);

    // PARSER

    const grammar_url = "http://localhost:8080/api/parsing/grammar";

    const grammar_data = JSON.stringify({
        variables: ["PROGRAM","FUNCTION_D","FUNCTION_BLOCK","CODE","STATEMENT"],
    terminals: ["KEYWORD","IDENTIFIER","OPERATOR","INTEGER","OPEN_BRACKET", "CLOSE_BRACKET", "OPEN_SCOPE", "CLOSE_SCOPE"],
    start: "PROGRAM",
    rules: [
        { "input": "PROGRAM", "output": ["FUNCTION_D", "PROGRAM"] },

        { "input": "FUNCTION_D", "output": ["KEYWORD","IDENTIFIER","CLOSE_BRACKET","OPEN_BRACKET","FUNCTION_BLOCK"] },
        { "input": "FUNCTION_BLOCK", "output": ["OPEN_SCOPE", "CODE", "CLOSE_SCOPE"] },

        { "input": "CODE", "output": ["STATEMENT", "CODE"] },
        { "input": "CODE", "output": ["STATEMENT"] },

        { "input": "STATEMENT", "output": ["KEYWORD","IDENTIFIER","KEYWORD"] },
        { "input": "STATEMENT", "output": ["IDENTIFIER","OPERATOR","INTEGER"] },
    ],

        project_name: project_name
    });

    const grammar_params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let grammar_res = http.post(grammar_url, grammar_data, grammar_params);
    check(grammar_res, {
      "status is 200": (r) => r.status === 200,
    });
    
    sleep(1);

    const delete_project_url = "http://localhost:8080/api/users/deleteProject";
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