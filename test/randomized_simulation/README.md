# Randomized Simulation

## Motivation
To quote the original paper: `When NetMARKS sees some application for the first time, it does not have any communication statistics so it refrains from voting by returning score zero for all nodes. That is why before the actual test had begun, our application was deployed to allow NetMARKS to collect communication metrics.`

This test is used to simulate how user will use our services.

## Setup
Aim: Soak test (Average-Load test)

Virtual User (VU) count: 200. This means a maximum of 200 parallel requests
Stages: 