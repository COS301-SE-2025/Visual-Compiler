import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    {duration: '20s', target: 5},
    {duration: '30', target: 5},
    {duration: '20s', target: 0},
  ]
}
export default function () {
  let res = http.post('http://localhost:8080/api/users/register', {
    email: "invalidemail",
    password: "invalidpass",
    username: "username"
  });
  check(res, {
    'status is 400': (r) => r.status === 400,
  });
  sleep(1);
}