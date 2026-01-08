Projet API-GO-CRUD (Gestion de Livres)
Ce projet est une API REST développée en Go permettant de gérer une bibliothèque de livres. Il utilise le routeur chi, interagit avec une base de données PostgreSQL via Docker, et respecte les standards de développement (validation, codes HTTP, logs).

1. Prérequis
Docker et Docker Compose.

Go (1.21+) pour l'exécution locale.

2. Installation et Lancement (DB + API)
Lancement de la base de données
À la racine du projet :

Bash

docker-compose up -d
Cela initialise automatiquement la table books via les migrations SQL.

Lancement de l'API
Bash

go run cmd/api/main.go
L'API tourne sur : http://localhost:8080.

3. Exemples de requêtes (curl)
Voici les commandes pour tester chaque point d'entrée de l'API :

Créer un livre (POST)
Bash

curl -X POST http://localhost:8080/books \
     -H "Content-Type: application/json" \
     -d '{"title":"Dune", "author":"Frank Herbert", "year":1965}'
Réponse attendue : {"id": 1}.

Lister les livres (GET)
Bash

curl http://localhost:8080/books
Réponse attendue : Un tableau JSON des livres.

Récupérer un livre par son ID (GET)
Bash

curl http://localhost:8080/books/1
Mettre à jour un livre (PUT)
Bash

curl -X PUT http://localhost:8080/books/1 \
     -H "Content-Type: application/json" \
     -d '{"title":"Dune Messiah", "author":"Frank Herbert", "year":1969}'
Supprimer un livre (DELETE)
Bash

curl -X DELETE http://localhost:8080/books/1
4. Livrables inclus
Code source complet sur ce dépôt Git.

Infrastructure : Fichier docker-compose.yml et scripts SQL fournis.

Documentation : Ce présent README.md avec instructions de lancement.

5. Fonctionnalités techniques
Validation : Les champs title, author et year sont obligatoires.

Logs : Affichage en temps réel des requêtes dans le terminal.

Codes HTTP : 201 (Création), 200 (OK), 204 (Suppression), 400/404 (Erreurs).