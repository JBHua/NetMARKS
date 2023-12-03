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

// const services = {
//     "netmarks-boat":   "http://127.0.0.1:50773",
//     "netmarks-beer":   "http://127.0.0.1:63583",
//     "netmarks-board":  "http://127.0.0.1:63586",
//     "netmarks-bread":  "http://127.0.0.1:63590",
//     "netmarks-coal":   "http://127.0.0.1:63592",
//     "netmarks-flour":  "http://127.0.0.1:63598",
//     "netmarks-gold":   "http://127.0.0.1:63600",
//     "netmarks-ironore":"http://127.0.0.1:63608",
//     "netmarks-log":    "http://127.0.0.1:63610",
//     "netmarks-meat":   "http://127.0.0.1:63612",
//     "netmarks-pig":    "http://127.0.0.1:63614",
//     "netmarks-coin"  : "http://127.0.0.1:63594",
//     "netmarks-iron"  : "http://127.0.0.1:63606",
//     "netmarks-sword"  : "http://127.0.0.1:63616",
//     "netmarks-tools"  : "http://127.0.0.1:63618",
// }

// Simulated user behavior
export default function () {
    let base_url = "http://127.0.0.1:63618/"
    let res = http.get(`${base_url}?quantity=${__ENV.Q}&response_size=${__ENV.RES}`);

    // Validate response status
    check(res, { "status was 200": (r) => r.status === 200 });
    sleep(1);
}
