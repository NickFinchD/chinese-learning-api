# Wojiao

Веб-приложение для изучения китайского языка: курсы с уроками (слова/квизы), словарь, тексты для чтения, тесты (грамматика по HSK + тренировка «слово → перевод») со spaced-repetition прогрессом.

Монорепозиторий: `backend/` (API) + `frontend/` (SPA).

## Стек

- **Backend:** Go 1.25, Gin, pgx/pgxpool (PostgreSQL), JWT, godotenv
- **Frontend:** Vue 3, TypeScript, Vite, Pinia, Vue Router, Tailwind CSS 4, Axios
- **БД:** PostgreSQL 17 (Docker или локально)

## Архитектура

**Backend** — package-by-feature: каждый домен лежит в `backend/internal/<domain>` и внутри разбит на слои `Handler → Service → Repository` (+ `model`, `routes`, `request`/`response`). Для тестируемости `Service` зависит от небольшого unexported-интерфейса `repository`, а не от конкретной структуры. Доступ к БД — через пул соединений `pgxpool.Pool` (не одиночный `*pgx.Conn` — тот не потокобезопасен под конкурентной нагрузкой).

Домены: `auth`, `users`, `courses`, `lessons`, `words`, `savedwords`, `quizzes`, `progress`, `learning`, `texts`.

**Frontend** — страницы/лейауты/компоненты + Pinia-сторы + сервисы-обёртки над axios, по одному сервису/стору на домен.

## Быстрый старт

### Требования

- Go 1.25+
- Node.js 18+
- Docker (для Postgres) — либо локальный PostgreSQL 17

### 1. База данных

```bash
cd backend
docker-compose up -d
```

Поднимет Postgres 17 на `localhost:5432` (db `chinese_learning`, user/password `postgres`/`postgres`).

### 2. Миграции

Миграции лежат в `backend/migrations` (формат [golang-migrate](https://github.com/golang-migrate/migrate)):

```bash
migrate -path backend/migrations \
  -database "postgres://postgres:postgres@localhost:5432/chinese_learning?sslmode=disable" \
  up
```

### 3. Backend

```bash
cd backend
go run ./cmd/server
```

Читает конфиг из `backend/.env`:

```
APP_NAME=
APP_PORT=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
JWT_SECRET=
```

Сервер поднимается на `:$APP_PORT` (по умолчанию 8080).

Наполнение курса «HSK 1» уроками 7–50 (тематическая практика/повторение поверх ручных уроков 1–6) — отдельный одноразовый сидер, не миграция:

```bash
go run ./cmd/seed
```

Destructive-но-идемпотентен: сначала удаляет уроки 7–50 (если уже были сгенерированы), потом строит заново — безопасно перезапускать при правке тем/грамматики в самом скрипте.

### 4. Frontend

```bash
cd frontend
npm install
npm run dev
```

Открывается на `http://localhost:5173`. Адрес API берётся из `frontend/.env`:

```
VITE_API_URL=http://localhost:8080/api/v1
```

## Backend

### Модули (`backend/internal/`)

| Модуль | Назначение |
|---|---|
| `auth` | регистрация/логин/logout, JWT (кука), middleware авторизации |
| `users` | модель пользователя |
| `courses` | список курсов, детали курса со списком уроков |
| `lessons` | уроки с полиморфными шагами (`word` / `quiz`) |
| `words` | словарь: поиск, фильтр по HSK |
| `savedwords` | сохранённые пользователем слова |
| `quizzes` | квизы (словарные и грамматические) с фильтром по HSK и серверной проверкой ответа |
| `progress` | прогресс по уроку (start/resume/step/complete) и по курсу |
| `learning` | spaced-repetition прогресс слов через тренировку «слово → перевод» (используется вкладками «Изучено»/«На изучении» в словаре) |
| `texts` | тексты для чтения разного уровня HSK |

### API

Базовый префикс `/api/v1` (плюс отдельный `GET /health` вне префикса).

**Публичные:**

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/logout`

**Требуют JWT (кука, выданная при логине):**

- `GET /me`
- `GET /words`, `GET /words/:id`
- `GET /words/saved`, `POST /words/:id/save`, `DELETE /words/:id/save`
- `GET /courses`, `GET /courses/:id`
- `GET /lessons/:id`
- `POST /lessons/:id/start`, `GET /lessons/:id/progress`, `POST /lessons/:id/step`, `POST /lessons/:id/complete`
- `GET /quizzes/` (опционально `?hsk=N`), `GET /quizzes/:id`, `POST /quizzes/`, `POST /quizzes/:id/check`
- `GET /learning/`, `GET /learning/learned`, `GET /learning/in-progress`, `POST /learning/:id/answer`
- `GET /texts/` (опционально `?hsk=N`), `GET /texts/:id`

### Миграции

19 миграций в `backend/migrations`. Основная схема: `users` → `courses` → `lessons` → `lesson_steps` → `words` → `saved_words` → `user_lesson_progress` → `quizzes` (+ `hsk_level`) → `user_course_progress` → `word_learning_progress` → `texts` (+ сид 10 текстов). Таблица `user_word_progress` (использовалась только удалённым модулем `review`) удалена.

### Тесты

```bash
cd backend
go test ./...
```

Покрыто: JWT-утилиты, `auth.Service`, `quizzes.Service`, `progress.Service` (в т.ч. что курс не обновляется при неудачном завершении урока), `learning.Service` (расписание повторений), `texts.Service`.

## Frontend

### Структура (`frontend/src/`)

- `layouts/` — `AuthLayout` (гостевые страницы), `DefaultLayout` (шапка + сайдбар для авторизованных)
- `pages/` — `HomePage`, `CoursesPage`, `CoursePage`, `LessonPage`, `VocabularyPage`, `TextsPage`, `TextPage`, `TestsPage`, `GrammarTestPage`, `WordTrainingPage`, `SettingsPage`, `LoginPage`
- `components/base` — `BaseButton`, `BaseCard`, `BaseInput`
- `components/layout` — `AppHeader`, `AppSidebar`
- `components/lesson` — `LessonStepRenderer`, `WordStep`, `QuizStep`
- `stores/` — Pinia: `auth`, `courses`, `lessons`, `vocabulary`, `savedWords`, `learning`, `texts`, `grammarTest`, `wordTraining`
- `services/` — axios-обёртки под каждый домен (`client.ts` — общий инстанс с `withCredentials`)
- `router/` — маршруты с гвардами `requiresAuth` / `guest`

### Роутинг

- `/login` — только для гостей
- `/app`, `/app/courses`, `/app/courses/:id`, `/app/lessons/:id`, `/app/vocabulary`, `/app/texts`, `/app/texts/:id`, `/app/tests`, `/app/tests/grammar`, `/app/tests/words`, `/app/settings` — только для авторизованных

### Сборка / типы

```bash
cd frontend
npm run build   # vue-tsc -b && vite build
```

## Известные ограничения

- нет тестов на фронтенде
- нет CI
