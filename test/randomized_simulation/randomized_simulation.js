import http from "k6/http";
import { check, sleep } from "k6";

// Test configuration
export const options = {
    stages: [
        { duration: "10s", target: 10 },
        { duration: "60s", target: 10 },
        { duration: '10s', target: 0 }, // ramp-down to 0 users
    ],
};

// Step 0: Init Code
// regular services
const services = {
    "netmarks-boat":   "http://127.0.0.1:50773",
    "netmarks-beer":   "http://127.0.0.1:63583",
    "netmarks-board":  "http://127.0.0.1:63586",
    "netmarks-bread":  "http://127.0.0.1:63590",
    "netmarks-coal":   "http://127.0.0.1:63592",
    "netmarks-flour":  "http://127.0.0.1:63598",
    "netmarks-gold":   "http://127.0.0.1:63600",
    "netmarks-ironore":"http://127.0.0.1:63608",
    "netmarks-log":    "http://127.0.0.1:63610",
    "netmarks-meat":   "http://127.0.0.1:63612",
    "netmarks-pig":    "http://127.0.0.1:63614",
}

// services with more than 20 dependencies
const complex_service = {
    "netmarks-coin"  : "http://127.0.0.1:63594",
    "netmarks-iron"  : "http://127.0.0.1:63606",
    "netmarks-sword"  : "http://127.0.0.1:63616",
    "netmarks-tools"  : "http://127.0.0.1:63618",
}

const service_quantity = 1
const complex_service_quantity = 1

// Simulated user behavior
export default function () {
    for (let key in services) {
        let addr = services[key]
        let res = http.get(`${addr}?quantity=${service_quantity}&response_size=${__ENV.RES}`);
        check(res, { "status was 200": (r) => r.status === 200 });
    }
    sleep(1);

    for (let key in complex_service) {
        let addr = complex_service[key]
        let res = http.get(`${addr}?quantity=${complex_service_quantity}&response_size=${__ENV.RES}`);
        check(res, { "status was 200": (r) => r.status === 200 });
    }
    sleep(1);
}
