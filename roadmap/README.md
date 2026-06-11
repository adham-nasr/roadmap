# Roadmap Backend Service

A simple NestJS-based roadmap API built with MongoDB and Mongoose. The service manages roadmaps, topics, and resources, and exposes REST endpoints for the frontend to consume.

## Overview

This backend provides core entities for a roadmap visualization and learning tracker:
- `Roadmap`
- `Topic`
- `Resource`

The application uses MongoDB at `mongodb://localhost/roadmap` by default.

> Note: The ERD also includes a `UserProfile` entity, but the current implementation exposes only `Roadmap`, `Topic`, and `Resource` APIs.

## Architecture

- NestJS application
- MongoDB + Mongoose for persistence
- Modules:
  - `RoadmapModule`
  - `TopicModule`
  - `ResourceModule`

## Data Model

### Roadmap
- `id` (string)
- `name` (string)
- `description` (string)

### Topic
- `id` (string)
- `repoTopicid` (string)
- `name` (string)
- `description` (string)
- `type` (string)
- `x_axis` (number)
- `y_axis` (number)
- `roadmap_id` (string) — relation to `Roadmap`
- `parent_topic_id` (string) — self-referencing topic parent

### Resource
- `id` (string)
- `link` (string)
- `type` (string)
- `topic_id` (string) — relation to `Topic`

### Relationships

- A `Roadmap` has many `Topic` records.
- A `Topic` can belong to one `Roadmap`.
- A `Topic` may have many child `Topic` entries via `parent_topic_id`.
- A `Topic` has many `Resource` records.
- A `Resource` belongs to one `Topic`.

## API Endpoints

Base URL: `http://localhost:3000`

### Roadmaps

- `POST /roadmaps`
  - Create a roadmap.
  - Body example:
    ```json
    {
      "name": "Frontend Roadmap",
      "description": "A path for learning UI and state management"
    }
    ```

- `GET /roadmaps`
  - Returns all roadmaps.

- `GET /roadmaps/:id`
  - Returns a roadmap by ID.

- `PATCH /roadmaps/:id`
  - Update roadmap fields.
  - Body example:
    ```json
    {
      "description": "Updated roadmap description"
    }
    ```

- `DELETE /roadmaps/:id`
  - Delete a roadmap.

- `GET /roadmaps/:id/topics`
  - Returns all topics for the selected roadmap.

### Topics

- `POST /topics`
  - Create a topic.
  - Body example:
    ```json
    {
      "name": "React Basics",
      "description": "Learn components and hooks",
      "type": "frontend",
      "x_axis": 120,
      "y_axis": 80,
      "repoTopicid": "react-basics",
      "roadmap_id": "648c5f4a2f4f9d0012345678",
      "parent_topic_id": null
    }
    ```

- `GET /topics`
  - Returns all topics.

- `GET /topics/:id`
  - Returns a topic by ID.

- `PATCH /topics/:id`
  - Update topic fields.
  - Body example:
    ```json
    {
      "name": "Advanced React",
      "x_axis": 220
    }
    ```

- `DELETE /topics/:id`
  - Delete a topic.

- `GET /topics/:id/resources`
  - Returns all resources attached to a topic.

### Resources

- `POST /resources`
  - Create a resource.
  - Body example:
    ```json
    {
      "link": "https://example.com/article",
      "type": "article",
      "topic_id": "648c5f4a2f4f9d0012345679"
    }
    ```

- `GET /resources`
  - Returns all resources.

- `GET /resources/:id`
  - Returns a resource by ID.

- `PATCH /resources/:id`
  - Update resource fields.
  - Body example:
    ```json
    {
      "type": "video"
    }
    ```

- `DELETE /resources/:id`
  - Delete a resource.

## Frontend Integration Notes

- Use `roadmap_id` when creating or filtering topics by roadmap.
- Use `topic_id` when creating resources or grouping resources for a topic.
- Use `parent_topic_id` to build nested topic trees or hierarchical relationships.
- `repoTopicid` can be used as a stable external topic identifier for UI state or repository mapping.
- The service currently does not expose authentication or user profile endpoints.

## Running the Project

```bash
cd roadmap
npm install
npm run start:dev
```

The server starts on port `3000` by default.

## Notes

- MongoDB must be running locally for the service to connect successfully.
- The source is in `src/` and includes three main modules: `roadmap`, `topic`, and `resource`.
- The ERD in `diagrams/roadmapErd.png` shows the expected model relationships.
