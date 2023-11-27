# Randomized Simulation

## Motivation
To quote the original paper: `When NetMARKS sees some application for the first time, it does not have any communication statistics so it refrains from voting by returning score zero for all nodes. That is why before the actual test had begun, our application was deployed to allow NetMARKS to collect communication metrics.`

This test is used to simulate how user will use our services.

## Setup
Aim: Soak test (Average-Load test)

Virtual User (VU) count: 200. This means a maximum of 200 parallel requests
Stages: 

socat TCP4-LISTEN:54508,fork,reuseaddr TCP4:127.0.0.1:54508 | \
socat TCP4-LISTEN:54511,fork,reuseaddr TCP4:127.0.0.1:54511 | \
socat TCP4-LISTEN:54513,fork,reuseaddr TCP4:127.0.0.1:54513 | \
socat TCP4-LISTEN:54515,fork,reuseaddr TCP4:127.0.0.1:54515 | \
socat TCP4-LISTEN:54517,fork,reuseaddr TCP4:127.0.0.1:54517 | \
socat TCP4-LISTEN:54519,fork,reuseaddr TCP4:127.0.0.1:54519 | \
socat TCP4-LISTEN:54521,fork,reuseaddr TCP4:127.0.0.1:54521 | \
socat TCP4-LISTEN:54523,fork,reuseaddr TCP4:127.0.0.1:54523 | \
socat TCP4-LISTEN:54525,fork,reuseaddr TCP4:127.0.0.1:54525 | \
socat TCP4-LISTEN:54527,fork,reuseaddr TCP4:127.0.0.1:54527 | \
socat TCP4-LISTEN:54529,fork,reuseaddr TCP4:127.0.0.1:54529 | \
socat TCP4-LISTEN:54531,fork,reuseaddr TCP4:127.0.0.1:54531 | \
socat TCP4-LISTEN:54533,fork,reuseaddr TCP4:127.0.0.1:54533 | \
socat TCP4-LISTEN:54535,fork,reuseaddr TCP4:127.0.0.1:54535 | \
socat TCP4-LISTEN:54537,fork,reuseaddr TCP4:127.0.0.1:54537 | \
socat TCP4-LISTEN:54539,fork,reuseaddr TCP4:127.0.0.1:54539 | \
socat TCP4-LISTEN:54541,fork,reuseaddr TCP4:127.0.0.1:54541 | \
socat TCP4-LISTEN:54543,fork,reuseaddr TCP4:127.0.0.1:54543 | \
socat TCP4-LISTEN:54545,fork,reuseaddr TCP4:127.0.0.1:54545