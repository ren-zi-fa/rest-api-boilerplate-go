import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
    vus: 10,
    duration: "30s"
};

const API_BASE = "http://localhost:8080/api/v1";
const AUTH_TOKEN = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoiYWRtaW4iLCJzdWIiOiJhY2Nlc3NfdG9rZW4iLCJleHAiOjE3NTkwNTExNDIsImlhdCI6MTc1OTA1MDI0Mn0.MLXWo51QRhTdNLlo83N98oTtIcVSJOEdKd7kgWQnNfc";

export default function () {
    // GET all posts
    let getRes = http.get(`${API_BASE}/posts`, {
        headers: { "Accept": "application/json", "Authorization": AUTH_TOKEN }
    });
    check(getRes, { "GET /posts status is 200": (r) => r.status === 200 });

    // GET post by ID
    let getByIdRes = http.get(`${API_BASE}/posts/1`, {
        headers: { "Accept": "application/json", "Authorization": AUTH_TOKEN }
    });
    check(getByIdRes, { "GET /posts/1 status is 200": (r) => r.status === 200 });

    // POST new post
    let payload = JSON.stringify({
        "author": "john doe ",
        "title": "New Post",
        "content": "This is a new post."
    });
    let params = {
        headers: {
            "Content-Type": "application/json",
            "Authorization": AUTH_TOKEN
        }
    };
    let postRes = http.post(`${API_BASE}/posts`, payload, params);
    check(postRes, { "POST /posts status is 201": (r) => r.status === 201 });

    // DELETE post by ID
    let deleteRes = http.del(`${API_BASE}/posts/340`, null, {
        headers: { "Authorization": AUTH_TOKEN }
    });
    check(deleteRes, {
        "DELETE /posts/2 status is 200 or 204": (r) => r.status === 200 || r.status === 204
    });

    sleep(1);
}
