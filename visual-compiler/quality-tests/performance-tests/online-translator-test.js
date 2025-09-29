import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
      {duration: '20s', target: 10},
      {duration: '20s', target: 10},
      {duration: '20s', target: 0},
    ],
    thresholds: {
      http_req_duration: ["p(95)<3000"],
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
    const source_code_url = "https://www.visual-compiler.co.za/api/lexing/code";

    const source_code_data = JSON.stringify({
            source_code: `int blue = 13;

int new(int red)
{
    red = red + 1;
    return red;
}

int _i = 0;
for _i range(12)
{
    blue = new(blue);
    print(blue);
}

int change(int r)
{
    r = r + 1;
    return r;
}
int _j = 0;
for _j range(12)
{
    blue = change(blue);
    print(blue);
}
    
int next(int re)
{
    re = re + 1;
    return re;
}
int _k = 0;
for _k range(12)
{
    blue = next(blue);
    print(blue);
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

    sleep(2);

    // LEXER

    const lexer_rules_url = "https://www.visual-compiler.co.za/api/lexing/rules";

    const lexer_rules_data = JSON.stringify({
        pairs: [
            { Type: 'KEYWORD', Regex: 'int|return|print'},
            { Type: 'CONTROL', Regex: 'for|range'},
            { Type: 'IDENTIFIER', Regex: '[a-zA-Z_]+'},
            { Type: 'INTEGER', Regex: '[0-9]+'},
            { Type: 'ASSIGNMENT', Regex: '='},
            { Type: 'OPERATOR', Regex: '[+\\-*/%]' },
            { Type: 'DELIMITER', Regex: ';'},
            { Type: 'OPEN_BRACKET', Regex: '\\('},
            { Type: 'CLOSE_BRACKET', Regex: '\\)'},
            { Type: 'OPEN_SCOPE', Regex: '\{'},
            { Type: 'CLOSE_SCOPE', Regex: '\}'}
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
    
    sleep(2);

    const lexer_url = "https://www.visual-compiler.co.za/api/lexing/lexer";

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
    sleep(2);

    // PARSER

    const grammar_url = "https://www.visual-compiler.co.za/api/parsing/grammar";

    const grammar_data = JSON.stringify({
        variables: ["PROGRAM", "STATEMENT", "FUNCTION", "ITERATION", "DECLARATION", "ELEMENT", "TYPE", "EXPRESSION", "FUNCTION_DEFINITION", "FUNCTION_BLOCK", "RETURN", "ITERATION_DEFINITION", "ITERATION_BLOCK", "PARAMETER", "PRINT"],
        terminals: ["KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "DELIMITER", "OPEN_BRACKET", "CLOSE_BRACKET", "OPEN_SCOPE", "CLOSE_SCOPE", "CONTROL"],
        start: "PROGRAM",
        rules: [
            { "input": "PROGRAM", "output": ["STATEMENT", "FUNCTION", "STATEMENT", "ITERATION", "FUNCTION", "STATEMENT", "ITERATION", "FUNCTION", "STATEMENT", "ITERATION"] },
            { "input": "STATEMENT", "output": ["DECLARATION", "DELIMITER"] },
            { "input": "DECLARATION", "output": ["TYPE", "IDENTIFIER", "ASSIGNMENT", "ELEMENT"] },
            { "input": "DECLARATION", "output": ["IDENTIFIER", "ASSIGNMENT", "EXPRESSION"] },
            { "input": "DECLARATION", "output": ["IDENTIFIER", "ASSIGNMENT", "IDENTIFIER", "PARAMETER"] },
            { "input": "TYPE", "output": ["KEYWORD"] },
            { "input": "EXPRESSION", "output": ["ELEMENT", "OPERATOR", "ELEMENT"] },
            { "input": "ELEMENT", "output": ["INTEGER"] },
            { "input": "ELEMENT", "output": ["IDENTIFIER"] },
            { "input": "FUNCTION", "output": ["FUNCTION_DEFINITION", "FUNCTION_BLOCK"] },
            { "input": "FUNCTION_DEFINITION", "output": ["TYPE", "IDENTIFIER", "PARAMETER"] },
            { "input": "FUNCTION_BLOCK", "output": ["OPEN_SCOPE", "STATEMENT", "RETURN", "CLOSE_SCOPE"] },
            { "input": "RETURN", "output": ["KEYWORD", "ELEMENT", "DELIMITER"] },
            { "input": "ITERATION", "output": ["ITERATION_DEFINITION", "ITERATION_BLOCK"] },
            { "input": "ITERATION_DEFINITION", "output": ["CONTROL", "IDENTIFIER", "CONTROL", "PARAMETER"] },
            { "input": "ITERATION_BLOCK", "output": ["OPEN_SCOPE", "STATEMENT", "PRINT", "CLOSE_SCOPE"] },
            { "input": "PARAMETER", "output": ["OPEN_BRACKET", "ELEMENT", "CLOSE_BRACKET"] },
            { "input": "PARAMETER", "output": ["OPEN_BRACKET", "TYPE", "IDENTIFIER", "CLOSE_BRACKET"] },
            { "input": "PRINT", "output": ["KEYWORD", "OPEN_BRACKET", "ELEMENT", "CLOSE_BRACKET", "DELIMITER"] }
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


    sleep(2);

    const parser_url = "https://www.visual-compiler.co.za/api/parsing/tree";

    const parser_data = JSON.stringify({
        project_name: project_name
    });

    const parser_params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let parser_res = http.post(parser_url, parser_data, parser_params);
    check(parser_res, {
      "status is 200": (r) => r.status === 200,
    });
    
    sleep(2);

    // ANALYSER

     const analyse_url = "https://www.visual-compiler.co.za/api/analysing/analyse";

    const analyse_data = JSON.stringify({
       scope_rules: [
            {"start": "\\{", "end":"\\}"}
        ],
       type_rules: [
            {"ResultData": "int", "Assignment":"=", "LHSData":"INTEGER"},
            {"ResultData": "int", "Assignment":"=", "LHSData":"int"},
            {"ResultData": "int", "Assignment":"=", "LHSData":"int", "OPERATOR":["+"], "RHSData":"INTEGER"}
       ],
       grammar_rules: {
            "TypeRule": "TYPE",
            "VariableRule":   "IDENTIFIER",
            "FunctionRule":   "FUNCTION_DEFINITION",
            "ParameterRule":  "PARAMETER",
            "AssignmentRule": "ASSIGNMENT",
            "OperatorRule":   "OPERATOR",
            "TermRule" :      "ELEMENT"
    },
        project_name: project_name
    });

    const analyse_params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let analyse_res = http.post(analyse_url, analyse_data, analyse_params);
    check(analyse_res, {
      "status is 200": (r) => r.status === 200,
    });

    sleep(2);

    // TRANSLATOR

    const translator_rule_url = "https://www.visual-compiler.co.za/api/translating/readRules";

    const translator_rule_data = JSON.stringify({
        translation_rules: [
            { "sequence": ["KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "DELIMITER"], "translation": ['add     rax, {INTEGER}', 'mov     [{IDENTIFIER}], rax'] },
            { "sequence": ["KEYWORD", "IDENTIFIER", "OPEN_BRACKET", "KEYWORD", "IDENTIFIER", "CLOSE_BRACKET", "OPEN_SCOPE", "IDENTIFIER", "ASSIGNMENT", "IDENTIFIER", "OPERATOR", "INTEGER", "DELIMITER", "KEYWORD", "IDENTIFIER", "DELIMITER", "CLOSE_SCOPE"], "translation": ['func {IDENTIFIER}:', '     mov     rbx, [{IDENTIFIER}]', '     add     rbx, {INTEGER}', '     mov     [{IDENTIFIER}], rbx', '     return'] },
            { "sequence": ["CONTROL", "IDENTIFIER", "CONTROL", "OPEN_BRACKET", "INTEGER", "CLOSE_BRACKET", "OPEN_SCOPE", "IDENTIFIER", "ASSIGNMENT", "IDENTIFIER", "OPEN_BRACKET", "IDENTIFIER", "CLOSE_BRACKET", "DELIMITER", "KEYWORD", "OPEN_BRACKET", "IDENTIFIER", "CLOSE_BRACKET", "DELIMITER", "CLOSE_SCOPE"], "translation": ['func {CONTROL}:', '     jump    [{IDENTIFIER}], {INTEGER}', '     param   rcx, [{IDENTIFIER}]', '     call    {IDENTIFIER}', '     print   [{IDENTIFIER}]'] },
        ],
        project_name: project_name
    });

    const translator_rule_params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let translator_rule_res = http.post(translator_rule_url, translator_rule_data, translator_rule_params);
    check(translator_rule_res, {
      "status is 200": (r) => r.status === 200,
    });


    sleep(2);

    const translator_url = "https://www.visual-compiler.co.za/api/translating/translate";

    const translator_data = JSON.stringify({
        project_name: project_name
    });

    const translator_params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let translator_res = http.post(translator_url, translator_data, translator_params);
    check(translator_res, {
      "status is 200": (r) => r.status === 200,
    });
    
    sleep(2);

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
    
    sleep(2);

}