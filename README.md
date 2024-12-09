*# Babel-Revolution-IA
IA component for Babel Revolution

## API
### Request for censorship `api/CheckMsg`

#### Resquest description
- `POST`
- `JSON` object

| Property  | Type        | Example                                                                     |
|-----------|-------------|-----------------------------------------------------------------------------|
| `title`   | `string`    | `"Instructions"`, ...                                                       |
| `text`    | `string`    | `"J'aime la compote."`,`"Prends l'objet pour se protéger de la pluie"`, ... |
| `author`  | `string`    | `"Mike"`, `"Corinne"`, ...                                                  |
| `parent`  | `string`    | see svelt-revolution specs                                                  |
| `session` | `string`    | see svelt-revolution specs                                                  |
| `side`    | `string`    | `"Pour"`,`"QG"`, ...                                                        |
| `isCensored` | `bool`   |                                                                             |
| `triggerNewEvent` | `bool` ||
| `qg` | `Node`||
| `terrain`| `Node`||


#### Response description
| HTTP Code   | Meaning                                           |
|-------------|---------------------------------------------------|
| `201`       | message analyzed successfully                     |
| `400`       | bad request                                       |
| `409`       | the session doesn't exists                        |
| `500` 	  | internal error while accessing external resources |

##### `JSON` object sent (if `201`)
Same as request.

### Request for new session `api/NewSession`

#### Resquest description
- `POST`
- `JSON` object

| Property     | Type        | Example                                                                                 |
|--------------|-------------|-----------------------------------------------------------------------------------------|
| `session`    | `string`    | see svelt-revolution specs                                                              |
| `bannedWords`| `[]string`  | `"clavier", "parapluie", "flaque", "écran", "machine", "IA", "SOPHIA", "détruire"`, ... |

#### Response description
| HTTP Code   | Meaning                                           |
|-------------|---------------------------------------------------|
| `201`       | session created successfully                      |
| `400`       | bad request                                       |
| `409`       | the session already exists                        |

##### `JSON` object sent (if `201`)
| Property     | Type        | Example                                                                                 |
|--------------|-------------|-----------------------------------------------------------------------------------------|
| `session`    | `string`    | see svelt-revolution specs                                                              |

## How to run with Docker

Download Docker [here](https://www.docker.com/products/docker-desktop)

This project is part of a bigger project, Babel Revolution. You can find the main project [here](https://github.com/KoroSensei10/svelte-revolution). You are more likely to run this project with the main project with the following instruction :

### Run with the main project

Run the following command at the root of the main project :
```bash
docker compose up --build
```

## Run this project only

### First method, docker compose

Run the following command at the root of this project :
```bash
docker compose up --build
```

### Second method, docker build and run

#### Build the image
```bash
docker build -t babel-revolution-ia .
```

#### Run the container
```bash
docker run -p 8000:8000 babel-revolution-ia
```

## Test the API
```bash
curl -X POST -H "Content-Type: application/json" -d '{"message":"Prends l objet pour se protéger de la pluie"}' http://localhost:8000/is_censored
```
