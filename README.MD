# Distributed In-Memory Key-Value Store (Nebula) Documentation

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Technologies Used](#technologies-used)
   <!-- 4. [Getting Started](#getting-started) -->
   <!-- 5. [Configuration](#configuration) -->
4. [Usage](#usage)
   <!-- 7. [API Reference](#api-reference) -->
   <!-- 8. [Contributing](#contributing) -->
5. [Author Information](#author-information)
6. [Contact](#contact)
<!-- 11. [License](#license) -->

## Introduction

This document provides an overview of the application, its features, and instructions for installation, configuration, and usage.

Nebula is a distributed in-memory key-value store that uses consistent hashing to distribute data across multiple stores, ensuring high availability and fault tolerance.

## Features

- Distributed architecture using consistent hashing.
- In-memory storage for fast data access.
- Support for GET, PUT, and DELETE operations.
  <!-- - Data replication for fault tolerance. -->
  <!-- - TLS encryption for secure communication. -->

## Technologies Used

DIM-KVS is built using the following technologies and libraries:

- Golang: The primary programming language.
- gRPC: Used for communication between components.
  <!-- - TLS: For secure communication. -->
  <!-- - (List any other technologies or libraries you used) -->

## Getting Started

Follow these steps to get started with NEBULA (DKVS):

1. Download and install Golang from the [official website](https://golang.org/dl/).

2. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/priyansh32/nebula.git
   ```

3. Install the dependencies:

   ```bash
    go mod download
   ```

4. Build the project components:

   ```bash
    chmod +x ./build.sh
    ./build.sh
   ```

5. Run the coordinator:

   ```bash
    ./bin/coordinator <port> <replication-factor>
   ```

6. Run the store:

   ```bash
    ./bin/store <port>
   ```

   you can run multiple stores by running the above command with different ports.

7. Run the CLI:

   ```bash
    ./bin/cli <coordinator-address>
   ```

   the coordinator address is the address of the coordinator in the format `host:port`.

## Usage

Application provides a CLI to interact with the coordinator. Below is a sample usage of CLI.  
To interact with the coordinator programmatically, you need to use gRPC with protobuf specification the `/internal/api/coordinator/coordinator.proto`

```bash

NEBULA> ADDSTORE localhost:50012 alpha-store
Added store:  alpha-store  with status:  OK

NEBULA> ADDSTORE localhost:50013 beta-store
Added store:  beta-store  with status:  OK

NEBULA> PUT email patidarpriyansh936@gmail.com
Status:  OK

NEBULA> PUT name priyanshpatidar
Status:  OK

NEBULA> GET name
Value:  priyanshpatidar

NEBULA> DELETE email
Status:  OK

NEBULA> GET email
CACHE MISS

NEBULA> EXIT
```

## **Author Information**

- Author: Priyansh Patidar
- GitHub: **[@priyansh32](https://github.com/priyansh32)**

## **Contact**

If you have questions or need assistance, feel free to reach out to the author:

- Email: **[patidarpriyansh936@gmail.com](mailto:patidarpriyansh936@gmail.com)**
- Twitter: **[@priyanshh32](https://twitter.com/priyanshh32)**
