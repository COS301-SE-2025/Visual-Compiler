import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
      {duration: '10s', target: 5},
      {duration: '20s', target: 5},
      {duration: '10s', target: 0},
    ],
    thresholds: {
      http_req_duration: ["p(95)<500"],
      http_req_failed: ["rate<0.01"],
    },
}

export default function () {
    const url = "http://localhost:8080/api/users/login";

    const login_data = JSON.stringify({
          login: "halfstack.guest@gmail.com",
          password: "guestUser13",
    });

    const params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer guestuser"
        },
    };

    let res = http.post(url, login_data, params);

    check(res, {
      "status is 200": (r) => r.status === 200,
    });
    
    sleep(1);
}