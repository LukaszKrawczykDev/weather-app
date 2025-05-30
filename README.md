## 🏷️ Tagowanie obrazów

Obrazy tagowane są automatycznie przy pomocy(https://github.com/docker/metadata-action):

- `latest` – zawsze wskazuje na najnowszą wersję aplikacji z gałęzi `main`.
- `sha-<hash>` – unikalny tag oparty na skrócie commit-u (`GIT_SHA`), np. `sha-1a2b3c4`.
- `semver` – wersje semantyczne (`v1.0.0`) dodawane przy tworzeniu tagów w repozytorium.

**Uzasadnienie**: taki schemat umożliwia jednoznaczną identyfikację wersji obrazu, łatwe wdrażanie (`latest`) i wersjonowanie produkcyjne (`semver`). (https://github.com/docker/metadata-action#tags-input)

## 💾 Tagowanie danych cache

Dane cache budowania są przechowywane w publicznym repozytorium Docker Hub:
```bash
docker.io/<DOCKERHUB_USERNAME>/weather-app:cache
```
Z wykorzystaniem eksportera `registry` w trybie `mode=max`, co pozwala na pełne współdzielenie cache między buildami.

**Uzasadnienie**: `registry` + `mode=max` to zalecany sposób buforowania w środowiskach CI/CD (https://docs.docker.com/build/cache/backends/#registry-cache-backend)).

## ⚙️ Etapy realizacji zadania

1. **Checkout kodu** – klonowanie repozytorium.
2. **Budowa obrazu `linux/amd64`** – tymczasowa budowa do analizy bezpieczeństwa.
3. **Skanowanie Trivy (CVE)** – sprawdzenie obrazu pod kątem podatności typu `HIGH/CRITICAL`, z `exit-code: 1`, co zatrzymuje pipeline w razie wykrycia zagrożeń. Sprawdzamy tylko jeden obraz (linux/amd64) poniewaz obie architektury są budowane z tego samego kodu źródłowego. (https://docs.docker.com/buildx/working-with-buildx/#building-multi-platform-images)
4. **Budowa multiarch (`linux/amd64` + `linux/arm64`)** – finalny obraz aplikacji.
5. **Publikacja do GHCR** – tylko jeśli obraz przeszedł skanowanie.
6. **Cache push/pull** – dane cache trafiają do publicznego repozytorium DockerHub.

---

