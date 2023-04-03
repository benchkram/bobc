<h1 align="center">
  bobc
</h1>

<h3 align="center">
   A Remote Cache for <a href="https://bob.build">Bob</a>
</h3>

<div align="center">
  <img  width="100%" src="https://user-images.githubusercontent.com/17600197/226882033-c9f741ca-4d42-4d75-8f43-956c404824b2.gif" />
</div>

<br/><br/>

**bobc** is a lightweight, open source implementation of the bob Cloud Platform - https://bob.build.
It is implemented in Go, and it uses an AWS S3-compatible storage backend for storing build artifacts, along with
a Postgres database for storing projects and artifact metadata.

<br/>

## Getting Started

To build and run bobc locally, you must have the following tools present on your system:
- Bob with Nix (https://bob.build/docs/getting-started/installation/)
- Docker-Compose (https://docs.docker.com/get-docker/)

Bob is used to build the Go binary and the container image.
Docker-compose is used to ramp up a local environment with a Postgres database, MinIO (S3-compatible object storage)
and Adminer to aid in inspecting the database contents.

First of all, clone the repository and `cd` into it:

```bash
git clone https://github.com/benchkram/bobc
cd bobc
```


### Building
To build the bobc container, run the following command:

```bash
bob build container
```

This command will install any build dependencies (Go, Docker, GolangCI-Lint), bootstrap the project, build the bobc
binary and subsequently build the container.

### Running
To set up the docker-compose environment and start the server run:

```bash
export API_KEY="example-api-key"
docker compose up -d
```

You should now see bobc running on port 8100.

Note: MinIO at localhost requires a host alias to be set up in order to work properly.
You should add the following to your `/etc/hosts` file:

```bash
127.0.0.1       minio
```

### Example: Creating a project and pushing artifacts to it

You must create a project to be able to sync artifacts to the server.
To do so, open a new terminal session and use the following `curl` command:

```bash
curl -X POST http://localhost:8100/api/projects \
   -H "Content-Type: application/json" \
   -H "Authorization: Bearer $API_KEY" \
   -d '{"name": "bobc-example"}'
```

This will create a project named `bobc-example`.

Next, you should configure an authentication context for bob. We'll use the same API_TOKEN as bearer token:

```bash
export API_KEY="example-api-key"
bob auth init --token=$API_KEY
```

You can verify artifact sync is working by typing
```bash
cd example
bob build --insecure --push
```

We have to use the `--insecure` flag since bobc is running over HTTP by default locally.
We also pass the `--push` flag to instruct bob to attempt to push artifacts to the artifact store. Bob will not push artifacts upstream by default.

For more information on how to use bob with bobc, please refer to the official documentation: https://bob.build/docs/remote-cache#pushing--pulling-artifacts
