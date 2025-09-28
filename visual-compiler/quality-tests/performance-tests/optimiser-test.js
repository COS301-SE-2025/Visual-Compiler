import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
      {duration: '20s', target: 50},
      {duration: '20s', target: 50},
      {duration: '20s', target: 0},
    ],
    thresholds: {
      http_req_duration: ["p(95)<5000"],
      http_req_failed: ["rate<0.01"],
    },
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


    const source_code_url = "http://localhost:8080/api/optimising/source_code";

    const source_code_data = JSON.stringify({
            source_code: `package main

import "fmt"

func main() {
    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }
    
    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }

    for blue := 1; blue < 5; blue++ {
        fmt.Println("[Block " + blue + "]")
        for red := 1; red <= blue; red++ {
            fmt.Println("	-> " + red)
            if blue == red {
                fmt.Println(blue + red)
            }
        }
    }
}
    
func nothing() (int) {
    var random int = 13
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


    const optimiser_url = "http://localhost:8080/api/optimising/optimise";

    const optimiser_data = JSON.stringify({
        project_name: project_name,
        constant_folding: true,
        dead_code: true,
        loop_unrolling: true
    });

    const optimiser__params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let optimiser_res = http.post(optimiser_url, optimiser_data, optimiser__params);
    check(optimiser_res, {
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