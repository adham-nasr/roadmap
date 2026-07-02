# Roadmaps Platform

> Serverless, event-driven ETL platform built with AWS SAM, Go, NestJS
> and LocalStack.

## Overview

This project consists of two major subsystems:

-   **NestJS Backend** deployed as a Lambda behind API Gateway.
-   **Event-driven ETL Pipeline** that synchronizes roadmap definitions
    from GitHub, transforms them into a normalized model, and loads them
    into MongoDB Atlas.

The platform is developed against **LocalStack** while remaining
compatible with AWS.

------------------------------------------------------------------------

# High-Level Architecture

``` mermaid
flowchart TB
    %% Force 3-column layout: Backend | ETL | Storage
    Backend ~~~ ETL ~~~ Storage

    %% =============================================
    %% COLUMN 1: BACKEND
    %% =============================================
    subgraph Backend["Backend System"]
        direction TB
        APIGW[API Gateway]
        Nest[NestJS Lambda]
        APIGW --> Nest
    end

    Client([Client Apps]) --> APIGW

    %% =============================================
    %% COLUMN 2: ETL + LONG BUS
    %% =============================================
    subgraph ETL["Event-Driven ETL Pipeline"]
        direction LR

        subgraph Pipeline[" "]
            direction TB
            Scheduler[⏱️ EventBridge Scheduler]
            Fetch[📥 Fetch Lambda]
            Queue[📨 SQS Queue]
            
            subgraph Workers["Parallel Download Workers"]
                direction LR
                D1[Worker 1] ~~~ D2[Worker 2] ~~~ D3[Worker 3]
            end
            
            Transform[🔄 Transform Lambda]
            Load[💾 Load Lambda]

            %% Explicit edges = perfect vertical stacking
            Scheduler --> Fetch --> Queue --> Workers
            Workers --> Transform
            Transform --> Load
        end

        %% The Long Bus - tall vertical orchestrator
        subgraph Bus["📡 EventBridge Bus"]
            direction TB
            B1[ ] ~~~ B2[ ] ~~~ B3[ ] ~~~ B4[ ]
        end

        %% Event handoffs with real event names
        Workers -.->|downloadComplete| Bus
        Bus -->|transformTrigger| Transform
        Transform -.->|transformComplete| Bus
        Bus -->|loadTrigger| Load
    end

    %% =============================================
    %% COLUMN 3: STORAGE
    %% =============================================
    subgraph Storage["Storage Layer"]
        direction TB
        Raw[(raw-bucket)]
        Output[(output-bucket)]
        DynamoDB[(DynamoDB Tables)]
        Raw ~~~ Output ~~~ DynamoDB
    end

    %% =============================================
    %% MINIMAL DATA FLOWS (Essential only)
    %% =============================================
    GitHub[(GitHub Repository)] -->|gets roadmap list| Fetch
    Fetch -->|enqueues URLs| Queue
    Fetch -->|creates run| DynamoDB
    Workers -->|store raw files| Raw
    Transform -->|write output| Output
    Load -->|loads| Mongo[(MongoDB Atlas)]

    %% =============================================
    %% STYLING
    %% =============================================
    classDef external fill:#FFEBEE,stroke:#C62828,stroke-width:2px,color:#000
    classDef backend fill:#E3F2FD,stroke:#1565C0,stroke-width:2px,color:#000
    classDef compute fill:#F3E5F5,stroke:#6A1B9A,stroke-width:2px,color:#000
    classDef storage fill:#E8F5E9,stroke:#2E7D32,stroke-width:2px,color:#000
    classDef trigger fill:#FFF3E0,stroke:#E65100,stroke-width:2px,color:#000
    classDef invisible fill:none,stroke:none,color:#fff

    class Client,GitHub,Mongo external
    class APIGW,Nest backend
    class Scheduler,Queue,Bus trigger
    class Fetch,Transform,Load,D1,D2,D3 compute
    class Raw,Output,DynamoDB storage
    class B1,B2,B3,B4 invisible
```

# Components

## Fetch Lambda

-   Scans GitHub repository
-   Detects changed roadmaps
-   Creates pipeline run
-   Enqueues download jobs in SQS

## Download Lambda

-   Triggered by SQS
-   Downloads roadmap assets
-   Stores raw content in S3
-   Updates DynamoDB state
-   Emits **RunComplete** on the EventBridge bus

## Transform Lambda

-   Reads raw artifacts
-   Parses and normalizes roadmap data
-   Produces `roadmaps.json` and `topics.json`
-   Emits **TransformComplete**

## Load Lambda

-   Reads generated artifacts
-   Upserts MongoDB Atlas collections

# AWS Resources

  Service         Purpose
  --------------- -----------------------------
  API Gateway     Backend entrypoint
  Lambda          Compute
  EventBridge     Event routing
  SQS             Download work queue
  S3              Raw & transformed artifacts
  DynamoDB        Pipeline state
  MongoDB Atlas   Application database
  CloudWatch      Logging

# DynamoDB Tables

  Table               Purpose
  ------------------- -------------------------------------
  SyncState           Roadmap fingerprint synchronization
  PipelineRun         Current ETL execution
  CompletedRoadmaps   Completed downloads for a run

# S3 Buckets

  Bucket          Contents
  --------------- -----------------------------------------
  raw-bucket      Downloaded roadmap assets
  output-bucket   Generated roadmaps.json and topics.json

# Event Flow

``` text
Scheduler
   │
Fetch
   │
SQS
   │
Download Workers
   │
EventBridge (RunComplete)
   │
Transform
   │
EventBridge (TransformComplete)
   │
Load
   │
MongoDB
```

# Local Development

``` bash
sam build
sam deploy
```

Local endpoint:

    http://localhost:4566/_aws/execute-api/<api-id>/Prod/

# Repository Layout

``` text
roadmap/          NestJS API
etl/              Go Lambdas
stateMachine/     Legacy workflow definitions
template.yaml     AWS SAM infrastructure
README.md
```

# Future Improvements

-   Parallel transform workers
-   Event versioning
-   CloudWatch dashboards
-   OpenTelemetry
-   CI/CD deployment
-   Infrastructure tests

# Technology Stack

-   Go
-   NestJS
-   AWS SAM
-   AWS Lambda
-   API Gateway
-   EventBridge
-   Amazon SQS
-   Amazon S3
-   DynamoDB
-   MongoDB Atlas
-   LocalStack
