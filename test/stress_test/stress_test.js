import http from "k6/http";
import { check, sleep } from "k6";

// Test configuration
export const options = {
    thresholds: {
        // Assert that 99% of requests finish within 3000ms.
        http_req_duration: ["p(99) < 3000"],
    },
    // Ramp the number of virtual users up and down
    stages: [
        // { duration: "10s", target: 200 },
        { duration: "120s", target: 2000 },
        // { duration: '10s', target: 0 }, // ramp-down to 0 users
    ],
};

// Simulated user behavior
export default function () {
    let base_url = "http://127.0.0.1:8080"
    let res = http.get(`${base_url}?quantity=${__ENV.Q}&response_size=${__ENV.RES}`);

    // Validate response status
    check(res, { "status was 200": (r) => r.status === 200 });
    sleep(1);
}
