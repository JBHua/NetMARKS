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
        { duration: "10s", target: 10 },
    ],
};

const services = {
    "boat":   "http://127.0.0.1:50773",
    "beer":   "http://127.0.0.1:63583",
    "board":  "http://127.0.0.1:63586",
    "bread":  "http://127.0.0.1:63590",
    "coal":   "http://127.0.0.1:63592",
    "flour":  "http://127.0.0.1:63598",
    "gold":   "http://127.0.0.1:63600",
    "ironore":"http://127.0.0.1:63608",
    "log":    "http://127.0.0.1:63610",
    "meat":   "http://127.0.0.1:63612",
    "pig":    "http://127.0.0.1:63614",
    "coin"  : "http://127.0.0.1:63594",
    "iron"  : "http://127.0.0.1:63606",
    "sword"  : "http://127.0.0.1:63616",
    "tools"  : "http://127.0.0.1:63618",
}

// Simulated user behavior
export default function () {
    let serviceName = `${__ENV.SVC_NAME}`
    let base_url = services[serviceName]
    let res = http.get(`${base_url}?quantity=${__ENV.Q}&response_size=${__ENV.RES}`);

    // Validate response status
    check(res, { "status was 200": (r) => r.status === 200 });
    sleep(1);
}
