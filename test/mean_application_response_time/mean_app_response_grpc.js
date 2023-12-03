import grpc from 'k6/net/grpc'
import { check, sleep } from "k6";

export const options = {
    thresholds: {
        // Assert that 99% of requests finish within 3000ms.
        http_req_duration: ["p(99) < 3000"],
    },
    // Ramp the number of virtual users up and down
    stages: [
        // { duration: "60s", target: 10 },
        { duration: "10s", target: 10 },
    ],
};

const services = {
    "boat":   "127.0.0.1:50773",
    "beer":   "127.0.0.1:63583",
    "board":  "127.0.0.1:63586",
    "bread":  "127.0.0.1:63590",
    "coal":   "127.0.0.1:63592",
    "flour":  "127.0.0.1:63598",
    "gold":   "127.0.0.1:63600",
    "ironore":"127.0.0.1:63608",
    "log":    "127.0.0.1:63610",
    "meat":   "127.0.0.1:63612",
    "pig":    "127.0.0.1:63614",
    "coin"  : "127.0.0.1:63594",
    "iron"  : "127.0.0.1:63606",
    "sword"  : "127.0.0.1:63616",
    "tools"  : "127.0.0.1:63618",
}

const client = new grpc.Client();
export default function () {
    let serviceName = `${__ENV.SVC_NAME}`
    let base_url = services[serviceName]

    client.connect(base_url, {
        plaintext: true,
        timeout: "3s",
        reflect: true
    });

    let SVC_NAME = `${__ENV.SVC_NAME}`
    let package_name = `netmarks_${__ENV.SVC_NAME}`
    let service_name = SVC_NAME.charAt(0).toUpperCase() + SVC_NAME.slice(1);

    const response = client.invoke(`${package_name}.${service_name}/Produce`, {
        quantity: `${__ENV.Q}`,
        response_size: `${__ENV.RES}`,
    })

    check(response, {
        "status is OK": (r) => r && r.status === grpc.StatusOK,
    });

    sleep(1);
}
