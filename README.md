# go-boilerplate-rest-api-chi

## Présentation générale

Ce projet est un **boilerplate d'API REST** en Go, utilisant le framework [Chi](https://github.com/go-chi/chi), pensé pour servir de base solide à tout projet backend. Il adopte une **architecture monolithique modulaire** : chaque domaine métier est isolé, ce qui facilite la maintenance, l’évolutivité et la clarté du code.

---

## Table des matières

- [go-boilerplate-rest-api-chi](#go-boilerplate-rest-api-chi)
  - [Présentation générale](#présentation-générale)
  - [Table des matières](#table-des-matières)
  - [Architecture du projet](#architecture-du-projet)
  - [Technologies utilisées](#technologies-utilisées)
  - [Gestion de la configuration](#gestion-de-la-configuration)
  - [Fichiers d'environnement](#fichiers-denvironnement)
  - [Utilisation de Docker](#utilisation-de-docker)
  - [Automatisation avec Task](#automatisation-avec-task)
  - [Tests et qualité](#tests-et-qualité)
  - [Intégration Continue (CI)](#intégration-continue-ci)
  - [Documentation API (Swagger \& Scalar)](#documentation-api-swagger--scalar)
  - [Collections Bruno](#collections-bruno)
  - [Démarrage rapide](#démarrage-rapide)

---

## Architecture du projet

L’architecture est **monolithique modulaire** :

- Chaque domaine (ex : `author`, `book`) possède son propre dossier dans `internal/`, avec ses handlers, services, repositories, DTO, erreurs, etc.
- Les dépendances sont injectées explicitement, ce qui facilite les tests et la maintenance.
- La configuration, la base de données, le logger, la validation, la gestion des réponses et les mocks sont tous séparés dans des modules dédiés.
- Les routes sont centralisées dans `internal/api`.
- Le point d’entrée de l’application se trouve dans `cmd/go-boilerplate-rest-api-chi/main.go`.

**Exemple d’organisation :**

```
internal/
	author/
		handler.go
		service.go
		repository.go
		dto/
	book/
		...
	api/
	config/
	database/
	entity/
	logger/
	mocks/
	response/
	test-utils/
	validator/
```

---

## Technologies utilisées

- **Go** (>=1.24)
- **Chi** : router HTTP léger et performant
- **GORM** : ORM pour la gestion de la base de données (MySQL/MariaDB par défaut)
- **Zerolog** : logging structuré, performant et lisible
- **Go-Playground Validator** : validation des entrées utilisateur (DTO)
- **Swaggo** : génération automatique de documentation Swagger
- **Scalar** : UI moderne pour la doc Swagger
- **Bruno** : gestionnaire de collections de requêtes API (voir [`bruno-collection/`](bruno-collection/))
- **Task** : automatisation des tâches de développement
- **Docker & Docker Compose** : conteneurisation de l’API et de la base de données
- **Testify, GoMock** : tests unitaires et mocks

---

## Gestion de la configuration

La configuration est centralisée dans [`internal/config`](internal/config/). Elle est chargée à partir de fichiers d’environnement (`.env`) et de variables système. Les paramètres de connexion à la base de données, au serveur, au logger, etc., sont tous configurables.

- **.env.example** : modèle à copier pour créer votre propre `.env`.
- **.env** : contient les variables locales
- Les valeurs sont lues automatiquement au démarrage.

**Exemple de variables :**

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=secret
DB_NAME=go_boilerplate
API_PORT=8080
LOG_LEVEL=debug
```

---

## Fichiers d'environnement

- **.env.example** : référence de toutes les variables nécessaires.
- **.env** : à personnaliser selon votre environnement local ou de production.
- **Sécurité** : `.env` est ignoré par git pour éviter toute fuite de secrets.

---

## Utilisation de Docker

Le projet est prêt à l’emploi avec Docker :

- **docker-compose.yml** orchestre l’API et la base de données (MariaDB).
- **Dockerfile** construit une image Go statique, légère et sécurisée (distroless).
- Les volumes assurent la persistance des données.
- Les variables d’environnement sont injectées automatiquement.

**Commandes utiles :**

```sh
docker compose up --build
```

---

## Automatisation avec Task

Le projet utilise [Task](https://taskfile.dev) pour automatiser les tâches courantes :

- **Formatage** : `task format` (go fmt, goimports)
- **Lint** : `task lint` (staticcheck, golangci-lint)
- **Tests** : `task test` (unitaires), `task test-cover` (avec couverture)
- **Génération de documentation** : `task doc` (Swagger)
- **Génération des mocks** : `task generate`
- **Démarrage complet (API + DB)** : `task dev`

**Exemple de workflow développeur :**

```sh
cp .env.example .env
task dev
```

---

## Tests et qualité

- **Tests unitaires** : chaque module métier possède ses propres tests (`*_test.go`).
- **Mocks** : générés automatiquement pour isoler les dépendances (voir `internal/mocks/`).
- **Test d’intégration** : possible via des utilitaires dans `internal/test-utils/`.
- **Couverture** : mesurée et exportée dans `coverage.txt`.
- **Commandes** :
	- `task test` : lance tous les tests
	- `task test-cover` : lance les tests avec rapport de couverture

---

## Intégration Continue (CI)

Le projet est prêt pour la CI avec **GitHub Actions** (voir `.github/workflows/ci.yml`) :

- **Lint** : vérification du style et des erreurs statiques
- **Tests** : exécution automatique à chaque push/PR
- **Couverture** : rapport publié
- **Build Docker** : vérification que l’image se construit correctement

---

## Documentation API (Swagger & Scalar)

- **Swagger** : la documentation OpenAPI est générée automatiquement à partir des annotations dans le code (voir `docs/`).
- **Scalar** : une UI moderne pour explorer et tester l’API, accessible sur `/api/docs` en local.
- **Mise à jour** : `task doc` régénère la documentation après modification des routes ou des schémas.

---

## Collections Bruno

Le dossier [`bruno-collection/`](bruno-collection/) contient des collections de requêtes prêtes à l’emploi pour [Bruno](https://www.usebruno.com/), un outil open-source pour tester et documenter les APIs :

- **CRUD complet** sur les entités (auteur, livre, etc.)
- **Healthcheck**
- **Organisation par dossier**
- **Environnements** (local, prod, etc.)

Importez la collection dans Bruno pour tester rapidement tous les endpoints de l’API.

---

## Démarrage rapide

1. Copier le fichier d’exemple d’environnement :
	 ```sh
	 cp .env.example .env
	 ```
2. Lancer l’environnement de dev (API + DB) :
	 ```sh
	 task dev
	 ```
3. Accéder à la documentation interactive :
	 - [http://localhost:8080/api/docs](http://localhost:8080/api/docs)
