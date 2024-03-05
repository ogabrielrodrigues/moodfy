# Moodfy
### O back-end do Moodfy é construido completamente em [Go](https://go.dev/).

### Rotas disponíveis:
- `POST /artist` – Para criar um artista.
- `GET /artist` – Para listar todos os artistas.
- `POST /style` – Para criar um estilos.
- `GET /style` – Para listar todos os estilos.
- `POST /music` – Para criar uma música.
- `GET /music?artist=[:nome_do_artista]&style=[:nome_do_estilo]` – Para listar todos as músicas ou filtrar.

### Criação de Artista
`POST /artist`

| atributo | descrição |
| --- | --- |
| **name** | obrigatório, string de no mínimo 3 caracteres e no máximo 100 caracteres. |

Exemplo de request:
```json
{
  "name": "Don Toliver"
}
```

### Listagem de Artista
`GET /artist`

Exemplo de response:
```json
[
  {
    "id": 1,
    "name": "Don Toliver"
  }
]
```

### Criação de Estilo
`POST /style`

| atributo | descrição |
| --- | --- |
| **name** | obrigatório, string de no mínimo 3 caracteres e no máximo 50 caracteres. |

Exemplo de request:
```json
{
    "name": "Trap"
}
```

### Listagem de Estilo
`GET /style`

Exemplo de response:
```json
[
  {
    "id": 1,
    "name": "Trap"
  }
]
```

### Criação de Música
`POST /music`

| atributo | descrição |
| --- | --- |
| **artist_id** | obrigatório, inteiro, sendo um ID de um artista válido. |
| **name** | obrigatório, string de no mínimo 3 caracteres e no máximo 50 caracteres. |
| **spotify_link** | obrigatório, string, sendo uma URL válida e de no máximo 200 caracteres. |
| **styles** | opcional, array de inteiro, sendo os ID's dos estilos. |

Exemplo de request:
```json
{
  "artist_id": 1,
  "name": "No Idea",
  "spotify_link": "https://open.spotify.com/intl-pt/track/7AzlLxHn24DxjgQX73F9fU?si=b5a6e2cffec9405b",
  "styles": [1]
}
```

### Listagem de Música
`GET /music`

### Filtragem:
 - Para filtrar por artista, basta imbutir na url: `?artist=[:nome_do_artista]`
 - Para filtrar por estilo, basta imbutir na url: `?style=[:nome_do_estilo]`
 - Para filtrar por artista e estilo, basta imbutir na url: `?artist=[:nome_do_artista]&style=[:nome_do_estilo]`

Exemplo de response:
```json
[
  {
    "id": 1,
    "music": "No Idea",
    "artist": "Don Toliver",
    "cover_image": "https://i.scdn.co/image/ab67616d0000b27345190a074bef3e8ce868b60c",
    "spotify_link": "https://open.spotify.com/intl-pt/track/7AzlLxHn24DxjgQX73F9fU?si=b5a6e2cffec9405b",
    "styles": [
      "Trap"
    ]
  }
]
```

<br>
<h3 align="center" style="font-size: 20px; text-decoration:none;">Made with ❤️ by <a href="https://github.com/ogabrielrodrigues">ogabrielrodrigues</a></h3>
