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

const client = new grpc.Client();
export default function () {
    let base_url = `127.0.0.1:${__ENV.PORT}`
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
