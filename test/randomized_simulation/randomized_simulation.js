import http from "k6/http";
import { check, sleep } from "k6";

// Test configuration
export const options = {
    stages: [
        { duration: "10s", target: 50 },
        { duration: "120s", target: 50 },
        { duration: '10s', target: 0 }, // ramp-down to 0 users
    ],
};

// Step 0: Init Code
// regular services
const services = {
    "netmarks-boat":   "http://127.0.0.1:54513",
    "netmarks-beer":   "http://127.0.0.1:54508",
    "netmarks-board":  "http://127.0.0.1:54511",
    "netmarks-bread":  "http://127.0.0.1:54515",
    "netmarks-coal":   "http://127.0.0.1:54517",
    "netmarks-flour":  "http://127.0.0.1:54523",
    "netmarks-gold":   "http://127.0.0.1:54525",
    "netmarks-ironore":"http://127.0.0.1:54531",
    "netmarks-log":    "http://127.0.0.1:54533",
    "netmarks-meat":   "http://127.0.0.1:54535",
    "netmarks-pig":    "http://127.0.0.1:54537",
}

// services with no 20 dependencies
const base_service = {
    "netmarks-fish"  : "http://127.0.0.1:54521",
    "netmarks-grain" : "http://127.0.0.1:54527",
    "netmarks-tree"  : "http://127.0.0.1:54543",
    "netmarks-water" : "http://127.0.0.1:54545",
}

// services with more than 20 dependencies
const complex_service = {
    "netmarks-coin"  : "http://127.0.0.1:54519",
    "netmarks-iron"  : "http://127.0.0.1:54529",
    "netmarks-sword"  : "http://127.0.0.1:54539",
    "netmarks-tools"  : "http://127.0.0.1:54541",
}

const service_quantity = 3
const base_service_quantity = 5
const complex_service_quantity = 1

// Simulated user behavior
export default function () {
    for (let key in complex_service) {
        let addr = complex_service[key]
        let res = http.get(`${addr}?quantity=${complex_service_quantity}&response_size=${__ENV.RES}`);
        check(res, { "status was 200": (r) => r.status === 200 });
    }
    sleep(1);

    for (let key in base_service) {
        let addr = base_service[key]
        let res = http.get(`${addr}?quantity=${base_service_quantity}&response_size=${__ENV.RES}`);
        check(res, { "status was 200": (r) => r.status === 200 });
    }
    sleep(1);

    for (let key in services) {
        let addr = services[key]
        let res = http.get(`${addr}?quantity=${service_quantity}&response_size=${__ENV.RES}`);
        check(res, { "status was 200": (r) => r.status === 200 });
    }
    sleep(1);
}
