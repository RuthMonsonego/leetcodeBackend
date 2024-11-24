# LeetCode Backend API

## Overview
This project provides a backend API for managing coding questions, with CRUD operations using MySQL and Gin framework. The execution environment supports code written in Go and Python.

## Prerequisites
- Docker and Docker Compose installed

## Setup Instructions
1. Clone the repository:
    ```bash
    git clone https://github.com/RuthMonsonego/leetcodeBackend
    cd leetcodeBackend
    ```

2. Build and run the project using Docker Compose:
    ```bash
    docker-compose up --build
    ```

3. The API will be available at `http://localhost:8080`.

## Supported Languages
- Go
- Python

## Supported Data Types
- `int`
- `double`
- `float`
- `string`
- `char`
- `bool`
- `int[]`
- `double[]`
- `float[]`
- `string[]`
- `char[]`
- `bool[]`

## API Endpoints
- `GET /questions`: Fetch all questions.
- `POST /questions`: Add a new question.
- `GET /questions/:code`: Fetch a question by code.
- `PUT /questions/:code`: Update a question by code.
- `DELETE /questions/:code`: Delete a question by code.