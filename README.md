# Babel-Revolution-IA
IA component for Babel Revolution

## API
### Request for censorship `/is_censored`

#### Resquest description
- `POST`
- `JSON` object

| Property  | Type        | Example                                                                     |
|-----------|-------------|-----------------------------------------------------------------------------|
| `message` | `string`    | `"J'aime la compote."`,`"Prends l'objet pour se protéger de la pluie"`, ... |

#### Response description
| HTTP Code   | Meaning                                           |
|-------------|---------------------------------------------------|
| `201`       | message analyzed successfully                     |
| `400`       | bad request                                       |
| `500` 	  | internal error while accessing external resources |

##### `JSON` object sent (if `201`)
| Property    | Type    | Example    |
|-------------|---------|------------|
|`is_censored`| `bool`  | `true`     |