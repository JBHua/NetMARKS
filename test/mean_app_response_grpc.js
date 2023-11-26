import grpc from 'k6/net/grpc'
import { check, sleep } from "k6";

const client = new grpc.Client();
// client.load(['../services/'], '/beer/proto/beer.proto');

// Test configuration
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

// Simulated user behavior
export default function () {
    let base_url = "127.0.0.1:62792"

    client.connect(base_url, {
        plaintext: true,
        timeout: "2s",
        reflect: true
    });

    const response = client.invoke('netmarks_coal.Coal/Produce', {
        quantity: 1,
        response_size: "512b"
    })

    console.log(JSON.stringify(response.message));

    check(response, {
        "status is OK": (r) => r && r.status === grpc.StatusOK,
    });

    sleep(1);
}
